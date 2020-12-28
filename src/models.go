package ttrn

import (
	"fmt"

	"gorm.io/gorm"
)

// User is a User of the system.
type User struct {
	gorm.Model
	Session      string
	Name         string
	IsAdmin      bool
	IsBanned     bool
	DrivedBefore int
}

// State represents the current state of the application.
type State struct {
	TrainSpeeds      [2]int
	TurnoutPositions [5]int
}

// DefaultState returns the initial state.
func DefaultState() State {
	return State{
		TrainSpeeds:      [2]int{0, 0},
		TurnoutPositions: [5]int{1, 1, 1, 1, 1},
	}
}

// ChangeTrainSpeed changes the speed of a train by a given delta and returns the new
// speed value.
func (s *State) ChangeTrainSpeed(id int, delta int) (int, error) {
	if id < 1 || id > 2 {
		return 0, fmt.Errorf("got invalid train id %d", id)
	}
	cSpeed := s.TrainSpeeds[id-1]
	nSpeed := cSpeed + delta
	if nSpeed < -4 || nSpeed > 4 {
		return cSpeed, nil
	}
	s.TrainSpeeds[id-1] = nSpeed
	return nSpeed, nil
}

// SwitchTurnout switches the turnout to the other position and returns the new state.
func (s *State) SwitchTurnout(id int) (int, error) {
	if id < 0 || id > 4 {
		return 0, fmt.Errorf("got invalid turnout id %d", id)
	}
	nPos := -1
	if s.TurnoutPositions[id] == -1 {
		nPos = 1
	}
	s.TurnoutPositions[id] = nPos
	return nPos, nil
}

// FrontendState represents the state for the frontend.
type FrontendState struct {
	Train1Speed      int `json:"train_1_speed"`
	Train2Speed      int `json:"train_2_speed"`
	Turnout0Position int `json:"turnout_1_position"`
	Turnout1Position int `json:"turnout_2_position"`
	Turnout2Position int `json:"turnout_3_position"`
	Turnout3Position int `json:"turnout_4_position"`
	Turnout4Position int `json:"turnout_5_position"`
}

// FromState returns a FrontendState instance from a State instance.
func FromState(s State) FrontendState {
	return FrontendState{
		Train1Speed:      s.TrainSpeeds[0],
		Train2Speed:      s.TrainSpeeds[1],
		Turnout0Position: s.TurnoutPositions[0],
		Turnout1Position: s.TurnoutPositions[1],
		Turnout2Position: s.TurnoutPositions[2],
		Turnout3Position: s.TurnoutPositions[3],
		Turnout4Position: s.TurnoutPositions[4],
	}
}

// TurnoutPositionEvent is a event on the position change of a turnaout.
type TurnoutPositionEvent struct {
	// Id of the turnout.
	Id int
	// NewPosition is the new position of the turnout. O = straight, 1 = diverging.
	NewPosition int
}

// TrainSpeedEvent is a event on the speed change of a train.
type TrainSpeedEvent struct {
	// Id of the train.
	Id int
	// NewSpeed is the new speed of the train ranging from -4 to 4.
	NewSpeed int
}

// TrainPositionEvent is a event on the position change of a train.
type TrainPositionEvent struct {
	// Id of the train.
	Id int
	// NewSpeed is the new speed of the train ranging from 0 to 2.
	NewPosition int
}
