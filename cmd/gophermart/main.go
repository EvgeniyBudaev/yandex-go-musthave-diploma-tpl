package main

import (
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/config"
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/db"
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/server"
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/storage"
)

func main() {
	config.Init()
	db.RunMigrations(config.GetDBURI())
	storageForHandler := storage.GetStorage(config.GetDBURI())
	serverToRun := server.CreateServer(storageForHandler)
	serverToRun.ListenAndServe()
}
