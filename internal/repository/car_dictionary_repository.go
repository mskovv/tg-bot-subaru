package repository

import (
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"gorm.io/gorm"
)

type CarDictionaryRepository struct {
	db *gorm.DB
}

func NewCarDictionaryRepository(db *gorm.DB) *CarDictionaryRepository {
	return &CarDictionaryRepository{db: db}
}

func (r *CarDictionaryRepository) GetAllModelsByMark(mark string) ([]models.CarDictionary, error) {
	var ms []models.CarDictionary
	err := r.db.Where("mark = ?", mark).Select("car_model").Find(&ms).Error
	return ms, err
}

func (r *CarDictionaryRepository) GetAllMarks() ([]string, error) {
	var ms []string
	err := r.db.Model(models.CarDictionary{}).Group("mark").Pluck("mark", &ms).Error
	return ms, err
}
