package storage

import (
	"context"
	"database/sql"
	wrapError "github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/errors"
	_ "github.com/jackc/pgx/v5/stdlib"
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

type AccrualDto struct {
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

func NewStorage(dbDSN string) Storage {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		log.Printf("error while starting db %s", err.Error())
		return nil
	}
	return &DBStorage{db}
}

type Storage interface {
	Register(ctx context.Context, registerData Auth, passwordHash string) (string, error)
	GetUserByLogin(ctx context.Context, authData Auth) (Auth, error)
	GetOrdersByUser(ctx context.Context, userID string) ([]Order, error)
	AddOrderForUser(ctx context.Context, externalOrderID string, userID string) error
	GetUserBalance(ctx context.Context, userID string) (UserBalance, error)
	AddWithdrawalForUser(ctx context.Context, userID string, withdrawal Withdrawal) error
	GetWithdrawalsForUser(ctx context.Context, userID string) ([]Withdrawal, error)
	GetOrdersInProgress(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, order AccrualDto) error
}

type DBStorage struct {
	db *sql.DB
}

func (s *DBStorage) Register(ctx context.Context, a Auth, p string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id FROM \"user\" WHERE \"login\" = $1", a.Login)
	var userID sql.NullString
	err := row.Scan(&userID)
	if err != nil && userID.Valid {
		log.Printf("Got error %s", err.Error())
		return "", err
	}
	if userID.Valid {
		log.Printf("user with login %s is existing", a.Login)
		return "", err
	}
	row = s.db.QueryRowContext(ctx, "INSERT INTO \"user\" (\"login\", password_hash) VALUES ($1, $2) RETURNING id", a.Login, p)
	if err := row.Scan(&userID); err != nil {
		log.Printf("error %s", err.Error())
		return "", err
	}
	log.Printf("userID %v", userID)
	if userID.Valid {
		userIDValue := userID.String
		log.Printf("new userID %s", userIDValue)
		return userIDValue, nil
	}
	return "", err
}

func (s *DBStorage) GetUserByLogin(ctx context.Context, a Auth) (Auth, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, login, password_hash FROM \"user\" WHERE login = $1", a.Login)
	var userData Auth
	err := row.Scan(&userData.UserID, &userData.Login, &userData.Password)
	if err != nil {
		log.Printf("could not get user data for login %s", a.Login)
		return userData, err
	}
	return userData, nil
}

func (s *DBStorage) AddOrderForUser(ctx context.Context, id string, u string) error {
	row := s.db.QueryRowContext(ctx, "SELECT user_id FROM \"order\" WHERE external_id = $1", id)
	var orderUserID sql.NullString
	err := row.Scan(&orderUserID)
	if err != nil && orderUserID.Valid {
		log.Printf("error while querying %s", err.Error())
		return err
	}
	if orderUserID.Valid {
		if orderUserID.String == u {
			log.Printf("same userID %s for orderID %s", u, id)
			return wrapError.NewOrderIsExistThisUserError("this order is exist the user", err)
		} else {
			log.Printf("another userID %s (instead of %s) for orderID %s", orderUserID.String, u, id)
			return wrapError.NewOrderIsExistAnotherUserError("this order is exist another user", err)
		}
	}
	log.Printf("order with id %v not found in database", id)
	row = s.db.QueryRowContext(
		ctx,
		"INSERT INTO \"order\" (user_id, status, external_id) VALUES ($1, $2, $3) RETURNING id",
		u, "NEW", id,
	)
	var orderID string
	err = row.Scan(&orderID)
	if err != nil {
		log.Printf("error while adding new order: %s", err.Error())
		return err
	}
	log.Printf("new order with id %s added", orderID)
	return nil
}

func (s *DBStorage) GetOrdersByUser(ctx context.Context, u string) ([]Order, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT external_id, status, amount, registered_at FROM \"order\" WHERE user_id = $1", u)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	orderList := make([]Order, 0)
	for rows.Next() {
		var orderFromDBVal OrderFromDB
		err := rows.Scan(&orderFromDBVal.Number, &orderFromDBVal.Status, &orderFromDBVal.Accrual, &orderFromDBVal.UploadedAt)
		if err != nil {
			log.Printf("error: %s", err.Error())
			return nil, err
		}
		order := Order{Number: orderFromDBVal.Number, Status: orderFromDBVal.Status, UploadedAt: orderFromDBVal.UploadedAt}
		if orderFromDBVal.Accrual.Valid {
			order.Accrual = orderFromDBVal.Accrual.Float64
		}
		orderList = append(orderList, order)
	}
	if err := rows.Err(); err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}
	return orderList, nil
}

