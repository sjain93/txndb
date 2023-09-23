package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sjain93/userservice/api/user"
	"github.com/sjain93/userservice/config"
	"github.com/sjain93/userservice/migrations"
	"github.com/sjain93/userservice/routes"
)

func main() {
	err := config.LoadEnvVars()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
	config.ConnectDatabase()
	migrations.AutoMigrate(config.DB)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userRepository := user.NewUserRepository(config.DB)
	userService := user.NewUserService(userRepository)

	routes.SetupRoutes(e, userService)

	e.Start(":8080")
}
