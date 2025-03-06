package main

import (
	"fmt"
	"os"

	"github.com/gistsapp/api/auth/config"
	"github.com/gistsapp/api/auth/repositories"
	"github.com/gistsapp/api/types"
)

func main() {
	config.LoadConfig() // reads the config file

	conf := config.GetConfig()

	// just testing thhings for now

	db, error:= repositories.NewPgDatabase(conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Database)
	fmt.Println(db)

	if error != nil {
		panic(error)
	}

	if os.Getenv("BOOTSTRAP") == "true" {
		err := db.Bootstrap()
		if err != nil {
			panic(err)
		}
	}

	user, error := db.CreateUser(&types.User{
		Username: "mihai",
		Email:    "mihai@example.com",
		Picture:  "https://avatars.githubusercontent.com/u/123456?v=4",
	})

	fmt.Println(user)
}
