package service

import (
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"github.com/mskovv/tg-bot-subaru96/internal/repository"
)

type AppointmentService struct {
	repo *repository.AppointmentRepository
}

func NewAppointmentService(repo *repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{repo: repo}
}

func (s *AppointmentService) UpdateAppointment(appointment *models.Appointment) error {
	return s.repo.UpdateAppointment(appointment)
}

func (s *AppointmentService) GetAppointmentsById(id int) (*models.Appointment, error) {
	return s.repo.GetAppointmentById(id)
}

func (s *AppointmentService) RemoveAppointment(appointmentId uint) error {
	return s.repo.RemoveAppointment(appointmentId)
}

func (s *AppointmentService) CreateAppointment(appointment *models.Appointment) error {
	return s.repo.CreateAppointment(appointment)
}
