package appointment

import (
	"fmt"
	"github.com/mskovv/tg-bot-subaru96/internal/fsm"
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"strings"
	"time"
)

func (h *Handler) getWeekCalendar(weekOffset int) *telego.InlineKeyboardMarkup {
	now := time.Now()

	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1)
	startOfWeek = startOfWeek.AddDate(0, 0, 7*weekOffset)

	appointments, err := h.appointmentSrv.GetAppointmentsOnWeek(startOfWeek)
	if err != nil {
		log.Fatalln("Failed to get GetAppointmentsOnWeek: ", err)
	}

	appointmentsCount := make(map[string]int)
	for _, appointment := range appointments {
		dateStr := appointment.Date.Format("02.01.2006")
		appointmentsCount[dateStr]++
	}

	var weekButtons [][]telego.InlineKeyboardButton
	currentState := h.fsm.Current()
	for i := 0; i < 5; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		dateStr := day.Format("02.01.2006")
		count := appointmentsCount[dateStr]
		buttonText := fmt.Sprintf("%s - %d зап.", dateStr, count)
		weekButtons = append(weekButtons, []telego.InlineKeyboardButton{
			tu.InlineKeyboardButton(buttonText).
				WithCallbackData(currentState + ":" + dateStr),
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

	appointments, err := h.appointmentSrv.GetAppointmentsOnDate(h.appointment.Date)
	if err != nil {
		log.Fatalln("Failed to get GetAppointmentsOnWeek: ", err)
	}

	appointmentsCount := make(map[string]int)
	for _, appointment := range appointments {
		timeStr := appointment.Time.Format("15:04")
		appointmentsCount[timeStr]++
	}

	timeSlots := h.createTimesSlots()
	for _, tm := range timeSlots {
		count := appointmentsCount[tm]
		buttonText := fmt.Sprintf("%s - %d зап.", tm, count)
		buttons = append(buttons, tu.InlineKeyboardButton(buttonText).WithCallbackData(fsm.StateSelectTime+":"+tm))
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
		timeSlots = append(timeSlots, t.Format("15:04"))
	}

	return timeSlots
}

func (h *Handler) getCarMarkSelection() *telego.InlineKeyboardMarkup {
	carMarks, err := h.carDictionarySrv.GetAllMarks()
	if err != nil {
		log.Fatalln("Failed to get carMarks: ", err)
	}

	var buttons []telego.InlineKeyboardButton
	for _, cm := range carMarks {
		buttons = append(buttons, tu.InlineKeyboardButton(cm).WithCallbackData(fsm.StateEnterCarMark+":"+cm))
	}

	//buttons = append(buttons, tu.InlineKeyboardButton("Другое(Пока не жмакать)").WithCallbackData(fsm.StateEnterCarMark+":other"))
	//TODO
	keyboard := tu.InlineKeyboardGrid(tu.InlineKeyboardCols(3, buttons...))
	return keyboard
}

func (h *Handler) getCarModelSelection() *telego.InlineKeyboardMarkup {
	var buttons []telego.InlineKeyboardButton

	carModels, err := h.carDictionarySrv.GetAllModelsByMark(h.appointment.CarMark)
	if err != nil {
		log.Fatalln("Failed to get carModels: ", err)
	}

	for _, frames := range carModels {
		buttons = append(buttons, tu.InlineKeyboardButton(frames.CarModel).WithCallbackData(fsm.StateEnterCarModel+":"+frames.CarModel))
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

func (h *Handler) FormatAppointmentsOnDate(appointments []models.Appointment) string {
	if len(appointments) == 0 {
		return "Нет записей на выбранный день."
	}

	var sb strings.Builder
	for i, appointment := range appointments {
		sb.WriteString(fmt.Sprintf(
			"Время: %s\nМарка: %s\nМодель: %s\nОписание: %s\n",
			appointment.Time.Format("15:04"),
			appointment.CarMark,
			appointment.CarModel,
			appointment.Description,
		))

		if i < len(appointments)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
