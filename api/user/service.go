package user

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/sjain93/userservice/api/common"
)

var (
	once     sync.Once
	instance *userService
)

// service errors
var (
	ErrSvcUserExists   = errors.New("target credential(s) already exist")
	ErrSvcUserNotFound = errors.New("target was not found")
)

type UserServiceManager interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	GetUser(id string) (User, error)
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
	if err != nil && errors.Is(err, common.ErrUniqueKeyViolated) {
		return user, ErrSvcUserExists
	} else if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) GetUser(id string) (User, error) {
	u := User{ID: id}

	user, err := s.userRepo.GetUser(&u)
	if err != nil && errors.Is(err, common.ErrRecordNotFound) {
		return user, ErrSvcUserNotFound
	} else if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]User, error) {
	return s.userRepo.GetAllUsers()
}

// Implement other service methods (Read, Update, Delete) here
