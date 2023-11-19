package server

import (
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/config"
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/routes"
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/storage"
	"net/http"
)

func CreateServer(s storage.Storage) *http.Server {
	router := routes.InitRouter(s)
	server := &http.Server{
		Addr:    config.ServerAddr,
		Handler: router,
	}
	return server
}
