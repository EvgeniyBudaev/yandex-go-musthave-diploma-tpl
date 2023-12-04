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
	SecretKeyToUserID string `env:"SECRET_KEY_TO_USER_ID"`
}

var config Config

func Init() *Config {
	flag.StringVar(&config.ServerAddr, "a", "", "GopherMart server address")
	flag.StringVar(&config.DBURI, "d", "", "GopherMart database address")
	flag.StringVar(&config.AccrualSysAddr, "r", "", "Accrual system address")
	flag.StringVar(&config.MigrateSourceURL, "m", "file://internal/db/migrations",
		"Migrate source URL")
	flag.StringVar(&config.UserCookie, "uc", "UserCookie", "User cookie")
	flag.StringVar(&config.SecretKeyToUserID, "sk", "SecretKeyToUserID", "Secret key to user ID")
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
	SecretKeyToUserIDEnv := os.Getenv("SECRET_KEY_TO_USER_ID")
	if SecretKeyToUserIDEnv != "" {
		config.SecretKeyToUserID = SecretKeyToUserIDEnv
	}

	config := Config{
		ServerAddr:        config.ServerAddr,
		DBURI:             config.DBURI,
		AccrualSysAddr:    config.AccrualSysAddr,
		MigrateSourceURL:  config.MigrateSourceURL,
		UserCookie:        config.UserCookie,
		SecretKeyToUserID: config.SecretKeyToUserID,
	}

	log.Printf("Got ServerAddr %s, DBURI %s, AccrualSysAddr %s to run GopherMart", config.ServerAddr,
		config.DBURI, config.AccrualSysAddr)
	return &config
}

func (c *Config) GetServerAddr() string {
	return config.ServerAddr
}

func (c *Config) GetDBURI() string {
	return config.DBURI
}

func (c *Config) GetAccrualSysAddr() string {
	return config.AccrualSysAddr
}

func (c *Config) GetMigrateSourceURL() string {
	return config.MigrateSourceURL
}

func (c *Config) GetUserCookie() string {
	return config.UserCookie
}

func (c *Config) GetSecretKeyToUserID() string {
	return config.SecretKeyToUserID
}
