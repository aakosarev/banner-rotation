package model

import "github.com/google/uuid"

type Slot struct {
	ID          uuid.UUID
	Description string
}
