package handler

import (
	"github.com/mskovv/tg-bot-subaru96/internal/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"time"
)

type AppointmentHandler struct {
	srv *service.AppointmentService
	bot *telego.Bot
}

func NewAppointmentHandler(srv *service.AppointmentService, bot *telego.Bot) *AppointmentHandler {
	return &AppointmentHandler{srv: srv, bot: bot}
}

func (h *AppointmentHandler) UpdateAppointment(message *telego.Message) {
	//args := message.Text[len("/update_appointment "):]
	//parts := strings.SplitN(args, " ", 2)
	//
	//if len(parts) < 2 {
	//	h.sendMessage(message.Chat.ID, "Пожалуйста, укажите ID записи и новое описание.")
	//	return
	//}
	//
	//appointmentID, err := strconv.Atoi(parts[0])
	//if err != nil {
	//	h.sendMessage(message.Chat.ID, "Неверный формат ID.")
	//	return
	//}
	//
	//newDescription := parts[1]
	//
	//appointment := &models.Appointment{
	//	ID:          uint(appointmentID),
	//	Description: newDescription,
	//}
	//
	//if err := h.srv.UpdateAppointment(appointment); err != nil {
	//	h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка обновления: %v", err))
	//	return
	//}
	//
	//h.sendMessage(message.Chat.ID, "Запись успешно обновлена!")
}

func (h *AppointmentHandler) RemoveAppointment(message *telego.Message) {
	//args := message.Text[len("/remove_appointment "):]
	//appointmentID, err := strconv.Atoi(args)
	//if err != nil {
	//	h.sendMessage(message.Chat.ID, "Неверный формат ID.")
	//	return
	//}
	//
	//if err := h.srv.RemoveAppointment(uint(appointmentID)); err != nil {
	//	h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка удаления: %v", err))
	//	return
	//}
	//
	//h.sendMessage(message.Chat.ID, "Запись успешно удалена!")
}

func (h *AppointmentHandler) GetAppointmentBuId(message *telego.Message) {}

func (h *AppointmentHandler) isDateAvailable(date time.Time) bool {
	return true
}

func (h *AppointmentHandler) CreateAppointment(update telego.Update) {
	startDate := time.Now()
	if startDate.Weekday() != time.Monday {
		for startDate.Weekday() != time.Monday {
			startDate = startDate.AddDate(0, 0, -1)
		}
	}

	var daysButtons []telego.KeyboardButton
	for i := 0; i < 5; i++ {
		date := startDate.AddDate(0, 0, i)
		if h.isDateAvailable(date) { // Проверяем доступность даты
			buttonText := date.Format("02 Января")
			daysButtons = append(daysButtons, tu.KeyboardButton(buttonText))
		}
	}

	keyboard := tu.Keyboard(
		tu.KeyboardRow(daysButtons...),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Выберите дату")

	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Выберите свободную дату для записи:",
	).WithReplyMarkup(keyboard))

	if err != nil {
		return
	}
}

func (h *AppointmentHandler) SendStartMessage(update telego.Update) {
	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Создать запись"),
			tu.KeyboardButton("Обновить запись"),
			tu.KeyboardButton("Удалить запись"),
		),
	).WithResizeKeyboard()

	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Добро пожаловать! Выберите команду для работы с записями:",
	).WithReplyMarkup(keyboard))

	if err != nil {
		// Обрабатываем ошибку отправки сообщения
		err = h.SendMessage(update, "Произошла ошибка при отправке сообщения. Попробуйте снова.")
		if err != nil {
			return
		}
	}
}

func (h *AppointmentHandler) SendMessage(update telego.Update, message string) error {
	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(update.Message.Chat.ID),
		message,
	))

	return err
}
