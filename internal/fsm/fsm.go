package fsm

import (
	"github.com/looplab/fsm"
)

const (
	StateStart            = "start"
	StateSelectDate       = "select_date"
	StateSelectTime       = "select_time"
	StateEnterCarModel    = "enter_car_model"
	StateEnterCarMark     = "enter_car_mark"
	StateEnterDescription = "enter_description"
	StateConfirmation     = "confirmation"
)

const (
	EventChoseDate        = "chose_date"
	EventChoseTime        = "chose_time"
	EventChoseCarModel    = "enter_car_model"
	EventChoseCarMark     = "enter_car_mark"
	EventChoseDescription = "enter_description"
	EventConfirm          = "confirm"
	EventReset            = "reset"
)

func NewAppointmentFSM() *fsm.FSM {
	return fsm.NewFSM(
		StateStart,
		fsm.Events{
			{Name: EventChoseDate, Src: []string{StateStart}, Dst: StateSelectDate},
			{Name: EventChoseTime, Src: []string{StateStart, StateSelectDate}, Dst: StateSelectTime},
			{Name: EventChoseCarModel, Src: []string{StateSelectTime}, Dst: StateEnterCarModel},
			{Name: EventChoseCarMark, Src: []string{StateEnterCarModel}, Dst: StateEnterCarMark},
			{Name: EventChoseDescription, Src: []string{StateEnterCarMark}, Dst: StateEnterDescription},
			{Name: EventConfirm, Src: []string{StateEnterDescription}, Dst: StateConfirmation},
			{Name: EventReset, Src: []string{"*"}, Dst: StateStart},
		},
		nil,
	)
}
