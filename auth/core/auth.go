package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gistsapp/api/auth/config"
	"github.com/gistsapp/api/auth/repositories"
	"github.com/gistsapp/api/auth/utils"
	"github.com/gistsapp/api/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

type AuthService interface {
	RegisterProviders() //done
	IsAuthenticated(token string) (*types.JWTClaims, error) //done
	Renew(token string) (*types.AuthTokens, error) //done
	AuthenticateWithRedirect(c *fiber.Ctx) error //done
	AuthenticateWithCode(email string) (*types.VerificationToken, error) //done
	VerifyAuthToken(code string, email string) (*types.AuthTokens, error)
	Callback(c *fiber.Ctx) (*types.AuthTokens, error) //done
	RegisterUser(options *RegistrationOptions) (*types.User, error) //done
}

type RegistrationOptions struct {
	GothUser goth.User   // user data from federated identity provider
	User     *types.User // user data provided by the user through the registration form
}

type authService struct {
	providers    []goth.Provider
	jwtService   JWTService
	userService  UserService
	database     repositories.Database
	emailService repositories.EmailService
}

func NewAuthService(providers_config config.AuthProviders, jwtService JWTService, userService UserService, database repositories.Database, emailService repositories.EmailService) AuthService {
	providers := []goth.Provider{}
	
	for _, provider := range providers_config {
		switch provider.Name {
		case "github":
			provider := github.New(provider.ClientID, provider.ClientSecret, provider.RedirectURI)
			providers = append(providers, provider)
			break
		case "google":
			provider := google.New(provider.ClientID, provider.ClientSecret, provider.RedirectURI)
			providers = append(providers, provider)
			break
		}
	}

	return &authService{
		providers:    providers,
		jwtService:   jwtService,
		userService:  userService,
		database:     database,
		emailService: emailService,
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

	if provider != "github" && provider != "google" {
		return nil, UnkownProvider
	}


	auth_user, err := goth_fiber.CompleteUserAuth(c)

	if err != nil {
		return nil, ErrCantCompleteAuth
	}

	user, err := a.userService.GetUserByID(auth_user.UserID)

	if err == types.ErrNotFound { // user not found, create user
		user, err = a.firstLogin(auth_user, provider)
		if err != nil {
			return nil, err
		}
		return a.generateTokens(user)
	} else {
		return a.generateTokens(user)
	}
}

func (a *authService) firstLogin(auth_user goth.User, provider string) (*types.User, error) {
	var user_identity *types.User
	var err error
	if provider == "local" {
		user_identity, err = a.RegisterUser(withEmailPrefix(auth_user))
	} else {
		user_identity, err = a.RegisterUser(withOIDCUsername(auth_user))
	}
	if err != nil {
		return nil, err
	}
	return user_identity, nil
}

func (a *authService) generateTokens(user *types.User) (*types.AuthTokens, error) {
	access_token, err := a.jwtService.CreateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refresh_token, err := a.jwtService.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &types.AuthTokens{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func withOIDCUsername(user goth.User) *RegistrationOptions {
	return &RegistrationOptions{
		GothUser: user,
		User: &types.User{
			ID:       user.UserID,
			Email:    user.Email,
			Username: user.NickName,
			Picture:  user.AvatarURL,
		},
	}
}

func withEmailPrefix(user goth.User) *RegistrationOptions {
	return &RegistrationOptions{
		GothUser: user,
		User: &types.User{
			ID:       user.UserID,
			Email:    user.Email,
			Username: strings.Split(user.Email, "@")[0],
			Picture:  user.AvatarURL,
		},
	}
}

func (a *authService) RegisterUser(options *RegistrationOptions) (*types.User, error) {
	auth_user := options.GothUser
	data, err := json.Marshal(auth_user)
	if err != nil {
		return nil, err
	}
	user_data, err := a.userService.CreateUser(options.User)

	if err != nil {
		return nil, err
	}
	federated_identity := types.FederatedIdentity{
		ID:         auth_user.UserID,
		Data:       string(data),
		Provider:   auth_user.Provider,
		UserID:     user_data.ID,
	}

	_, err = a.database.CreateFederatedIdentity(&federated_identity)
	return user_data, err
}

func (a *authService) AuthenticateWithCode(email string) (*types.VerificationToken, error) {
	token_value := utils.GenToken(6)
	token_data := types.VerificationToken{
		Token: token_value,
		Email: email,
	}

	token, err := a.database.CreateVerificationToken(&token_data)
	if err != nil { // need to check if the token already exists
		token, err = a.database.GetVerificationTokenByEmail(email)
		if err != nil {
			return nil, err
		}
		err = a.database.DeleteVerificationToken(email, token.Token)
		if err != nil {
			return nil, err
		}

		// retry
		return a.AuthenticateWithCode(email)
	}

	a.emailService.SendVerificationEmail(email, token_value)

	return token, nil
}

func (a *authService) AuthenticateWithRedirect(c *fiber.Ctx) error {
	if _, err := goth_fiber.CompleteUserAuth(c); err == nil {
		return nil
	} else {
		return goth_fiber.BeginAuthHandler(c)
	}
}

func (a *authService) VerifyAuthToken(code string, email string) (*types.AuthTokens, error) {
	err := a.database.DeleteVerificationToken(email, code)
	if err == sql.ErrNoRows {
		return nil, ErrCantCompleteAuth
	}

	//now we finish user registration

	goth_user := goth.User{
		UserID:   email,
		Email:    email,
		Provider: "local",
		AvatarURL:  "https://vercel.com/api/www/avatar/?u=" + email + "&s=80",
	}

	federated_identity, err := a.database.GetFederatedIdentityByID(goth_user.UserID)
	var user *types.User
	if err == sql.ErrNoRows { // user not found, create user because first connection
		user, err = a.RegisterUser(withEmailPrefix(goth_user))
	}else {
		user, err = a.userService.GetUserByID(federated_identity.UserID)
	}
	return a.generateTokens(user)
}

var ErrCantCompleteAuth error = errors.New("Couldn't complete auth")
var ErrInvalidCode error = errors.New("Invalid verification code")
var UnkownProvider error = errors.New("Unkown provider")
