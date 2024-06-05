package models

import "github.com/google/uuid"

type Message struct {
	Id       uuid.UUID `json:"id"`
	Command  string    `json:"command"`
	Schedule string    `json:"scheduler"`
}
