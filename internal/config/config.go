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

func Init() *Config {
	var (
		serverAddr        string
		dBURI             string
		accrualSysAddr    string
		migrateSourceURL  string
		userCookie        string
		userID            string
		secretKeyToUserID string
	)
	flag.StringVar(&serverAddr, "a", "", "GopherMart server address")
	flag.StringVar(&dBURI, "d", "", "GopherMart database address")
	flag.StringVar(&accrualSysAddr, "r", "", "Accrual system address")
	flag.StringVar(&migrateSourceURL, "m", "file://internal/db/migrations",
		"Migrate source URL")
	flag.StringVar(&userCookie, "uc", "UserCookie", "User cookie")
	flag.StringVar(&userID, "uid", "UserID", "User ID")
	flag.StringVar(&secretKeyToUserID, "sk", "SecretKeyToUserID", "Secret key to user ID")
	flag.Parse()

	config := &Config{
		ServerAddr:        serverAddr,
		DBURI:             dBURI,
		AccrualSysAddr:    accrualSysAddr,
		MigrateSourceURL:  migrateSourceURL,
		UserCookie:        userCookie,
		UserID:            userID,
		SecretKeyToUserID: secretKeyToUserID,
	}

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
	return config
}

func (c *Config) GetServerAddr() string {
	return c.ServerAddr
}

func (c *Config) GetDBURI() string {
	return c.DBURI
}

func (c *Config) GetAccrualSysAddr() string {
	return c.AccrualSysAddr
}

func (c *Config) GetMigrateSourceURL() string {
	return c.MigrateSourceURL
}

func (c *Config) GetUserCookie() string {
	return c.UserCookie
}

func (c *Config) GetUserID() string {
	return c.UserID
}

func (c *Config) GetSecretKeyToUserID() string {
	return c.SecretKeyToUserID
}
