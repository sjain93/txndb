package user

import (
	"net/http"

	"github.com/pkg/errors"

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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	valid, err := govalidator.ValidateStruct(userReq)
	if !valid {
		c.Logger().Errorf("invalid request parameters w/ error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, errors.
			Wrapf(err, "invalid request parameters").
			Error())
	}
	if err != nil {
		c.Logger().Errorf("failure to validate request parameters w/ error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	uRes, err := h.uService.CreateUser(userReq.mapToUser())
	if err != nil {
		switch {
		case errors.Is(err, ErrSvcUserExists):
			return echo.NewHTTPError(http.StatusConflict, errors.
				Wrapf(err, "failed to create user").
				Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, uRes.mapToResponse())
}

func (h *UserHandler) GetUser(c echo.Context) error {
	paramIDs := c.ParamValues()
	if len(paramIDs) == 0 || len(paramIDs) > 1 {
		c.Logger().Error("invalid number of ID parameters")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid number of parameters")
	}

	uRes, err := h.uService.GetUser(paramIDs[0])
	if err != nil {
		switch {
		case errors.Is(err, ErrSvcUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, errors.
				Wrapf(err, "failed to find user with given ID").
				Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusFound, uRes.mapToResponse())
}

func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.uService.GetAllUsers()
	if err != nil {
		c.Logger().Errorf("failure to get all users w/ error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	uRes := []HTTPUserResponse{}
	for _, u := range users {
		uRes = append(uRes, u.mapToResponse())
	}

	return c.JSON(http.StatusOK, uRes)
}
