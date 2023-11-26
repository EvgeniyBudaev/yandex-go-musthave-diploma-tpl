package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/config"
	"log"
)

func RunMigrations(dbURI string) error {
	if dbURI == "" {
		return fmt.Errorf("got empty dbURI")
	}
	m, err := migrate.New(config.GetMigrateSourceURL(), dbURI)
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
