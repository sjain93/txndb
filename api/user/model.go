package user

import (
	"time"
)

type HTTPUserRequest struct {
	Username      string `json:"username" valid:"required"`
	AccountNumber string `json:"account_number" valid:"required"`
	Email         string `json:"email" valid:"email,required"`
}

type HTTPUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Shared model for Service and Repository layer
type User struct {
	ID            string `gorm:"primaryKey"`
	Username      string `gorm:"unique"`
	Email         string `gorm:"unique"`
	AccountNumber string `gorm:"unique"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *User) mapToResponse() HTTPUserResponse {
	return HTTPUserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

func (uR *HTTPUserRequest) mapToUser() User {
	return User{
		Username:      uR.Username,
		AccountNumber: uR.AccountNumber,
		Email:         uR.Email,
	}
}
