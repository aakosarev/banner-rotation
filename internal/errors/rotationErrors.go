package errors

import "errors"

var (
	ErrNoOneBannerFoundForSlot   = errors.New("no banner was found for this slot")
	ErrBannerAlreadyLinkedToSlot = errors.New("banner is already linked to this slot")
	ErrBannerNotFound            = errors.New("banner not found")
	ErrSlotNotFound              = errors.New("slot not found")
	ErrSocialGroupNotFound       = errors.New("social group not found")
)
