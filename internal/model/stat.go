package model

import "github.com/google/uuid"

type Stat struct {
	BannerID uuid.UUID
	SlotID   uuid.UUID
	GroupID  uuid.UUID
	Shows    int
	Clicks   int
}
