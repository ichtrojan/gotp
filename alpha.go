package gotp

import (
	"math/rand"
	"time"
)

func generateAlphaToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, length)

	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}

	return string(token)
}
