package user

import "sync"

var (
	once     sync.Once
	instance *userService
)

type UserServiceManager interface {
	CreateUser(user *User) error
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

func (s *userService) CreateUser(user *User) error {
	// check if user exists

	return s.userRepo.Create(user)
}

// Implement other service methods (Read, Update, Delete) here
