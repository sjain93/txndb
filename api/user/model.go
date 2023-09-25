package user

import (
	"time"

	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Username string `json:"username" valid:"required"`
	Email    string `json:"email" valid:"email,required"`
}

type CreateUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type User struct {
	ID        string `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Email     string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) mapToResponse() CreateUserResponse {
	return CreateUserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

func (uR *CreateUserRequest) mapToUser() User {
	return User{
		Username: uR.Username,
		Email:    uR.Email,
	}
}
