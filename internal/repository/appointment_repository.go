package repository

import (
	"errors"
	"github.com/mskovv/tg-bot-subaru96/internal/database"
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"gorm.io/gorm"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) CreateAppointment(appointment *models.Appointment) error {
	return r.db.Create(&appointment).Error
}

func (r *AppointmentRepository) GetAppointmentById(id int) (*models.Appointment, error) {
	var ap models.Appointment
	err := r.db.First(&ap, id).Error
	return &ap, err
}

func (r *AppointmentRepository) RemoveAppointment(appointmentId uint) error {
	if err := database.DB.First(&models.Appointment{}, appointmentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("пользователь не найден")
		}
		return err
	}
	return r.db.Delete(&models.Appointment{}, appointmentId).Error
}

func (r *AppointmentRepository) UpdateAppointment(appointment *models.Appointment) error {
	if err := database.DB.First(&models.Appointment{}, appointment.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("пользователь не найден")
		}
		return err
	}

	return r.db.First(&models.Appointment{}, appointment.ID).Updates(&appointment).Error
}
