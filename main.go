package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sjain93/userservice/api/user"
	"github.com/sjain93/userservice/config"
	"github.com/sjain93/userservice/migrations"
	"github.com/sjain93/userservice/routes"
)

func main() {

	var (
		userRepository user.UserRepoManager
		err            error
	)

	noDB := flag.Bool("noDB", false, "Bool if the server should init in memory store")

	if *noDB {
		config.InitInMemoryStore()
		userRepository, err = user.NewUserRepository(nil, config.InMemDB)
		if err != nil {
			log.Fatalf("Error initializing in memory datastore: %v", err.Error())
		}

	} else {
		err = config.LoadEnvVars()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}

		config.ConnectDatabase()
		migrations.AutoMigrate(config.DB)

		userRepository, err = user.NewUserRepository(config.DB, nil)
		if err != nil {
			log.Fatalf("Error initializing postgres datastore: %v", err.Error())
		}
	}

	userService := user.NewUserService(userRepository)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.SetupRoutes(e, userService)

	e.Start(":8080")
}
