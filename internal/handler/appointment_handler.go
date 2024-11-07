package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"github.com/mskovv/tg-bot-subaru96/internal/service"
	"net/http"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AppointmentHandler struct {
	srv *service.AppointmentService
	bot *tgbotapi.BotAPI
}

func NewAppointmentHandler(srv *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{srv: srv}
}

func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := h.srv.UpdateAppointment(&appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

func (h *AppointmentHandler) RemoveAppointment(update tgbotapi.Update) {
	args := update.Message.CommandArguments()
	appointmentId, err := strconv.Atoi(args)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат ID.")
		h.bot.Send(msg)
		return
	}

	if err = h.srv.RemoveAppointment(uint(appointmentId)); err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Запись успешно удалена!")
	h.bot.Send(msg)

}

func (h *AppointmentHandler) GetAppointmentBuId(c *gin.Context) {}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := h.srv.CreateAppointment(&appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ok"})
}
