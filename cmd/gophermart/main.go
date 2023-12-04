package main

import (
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/db"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/server"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
)

func main() {
	config.Init()
	db.RunMigrations(config.GetDBURI())
	storageForHandler := storage.NewStorage(config.GetDBURI())
	serverToRun := server.CreateServer(storageForHandler)
	serverToRun.ListenAndServe()
}
