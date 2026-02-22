package gotp

import (
	"testing"
	"time"
	"unicode"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestGenerateNumericToken(t *testing.T) {
	tests := []struct {
		length int
	}{
		{length: 4},
		{length: 6},
		{length: 8},
	}

	for _, test := range tests {
		t.Run("GenerateNumericToken", func(t *testing.T) {
			token, err := generateNumericToken(test.length)

			assert.NoError(t, err)
			assert.Equal(t, test.length, len(token), "Token length should match the requested length")

			for _, char := range token {
				assert.True(t, unicode.IsDigit(char), "Token should contain only numeric characters")
			}
		})
	}
}

func TestGenerateAlphaNumericToken(t *testing.T) {
	tests := []struct {
		length int
	}{
		{length: 4},
		{length: 6},
		{length: 8},
	}

	for _, test := range tests {
		t.Run("GenerateAlphaNumericToken", func(t *testing.T) {
			token, err := generateAlphaNumericToken(test.length)

			assert.NoError(t, err)
			assert.Equal(t, test.length, len(token), "Token length should match the requested length")

			for _, char := range token {
				assert.True(t, unicode.IsDigit(char) || unicode.IsLower(char), "Token should contain only alphanumeric characters")
			}
		})
	}
}

func TestGenerateAlphaToken(t *testing.T) {
	tests := []struct {
		length int
	}{
		{length: 4},
		{length: 6},
		{length: 8},
	}

	for _, test := range tests {
		t.Run("GenerateAlphaToken", func(t *testing.T) {
			token, err := generateAlphaToken(test.length)

			assert.NoError(t, err)
			assert.Equal(t, test.length, len(token), "Token length should match the requested length")

			for _, char := range token {
				assert.True(t, unicode.IsLower(char), "Token should contain only lowercase alphabetic characters")
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	db, mock := redismock.NewClientMock()
	config := Config{Redis: db}


	tests := []struct {
		name      string
		payload   Generate
		expectErr bool
	}{
		{
			name: "Length too short",
			payload: Generate{
				Format:     Alpha,
				Length:     3,
				Identifier: "testIdentifier",
				Expires:    10 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "Length too long",
			payload: Generate{
				Format:     Alpha,
				Length:     11,
				Identifier: "testIdentifier",
				Expires:    10 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "Empty Identifier",
			payload: Generate{
				Format:     Alpha,
				Length:     6,
				Identifier: "",
				Expires:    10 * time.Minute,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := config.Generate(tt.payload)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGenerateSuccess(t *testing.T) {
	db, mock := redismock.NewClientMock()
	config := Config{Redis: db}


	payload := Generate{
		Format:     Numeric,
		Length:     6,
		Identifier: "testUser",
		Expires:    5 * time.Minute,
	}

	mock.Regexp().ExpectSet(prefix+payload.Identifier, `^\d{6}$`, payload.Expires).SetVal("OK")

	token, err := config.Generate(payload)

	assert.NoError(t, err)
	assert.Len(t, token, 6)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerifySuccess(t *testing.T) {
	db, mock := redismock.NewClientMock()
	config := Config{Redis: db}


	identifier := "testUser"
	token := "123456"

	mock.ExpectGet(prefix + identifier).SetVal(token)
	mock.ExpectDel(prefix + identifier).SetVal(1)

	valid, err := config.Verify(Verify{
		Token:      token,
		Identifier: identifier,
	})

	assert.NoError(t, err)
	assert.True(t, valid)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerifyInvalidToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	config := Config{Redis: db}


	identifier := "testUser"

	mock.ExpectGet(prefix + identifier).SetVal("123456")

	valid, err := config.Verify(Verify{
		Token:      "wrongtoken",
		Identifier: identifier,
	})

	assert.NoError(t, err)
	assert.False(t, valid)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerifyNotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	config := Config{Redis: db}


	identifier := "nonexistent"

	mock.ExpectGet(prefix + identifier).RedisNil()

	valid, err := config.Verify(Verify{
		Token:      "123456",
		Identifier: identifier,
	})

	assert.NoError(t, err)
	assert.False(t, valid)
	assert.NoError(t, mock.ExpectationsWereMet())
}
