package user

import (
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uService UserServiceManager
}

func NewUserHandler(service UserServiceManager) *UserHandler {
	return &UserHandler{uService: service}
}

func (h *UserHandler) Create(c echo.Context) error {
	userReq := new(HTTPUserRequest)
	if err := c.Bind(userReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	valid, err := govalidator.ValidateStruct(userReq)
	if !valid {
		c.Logger().Errorf("invalid request parameters w/ error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err != nil {
		c.Logger().Errorf("failure to validate request parameters w/ error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	uRes, err := h.uService.CreateUser(userReq.mapToUser())
	if err != nil {
		switch {
		case errors.Is(err, ErrSvcUserExists):
			return echo.NewHTTPError(http.StatusConflict)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return c.JSON(http.StatusCreated, uRes.mapToResponse())
}
