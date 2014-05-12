package models

import (
	"time"
)

type Task struct {
	ID           string    `json:"id"`
	ParentID     string    `json:"parentID"`
	VersionID    string    `json:"versionID"`
	IsPaid       bool      `json:"isPaid"`
	Hours        float64   `json:"hours"`
	SubTasks     []*Task   `json:"subTasks"`
	Title        string    `json:"title"`
	DateExpected time.Time `json:"dateExpected"`
	LastAction   *Action   `json:"lastAction"`
}

type Action struct {
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	UserID string    `json:"userID"`
	Type   string    `json:"type,omitempty"`
}
