package ttrn

import (
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
	OperatorTrain1 *User
	OperatorTrain2 *User
	OperatorTrain3 *User
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
