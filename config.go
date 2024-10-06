package gotp

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis *redis.Client
}
