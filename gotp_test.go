package gotp

import (
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"unicode"
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
			// Generate the token
			token := generateNumericToken(test.length)

			// Check if the token has the correct length
			assert.Equal(t, test.length, len(token), "Token length should match the requested length")

			// Check if the token contains only numeric characters
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
			// Generate the token
			token := generateAlphaNumericToken(test.length)

			// Check if the token has the correct length
			assert.Equal(t, test.length, len(token), "Token length should match the requested length")

			// Check if the token contains only alphanumeric characters
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
			// Generate the token
			token := generateAlphaToken(test.length)

			// Check if the token has the correct length
			assert.Equal(t, test.length, len(token), "Token length should match the requested length")

			// Check if the token contains only lowercase alphabetic characters
			for _, char := range token {
				assert.True(t, unicode.IsLower(char), "Token should contain only lowercase alphabetic characters")
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	// Create a new mock Redis client
	db, mock := redismock.NewClientMock()
	config := Config{Redis: db}

	tests := []struct {
		name         string
		payload      Generate
		mockReturn   string
		mockError    error
		expectErr    bool
		validateType bool
	}{
		{
			name: "Length too short",
			payload: Generate{
				Format:     ALPHA,
				Length:     3,
				Identifier: "testIdentifier",
				Expires:    10 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "Length too long",
			payload: Generate{
				Format:     ALPHA,
				Length:     11,
				Identifier: "testIdentifier",
				Expires:    10 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "Empty Identifier",
			payload: Generate{
				Format:     ALPHA,
				Length:     6,
				Identifier: "",
				Expires:    10 * time.Minute,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockReturn != "" {
				mock.ExpectSet(prefix+tt.payload.Identifier, tt.mockReturn, tt.payload.Expires).SetVal(tt.mockReturn)
			}

			if tt.mockError != nil {
				mock.ExpectSet(prefix+tt.payload.Identifier, "", tt.payload.Expires).SetErr(tt.mockError)
			}

			token, err := config.Generate(tt.payload)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.mockReturn, token)
			}

			// Ensure that all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
