package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"hash"
)

func GenerateCookie(c *config.Config) hash.Hash {
	var CookieKey = []byte(c.GetSecretKeyToUserID())
	return hmac.New(sha256.New, CookieKey)
}
