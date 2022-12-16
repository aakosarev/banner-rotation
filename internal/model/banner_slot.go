package model

import "github.com/google/uuid"

type BannerSlot struct {
	BannerID uuid.UUID `json:"banner_id"`
	SlotID   uuid.UUID `json:"slot_id"`
}
