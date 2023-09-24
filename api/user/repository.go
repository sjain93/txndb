package user

import "gorm.io/gorm"

type UserRepoManager interface {
	Create(user *User) error
	GetAllUsers() ([]User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoManager {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) Create(user *User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	if err := r.DB.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

// Implement other CRUD operations (Read, Update, Delete) here
