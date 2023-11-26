package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	ServerAddr        string `env:"RUN_ADDRESS"`
	DBURI             string `env:"DATABASE_URI"`
	AccrualSysAddr    string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	MigrateSourceURL  string `env:"MIGRATE_SOURCE_URL"`
	UserCookie        string `env:"USER_COOKIE"`
	UserID            string `env:"USERID"`
	SecretKeyToUserID string `env:"SECRET_KEY_TO_USER_ID"`
}

var config Config

func Init() {
	flag.StringVar(&config.ServerAddr, "a", "", "GopherMart server address")
	flag.StringVar(&config.DBURI, "d", "", "GopherMart database address")
	flag.StringVar(&config.AccrualSysAddr, "r", "", "Accrual system address")
	flag.StringVar(&config.MigrateSourceURL, "m", "file://internal/db/migrations",
		"Migrate source URL")
	flag.StringVar(&config.UserCookie, "r", "UserCookie", "User cookie")
	flag.StringVar(&config.UserID, "r", "UserID", "User ID")
	flag.StringVar(&config.SecretKeyToUserID, "r", "SecretKeyToUserID", "Secret key to user ID")
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
	MigrateSourceURLEnv := os.Getenv("MIGRATE_SOURCE_URL")
	if MigrateSourceURLEnv != "" {
		config.MigrateSourceURL = MigrateSourceURLEnv
	}
	UserCookieEnv := os.Getenv("USER_COOKIE")
	if UserCookieEnv != "" {
		config.UserCookie = UserCookieEnv
	}
	UserIDEnv := os.Getenv("USERID")
	if UserIDEnv != "" {
		config.UserID = UserIDEnv
	}
	SecretKeyToUserIDEnv := os.Getenv("SECRET_KEY_TO_USER_ID")
	if SecretKeyToUserIDEnv != "" {
		config.SecretKeyToUserID = SecretKeyToUserIDEnv
	}

	log.Printf("Got ServerAddr %s, DBURI %s, AccrualSysAddr %s to run GopherMart", config.ServerAddr,
		config.DBURI, config.AccrualSysAddr)
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

func GetMigrateSourceURL() string {
	return config.MigrateSourceURL
}

func GetUserCookie() string {
	return config.UserCookie
}

func GetUserID() string {
	return config.UserID
}

func GetSecretKeyToUserID() string {
	return config.SecretKeyToUserID
}
