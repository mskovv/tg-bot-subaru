package appointment

import (
	"fmt"
	"github.com/mskovv/tg-bot-subaru96/internal/fsm"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"time"
)

func (h *Handler) getWeekCalendar(weekOffset int) *telego.InlineKeyboardMarkup {
	now := time.Now()

	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1)
	startOfWeek = startOfWeek.AddDate(0, 0, 7*weekOffset)

	var weekButtons [][]telego.InlineKeyboardButton
	for i := 0; i < 5; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		dateStr := day.Format("02.01.2006")
		weekButtons = append(weekButtons, []telego.InlineKeyboardButton{
			tu.InlineKeyboardButton(dateStr).
				WithCallbackData(fsm.StateSelectDate + ":" + dateStr),
		})
	}

	navigationRow := []telego.InlineKeyboardButton{
		tu.InlineKeyboardButton("⬅️").WithCallbackData(fmt.Sprintf("nav_week:%d", weekOffset-1)),
		tu.InlineKeyboardButton("➡️").WithCallbackData(fmt.Sprintf("nav_week:%d", weekOffset+1)),
	}
	weekButtons = append(weekButtons, navigationRow)

	keyboard := tu.InlineKeyboard(
		weekButtons...,
	)
	return keyboard
}

func (h *Handler) getTimeSelection() *telego.InlineKeyboardMarkup {
	var buttons []telego.InlineKeyboardButton
	timeSlots := h.createTimesSlots()
	for _, tm := range timeSlots {
		buttons = append(buttons, tu.InlineKeyboardButton(tm).WithCallbackData(fsm.StateSelectTime+":"+tm))
	}

	keyboard := tu.InlineKeyboardGrid(
		tu.InlineKeyboardCols(2, buttons...),
	)
	return keyboard
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

func (h *Handler) getCarMarkSelection() *telego.InlineKeyboardMarkup {
	var buttons []telego.InlineKeyboardButton
	carModels := []string{"Subaru", "Toyota", "Suzuki", "Другое"}
	for _, cm := range carModels {
		buttons = append(buttons, tu.InlineKeyboardButton(cm).WithCallbackData(fsm.StateEnterCarMark+":"+cm))
	}

	keyboard := tu.InlineKeyboardGrid(tu.InlineKeyboardCols(3, buttons...))
	return keyboard
}

func (h *Handler) getCarModelSelection() *telego.InlineKeyboardMarkup {
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

	for model, frames := range carModelsMarks[h.appointment.CarMark] {
		for _, frame := range frames {
			buttons = append(buttons, tu.InlineKeyboardButton(model+" "+frame).WithCallbackData(fsm.StateEnterCarModel+":"+model+" "+frame))
		}
	}

	keyboard := tu.InlineKeyboardGrid(tu.InlineKeyboardCols(3, buttons...))
	return keyboard
}

func (h *Handler) getConfirmation() *telego.InlineKeyboardMarkup {
	var buttons []telego.InlineKeyboardButton
	buttons = append(buttons, tu.InlineKeyboardButton("Подтверждаю ✅").WithCallbackData(fsm.StateConfirmation+":"+"yes"))
	buttons = append(buttons, tu.InlineKeyboardButton("Отмена ✖️").WithCallbackData(fsm.StateConfirmation+":"+"no"))

	keyboard := tu.InlineKeyboardGrid(tu.InlineKeyboardRows(2, buttons...))

	return keyboard
}
