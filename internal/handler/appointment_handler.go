package handler

import (
	"context"
	"fmt"
	fsmstate "github.com/looplab/fsm"
	"github.com/mskovv/tg-bot-subaru96/internal/fsm"
	"github.com/mskovv/tg-bot-subaru96/internal/service"
	"github.com/mskovv/tg-bot-subaru96/internal/storage"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"strings"
	"time"
)

type AppointmentHandler struct {
	srv     *service.AppointmentService
	storage *storage.RedisStorage
	bot     *telego.Bot
	fsm     *fsmstate.FSM
}

func NewAppointmentHandler(srv *service.AppointmentService, storage *storage.RedisStorage, bot *telego.Bot, fsm *fsmstate.FSM) *AppointmentHandler {
	return &AppointmentHandler{
		srv:     srv,
		storage: storage,
		bot:     bot,
		fsm:     fsm,
	}
}

func (h *AppointmentHandler) HandleMessage(update telego.Update) {
	userID := update.Message.From.ID
	command := update.Message.Text

	switch command {
	case "/create_appointment":
		h.StartAppointmentCreation(update)
	case "update":
		//h.StartAppointmentUpdate(update)
	case "delete":
		//h.StartAppointmentDeletion(update)
	default:
		h.bot.SendMessage(tu.Message(
			tu.ID(userID),
			"Неизвестная команда. Попробуйте снова.",
		))
	}
}

func (h *AppointmentHandler) StartAppointmentCreation(update telego.Update) {
	userID := update.Message.Chat.ID

	currentState := h.fsm.Current()

	if currentState == "" {
		h.fsm.SetState(fsm.StateStart)
		h.setRedisState(update.Context(), userID)
	}

	switch currentState {
	case fsm.StateStart:
		h.ShowCalendar(update)
		h.fsm.Event(update.Context(), fsm.EventChoseDate)
		h.setRedisState(update.Context(), userID)
	case fsm.StateSelectDate:
		h.ShowCalendar(update)
		h.fsm.Event(update.Context(), fsm.EventChoseTime)
		h.setRedisState(update.Context(), userID)
	case fsm.StateSelectTime:
		h.handleTimeSelection(update.Message.Chat.ID)
		h.fsm.Event(update.Context(), fsm.EventChoseCarModel)
		h.setRedisState(update.Context(), userID)
	case fsm.StateEnterCarMark:
		h.handleCarMarkSelection(update.Message.Chat.ID)
		h.fsm.Event(update.Context(), fsm.EventChoseCarModel)
		h.setRedisState(update.Context(), userID)
	case fsm.StateEnterCarModel:
		//h.handleCarModelSelection(update.Message.Chat.ID)
	//	TODO
	case fsm.StateEnterDescription:
	//	TODO
	case fsm.StateConfirmation:
		//	TODO
	default:
		_, err := h.bot.SendMessage(tu.Message(
			tu.ID(userID),
			"Неизвестное состояние. Начинаем сначала.",
		))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		h.resetState(update.Context(), userID)
	}
}

func (h *AppointmentHandler) HandleCallback(callback telego.CallbackQuery) {
	state, payload := parseCallbackData(callback.Data)
	fmt.Println("state: ", state, "payload: ", payload)
	ctx := context.Background()
	userId := callback.Message.GetChat().ID

	switch state {
	case fsm.StateSelectDate:
		err, _ := h.bot.SendMessage(tu.Message(tu.ID(userId), "Выбранная дата: "+payload))
		if err != nil {
			log.Println("Error sending callback:", err)
		}
		h.handleTimeSelection(callback.Message.GetChat().ID)

		err2 := h.fsm.Event(ctx, fsm.EventChoseTime)
		if err2 != nil {
			log.Println("Error state.Event:", err2)
			h.resetState(ctx, userId)
		}
		h.setRedisState(ctx, userId)
	case fsm.StateSelectTime:
		err, _ := h.bot.SendMessage(tu.Message(tu.ID(userId), "Выбранное время: "+payload))
		if err != nil {
			log.Println("Error sending callback:", err)
		}

		err2 := h.fsm.Event(ctx, fsm.EventChoseCarModel)
		if err2 != nil {
			log.Println("Error state.Event:", err2)
			h.resetState(ctx, userId)
		}
		h.setRedisState(ctx, userId)
		h.handleCarMarkSelection(callback.Message.GetChat().ID)
	case fsm.StateEnterCarMark:
		err, _ := h.bot.SendMessage(tu.Message(tu.ID(userId), "Выбранная марка: "+payload))
		if err != nil {
			log.Println("Error sending callback:", err)
		}
		err2 := h.fsm.Event(ctx, fsm.EventChoseCarMark)
		if err2 != nil {
			log.Println("Error state.Event:", err2)
			log.Println("Current state", h.fsm.Current())
			h.resetState(ctx, userId)
		}
		h.setRedisState(ctx, userId)
		h.handleCarModelSelection(userId, payload)

	case fsm.StateEnterCarModel:

	case "confirm_details":
		// TODO
		//h.HandleConfirmation(callback)
	default:
		h.bot.AnswerCallbackQuery(tu.CallbackQuery(callback.ID).WithText("Неизвестное действие"))
	}
}

func (h *AppointmentHandler) setRedisState(ctx context.Context, userId int64) {
	err := h.storage.SetState(ctx, userId, h.fsm.Current())
	if err != nil {
		log.Println("Error set state:", err)
	}
}

