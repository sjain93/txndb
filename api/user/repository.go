package user

import (
	"errors"

	"github.com/sjain93/userservice/config"
	"gorm.io/gorm"
)

var ErrNoDatastore = errors.New("no datastore provided")

type UserRepoManager interface {
	Create(user *User) error
	GetAllUsers() ([]User, error)
}

type userRepository struct {
	DB       *gorm.DB
	memstore config.MemoryStore
}

func NewUserRepository(db *gorm.DB, inMemStore config.MemoryStore) (UserRepoManager, error) {
	if db != nil {
		return &userRepository{
			DB: db,
		}, nil
	} else if inMemStore != nil {
		return &userRepository{
			memstore: inMemStore,
		}, nil
	}

	return &userRepository{}, ErrNoDatastore

}

func (r *userRepository) Create(user *User) error {
	if r.DB != nil {
		if err := r.DB.Create(user).Error; err != nil {
			return err
		}
	} else {
		r.memstore[user.ID] = *user
	}

	return nil
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	if r.DB != nil {
		if err := r.DB.Find(&users).Error; err != nil {
			return users, err
		}
	} else {
		for _, u := range r.memstore {
			user, ok := u.(User)
			if !ok {
				return users, errors.New("invalid user data type")
			}
			users = append(users, user)
		}
	}

	return users, nil
}

// Implement other CRUD operations (Read, Update, Delete) here
