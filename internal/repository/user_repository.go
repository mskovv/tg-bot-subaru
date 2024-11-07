package repository

import (
	"errors"
	"github.com/mskovv/tg-bot-subaru96/internal/database"
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) CreateUser(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *UserRepository) GetUserById(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) RemoveUser(id int) error {
	if err := database.DB.First(&models.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("пользователь не найден")
		}
		return err
	}
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) UpdateUser(id int, user models.User) error {
	if err := database.DB.First(&models.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("пользователь не найден")
		}
		return err
	}

	return r.db.First(&models.User{}, id).Updates(&user).Error
}
