package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	ServerAddr     string `env:"RUN_ADDRESS"`
	DBURI          string `env:"DATABASE_URI"`
	AccrualSysAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

var config Config

func Init() {
	flag.StringVar(&config.ServerAddr, "a", "", "GopherMart server address")
	flag.StringVar(&config.DBURI, "d", "", "GopherMart database address")
	flag.StringVar(&config.AccrualSysAddr, "r", "", "Accrual system address")
	flag.Parse()

	ServerAddrEnv := os.Getenv("RUN_ADDRESS")
	if ServerAddrEnv != "" {
		config.ServerAddr = ServerAddrEnv
	}

	DBURIEnv := os.Getenv("DATABASE_URI")
	if DBURIEnv != "" {
		config.DBURI = DBURIEnv
	}

	AccrualSysAddrEnv := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if AccrualSysAddrEnv != "" {
		config.AccrualSysAddr = AccrualSysAddrEnv
	}

	log.Printf("Got ServerAddr %s, DBURI %s, AccrualSysAddr %s to run GopherMart", &config.ServerAddr,
		&config.DBURI, &config.AccrualSysAddr)
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
