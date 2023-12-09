package main

import (
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/clients"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/db"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/server"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/services"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"github.com/go-resty/resty/v2"
	"time"
)

func main() {
	const (
		countWorker             = 5
		retryTimeCheckNewOrders = 5 * time.Second
	)
	configInit := config.Init()
	db.RunMigrations(configInit)
	storageForHandler := storage.NewStorage(configInit.GetDBURI())
	service := services.NewUserService(storageForHandler)
	client := resty.New()
	accrual := clients.NewClientAccrual(client, configInit.AccrualSysAddr)
	ticker := time.NewTicker(retryTimeCheckNewOrders)
	worker := services.NewPoolWorker(accrual, service)
	go func() {
		worker.StarIntegration(countWorker, ticker)
	}()
	serverToRun := server.CreateServer(storageForHandler, configInit)
	serverToRun.ListenAndServe()
}
