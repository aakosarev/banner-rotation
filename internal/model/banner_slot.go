package model

import "github.com/google/uuid"

type BannerSlot struct {
	BannerID uuid.UUID
	SlotID   uuid.UUID
}
