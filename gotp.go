package gotp

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Generate struct {
	Format     format
	Length     int
	Identifier string
	Expires    time.Duration
}

type Verify struct {
	Token      string
	Identifier string
}

var ctx = context.Background()

const prefix = "gotp_"

func New(c Config) (Config, error) {
	err := c.Redis.Ping(ctx).Err()

	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func (c Config) Generate(payload Generate) (token string, err error) {
	if payload.Length < 4 || payload.Length > 10 {
		return "", errors.New("length must be between 4 and 10")
	}

	if payload.Identifier == "" {
		return "", errors.New("identifier is required")
	}

	switch payload.Format {
	case ALPHA:
		token = generateAlphaToken(payload.Length)
		break
	case ALPHA_NUMERIC:
		token = generateAlphaNumericToken(payload.Length)
		break
	case NUMERIC:
		token = generateNumericToken(payload.Length)
		break
	}

	err = c.Redis.Set(ctx, prefix+payload.Identifier, token, payload.Expires).Err()

	if err != nil {
		return "", err
	}

	return token, nil
}

func (c Config) Verify(payload Verify) (valid bool, err error) {
	storedToken, err := c.Redis.Get(ctx, prefix+payload.Identifier).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if storedToken == payload.Token {
		err = c.Redis.Del(ctx, prefix+payload.Identifier).Err()

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
