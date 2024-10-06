package gotp

import (
	"math/rand"
	"time"
)

func generateAlphaNumericToken(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, length)

	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}

	return string(token)
}
