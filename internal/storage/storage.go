package storage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	UserID   string `json:"userID,omitempty"`
}

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type OrderFromDB struct {
	Number     string
	Status     string
	Accrual    sql.NullFloat64
	UploadedAt time.Time
}

type OrderFromBlackBox struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type UserBalance struct {
	Orders    float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

func GetStorage(dbDSN string) Storage {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		log.Printf("Got error while starting db %s", err.Error())
		return nil
	}
	return &DBStorage{db}
}

type Storage interface {
	Register(registerData Auth) (string, int)
	GetUserByLogin(authData Auth) (Auth, int)
	GetOrdersByUser(userID string) ([]Order, int)
	AddOrderForUser(externalOrderID string, userID string) int
	GetUserBalance(userID string) (UserBalance, int)
	AddWithdrawalForUser(userID string, withdrawal Withdrawal) int
	GetWithdrawalsForUser(userID string) ([]Withdrawal, int)
	GetOrdersInProgress() ([]Order, int)
	UpdateOrder(order OrderFromBlackBox) int
}

type DBStorage struct {
	db *sql.DB
}

func (s *DBStorage) Register(a Auth) (string, int) {
	row := s.db.QueryRow("SELECT id FROM \"user\" WHERE \"login\" = $1", a.Login)
	var userID sql.NullString
	err := row.Scan(&userID)
	if err != nil && userID.Valid {
		log.Printf("Got error %s", err.Error())
		return "", http.StatusInternalServerError
	}
	if userID.Valid {
		log.Printf("Got existing user with login %s", a.Login)
		return "", http.StatusFailedDependency
	}
	h := sha256.New()
	h.Write([]byte(a.Password))
	passwordHash := hex.EncodeToString(h.Sum(nil))
	log.Printf("Got password hash %s", passwordHash)
	row = s.db.QueryRow("INSERT INTO \"user\" (\"login\", password_hash) VALUES ($1, $2) RETURNING id", a.Login, passwordHash)
	if err := row.Scan(&userID); err != nil {
		log.Printf("Error %s", err.Error())
		return "", http.StatusInternalServerError
	} else {
		log.Printf("Got userID %v", userID)
		if userID.Valid {
			userIDValue := userID.String
			log.Printf("Got new userID %s", userIDValue)
			return userIDValue, http.StatusOK
		}
	}
	return "", http.StatusInternalServerError
}

func (s *DBStorage) GetUserByLogin(a Auth) (Auth, int) {
	row := s.db.QueryRow("SELECT id, login, password_hash FROM \"user\" WHERE login = $1", a.Login)
	var userData Auth
	err := row.Scan(&userData.UserID, &userData.Login, &userData.Password)
	if err != nil {
		log.Printf("Could not get user data for login %s", a.Login)
		return userData, http.StatusUnauthorized
	}
	return userData, http.StatusOK
}

func (s *DBStorage) AddOrderForUser(id string, u string) int {
	row := s.db.QueryRow("SELECT user_id FROM \"order\" WHERE external_id = $1", id)
	var orderUserID sql.NullString
	err := row.Scan(&orderUserID)
	if err != nil && orderUserID.Valid {
		log.Printf("Got error while querying %s", err.Error())
		return http.StatusInternalServerError
	}
	if orderUserID.Valid {
		if orderUserID.String == u {
			log.Printf("Got same userID %s for orderID %s", u, id)
			return http.StatusOK
		} else {
			log.Printf("Got another userID %s (instead of %s) for orderID %s", orderUserID.String, u, id)
			return http.StatusConflict
		}
	}
	log.Printf("Order with id %v not found in DB, should add it", id)
	row = s.db.QueryRow(
		"INSERT INTO \"order\" (user_id, status, external_id) VALUES ($1, $2, $3) RETURNING id",
		u, "NEW", id,
	)
	var orderID string
	err = row.Scan(&orderID)
	if err != nil {
		log.Printf("Smth went wrong while adding new order: %s", err.Error())
		return http.StatusInternalServerError
	}
	log.Printf("New order with id %s added", orderID)
	return http.StatusAccepted
}

