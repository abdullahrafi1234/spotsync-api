package repository

import (
	"spotsync-api/models"

	"gorm.io/gorm"
)

// UserRepository defines what database operations we can do on users.
// Using an interface lets us swap implementations easily (e.g. for testing).
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}

// userRepository is the concrete (real) implementation using GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository is a "constructor" function that creates a UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}