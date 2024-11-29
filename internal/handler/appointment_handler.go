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
}

func NewAppointmentHandler(srv *service.AppointmentService, storage *storage.RedisStorage, bot *telego.Bot) *AppointmentHandler {
	return &AppointmentHandler{
		srv:     srv,
		storage: storage,
		bot:     bot,
	}
}

func (h *AppointmentHandler) HandleMessage(update telego.Update) {
	userID := update.Message.From.ID
	command := update.Message.Text
	log.Println(command)

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
		h.resetState(update.Context(), userID)
	}
}

func (h *AppointmentHandler) StartAppointmentCreation(update telego.Update) {
	userID := update.Message.Chat.ID

	currentState, err := h.storage.GetState(update.Context(), userID)
	if err != nil {
		log.Println("Error getting fsm:", err)
		return
	}

	fmt.Println("current state:", currentState)
	if currentState == "" {
		currentState = fsm.StateStart
		err = h.storage.SetState(update.Context(), userID, currentState)

		if err != nil {
			log.Println("Error setting fsm:", err)
			return
		}
	}

	stateMachine := fsm.NewAppointmentFSM()
	stateMachine.SetState(currentState)
	fmt.Println("current state2:", stateMachine.Current())

	switch stateMachine.Current() {
	case fsm.StateStart:
		h.ShowCalendar(update, stateMachine)
	case fsm.StateSelectDate:
		h.HandleTimeSelection(update)
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
		h.resetState(update.Context(), userID)
	}
}

func (h *AppointmentHandler) resetState(ctx context.Context, userId int64) {
	stateMachine := fsm.NewAppointmentFSM()
	err := stateMachine.Event(ctx, "reset")
	if err != nil {
		log.Println("Error state.reset:", err)
	}

	err = h.storage.SetState(ctx, userId, stateMachine.Current())
	if err != nil {
		log.Println("Error SetState:", err)
		return
	}
}

func (h *AppointmentHandler) ShowCalendar(update telego.Update, state *fsmstate.FSM) {
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

	err = state.Event(update.Context(), "chose_date")
	if err != nil {
		log.Println("Error state.Event:", err)
		h.resetState(update.Context(), userID)
		return
	}
	err = h.storage.SetState(update.Context(), userID, state.Current())
	if err != nil {
		log.Println("Error SetState:", err)
		h.resetState(update.Context(), userID)
		return
	}
}

func (h *AppointmentHandler) SendStartMessage(update telego.Update) {
	userID := update.Message.Chat.ID
	currentState, err := h.storage.GetState(update.Context(), userID)

	if err != nil {
		log.Println("Error getting fsm:", err)
		return
	}

	if currentState != "" {
		h.resetState(update.Context(), userID)
	}
	h.storage.SetState(update.Context(), userID, "start")
	err = h.SendMessage(update, "Состояние сброшено, выберите команду из меню")
	if err != nil {
		err = h.SendMessage(update, "Произошла ошибка при отправке сообщения. Попробуйте снова.")
		if err != nil {
			return
		}
	}
}

func (h *AppointmentHandler) HandleTimeSelection(update telego.Update) {
	userID := update.Message.Chat.ID

	h.bot.SendMessage(tu.Message(
		tu.ID(userID),
		"Выбор времени",
	))
}

func (h *AppointmentHandler) HandleCallback(callback telego.CallbackQuery) {
	state, payload := parseCallbackData(callback.Data)

	switch state {
	case "select_date":
		// TODO
		h.bot.AnswerCallbackQuery(tu.CallbackQuery(callback.ID).WithText("Выбранная дата: " + payload))
	case "select_time":
		//h.HandleTimeSelection(callback, payload)
		// TODO

	case "confirm_details":
		// TODO
		//h.HandleConfirmation(callback)
	default:
		h.bot.AnswerCallbackQuery(tu.CallbackQuery(callback.ID).WithText("Неизвестное действие"))
	}
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
