package handler

import (
	"context"
	fsmstate "github.com/looplab/fsm"
	"github.com/mskovv/tg-bot-subaru96/internal/fsm"
	"github.com/mskovv/tg-bot-subaru96/internal/service"
	"github.com/mskovv/tg-bot-subaru96/internal/storage"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"time"
)

type AppointmentHandler struct {
	srv     *service.AppointmentService
	storage *storage.RedisStorage
	bot     *telego.Bot
}

func NewAppointmentHandler(srv *service.AppointmentService, storage *storage.RedisStorage, bot *telego.Bot) *AppointmentHandler {
	return &AppointmentHandler{
		srv:     srv,
		storage: storage,
		bot:     bot,
	}
}

func (h *AppointmentHandler) HandleMessage(ctx context.Context, update telego.Update) {
	userID := update.Message.Chat.ID

	currentState, err := h.storage.GetState(ctx, userID)
	if err != nil {
		log.Println("Error getting fsm:", err)
		return
	}

	if currentState == "" {
		currentState = fsm.StateStart
		err = h.storage.SetState(ctx, userID, currentState)

		if err != nil {
			log.Println("Error setting fsm:", err)
			return
		}
	}

	stateMachine := fsm.NewAppointmentFSM()
	stateMachine.SetState(currentState)

	switch stateMachine.Current() {
	case fsm.StateStart:
		h.ShowCalendar(ctx, userID, stateMachine)
	case fsm.StateSelectDate:
	//	TODO
	case fsm.StateSelectTime:
	//	TODO
	case fsm.StateEnterCarModel:
	//	TODO
	case fsm.StateEnterDescription:
	//	TODO
	case fsm.StateConfirmation:
		//	TODO
	default:
		_, err = h.bot.SendMessage(tu.Message(
			tu.ID(userID),
			"Неизвестное состояние. Начинаем сначала.",
		))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		h.resetState(ctx, userID, stateMachine)
	}
}

func (h *AppointmentHandler) resetState(ctx context.Context, userId int64, state *fsmstate.FSM) {
	state.Event(ctx, "reset")
	h.storage.SetState(ctx, userId, state.Current())
}

func (h *AppointmentHandler) ShowCalendar(ctx context.Context, userId int64, state *fsmstate.FSM) {
	startDate := time.Now()
	if startDate.Weekday() != time.Monday {
		for startDate.Weekday() != time.Monday {
			startDate = startDate.AddDate(0, 0, -1)
		}
	}

	var buttons []telego.KeyboardButton
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, i).Format("02.01.2006")
		buttons = append(buttons, tu.KeyboardButton(date))
	}

	keyboard := tu.Keyboard(
		tu.KeyboardRow(buttons...),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Выберите дату")

	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(userId),
		"Выберите свободную дату для записи:",
	).WithReplyMarkup(keyboard))

	if err != nil {
		log.Println("Error sending calendar:", err)
		return
	}

	state.Event(ctx, "chose_date")
	h.storage.SetState(ctx, userId, state.Current())
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
