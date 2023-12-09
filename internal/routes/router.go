package routes

import (
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/clients"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/handlers"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/services"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"time"
)

func InitRouter(s storage.Storage, cfg *config.Config) *chi.Mux {
	const (
		countWorker             = 5
		retryTimeCheckNewOrders = 5 * time.Second
	)
	router := chi.NewRouter()
	handlerWithStorage := handlers.NewHandlerWithStorage(s, cfg)
	router.Use(handlerWithStorage.CheckAuth)

	service := services.NewUserService(s)
	client := resty.New()
	accrual := clients.NewClientAccrual(client, cfg.AccrualSysAddr)
	ticker := time.NewTicker(retryTimeCheckNewOrders)
	worker := services.NewPoolWorker(accrual, service)
	go func() {
		worker.StarIntegration(countWorker, ticker)
	}()

	//go handlerWithStorage.GetStatusesDaemon(cfg)
	router.Post("/api/user/register", handlerWithStorage.Register)
	router.Post("/api/user/login", handlerWithStorage.Login)
	router.Post("/api/user/orders", handlerWithStorage.AddOrder)
	router.Get("/api/user/orders", handlerWithStorage.GetOrders)
	router.Get("/api/user/balance", handlerWithStorage.GetBalance)
	router.Post("/api/user/balance/withdraw", handlerWithStorage.AddWithdrawal)
	router.Get("/api/user/withdrawals", handlerWithStorage.GetWithdrawals)
	return router
}
