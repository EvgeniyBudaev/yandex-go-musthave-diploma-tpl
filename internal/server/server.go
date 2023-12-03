package server

import (
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/routes"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"net/http"
)

func CreateServer(s storage.Storage) *http.Server {
	router := routes.InitRouter(s)
	server := &http.Server{
		Addr:    config.GetServerAddr(),
		Handler: router,
	}
	return server
}
