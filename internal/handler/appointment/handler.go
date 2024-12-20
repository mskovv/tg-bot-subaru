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
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	appointmentSrv   *service.AppointmentService
	carDictionarySrv *service.CarDictionaryService
	storage          *storage.RedisStorage
	bot              *telego.Bot
	fsm              *fsmstate.FSM
	appointment      *models.Appointment
}

func NewAppointmentHandler(
	appointmentSrv *service.AppointmentService,
	carDictionarySrv *service.CarDictionaryService,
	storage *storage.RedisStorage,
	bot *telego.Bot,
	fsm *fsmstate.FSM,
) *Handler {
	return &Handler{
		appointmentSrv:   appointmentSrv,
		carDictionarySrv: carDictionarySrv,
		storage:          storage,
		bot:              bot,
		fsm:              fsm,
	}
}

func (h *Handler) sendMessageWithReplyMarkup(userID int64, keyboard *telego.InlineKeyboardMarkup, text string) {
	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(userID),
		text,
	).WithReplyMarkup(keyboard))

	if err != nil {
		log.Println("Error sending calendar:", err)
		h.resetState(context.Background(), userID)
		return
	}
}

func (h *Handler) editReplyMarkupMessage(callback telego.CallbackQuery, keyboard *telego.InlineKeyboardMarkup, text string) {
	message := callback.Message
	chat := message.GetChat()

	var err error
	if text != "" {
		_, err = h.bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:      chat.ChatID(),
			MessageID:   message.GetMessageID(),
			Text:        text,
			ReplyMarkup: keyboard,
		})
	} else {
		_, err = h.bot.EditMessageReplyMarkup(&telego.EditMessageReplyMarkupParams{
			BusinessConnectionID: "",
			ChatID:               chat.ChatID(),
			MessageID:            message.GetMessageID(),
			InlineMessageID:      "",
			ReplyMarkup:          keyboard,
		})
	}
	if err != nil {
		log.Println("Error edit message:", err)
		h.resetState(context.Background(), chat.ID)
		return
	}
}

