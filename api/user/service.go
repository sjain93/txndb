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

	user.ID = common.GetMD5HashWithSum(user.Username + user.Email)
	err := s.userRepo.Create(&user)
	if err != nil {
		if errors.Is(err, ErrUniqueKeyViolated) {
			return user, ErrSvcUserExists
		} else {
			return user, err
		}
	}
	return user, nil
}

// Implement other service methods (Read, Update, Delete) here
