package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sjain93/userservice/api/user"
)

func SetupRoutes(e *echo.Echo, userService *user.UserService) {
	userHandler := user.NewUserHandler(userService)

	api := e.Group("/api")
	users := api.Group("/users")

	users.POST("", userHandler.Create)
	// Implement other routes (GET, PUT, DELETE) here
}
