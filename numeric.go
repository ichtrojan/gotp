package gotp

import (
	"fmt"
	"math/rand"
	"time"
)

func generateNumericToken(length int) (token string) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		token += fmt.Sprintf("%d", rand.Intn(10))
	}

	return token
}
