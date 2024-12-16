package appointment

import (
	"context"
	"fmt"
	fsmstate "github.com/looplab/fsm"
	"github.com/mskovv/tg-bot-subaru96/internal/fsm"
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"github.com/mskovv/tg-bot-subaru96/internal/service"
	"github.com/mskovv/tg-bot-subaru96/internal/storage"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"strings"
	"time"
)

type Handler struct {
	srv         *service.AppointmentService
	storage     *storage.RedisStorage
	bot         *telego.Bot
	fsm         *fsmstate.FSM
	appointment *models.Appointment
}

func NewAppointmentHandler(srv *service.AppointmentService, storage *storage.RedisStorage, bot *telego.Bot, fsm *fsmstate.FSM) *Handler {
	return &Handler{
		srv:     srv,
		storage: storage,
		bot:     bot,
		fsm:     fsm,
	}
}

func (h *Handler) HandleCommand(update telego.Update) {
	userID := update.Message.From.ID
	command := update.Message.Text
	h.appointment = &models.Appointment{}

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

func (h *Handler) StartAppointmentCreation(update telego.Update) {
	userID := update.Message.Chat.ID

	currentState := h.fsm.Current()

	if currentState == "" {
		h.fsm.SetState(fsm.StateStart)
		h.setRedisState(update.Context(), userID)
	}

	switch currentState {
	case fsm.StateStart:
		h.showCalendar(update)
		h.fsm.Event(update.Context(), fsm.EventChoseDate)
		h.setRedisState(update.Context(), userID)
	case fsm.StateSelectDate:
		h.showCalendar(update)
		h.fsm.Event(update.Context(), fsm.EventChoseTime)
		h.setRedisState(update.Context(), userID)
	case fsm.StateSelectTime:
		h.showTimeSelection(update.Message.Chat.ID)
		h.fsm.Event(update.Context(), fsm.EventChoseCarMark)
		h.setRedisState(update.Context(), userID)
	case fsm.StateEnterCarMark:
		h.showCarMarkSelection(update.Message.Chat.ID)
		h.fsm.Event(update.Context(), fsm.EventChoseCarModel)
		h.setRedisState(update.Context(), userID)
	case fsm.StateEnterCarModel:
		//h.showCarModelSelection(update.Message.Chat.ID)
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

func (h *Handler) HandleCallback(callback telego.CallbackQuery) {
	state, payload := parseCallbackData(callback.Data)
	fmt.Println("state: ", state, "payload: ", payload)
	ctx := context.Background()
	message := callback.Message
	chat := message.GetChat()
	userId := chat.ID

	switch state {
	case fsm.StateSelectDate:
		_, err := h.bot.SendMessage(tu.Message(tu.ID(userId), "Выбранная дата: "+payload))
		if err != nil {
			log.Println("Error sending result:", err)
		}
		h.showTimeSelection(chat.ID)

		err = h.fsm.Event(ctx, fsm.EventChoseTime)
		if err != nil {
			log.Println("Error state.Event EventChoseTime:", err)
			h.resetState(ctx, userId)
		}
		h.appointment.Date, err = time.Parse("02.01.2006", payload)
		if err != nil {
			h.bot.SendMessage(tu.Message(tu.ID(userId), "Произошла ошибка интерпритации даты"))
			log.Fatal("Error parse Date: ", err)
			return
		}
	case fsm.StateSelectTime:
		_, err := h.bot.SendMessage(tu.Message(tu.ID(userId), "Выбранное время: "+payload))
		if err != nil {
			log.Println("Error sending time payload:", err)
		}

		err2 := h.fsm.Event(ctx, fsm.EventChoseCarMark)
		if err2 != nil {
			log.Println("Error state.Event EventChoseCarMark:", err2)
			h.resetState(ctx, userId)
		}
		h.appointment.Time, _ = time.Parse("15:04", payload)

		h.showCarMarkSelection(chat.ID)
	case fsm.StateEnterCarMark:
		_, err := h.bot.SendMessage(tu.Message(tu.ID(userId), "Выбранная марка: "+payload))
		if err != nil {
			log.Println("Error sending car mark payload::", err)
		}
		err2 := h.fsm.Event(ctx, fsm.EventChoseCarModel)
		if err2 != nil {
			log.Println("Error state.Event EventChoseCarModel:", err2)
			log.Println("Current state", h.fsm.Current())
			h.resetState(ctx, userId)
		}
		h.appointment.CarMark = payload

		h.showCarModelSelection(userId, payload)
	case fsm.StateEnterCarModel:
		_, err := h.bot.SendMessage(tu.Message(tu.ID(userId), "Выбранная модель: "+payload))
		if err != nil {
			log.Println("Error sending Model payload:", err)
		}
		err2 := h.fsm.Event(ctx, fsm.EventChoseDescription)
		if err2 != nil {
			log.Println("Error state.Event EventChoseDescription:", err2)
			log.Println("Current state", h.fsm.Current())
			h.resetState(ctx, userId)
		}
		h.appointment.CarModel = payload

		h.handleEnterDescription(userId)
	case fsm.StateEnterDescription:
		_, _ = h.bot.SendMessage(tu.Message(tu.ID(userId), "Ввод описания "+payload))
	//	TODO NOTHING
	case fsm.StateConfirmation:
		// TODO
		//h.HandleConfirmation(callback)
	default:
		h.bot.AnswerCallbackQuery(tu.CallbackQuery(callback.ID).WithText("Неизвестное действие"))
	}

	h.setRedisState(ctx, userId)
	h.bot.DeleteMessage(&telego.DeleteMessageParams{MessageID: message.GetMessageID(), ChatID: chat.ChatID()})
}

func (h *Handler) HandleMessage(message telego.Message) {
	userId := message.GetChat().ID

	if message.ReplyToMessage != nil &&
		message.ReplyToMessage.Text == "Введите описание, ответив на это сообщение." {
		description := message.Text
		if description == "" {
			_, _ = h.bot.SendMessage(tu.Message(
				tu.ID(userId),
				"Описание не может быть пустым",
			))
			return
		}
		_, _ = h.bot.SendMessage(tu.Message(
			tu.ID(userId),
			"Описание успешно сохранено",
		))
		//h.fsm.Event(context.Background(), fsm.EventConfirm)
		//h.setRedisState(context.Background(), userId)
		err := h.fsm.Event(context.Background(), fsm.EventReset)
		if err != nil {
			log.Println("Error parse Date: ", err)
		}
		h.setRedisState(context.Background(), userId)

		h.appointment.Description = description
		err = h.srv.CreateAppointment(h.appointment)
		if err != nil {
			h.bot.SendMessage(tu.Message(tu.ID(userId), "Произошла ошибка cоздания записи"))
			log.Println("Error parse Date: ", err)
		}
		h.appointment = nil
	}
}

func (h *Handler) setRedisState(ctx context.Context, userId int64) {
	err := h.storage.SetState(ctx, userId, h.fsm.Current())
	if err != nil {
		log.Println("Error set state:", err)
	}
}

func (h *Handler) resetState(ctx context.Context, userId int64) {
	h.fsm = fsm.NewAppointmentFSM()
	h.setRedisState(ctx, userId)
}

func (h *Handler) SendStartMessage(update telego.Update) {
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

func (h *Handler) handleEnterDescription(userId int64) {
	_, _ = h.bot.SendMessage(tu.Message(
		tu.ID(userId),
		"Введите описание, ответив на это сообщение.",
	))
}

func parseCallbackData(data string) (state, payload string) {
	parts := strings.SplitN(data, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return data, ""
}

func (h *Handler) SendMessage(update telego.Update, message string) error {
	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(update.Message.Chat.ID),
		message,
	))

	return err
}
