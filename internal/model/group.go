package model

import "github.com/google/uuid"

type Group struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}
