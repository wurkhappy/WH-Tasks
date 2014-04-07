package models

import (
	"time"
)

type Action struct {
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	UserID string    `json:"userID"`
	Type   string    `json:"type,omitempty"`
}

var (
	ActionCreated   string = "created"
	ActionSubmitted string = "submitted"
	ActionCompleted string = "completed"
	ActionAccepted  string = "accepted"
)

func CreatedActionForUser(userID string) *Action {
	action := createActionForUser(userID)
	action.Name = ActionCreated
	return action
}

func SubmittedActionForUser(userID string) *Action {
	action := createActionForUser(userID)
	action.Name = ActionSubmitted
	return action
}
func CompletedActionForUser(userID string) *Action {
	action := createActionForUser(userID)
	action.Name = ActionCompleted
	return action
}
func AcceptedActionForUser(userID string) *Action {
	action := createActionForUser(userID)
	action.Name = ActionAccepted
	return action
}

func createActionForUser(userID string) *Action {
	action := new(Action)
	action.Date = time.Now()
	action.UserID = userID
	return action
}