func (s *DBStorage) GetOrdersByUser(u string) ([]Order, int) {
	rows, err := s.db.Query("SELECT external_id, status, amount, registered_at FROM \"order\" WHERE user_id = $1", u)
	if err != nil {
		log.Printf("Got error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}
	defer rows.Close()
	orders := make([]Order, 0)
	for rows.Next() {
		var orderFromDBVal OrderFromDB
		err := rows.Scan(&orderFromDBVal.Number, &orderFromDBVal.Status, &orderFromDBVal.Accrual, &orderFromDBVal.UploadedAt)
		if err != nil {
			log.Printf("Got error: %s", err.Error())
			return nil, http.StatusInternalServerError
		}
		order := Order{Number: orderFromDBVal.Number, Status: orderFromDBVal.Status, UploadedAt: orderFromDBVal.UploadedAt}
		if orderFromDBVal.Accrual.Valid {
			order.Accrual = orderFromDBVal.Accrual.Float64
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Got error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}
	return orders, http.StatusOK
}

func (s *DBStorage) GetUserBalance(u string) (UserBalance, int) {
	log.Printf("Got userID %s", u)
	sumOrdersRow := s.db.QueryRow("SELECT sum(amount) FROM \"order\" WHERE user_id = $1", u)
	sumWithdrawalsRow := s.db.QueryRow("SELECT sum(amount) FROM withdrawal WHERE user_id = $1", u)
	var sumOrders sql.NullFloat64
	var sumWithdrawals sql.NullFloat64
	err := sumOrdersRow.Scan(&sumOrders)
	if err != nil && sumOrders.Valid {
		log.Printf("Could not get sumOrders: %s", err.Error())
		return UserBalance{0, 0}, http.StatusInternalServerError
	}
	var resultBalance UserBalance
	if !sumOrders.Valid {
		log.Printf("Got empty resultBalance")
		resultBalance.Orders = 0
	} else {
		resultBalance.Orders = sumOrders.Float64
	}
	err = sumWithdrawalsRow.Scan(&sumWithdrawals)

	if err != nil && sumWithdrawals.Valid {
		log.Printf("Could not get sumWithdrawals: %s", err.Error())
		return UserBalance{0, 0}, http.StatusInternalServerError
	}
	if !sumWithdrawals.Valid {
		log.Printf("Got empty resultBalance")
		resultBalance.Withdrawn = 0
	} else {
		resultBalance.Withdrawn = sumWithdrawals.Float64
	}
	resultBalance.Orders -= resultBalance.Withdrawn
	log.Printf("Got balance %v", resultBalance)
	return resultBalance, http.StatusOK
}

func (s *DBStorage) AddWithdrawalForUser(u string, w Withdrawal) int {
	userBalance, errCode := s.GetUserBalance(u)
	if errCode != http.StatusOK {
		log.Printf("Got error while getting status %v", errCode)
		return errCode
	}
	if userBalance.Orders < w.Sum {
		log.Printf("Got less bonus points %v than expected %v", userBalance.Orders, w.Sum)
		return http.StatusPaymentRequired
	}
	var withdrawalID string
	row := s.db.QueryRow(
		"INSERT INTO withdrawal (user_id, amount, external_id) VALUES ($1, $2, $3) RETURNING id",
		u, w.Sum, w.Order,
	)
	err := row.Scan(&withdrawalID)
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return http.StatusInternalServerError
	}
	log.Printf("Got new withdrawal %s", withdrawalID)
	return http.StatusOK
}

func (s *DBStorage) GetWithdrawalsForUser(u string) ([]Withdrawal, int) {
	rows, err := s.db.Query("SELECT external_id, amount, registered_at FROM withdrawal WHERE user_id = $1", u)
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return make([]Withdrawal, 0), http.StatusInternalServerError
	}
	defer rows.Close()
	withdrawals := make([]Withdrawal, 0)
	for rows.Next() {
		var withdrawal Withdrawal
		err = rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			log.Printf("Got error %s", err.Error())
			return make([]Withdrawal, 0), http.StatusInternalServerError
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Got error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}
	return withdrawals, http.StatusOK
}

func (s *DBStorage) GetOrdersInProgress() ([]Order, int) {
	rows, err := s.db.Query("SELECT external_id, status, amount from \"order\" where status not in ('INVALID', 'PROCESSED')")

	if err != nil {
		log.Printf("Got error %s", err.Error())
		return make([]Order, 0), http.StatusInternalServerError
	}
	defer rows.Close()
	orders := make([]Order, 0)
	for rows.Next() {
		var orderFromDBVal OrderFromDB
		err = rows.Scan(&orderFromDBVal.Number, &orderFromDBVal.Status, &orderFromDBVal.Accrual)
		if err != nil {
			log.Printf("Got error %s", err.Error())
			return make([]Order, 0), http.StatusInternalServerError
		}
		order := Order{Number: orderFromDBVal.Number, Status: orderFromDBVal.Status}
		if orderFromDBVal.Accrual.Valid {
			order.Accrual = orderFromDBVal.Accrual.Float64
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Got error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}
	log.Printf("Got orders %v", orders)
	return orders, http.StatusOK
}

func (s *DBStorage) UpdateOrder(order OrderFromBlackBox) int {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return http.StatusInternalServerError
	}
	URLstmt, err := tx.Prepare("UPDATE \"order\" SET status = $1, amount = $2 where external_id = $3")
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return http.StatusInternalServerError
	}
	defer URLstmt.Close()
	if _, err := URLstmt.Exec(order.Status, order.Accrual, order.Order); err != nil {
		if err = tx.Rollback(); err != nil {
			log.Fatalf("Insert to url, need rollback, %v", err)
			return http.StatusInternalServerError
		}
		return http.StatusInternalServerError
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Unable to commit: %v", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
