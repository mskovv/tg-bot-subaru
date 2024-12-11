package appointment

import (
	"github.com/mskovv/tg-bot-subaru96/internal/fsm"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"time"
)

func (h *Handler) showCalendar(update telego.Update) {
	userID := update.Message.Chat.ID

	var buttons []telego.InlineKeyboardButton
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, i).Format("02.01.2006")
		buttons = append(buttons, tu.InlineKeyboardButton(date).WithCallbackData(fsm.StateSelectDate+":"+date))
	}

	keyboard := tu.InlineKeyboardGrid(
		tu.InlineKeyboardCols(1, buttons...),
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

func (h *Handler) showTimeSelection(userId int64) {
	var buttons []telego.InlineKeyboardButton
	timeSlots := h.createTimesSlots()
	for _, tm := range timeSlots {
		buttons = append(buttons, tu.InlineKeyboardButton(tm).WithCallbackData(fsm.StateSelectTime+":"+tm))
	}

	keyboard := tu.InlineKeyboardGrid(
		tu.InlineKeyboardCols(2, buttons...),
	)

	_, err := h.bot.SendMessage(tu.Message(
		tu.ID(userId),
		"Выберите свободное время для записи:",
	).WithReplyMarkup(keyboard))

	if err != nil {
		log.Println("Error sending time selection:", err)
	}
}

// Рабочее время с 10 до 19:30. Раннее время - 10:00, крайнее 19:00. По полчаса секции
func (h *Handler) createTimesSlots() []string {
	startTime := time.Date(1, 1, 1, 10, 0, 0, 0, time.UTC) // Начальное время 10:00
	endTime := time.Date(1, 1, 1, 19, 0, 0, 0, time.UTC)   // Конечное время 19:00
	interval := 30 * time.Minute                           // Интервал в полчаса

	var timeSlots []string

	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(interval) {
		timeSlots = append(timeSlots, t.Format("15:04")) // Форматируем время в строку "HH:MM"
	}

	return timeSlots
}

func (h *Handler) showCarMarkSelection(userId int64) {
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

func (h *Handler) showCarModelSelection(userId int64, mark string) {
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
			buttons = append(buttons, tu.InlineKeyboardButton(model+" "+frame).WithCallbackData(fsm.StateEnterCarModel+":"+model+" "+frame))
		}
	}

	keyboard := tu.InlineKeyboardGrid(tu.InlineKeyboardCols(3, buttons...))
	_, _ = h.bot.SendMessage(tu.Message(
		tu.ID(userId),
		"Выберите модель "+mark+":",
	).WithReplyMarkup(keyboard))
}

func (h *Handler) showConfirmation(userId int64) {
	var buttons []telego.InlineKeyboardButton
	buttons = append(buttons, tu.InlineKeyboardButton("Подтверждаю").WithCallbackData(fsm.StateConfirmation+":"+"yes"))
	buttons = append(buttons, tu.InlineKeyboardButton("Отмена").WithCallbackData(fsm.StateConfirmation+":"+"no"))

	keyboard := tu.InlineKeyboardGrid(tu.InlineKeyboardCols(2, buttons...))

	_, _ = h.bot.SendMessage(tu.Message(
		tu.ID(userId),
		"Проверьте информацию и подтвердите создание записи",
	).WithReplyMarkup(keyboard))
}
