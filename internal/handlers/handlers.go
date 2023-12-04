package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type userCtxName string

type orderStatus string

const (
	orderStatusInvalid   orderStatus = "INVALID"
	orderStatusProcessed orderStatus = "PROCESSED"
)

var UserID = userCtxName(config.GetUserID())

type HandlerWithStorage struct {
	storage         storage.Storage
	client          http.Client
	ordersToProcess chan string
}

func GetHandlerWithStorage(storage storage.Storage) *HandlerWithStorage {
	return &HandlerWithStorage{storage: storage, client: http.Client{}, ordersToProcess: make(chan string, 10)}
}

func ValidateOrder(order string) (uint, int, error) {
	orderNum, err := strconv.Atoi(order)
	if err != nil || orderNum < 0 {
		return 0, http.StatusBadRequest, err
	}
	sum := 0
	if len(order)%2 == 0 {
		for i, num := range []rune(order) {
			if i%2 == 0 {
				if int(num-'0')*2 > 9 {
					sum += int(num-'0')*2 - 9
				} else {
					sum += int(num-'0') * 2
				}
			} else {
				sum += int(num - '0')
			}
		}
	} else {
		for i, num := range []rune(order) {
			if i%2 == 1 {
				if int(num-'0')*2 > 9 {
					sum += int(num-'0')*2 - 9
				} else {
					sum += int(num-'0') * 2
				}
			} else {
				sum += int(num - '0')
			}
		}
	}
	if sum%10 == 0 {
		return uint(orderNum), http.StatusOK, nil
	} else {
		return 0, http.StatusUnprocessableEntity, err
	}
}

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/register" || r.URL.Path == "/api/user/login" {
			log.Printf("get %s url, skip check", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := (*r).Cookie(config.GetUserCookie())
		if cookie != nil && err != nil {
			log.Println(err.Error())
			http.Error(w, "could not auth user", http.StatusUnauthorized)
			return
		}
		if cookie == nil {
			log.Println("null value in Cookie for UserID")
			http.Error(w, "error auth user", http.StatusUnauthorized)
			return
		}
		data, err := hex.DecodeString(cookie.Value)
		log.Printf("cookie: %s", cookie.Value)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "error auth user", http.StatusUnauthorized)
			return
		}
		h := utils.GenerateCookie()
		h.Write(data[:36])
		sign := h.Sum(nil)
		if hmac.Equal(sign, data[36:]) {
			ctx := context.WithValue(r.Context(), UserID, string(data[:36]))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		log.Println("get not equal sign for UserID")
		http.Error(w, "error auth user", http.StatusUnauthorized)
	})
}

func (strg *HandlerWithStorage) GetStatusesDaemon() {
	ctx := context.Background()
	for orderNumber := range strg.ordersToProcess {
		log.Printf("order %s to process", orderNumber)
		response, err := strg.client.Get(config.GetAccrualSysAddr() + "/api/orders/" + orderNumber)
		if err != nil {
			log.Printf("error %s", err.Error())
			continue
		}
		if response.StatusCode == http.StatusOK {
			var newOrder storage.OrderFromBlackBox
			data, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("error %s", err.Error())
				continue
			}
			err = json.Unmarshal(data, &newOrder)
			if err != nil {
				log.Printf("error %s", err.Error())
				continue
			}
			newOrder.Order = orderNumber
			err = strg.storage.UpdateOrder(ctx, newOrder)
			if err != nil {
				return
			}
			if newOrder.Status != string(orderStatusInvalid) && newOrder.Status != string(orderStatusProcessed) {
				go func(orderNumber string) {
					strg.ordersToProcess <- orderNumber
				}(orderNumber)
			}
			response.Body.Close()
		} else {
			if response.StatusCode == http.StatusTooManyRequests {
				log.Printf("get 429 StatusTooManyRequests, need to sleep a bit")
				time.Sleep(1 * time.Second)
			}
			log.Printf("bad status code %v for order %s", response.StatusCode, orderNumber)
			go func(orderNumber string) {
				strg.ordersToProcess <- orderNumber
			}(orderNumber)
		}
	}
	close(strg.ordersToProcess)
}

func (strg *HandlerWithStorage) Register(w http.ResponseWriter, r *http.Request) {
	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error while reading body: %s", err.Error())
		http.Error(w, "error while reading body", http.StatusBadRequest)
		return
	}
	var authData storage.Auth
	log.Printf("data: %v", r.Body)
	err = json.Unmarshal(jsonBody, &authData)
	if err != nil {
		log.Printf("error unmarshal body: %s", err.Error())
		http.Error(w, "error unmarshal body", http.StatusBadRequest)
		return
	}
	h := sha256.New()
	h.Write([]byte(authData.Password))
	passwordHash := hex.EncodeToString(h.Sum(nil))
	userID, err := strg.storage.Register(r.Context(), authData, passwordHash)
	if err != nil {
		log.Println("error register user")
		http.Error(w, "error register user", http.StatusInternalServerError)
		return
	}
	h = utils.GenerateCookie()
	h.Write([]byte(userID))
	sign := h.Sum(nil)
	newCookie := http.Cookie{Name: config.GetUserCookie(), Value: hex.EncodeToString(append([]byte(userID)[:], sign[:]...))}
	log.Printf("sign %v, cookie %v", sign, []byte(userID))
	http.SetCookie(w, &newCookie)
	w.WriteHeader(http.StatusOK)
	w.Write(make([]byte, 0))
}

