package user



import "gorm.io/gorm"

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *User) error {
	result := r.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Implement other CRUD operations (Read, Update, Delete) here