func (h *AppointmentHandler) resetState(ctx context.Context, userId int64) {
	h.fsm = fsm.NewAppointmentFSM()
	//err := h.fsm.Event(ctx, "reset")
	//if err != nil {
	//	log.Println("Error state.reset:", err)
	//}
	h.setRedisState(ctx, userId)
}

func (h *AppointmentHandler) ShowCalendar(update telego.Update) {
	userID := update.Message.Chat.ID

	var buttons []telego.InlineKeyboardButton
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, i).Format("02.01.2006")
		buttons = append(buttons, tu.InlineKeyboardButton(date).WithCallbackData("select_date:"+date))
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardCols(1, buttons...)...,
	)

	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(userID),
		"Выберите свободную дату для записи:",
	).WithReplyMarkup(keyboard))

	if err != nil {
		log.Println("Error sending calendar:", err)
		h.resetState(update.Context(), userID)
		return
	}
}

func (h *AppointmentHandler) SendStartMessage(update telego.Update) {
	userID := update.Message.Chat.ID
	err := h.fsm.Event(update.Context(), fsm.EventReset)
	if err != nil {
		log.Println("Error state.Event:", err)
		h.resetState(update.Context(), userID)
	}
	h.setRedisState(update.Context(), userID)
	err = h.SendMessage(update, "Состояние сброшено, выберите команду из меню")
	if err != nil {
		err = h.SendMessage(update, "Произошла ошибка при отправке сообщения. Попробуйте снова.")
		if err != nil {
			return
		}
	}
}

func (h *AppointmentHandler) handleTimeSelection(userId int64) {
	var buttons []telego.InlineKeyboardButton
	timeSlots := h.createTimesSlots()
	for _, tm := range timeSlots {
		buttons = append(buttons, tu.InlineKeyboardButton(tm).WithCallbackData("select_time:"+tm))
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardCols(2, buttons...)...,
	)

	_, _ = h.bot.SendMessage(tu.Message(
		tu.ID(userId),
		"Выберите свободное время для записи:",
	).WithReplyMarkup(keyboard))
}

// Рабочее время с 10 до 19:30. Раннее время - 10:00, крайнее 19:00. По полчаса секции
func (h *AppointmentHandler) createTimesSlots() []string {
	startTime := time.Date(1, 1, 1, 10, 0, 0, 0, time.UTC) // Начальное время 10:00
	endTime := time.Date(1, 1, 1, 19, 0, 0, 0, time.UTC)   // Конечное время 19:00
	interval := 30 * time.Minute                           // Интервал в полчаса

	var timeSlots []string

	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(interval) {
		timeSlots = append(timeSlots, t.Format("15:04")) // Форматируем время в строку "HH:MM"
	}

	return timeSlots
}

func (h *AppointmentHandler) handleCarMarkSelection(userId int64) {
	var buttons []telego.InlineKeyboardButton
	carModels := []string{"Subaru", "Toyota", "Suzuki", "Другое"}
	for _, cm := range carModels {
		buttons = append(buttons, tu.InlineKeyboardButton(cm).WithCallbackData(fsm.StateEnterCarMark+":"+cm))
	}

	keyboard := tu.InlineKeyboardGrid(tu.InlineKeyboardCols(3, buttons...))

	_, _ = h.bot.SendMessage(tu.Message(
		tu.ID(userId),
		"Выберите марку:",
	).WithReplyMarkup(keyboard))
}

func (h *AppointmentHandler) handleCarModelSelection(userId int64, mark string) {

	fmt.Println("handleCarModelSelection")
	imprezaModels := []string{
		"GF/GC",
		"GG/GD",
		"GH/GЕ",
		"GJ",
		"GR/GV",
		"GP/GJ",
		"GP(XV)",
	}
	foresterModels := []string{
		"SF",
		"SG",
		"SH",
		"SJ",
		"SK",
	}
	outbackModels := []string{
		"BG",
		"BH-BHE",
		"BP9-BPE",
		"BM9-BR9",
		"BS",
	}

	subaruModels := make(map[string][]string)
	subaruModels["Impeza"] = imprezaModels
	subaruModels["Forester"] = foresterModels
	subaruModels["Outback"] = outbackModels

	carModelsMarks := make(map[string]map[string][]string)
	carModelsMarks["Subaru"] = subaruModels

	var buttons []telego.InlineKeyboardButton

	for model, frames := range carModelsMarks[mark] {
		for _, frame := range frames {
			buttons = append(buttons, tu.InlineKeyboardButton(model+" "+frame).WithCallbackData("select_model:"+model+" "+frame))
		}
		//for model, frame := range cm {
		//fmt.Println(model, frame)
		//buttons =
		//}
	}

	keyboard := tu.InlineKeyboardGrid(tu.InlineKeyboardCols(3, buttons...))
	//keyboard := tu.InlineKeyboard(
	//	tu.InlineKeyboardCols(2, buttons...)...,
	//)

	_, _ = h.bot.SendMessage(tu.Message(
		tu.ID(userId),
		"Выберите модель "+mark+":",
	).WithReplyMarkup(keyboard))
}

func parseCallbackData(data string) (state, payload string) {
	parts := strings.SplitN(data, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return data, ""
}

func (h *AppointmentHandler) SendMessage(update telego.Update, message string) error {
	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(update.Message.Chat.ID),
		message,
	))

	return err
}