func (strg *HandlerWithStorage) Login(w http.ResponseWriter, r *http.Request) {
	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error while reading body: %s", err.Error())
		http.Error(w, "error while reading body", http.StatusBadRequest)
		return
	}
	var authData storage.Auth
	err = json.Unmarshal(jsonData, &authData)
	if err != nil {
		log.Printf("error unmarshal body: %s", err.Error())
		http.Error(w, "error unmarshal body", http.StatusBadRequest)
		return
	}
	userData, err := strg.storage.GetUserByLogin(r.Context(), authData)
	if err != nil {
		log.Println("error get user by login")
		http.Error(w, "error get user by login", http.StatusUnauthorized)
		return
	}
	h := sha256.New()
	h.Write([]byte(authData.Password))
	pswdHash := hex.EncodeToString(h.Sum(nil))
	if pswdHash == userData.Password {
		h := utils.GenerateCookie()
		h.Write([]byte(userData.UserID))
		sign := h.Sum(nil)
		newCookie := http.Cookie{Name: config.GetUserCookie(), Value: hex.EncodeToString(append([]byte(userData.UserID)[:], sign[:]...))}
		http.SetCookie(w, &newCookie)
		w.WriteHeader(http.StatusOK)
		w.Write(make([]byte, 0))
	} else {
		log.Println("error login-password pair")
		http.Error(w, "error login-password pair", http.StatusUnauthorized)
	}
}

func (strg *HandlerWithStorage) AddOrder(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("error while reading body: ")
	}

	_, errCode, _ := ValidateOrder(string(data))
	if errCode != http.StatusOK {
		log.Printf("bad order number %s", data)
		http.Error(w, "bad order number", errCode)
		return
	}
	userID := r.Context().Value(UserID).(string)
	founded, err := strg.storage.AddOrderForUser(r.Context(), string(data), userID)
	if founded && err != nil {
		http.Error(w, "error add order into db", http.StatusConflict)
		return
	}
	if founded {
		http.Error(w, "error add order into db", http.StatusOK)
		return
	}
	if !founded && err != nil {
		http.Error(w, "error add order into db", http.StatusInternalServerError)
		return
	}
	if !founded {
		go func(orderNumber string) {
			strg.ordersToProcess <- orderNumber
		}(string(data))
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(make([]byte, 0))
}

func (strg *HandlerWithStorage) GetOrders(w http.ResponseWriter, r *http.Request) {
	log.Println("Got GetOrders request")
	userID := r.Context().Value(UserID).(string)
	orders, err := strg.storage.GetOrdersByUser(r.Context(), userID)
	if err != nil {
		log.Printf("error %v", err)
		http.Error(w, "bad status code", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		log.Println("orders is empty")
		http.Error(w, "no orders for this user", http.StatusNoContent)
		return
	}
	log.Printf("orders %v", orders)
	ordersMarshalled, err := json.Marshal(orders)
	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, "error while marshalling", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ordersMarshalled)
}

func (strg *HandlerWithStorage) GetBalance(w http.ResponseWriter, r *http.Request) {
	userBalance, err := strg.storage.GetUserBalance(r.Context(), r.Context().Value(UserID).(string))
	if err != nil {
		http.Error(w, "error get user balance", http.StatusInternalServerError)
		return
	}
	userBalanceMarshalled, err := json.Marshal(userBalance)
	if err != nil {
		log.Printf("error while marshalling: %s", err.Error())
		http.Error(w, "error while marshalling", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userBalanceMarshalled)
}

func (strg *HandlerWithStorage) AddWithdrawal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserID).(string)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error %s", err.Error())
		http.Error(w, "error while getting data", http.StatusInternalServerError)
		return
	}
	var withdrawal storage.Withdrawal
	err = json.Unmarshal(data, &withdrawal)
	if err != nil {
		log.Printf("error %s", err.Error())
		http.Error(w, "error while getting data", http.StatusInternalServerError)
		return
	}
	_, errCode, _ := ValidateOrder(withdrawal.Order)
	if errCode != http.StatusOK {
		log.Printf("bad order number %s", withdrawal.Order)
		http.Error(w, "bad order number", errCode)
		return
	}
	err = strg.storage.AddWithdrawalForUser(r.Context(), userID, withdrawal)
	if err != nil {
		log.Printf("errorCode %v", http.StatusInternalServerError)
		http.Error(w, "error from AddWithdrawalForUser", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(make([]byte, 0))
}

func (strg *HandlerWithStorage) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserID).(string)
	withdrawals, err := strg.storage.GetWithdrawalsForUser(r.Context(), userID)
	if err != nil {
		log.Printf("errCode %v", http.StatusInternalServerError)
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		http.Error(w, "no withdrawals for this user", http.StatusNoContent)
		return
	}
	withdrawalsMarshalled, err := json.Marshal(withdrawals)
	if err != nil {
		log.Printf("error %s", err.Error())
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(withdrawalsMarshalled)
}
