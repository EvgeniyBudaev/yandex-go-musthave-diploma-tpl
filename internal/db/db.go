package db

import (
	"fmt"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func RunMigrations(c *config.Config) error {
	if c.GetDBURI() == "" {
		return fmt.Errorf("got empty dbURI")
	}
	m, err := migrate.New(c.GetMigrateSourceURL(), c.GetDBURI())
	if err != nil {
		log.Printf("Got err %s", err.Error())
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Got err %s", err.Error())
		return err
	}
	return nil
}
