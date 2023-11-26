package config

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	ServerAddr     string `env:"SERVER_ADDRESS"`
	DBURI          string `env:"DATABASE_URI"`
	AccrualSysAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

var config Config

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error while Init. Error in Load")
		return
	}
	err := envconfig.Process("MYAPP", &config)
	if err != nil {
		log.Fatal("error while Init. Error in Process")
		return
	}

	flag.StringVar(&config.ServerAddr, "a", "", "GopherMart server address")
	flag.StringVar(&config.DBURI, "d", "", "GopherMart database address")
	flag.StringVar(&config.AccrualSysAddr, "r", "", "Accrual system address")
	flag.Parse()

	log.Printf("Got ServerAddr %s, DBURI %s, AccrualSysAddr %s to run GopherMart",
		&config.ServerAddr, &config.DBURI, &config.AccrualSysAddr)
}

func GetServerAddr() string {
	return config.ServerAddr
}

func GetDBURI() string {
	return config.DBURI
}

func GetAccrualSysAddr() string {
	return config.AccrualSysAddr
}