func (h *Handler) HandleCommand(update telego.Update) {
	userID := update.Message.From.ID
	command := update.Message.Text
	if h.appointment == nil {
		h.appointment = &models.Appointment{}
	}

	switch command {
	case "/create_appointment":
		h.StartAppointmentCreation(update)
	case "update":
		//h.StartAppointmentUpdate(update)
	case "delete":
		//h.StartAppointmentDeletion(update)
	default:
		_, err := h.bot.SendMessage(tu.Message(
			tu.ID(userID),
			"Неизвестная команда. Попробуйте снова.",
		))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}

func (h *Handler) StartAppointmentCreation(update telego.Update) {
	userID := update.Message.Chat.ID

	currentState := h.fsm.Current()

	if currentState == "" {
		h.fsm.SetState(fsm.StateStart)
		h.setRedisState(update.Context(), userID)
	}

	text := ""
	var keyboard *telego.InlineKeyboardMarkup

	switch currentState {
	case fsm.StateStart:
		keyboard = h.getWeekCalendar(0)
		text = "Выберите свободную дату для записи:"
	case fsm.StateSelectDate:
		keyboard = h.getWeekCalendar(0)
		text = "Выберите свободную дату для записи:"
	case fsm.StateSelectTime:
		keyboard = h.getTimeSelection()
		text = "Выберите свободное время для записи:"
	case fsm.StateEnterCarMark:
		keyboard = h.getCarMarkSelection()
		text = "Выберите марку:"
	case fsm.StateEnterCarModel:
		keyboard = h.getCarModelSelection()
		text = "Выберите модель " + h.appointment.CarMark + ":"
	case fsm.StateEnterDescription:
	//	TODO
	case fsm.StateConfirmation:
		keyboard = h.getConfirmation()
		text = fmt.Sprintln(h.appointment)
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

	h.sendMessageWithReplyMarkup(userID, keyboard, text)
	h.setRedisState(update.Context(), userID)
}

func (h *Handler) HandleCallback(callback telego.CallbackQuery) {
	state, payload := parseCallbackData(callback.Data)
	fmt.Println("state: ", state, "payload: ", payload)
	ctx := context.Background()
	message := callback.Message
	chat := message.GetChat()
	userId := chat.ID

	switch state {
	case "nav_week":
		weekOffset, err := strconv.Atoi(payload)
		if err != nil {
			log.Fatalln(err)
		}
		keyboard := h.getWeekCalendar(weekOffset)
		h.editReplyMarkupMessage(callback, keyboard, "")
		if err != nil {
			log.Println("Error edit message:", err)
			return
		}
	case fsm.StateSelectDate:
		var err error
		h.appointment.Date, err = time.Parse("02.01.2006", payload)
		if err != nil {
			_, err = h.bot.SendMessage(tu.Message(tu.ID(userId), "Произошла ошибка интерпритации даты"))
			if err != nil {
				log.Println("Error sending message:", err)
			}
			log.Fatal("Error parse Date: ", err)
			return
		}

		keyboard := h.getTimeSelection()
		h.editReplyMarkupMessage(callback, keyboard, "Выберите свободное время для записи")

		err = h.fsm.Event(ctx, fsm.EventChoseTime)
		if err != nil {
			log.Println("Error state.Event EventChoseTime:", err)
			h.resetState(ctx, userId)
		}
	case fsm.StateSelectTime:
		err := h.fsm.Event(ctx, fsm.EventChoseCarMark)
		if err != nil {
			log.Println("Error state.Event EventChoseCarMark:", err)
			h.resetState(ctx, userId)
		}
		h.appointment.Time, err = time.Parse("15:04", payload)
		if err != nil {
			log.Println("Error time.Parse:", err)
			h.resetState(ctx, userId)
		}
		keyboard := h.getCarMarkSelection()
		h.editReplyMarkupMessage(callback, keyboard, "Выберите марку:")
	case fsm.StateEnterCarMark:
		err := h.fsm.Event(ctx, fsm.EventChoseCarModel)
		if err != nil {
			log.Println("Error state.Event EventChoseCarModel:", err)
			h.resetState(ctx, userId)
		}
		h.appointment.CarMark = payload
		keyboard := h.getCarModelSelection()
		h.editReplyMarkupMessage(callback, keyboard, "Выберите модель "+h.appointment.CarMark+":")
	case fsm.StateEnterCarModel:
		err2 := h.fsm.Event(ctx, fsm.EventChoseDescription)
		if err2 != nil {
			log.Println("Error state.Event EventChoseDescription:", err2)
			log.Println("Current state", h.fsm.Current())
			h.resetState(ctx, userId)
		}
		h.appointment.CarModel = payload
		err := h.bot.DeleteMessage(&telego.DeleteMessageParams{MessageID: message.GetMessageID(), ChatID: chat.ChatID()})
		if err != nil {
			log.Println("Error delete message:", err)
			return
		}

		_, _ = h.bot.SendMessage(tu.Message(
			tu.ID(userId),
			"Отправьте описание необходимых действий",
		))

	//case fsm.StateEnterDescription: // UNUSED
	case fsm.StateConfirmation:
		if payload == "yes" {
			err := h.appointmentSrv.CreateAppointment(h.appointment)
			if err != nil {
				_, err = h.bot.SendMessage(tu.Message(tu.ID(userId), "Произошла ошибка cоздания записи"))
				if err != nil {
					log.Println("Error sending message:", err)
					return
				}
				log.Fatalln("Error create appointment: ", err)
			}

			_, err = h.bot.EditMessageReplyMarkup(&telego.EditMessageReplyMarkupParams{
				BusinessConnectionID: "",
				ChatID:               chat.ChatID(),
				MessageID:            message.GetMessageID(),
				InlineMessageID:      "",
				ReplyMarkup:          nil,
			})

			_, _ = h.bot.SendMessage(tu.Message(
				tu.ID(userId),
				"Запись успешно сохранена",
			))
			h.appointment = nil
			h.resetState(ctx, userId)
			return
		}
	default:
		err := h.bot.AnswerCallbackQuery(tu.CallbackQuery(callback.ID).WithText("Неизвестное действие"))
		if err != nil {
			log.Println("Error answer callback:", err)
			return
		}
		return
	}

	h.setRedisState(ctx, userId)
}

func (h *Handler) HandleMessage(message telego.Message) {
	chat := message.GetChat()
	userId := chat.ID

	currentState := h.fsm.Current()

	if currentState == fsm.StateEnterDescription || message.ReplyToMessage.Text == "Введите описание, ответив на это сообщение." {
		if message.Text == "" {
			_, _ = h.bot.SendMessage(tu.Message(tu.ID(userId), "Описание не может быть пустым"))
			return
		}

		err := h.fsm.Event(context.Background(), fsm.EventConfirm)
		if err != nil {
			log.Println("Error state.Event EventConfirm:", err)
			h.resetState(context.Background(), userId)
		}
		h.setRedisState(context.Background(), userId)

		h.appointment.Description = message.Text

		confirmKeyboard := h.getConfirmation()
		h.sendMessageWithReplyMarkup(userId, confirmKeyboard, fmt.Sprintln(h.appointment))
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
