package storage

import (
	"context"
	"errors"
	"github.com/aakosarev/banner-rotation/internal/model"
	"github.com/aakosarev/banner-rotation/pkg/client/postgresql"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	client postgresql.Client
}

func NewStorage(client postgresql.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) AddBannerToSlot(ctx context.Context, bannerID, slotID *uuid.UUID) error {
	query := `
		INSERT INTO banner_slot(banner_id, slot_id)
		VALUES ($1, $2);
	`
	_, err := s.client.Exec(ctx, query, bannerID, slotID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FindBannerSlot(ctx context.Context, bannerID, slotID *uuid.UUID) (*model.BannerSlot, error) {
	query := `
		SELECT banner_id, slot_id
		FROM banner_slot
		WHERE banner_id = $1 AND slot_id = $2;
	`

	var bannerSlot model.BannerSlot

	err := pgxscan.Get(ctx, s.client, &bannerSlot, query, bannerID, slotID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &bannerSlot, nil
}

func (s *Storage) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID *uuid.UUID) error {
	query := `
		DELETE FROM banner_slot
		WHERE banner_id = $1 AND slot_id = $2
	`
	_, err := s.client.Exec(ctx, query, bannerID, slotID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FindStatByParams(ctx context.Context, bannerID, slotID, socialGroupID *uuid.UUID) (*model.Stat, error) {
	query := `
		SELECT banner_id, slot_id, social_group_id, shows, clicks
		FROM stat
		WHERE banner_id = $1 AND slot_id = $2 AND social_group_id = $3
	`

	var stat model.Stat

	err := pgxscan.Get(ctx, s.client, &stat, query, bannerID, slotID, socialGroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &stat, nil
}

func (s *Storage) CreateStat(ctx context.Context, stat *model.Stat) error {
	query := `
		INSERT INTO stat(banner_id, slot_id, social_group_id, shows, clicks)
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err := s.client.Exec(ctx, query, stat.BannerID, stat.SlotID, stat.GroupID, stat.Shows, stat.Clicks)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddClickToStat(ctx context.Context, stat *model.Stat) error {
	query := `
		UPDATE stat
		SET clicks = clicks + 1
		WHERE banner_id = $1 AND slot_id = $2 AND social_group_id = $3
	`

	_, err := s.client.Exec(ctx, query, stat.BannerID, stat.SlotID, stat.GroupID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddShowToStat(ctx context.Context, stat *model.Stat) error {
	query := `
		UPDATE stat
		SET shows = shows + 1
		WHERE banner_id = $1 AND slot_id = $2 AND social_group_id = $3
	`

	_, err := s.client.Exec(ctx, query, stat.BannerID, stat.SlotID, stat.GroupID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FindStatsBySlotAndSocialGroup(ctx context.Context, slotID, socialGroupID *uuid.UUID) ([]*model.Stat, error) {
	query := `
		SELECT banner_id, slot_id, social_group_id, shows, clicks
		FROM stat
		WHERE slot_id = $1 AND social_group_id = $2
	`

	var stats []*model.Stat

	err := pgxscan.Select(ctx, s.client, &stats, query, slotID, socialGroupID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *Storage) FindBannersInSlot(ctx context.Context, slotID *uuid.UUID) ([]*uuid.UUID, error) {
	query := `
		SELECT banner_id
		FROM banner_slot
		WHERE slot_id = $1
	`

	var bannerIDs []*uuid.UUID

	err := pgxscan.Select(ctx, s.client, &bannerIDs, query, slotID)
	if err != nil {
		return nil, err
	}

	return bannerIDs, nil
}

func (s *Storage) FindBannerByID(ctx context.Context, bannerID *uuid.UUID) (*model.Banner, error) {
	query := `
		SELECT id, description
		FROM banner
		WHERE id = $1
	`
	var banner model.Banner

	err := pgxscan.Get(ctx, s.client, &banner, query, bannerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &banner, nil
}

func (s *Storage) FindSlotByID(ctx context.Context, slotID *uuid.UUID) (*model.Slot, error) {
	query := `
		SELECT id, description
		FROM slot
		WHERE id = $1
	`

	var slot model.Slot

	err := pgxscan.Get(ctx, s.client, &slot, query, slotID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &slot, nil
}

func (s *Storage) FindSocialGroupByID(ctx context.Context, socialGroupID *uuid.UUID) (*model.Group, error) {
	query := `
		SELECT id, description
		FROM social_group
		WHERE id = $1
	`

	var socialGroup model.Group

	err := pgxscan.Get(ctx, s.client, &socialGroup, query, socialGroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &socialGroup, nil
}
