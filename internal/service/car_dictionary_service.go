package service

import (
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"github.com/mskovv/tg-bot-subaru96/internal/repository"
)

type CarDictionaryService struct {
	repo *repository.CarDictionaryRepository
}

func NewCarDictionaryService(repo *repository.CarDictionaryRepository) *CarDictionaryService {
	return &CarDictionaryService{repo: repo}
}

func (s *CarDictionaryService) GetAllModelsByMark(mark string) ([]models.CarDictionary, error) {
	return s.repo.GetAllModelsByMark(mark)
}

func (s *CarDictionaryService) GetAllMarks() ([]string, error) {
	return s.repo.GetAllMarks()
}
