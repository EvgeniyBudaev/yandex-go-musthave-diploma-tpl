package main

import (
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/db"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/server"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
)

func main() {
	c := config.Init()
	db.RunMigrations(c.GetDBURI(), c)
	storageForHandler := storage.GetStorage(c.GetDBURI())
	serverToRun := server.CreateServer(storageForHandler, c)
	serverToRun.ListenAndServe()
}
