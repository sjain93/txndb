package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uService UserServiceManager
}

func NewUserHandler(service UserServiceManager) *UserHandler {
	return &UserHandler{uService: service}
}

func (h *UserHandler) Create(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := h.uService.CreateUser(user); err != nil {
		// (todo SJ) Add switch case on error for user collision scenario
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, user)
}
