package user

type UserService struct {
	Repository *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{Repository: repo}
}

func (s *UserService) CreateUser(user *User) error {
	return s.Repository.Create(user)
}

// Implement other service methods (Read, Update, Delete) here
