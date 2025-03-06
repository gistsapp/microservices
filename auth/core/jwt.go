package core

import "github.com/gistsapp/api/types"

type JWTService interface {
	CreateAccessToken(user *types.User) (string, error)
	CreateRefreshToken(user_id string) (string, error)
	VerifyAccessToken(token string) (*types.JWTClaims, error)
	VerifyRefreshToken(token string) (bool, error)
}

// for dorian -> implement this service bellow
