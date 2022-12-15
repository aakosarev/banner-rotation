package model

import "github.com/google/uuid"

type Banner struct {
	ID          uuid.UUID
	Description string
}
