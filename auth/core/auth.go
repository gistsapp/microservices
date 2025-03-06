package core

import (
	"errors"

	"github.com/gistsapp/api/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
)

type AuthService interface {
	RegisterProviders()
	IsAuthenticated(token string) (*types.JWTClaims, error)
	Renew(token string) (*types.AuthTokens, error)
	AuthenticateWithRedirect(c *fiber.Ctx) error
	AuthenticateWithCode(email string) (*types.OpaqueToken, error)
	Callback(c *fiber.Ctx) (*types.AuthTokens, error)
	RegisterUser(options *RegistrationOptions) (*types.User, error)
}

type RegistrationOptions struct {
	GothUser goth.User   // user data from federated identity provider
	User     *types.User // user data provided by the user through the registration form
}

type authService struct {
	providers  []goth.Provider
	jwtService JWTService
	userService UserService
}

func NewAuthService(providers []goth.Provider, jwtService JWTService, userService UserService) AuthService {
	return &authService{
		providers:  []goth.Provider{},
		jwtService: jwtService,
		userService: userService,
	}
}

func (a *authService) RegisterProviders() {
	goth.UseProviders(a.providers...)
}

func (a *authService) IsAuthenticated(token string) (*types.JWTClaims, error) {
	claims, err := a.jwtService.VerifyAccessToken(token)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (a *authService) Renew(token string) (*types.AuthTokens, error) {
	claims, err := a.jwtService.VerifyAccessToken(token)

	if err != jwt.ErrTokenExpired {
		return nil, err
	}

	refreshToken, err := a.jwtService.CreateRefreshToken(claims.UserID)

	if err != nil {
		return nil, err
	}

	user, err := a.userService.GetUserByID(claims.UserID)

	if err != nil {
		return nil, err
	}

	accessToken, err := a.jwtService.CreateAccessToken(user)

	return &types.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authService) Callback(c *fiber.Ctx) (*types.AuthTokens, error) {
	provider := c.Params("provider")
	auth_user, err := goth_fiber.CompleteUserAuth(c)

	if err != nil {
		return nil, ErrCantCompleteAuth
	}

	user, err := a.userService.GetUserByID(auth_user.)
}

var ErrCantCompleteAuth error = errors.New("Couldn't complete auth")
