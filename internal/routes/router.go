package routes

import (
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/handlers"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
)

func InitRouter(s storage.Storage, c *config.Config) *chi.Mux {
	router := chi.NewRouter()
	handlerWithStorage := handlers.GetHandlerWithStorage(s)
	router.Use(handlers.CheckAuth)
	go handlerWithStorage.GetStatusesDaemon(c)
	router.Post("/api/user/register", handlerWithStorage.Register)
	router.Post("/api/user/login", handlerWithStorage.Login)
	router.Post("/api/user/orders", handlerWithStorage.AddOrder)
	router.Get("/api/user/orders", handlerWithStorage.GetOrders)
	router.Get("/api/user/balance", handlerWithStorage.GetBalance)
	router.Post("/api/user/balance/withdraw", handlerWithStorage.AddWithdrawal)
	router.Get("/api/user/withdrawals", handlerWithStorage.GetWithdrawals)
	return router
}
