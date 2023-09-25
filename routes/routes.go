package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sjain93/userservice/api/user"
)

func SetupAPIRoutes(e *echo.Echo, userService user.UserServiceManager) {
	// API Groups
	api := e.Group("/api")
	users := api.Group("/users")

	// Service Handlers
	userHandler := user.NewUserHandler(userService)
	users.POST("", userHandler.Create)
	// Implement other routes (GET, PUT, DELETE) here
}
