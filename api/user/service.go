package user

import (
	"errors"
	"sync"

	"github.com/sjain93/userservice/api/common"
)

var (
	once     sync.Once
	instance *userService
)

// service errors
var ErrSvcUserExists = errors.New("target user already exists")

type UserServiceManager interface {
	CreateUser(user User) (User, error)
}

type userService struct {
	userRepo UserRepoManager
}

func NewUserService(r UserRepoManager) UserServiceManager {
	once.Do(func() {
		instance = &userService{
			userRepo: r,
		}
	})
	return instance
}

func (s *userService) CreateUser(user User) (User, error) {
	userID := common.GetMD5HashWithSum(user.Username + user.Email)
	// check if user exists
	_, err := s.userRepo.GetByID(userID)
	if errors.Is(err, ErrRecordNotFound) {
		user.ID = userID
		return user, s.userRepo.Create(&user)
	}

	return user, ErrSvcUserExists
}

// Implement other service methods (Read, Update, Delete) here
