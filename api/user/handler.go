package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) Create(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := h.Service.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, user)
}

// Implement other handlers (Read, Update, Delete) here
