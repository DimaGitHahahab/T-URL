package service

import (
	"crypto/sha256"
	"encoding/base64"
)

const maxShortKeySize = 5

type URLGenerator struct {
}

func NewURLGenerator() *URLGenerator {
	return &URLGenerator{}
}

func (g URLGenerator) GenerateKey(long string) string {
	hasher := sha256.New()
	hasher.Write([]byte(long))
	hashed := hasher.Sum(nil)
	shortURL := base64.URLEncoding.EncodeToString(hashed)

	return shortURL[:maxShortKeySize]
}
