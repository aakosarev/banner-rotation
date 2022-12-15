package model

import "github.com/google/uuid"

type Group struct {
	ID          uuid.UUID
	Description string
}
