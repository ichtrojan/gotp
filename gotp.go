package gotp

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Generate struct {
	Format     Format
	Length     int
	Identifier string
	Expires    time.Duration
}

type Verify struct {
	Token      string
	Identifier string
}

const prefix = "gotp_"

func New(c Config) (Config, error) {
	err := c.Redis.Ping(context.Background()).Err()
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func (c Config) Generate(payload Generate) (string, error) {
	ctx := context.Background()
	if payload.Length < 4 || payload.Length > 10 {
		return "", errors.New("length must be between 4 and 10")
	}

	if payload.Identifier == "" {
		return "", errors.New("identifier is required")
	}

	var token string
	var err error

	switch payload.Format {
	case Alpha:
		token, err = generateAlphaToken(payload.Length)
	case AlphaNumeric:
		token, err = generateAlphaNumericToken(payload.Length)
	case Numeric:
		token, err = generateNumericToken(payload.Length)
	default:
		return "", errors.New("invalid format")
	}

	if err != nil {
		return "", err
	}

	err = c.Redis.Set(ctx, prefix+payload.Identifier, token, payload.Expires).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c Config) Verify(payload Verify) (bool, error) {
	ctx := context.Background()
	storedToken, err := c.Redis.Get(ctx, prefix+payload.Identifier).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare([]byte(storedToken), []byte(payload.Token)) == 1 {
		err = c.Redis.Del(ctx, prefix+payload.Identifier).Err()
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
