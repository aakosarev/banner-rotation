package model

import "github.com/google/uuid"

type Stat struct {
	BannerID uuid.UUID `json:"banner_id" db:"banner_id"`
	SlotID   uuid.UUID `json:"slot_id" db:"slot_id"`
	GroupID  uuid.UUID `json:"group_id" db:"social_group_id"`
	Shows    int       `json:"shows" db:"shows"`
	Clicks   int       `json:"clicks" db:"clicks"`
}
