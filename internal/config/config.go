package config

import (
	"flag"
	"log"
	"os"
)

func Init() {
	var (
		ServerAddr     string
		DBURI          string
		AccrualSysAddr string
	)
	flag.StringVar(&ServerAddr, "a", "", "GopherMart server address")
	flag.StringVar(&DBURI, "d", "", "GopherMart database address")
	flag.StringVar(&AccrualSysAddr, "r", "", "Accrual system address")
	flag.Parse()

	ServerAddrEnv := os.Getenv("RUN_ADDRESS")
	if ServerAddrEnv != "" {
		ServerAddr = ServerAddrEnv
	}

	DBURIEnv := os.Getenv("DATABASE_URI")
	if DBURIEnv != "" {
		DBURI = DBURIEnv
	}

	AccrualSysAddrEnv := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if AccrualSysAddrEnv != "" {
		AccrualSysAddr = AccrualSysAddrEnv
	}

	log.Printf("Got ServerAddr %s, DBURI %s, AccrualSysAddr %s to run GopherMart", ServerAddr, DBURI, AccrualSysAddr)
}
