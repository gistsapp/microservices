package core

import (
	"errors"
	"time"

	"github.com/gistsapp/api/auth/repositories"
	"github.com/gistsapp/api/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService interface {
    CreateAccessToken(claims *types.JWTClaims) (string, error)
    CreateRefreshToken(userID string) (string, error)
    VerifyAccessToken(token string) (*types.JWTClaims, error)
    VerifyRefreshToken(token string) (string, error)
    InvalidateRefreshToken(token string) error
}

type jwtService struct {
    secretKey string
    db        repositories.Database
}

func NewJWTService(secretKey string, db repositories.Database) JWTService {
    return &jwtService{
        secretKey: secretKey,
        db:        db,
    }
}

func (j *jwtService) CreateAccessToken(claims *types.JWTClaims) (string, error) {
    claims.RegisteredClaims = jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        NotBefore: jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) CreateRefreshToken(userID string) (string, error) {
    tokenValue := uuid.New().String()
    
    expiresAt := time.Now().Add(time.Hour * 24 * 7)
    expiresAtStr := expiresAt.Format(time.RFC3339)
    
    opaqueToken := &types.OpaqueToken{
        ID:        uuid.New().String(),
        UserID:    userID,
        Token:     tokenValue,
        ExpiresAt: expiresAtStr,
    }
    
    _, err := j.db.CreateOpaqueToken(opaqueToken)
    if err != nil {
        return "", err
    }
    
    return tokenValue, nil
}

func (j *jwtService) VerifyAccessToken(tokenString string) (*types.JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(j.secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*types.JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrTokenMalformed
}

func (j *jwtService) VerifyRefreshToken(tokenString string) (string, error) {
    opaqueToken, err := j.db.GetOpaqueTokenByToken(tokenString)
    if err != nil {
        return "", err
    }

    expiresAt, err := time.Parse(time.RFC3339, opaqueToken.ExpiresAt)
    if err != nil {
        return "", err
    }
    
    if time.Now().After(expiresAt) {
        j.db.DeleteOpaqueToken(opaqueToken.ID)
        return "", errors.New("refresh token expired")
    }
    
    return opaqueToken.UserID, nil
}

func (j *jwtService) InvalidateRefreshToken(tokenString string) error {
    opaqueToken, err := j.db.GetOpaqueTokenByToken(tokenString)
    if err != nil {
        return err
    }
    
    return j.db.DeleteOpaqueToken(opaqueToken.ID)
}
