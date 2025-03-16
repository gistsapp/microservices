package main

import (
	"os"

	"github.com/gistsapp/api/auth/config"
	"github.com/gistsapp/api/auth/core"
	"github.com/gistsapp/api/auth/http"
	"github.com/gistsapp/api/auth/repositories"
	"github.com/gofiber/fiber/v2/log"
)

// @title			Gists auth service API
// @version		0.1
// @description	This is the API for the Gists auth service
// @contact.name	Courtcircuits
// @contact.url	https://github.com/courtcircuits
// @contact.email	tristan-mihai.radulescu@etu.umontpellier.fr
func main() {
	config.LoadConfig() // reads the config file
	conf := config.GetConfig()
	log.Info("Starting Gists auth service")
	// just testing things for now
	db, error := repositories.NewPgDatabase(conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Database)

	if error != nil {
		panic(error)
	}

	if os.Getenv("BOOTSTRAP") == "true" {
		err := db.Bootstrap()
		if err != nil {
			panic(err)
		}
	}
	user_service := core.NewUserService(db)
	jwt_service := core.NewJWTService(conf.JWTSecretKey, db)
	email_repository := repositories.NewEmailService(conf.EmailService)
	auth_service := core.NewAuthService(conf.AuthProviders, jwt_service, user_service, db, email_repository)

	auth_handler := http.NewAuthController(auth_service, &conf, jwt_service)
	docs_handler := http.NewDocsHandler()

	server := http.NewServer(conf.Port)
	server.Setup(auth_handler, docs_handler)
	server.Ignite()
}
