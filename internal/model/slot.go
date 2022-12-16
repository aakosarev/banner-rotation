package model

import "github.com/google/uuid"

type Slot struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}