func (s *DBStorage) GetUserBalance(ctx context.Context, u string) (UserBalance, error) {
	log.Printf("userID %s", u)
	sumOrdersRow := s.db.QueryRowContext(ctx, "SELECT sum(amount) FROM \"order\" WHERE user_id = $1", u)
	sumWithdrawalsRow := s.db.QueryRowContext(ctx, "SELECT sum(amount) FROM withdrawal WHERE user_id = $1", u)
	var sumOrders sql.NullFloat64
	var sumWithdrawals sql.NullFloat64
	err := sumOrdersRow.Scan(&sumOrders)
	if err != nil && sumOrders.Valid {
		log.Printf("error get sumOrders: %s", err.Error())
		return UserBalance{0, 0}, err
	}
	var resultBalance UserBalance
	if !sumOrders.Valid {
		log.Printf("resultBalance is empty")
		resultBalance.Orders = 0
	} else {
		resultBalance.Orders = sumOrders.Float64
	}
	err = sumWithdrawalsRow.Scan(&sumWithdrawals)

	if err != nil && sumWithdrawals.Valid {
		log.Printf("error get sumWithdrawals: %s", err.Error())
		return UserBalance{0, 0}, err
	}
	if !sumWithdrawals.Valid {
		log.Printf("resultBalance is empty")
		resultBalance.Withdrawn = 0
	} else {
		resultBalance.Withdrawn = sumWithdrawals.Float64
	}
	resultBalance.Orders -= resultBalance.Withdrawn
	log.Printf("balance %v", resultBalance)
	return resultBalance, nil
}

func (s *DBStorage) AddWithdrawalForUser(ctx context.Context, u string, w Withdrawal) error {
	userBalance, err := s.GetUserBalance(ctx, u)
	if err != nil {
		log.Printf("error while getting status %v", http.StatusInternalServerError)
		return err
	}
	if userBalance.Orders < w.Sum {
		log.Printf("got less bonus points %v than expected %v", userBalance.Orders, w.Sum)
		return err
	}
	var withdrawalID string
	row := s.db.QueryRowContext(ctx,
		"INSERT INTO withdrawal (user_id, amount, external_id) VALUES ($1, $2, $3) RETURNING id",
		u, w.Sum, w.Order,
	)
	err = row.Scan(&withdrawalID)
	if err != nil {
		log.Printf("error %s", err.Error())
		return err
	}
	log.Printf("new withdrawal %s", withdrawalID)
	return nil
}

func (s *DBStorage) GetWithdrawalsForUser(ctx context.Context, u string) ([]Withdrawal, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT external_id, amount, registered_at FROM withdrawal WHERE user_id = $1", u)
	if err != nil {
		log.Printf("error %s", err.Error())
		return make([]Withdrawal, 0), err
	}
	defer rows.Close()
	withdrawalList := make([]Withdrawal, 0)
	for rows.Next() {
		var withdrawal Withdrawal
		err = rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			log.Printf("error %s", err.Error())
			return make([]Withdrawal, 0), err
		}
		withdrawalList = append(withdrawalList, withdrawal)
	}
	if err := rows.Err(); err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}
	return withdrawalList, nil
}

func (s *DBStorage) GetOrdersInProgress(ctx context.Context) ([]Order, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT external_id, status, amount from \"order\" where status not in ('INVALID', 'PROCESSED')")

	if err != nil {
		log.Printf("error %s", err.Error())
		return make([]Order, 0), err
	}
	defer rows.Close()
	orderList := make([]Order, 0)
	for rows.Next() {
		var orderFromDBVal OrderFromDB
		err = rows.Scan(&orderFromDBVal.Number, &orderFromDBVal.Status, &orderFromDBVal.Accrual)
		if err != nil {
			log.Printf("error %s", err.Error())
			return make([]Order, 0), err
		}
		order := Order{Number: orderFromDBVal.Number, Status: orderFromDBVal.Status}
		if orderFromDBVal.Accrual.Valid {
			order.Accrual = orderFromDBVal.Accrual.Float64
		}
		orderList = append(orderList, order)
	}
	if err := rows.Err(); err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}
	log.Printf("orders %v", orderList)
	return orderList, nil
}

func (s *DBStorage) UpdateOrder(ctx context.Context, o AccrualDto) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error %s", err.Error())
		return err
	}
	stmt, err := tx.Prepare("UPDATE \"order\" SET status = $1, amount = $2 where external_id = $3")
	if err != nil {
		log.Printf("error %s", err.Error())
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(o.Status, o.Accrual, o.Order); err != nil {
		if err = tx.Rollback(); err != nil {
			log.Fatalf("insert to url, need rollback, %v", err)
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("unable to commit: %v", err)
		return err
	}
	return nil
}
