package main

import (
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/db"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/server"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
)

func main() {
	configInit := config.Init()
	db.RunMigrations(configInit)
	storageForHandler := storage.NewStorage(configInit.GetDBURI())
	serverToRun := server.CreateServer(storageForHandler, configInit)
	serverToRun.ListenAndServe()
}
