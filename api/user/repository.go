package user

import (
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sjain93/userservice/config"
	"gorm.io/gorm"
)

var (
	ErrNoDatastore       = errors.New("no datastore provided")
	ErrRecordNotFound    = errors.New("record not found")
	ErrInvalidDataType   = errors.New("invalid user data type")
	ErrUniqueKeyViolated = errors.New("duplicated key not allowed")
)

const (
	UniqueViolationErr = "23505"
)

type UserRepoManager interface {
	Create(user *User) error
	GetAllUsers() ([]User, error)
	GetUser(user *User) (User, error)
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

// Creates a new user for either memory store and enforces model constraints
func (r *userRepository) Create(user *User) error {
	if r.DB != nil {
		if err := r.DB.Create(user).Error; err != nil {
			// this is a GORM implementation detail
			var perr *pgconn.PgError
			if ok := errors.As(err, &perr); ok && perr.Code == UniqueViolationErr {
				return ErrUniqueKeyViolated
			} else {
				return err
			}
		}
	} else {
		_, ok := r.memstore[user.ID]
		if ok {
			return ErrUniqueKeyViolated
		}
		// ensuring unique constraint for email and username when using in mem
		u := r.userForMemStoreValues(user.Username, user.Email)
		if u != nil {
			return ErrUniqueKeyViolated
		}

		// save user when all unique constraints are met
		r.memstore[user.ID] = *user
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
				return *user, ErrRecordNotFound
			default:
				return *user, err
			}
		}
	} else {
		val, ok := r.memstore[user.ID]
		if !ok {
			u := r.userForMemStoreValues(user.Username, user.Email)
			if u == nil {
				return *user, ErrRecordNotFound
			}
			val = u
		}
		u, ok := val.(User)
		if !ok {
			return *user, ErrInvalidDataType
		}
		user = &u
	}
	return *user, nil
}

// In order to emulate the PQ constraint in the memory store
func (r *userRepository) userForMemStoreValues(vals ...string) interface{} {
	var foundID string
	for _, v := range vals {
		for id, storedUser := range r.memstore {
			sU := storedUser.(User)
			if sU.Email == v || sU.Username == v {
				foundID = id
				break
			}
		}
	}
	return r.memstore[foundID]
}

// Returns all users from chosen memory store
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
				return users, ErrInvalidDataType
			}
			users = append(users, user)
		}
	}

	return users, nil
}
