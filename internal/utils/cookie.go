package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"github.com/tank4gun/go-musthave-diploma-tpl/internal/config"
	"hash"
)

func GenerateCookie() hash.Hash {
	var CookieKey = []byte(config.GetSecretKeyToUserID())
	return hmac.New(sha256.New, CookieKey)
}
