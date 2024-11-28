package fsm

import "github.com/looplab/fsm"

const (
	StateStart            = "start"
	StateSelectDate       = "select_date"
	StateSelectTime       = "select_time"
	StateEnterCarModel    = "enter_car_model"
	StateEnterDescription = "enter_description"
	StateConfirmation     = "confirmation"
)

func NewAppointmentFSM() *fsm.FSM {
	return fsm.NewFSM(
		StateStart,
		fsm.Events{
			{Name: "chose_date", Src: []string{StateStart}, Dst: StateSelectDate},
			{Name: "choose_time", Src: []string{StateSelectDate}, Dst: StateSelectTime},
			{Name: "enter_car_model", Src: []string{StateSelectTime}, Dst: StateEnterCarModel},
			{Name: "enter_description", Src: []string{StateEnterCarModel}, Dst: StateEnterDescription},
			{Name: "confirm", Src: []string{StateEnterDescription}, Dst: StateConfirmation},
			{Name: "reset", Src: []string{"*"}, Dst: StateStart},
		},
		nil,
	)
}
