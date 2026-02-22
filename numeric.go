package gotp

import (
	"crypto/rand"
	"math/big"
)

func generateNumericToken(length int) (string, error) {
	token := make([]byte, length)

	for i := range token {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		token[i] = '0' + byte(n.Int64())
	}

	return string(token), nil
}
