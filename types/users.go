package types

import "github.com/golang-jwt/jwt"

type User struct {
	ID       string `json:"id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Picture  string `json:"picture" db:"picture"`
}

type FederatedIdentity struct {
	ID       string `db:"federated_identity_id"`
	UserID   string `db:"user_id"`
	Provider string `db:"provider"`
	Data     string `db:"data"`
}

// an opaque token is a token that tied to a user and stored in the database.
// it is used as a refresh token for example.
type OpaqueToken struct {
	ID       string `db:"token_id"`
	UserID   string `db:"user_id"`
	Token    string `db:"token"`
	ExpiresAt string `db:"expires_at"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type VerificationToken struct {
	Email    string `db:"email"`
	Token    string `db:"token"`
}
