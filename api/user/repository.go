package user

import (
	"github.com/pkg/errors"
	"github.com/sjain93/userservice/api/common"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRepoManager interface {
	Create(user *User) error
	GetAllUsers() ([]User, error)
	GetUser(user *User) (User, error)
	DeleteUser(user *User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) (UserRepoManager, error) {
	if db != nil {
		return &userRepository{
			DB: db,
		}, nil
	}

	return &userRepository{}, common.ErrNoDatastore
}

// Creates a new user for either memory store and enforces model constraints
func (r *userRepository) Create(user *User) error {
	if r.DB != nil {
		if err := r.DB.Create(user).Error; err != nil {
			// this is a GORM implementation detail
			var perr *pgconn.PgError
			if ok := errors.As(err, &perr); ok && perr.Code == common.UniqueViolationErr {
				return common.ErrUniqueKeyViolated
			} else {
				return err
			}
		}
	}

	return nil
}

// Returns a user for any field provided in the User struct
func (r *userRepository) GetUser(user *User) (User, error) {
	if r.DB != nil {
		err := r.DB.First(user).Error
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return *user, common.ErrRecordNotFound
			default:
				return *user, err
			}
		}
	}
	return *user, nil
}

// Returns all users from chosen memory store
func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	if r.DB != nil {
		if err := r.DB.Find(&users).Error; err != nil {
			return users, err
		}
	}

	return users, nil
}

func (r *userRepository) DeleteUser(user *User) error {
	if r.DB != nil {
		err := r.DB.Delete(&user).Error
		if err != nil {
			return err
		}

	}
	return nil
}
