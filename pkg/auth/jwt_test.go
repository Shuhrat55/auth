package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateAccessToken(t *testing.T) {
	userID := 123

	token, expiresAt, err := GenerateAccessToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Greater(t, expiresAt, time.Now().Unix())

	// Проверяем что токен можно распарсить
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestGenerateRefreshToken(t *testing.T) {
	userID := 456

	token, err := GenerateRefreshToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверяем что токен можно распарсить
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestValidateToken_ValidToken(t *testing.T) {
	userID := 789

	token, _, err := GenerateAccessToken(userID)
	assert.NoError(t, err)

	claims, err := ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.NotNil(t, claims.ExpiresAt)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	tests := []struct {
		name        string
		tokenString string
	}{
		{
			name:        "empty token",
			tokenString: "",
		},
		{
			name:        "invalid format token",
			tokenString: "invalid.token.here",
		},
		{
			name:        "wrong signature token",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImV4cCI6MTkwMDAwMDAwMH0.wrongsignature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.tokenString)

			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Создаем истекший токен вручную
	claims := &Claims{
		UserID: 123,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	assert.NoError(t, err)

	parsedClaims, err := ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestTokenTypeDifference(t *testing.T) {
	userID := 111

	accessToken, accessExpires, err := GenerateAccessToken(userID)
	assert.NoError(t, err)

	refreshToken, err := GenerateRefreshToken(userID)
	assert.NoError(t, err)

	// Проверяем что оба токена разные
	assert.NotEqual(t, accessToken, refreshToken)

	// Проверяем что срок действия access меньше чем у refresh
	accessClaims, err := ValidateToken(accessToken)
	assert.NoError(t, err)

	refreshClaims, err := ValidateToken(refreshToken)
	assert.NoError(t, err)

	assert.True(t, refreshClaims.ExpiresAt.Unix() > accessClaims.ExpiresAt.Unix())
}