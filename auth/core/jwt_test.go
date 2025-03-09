package core

import (
	"testing"
	"time"

	"github.com/gistsapp/api/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessToken(t *testing.T) {
	secretKey := "peuxpasdire"
	jwtService := NewJWTService(secretKey)
	claims := &types.JWTClaims{
		UserID:   "tristianoronaldo",
	}

	token, err := jwtService.CreateAccessToken(claims)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.ParseWithClaims(token, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	parsedClaims, ok := parsedToken.Claims.(*types.JWTClaims)
	assert.True(t, ok)
	assert.Equal(t, claims.UserID, parsedClaims.UserID)
	
	expectedExpiry := time.Now().Add(time.Hour * 24)
	tokenExpiry := parsedClaims.ExpiresAt.Time
	diff := expectedExpiry.Sub(tokenExpiry)
	assert.LessOrEqual(t, diff.Abs(), time.Second*5)
}

func TestCreateRefreshToken(t *testing.T) {
	secretKey := "peuxpasdire"
	jwtService := NewJWTService(secretKey)
	claims := &types.JWTClaims{
		UserID:   "tristianoronaldo",
	}

	token, err := jwtService.CreateRefreshToken(claims)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.ParseWithClaims(token, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	parsedClaims, ok := parsedToken.Claims.(*types.JWTClaims)
	assert.True(t, ok)
	assert.Equal(t, claims.UserID, parsedClaims.UserID)
	
	expectedExpiry := time.Now().Add(time.Hour * 24 * 7)
	tokenExpiry := parsedClaims.ExpiresAt.Time
	diff := expectedExpiry.Sub(tokenExpiry)
	assert.LessOrEqual(t, diff.Abs(), time.Second*5)
}

func TestVerifyAccessToken_Valid(t *testing.T) {
	secretKey := "peuxpasdire"
	jwtService := NewJWTService(secretKey)
	claims := &types.JWTClaims{
		UserID:   "tristianoronaldo",
	}

	token, err := jwtService.CreateAccessToken(claims)
	assert.NoError(t, err)

	parsedClaims, err := jwtService.VerifyAccessToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, parsedClaims)
	assert.Equal(t, claims.UserID, parsedClaims.UserID)
}

func TestVerifyAccessToken_Invalid(t *testing.T) {
	jwtService := NewJWTService("peuxpasdire")
	invalidTokens := []string{
		"invalid.token.format",
		"",
		createTokenWithDifferentKey(),
	}

	for _, invalidToken := range invalidTokens {
		claims, err := jwtService.VerifyAccessToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
	}
}

func TestVerifyRefreshToken_Valid(t *testing.T) {
	secretKey := "peuxpasdire"
	jwtService := NewJWTService(secretKey)
	claims := &types.JWTClaims{
		UserID:   "tristianoronaldo",
	}

	token, err := jwtService.CreateRefreshToken(claims)
	assert.NoError(t, err)

	valid, err := jwtService.VerifyRefreshToken(token)

	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestVerifyRefreshToken_Invalid(t *testing.T) {
	jwtService := NewJWTService("peuxpasdire")
	invalidTokens := []string{
		"invalid.token.format",
		"",
		createTokenWithDifferentKey(),
	}

	for _, invalidToken := range invalidTokens {
		valid, err := jwtService.VerifyRefreshToken(invalidToken)

		assert.Error(t, err)
		assert.False(t, valid)
	}
}

func TestExpiredToken(t *testing.T) {
	secretKey := "peuxpasdire"
	jwtService := NewJWTService(secretKey)
	
	claims := &types.JWTClaims{
		UserID:   "tristianoronaldo",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expiré il y a 1 heure
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString([]byte(secretKey))
	assert.NoError(t, err)
	
	parsedClaims, err := jwtService.VerifyAccessToken(expiredToken)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	
	valid, err := jwtService.VerifyRefreshToken(expiredToken)
	assert.Error(t, err)
	assert.False(t, valid)
}

func createTokenWithDifferentKey() string {
	claims := &types.JWTClaims{
		UserID:   "tristianoronaldo",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	differentKey := "peuxpasdiretopsecret"
	tokenString, _ := token.SignedString([]byte(differentKey))
	return tokenString
}