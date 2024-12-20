package repository

import (
	"errors"
	"github.com/mskovv/tg-bot-subaru96/internal/database"
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"gorm.io/gorm"
	"time"
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

func (r *AppointmentRepository) GetAppointmentsOnWeek(startWeek time.Time) ([]models.Appointment, error) {
	var ap []models.Appointment
	endWeek := startWeek.AddDate(0, 0, 5)
	err := r.db.Where("date >= ? AND date < ?", startWeek, endWeek).Find(&ap).Error
	return ap, err
}

func (r *AppointmentRepository) GetAppointmentsOnDate(date time.Time) ([]models.Appointment, error) {
	var ap []models.Appointment
	err := r.db.Select("date, time AT TIME ZONE 'UTC' AS time, car_mark, car_model, description").
		Where("DATE(date) = ?", date).Find(&ap).Error
	return ap, err
}

func (r *AppointmentRepository) RemoveAppointment(appointmentId uint) error {
	if err := database.DB.First(&models.Appointment{}, appointmentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("запись не найдена")
		}
		return err
	}
	return r.db.Delete(&models.Appointment{}, appointmentId).Error
}

func (r *AppointmentRepository) UpdateAppointment(appointment *models.Appointment) error {
	if err := database.DB.First(&models.Appointment{}, appointment.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("запись не найдена")
		}
		return err
	}

	return r.db.First(&models.Appointment{}, appointment.ID).Updates(&appointment).Error
}
