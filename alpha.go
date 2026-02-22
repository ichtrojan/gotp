package gotp

import (
	"crypto/rand"
	"math/big"
)

func generateAlphaToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	token := make([]byte, length)

	for i := range token {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		token[i] = charset[n.Int64()]
	}

	return string(token), nil
}
