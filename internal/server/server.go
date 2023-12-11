package server

import (
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/routes"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"net/http"
)

func CreateServer(s storage.Storage, c *config.Config) *http.Server {
	router := routes.InitRouter(s, c)
	server := &http.Server{
		Addr:    c.GetServerAddr(),
		Handler: router,
	}
	return server
}
