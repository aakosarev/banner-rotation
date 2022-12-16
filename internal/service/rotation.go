package service

import (
	"context"
	"github.com/aakosarev/banner-rotation/internal/errors"
	"github.com/aakosarev/banner-rotation/internal/mab"
	"github.com/aakosarev/banner-rotation/internal/model"
	"github.com/google/uuid"
	"sync"
)

type storage interface {
	AddBannerToSlot(ctx context.Context, bannerID, slotID *uuid.UUID) error
	FindBannerSlot(ctx context.Context, bannerID, slotID *uuid.UUID) (*model.BannerSlot, error)
	RemoveBannerFromSlot(ctx context.Context, bannerID, slotID *uuid.UUID) error
	FindStatByParams(ctx context.Context, bannerID, slotID, socialGroupID *uuid.UUID) (*model.Stat, error)
	CreateStat(ctx context.Context, stat *model.Stat) error
	AddClickToStat(ctx context.Context, stat *model.Stat) error
	AddShowToStat(ctx context.Context, stat *model.Stat) error
	FindStatsBySlotAndSocialGroup(ctx context.Context, slotID, socialGroupID *uuid.UUID) ([]*model.Stat, error)
	FindBannersInSlot(ctx context.Context, slotID *uuid.UUID) ([]*uuid.UUID, error)
	FindBannerByID(ctx context.Context, bannerID *uuid.UUID) (*model.Banner, error)
	FindSlotByID(ctx context.Context, slotID *uuid.UUID) (*model.Slot, error)
	FindSocialGroupByID(ctx context.Context, socialGroupID *uuid.UUID) (*model.Group, error)
}

type Service struct {
	storage storage
}

func NewService(storage storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) checkBannerAndSlotExists(ctx context.Context, bannerID, slotID *uuid.UUID) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var (
		bannerErr error
		slotErr   error
	)

	go func() {
		defer wg.Done()

		var banner *model.Banner

		banner, bannerErr = s.storage.FindBannerByID(ctx, bannerID)
		if bannerErr != nil {
			return
		}

		if banner == nil {
			bannerErr = errors.ErrBannerNotFound
		}
	}()

	go func() {
		defer wg.Done()

		var slot *model.Slot

		slot, slotErr = s.storage.FindSlotByID(ctx, slotID)
		if slotErr != nil {
			return
		}

		if slot == nil {
			slotErr = errors.ErrSlotNotFound
		}
	}()

	wg.Wait()

	if bannerErr != nil {
		return bannerErr
	}

	if slotErr != nil {
		return slotErr
	}

	return nil
}

func (s *Service) AddBannerToSlot(ctx context.Context, bannerID, slotID *uuid.UUID) error {
	err := s.checkBannerAndSlotExists(ctx, bannerID, slotID)
	if err != nil {
		return err
	}

	bannerSlot, err := s.storage.FindBannerSlot(ctx, bannerID, slotID)
	if err != nil {
		return err
	}

	if bannerSlot != nil {
		return errors.ErrBannerAlreadyLinkedToSlot
	}

	return s.storage.AddBannerToSlot(ctx, bannerID, slotID)
}

func (s *Service) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID *uuid.UUID) error {
	err := s.checkBannerAndSlotExists(ctx, bannerID, slotID)
	if err != nil {
		return err
	}

	return s.storage.RemoveBannerFromSlot(ctx, bannerID, slotID)
}

func (s *Service) checkSlotAndSocialGroupExists(ctx context.Context, slotID, socialGroupID *uuid.UUID) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var (
		slotErr        error
		socialGroupErr error
	)

	go func() {
		defer wg.Done()

		var slot *model.Slot

		slot, slotErr = s.storage.FindSlotByID(ctx, slotID)
		if slotErr != nil {
			return
		}

		if slot == nil {
			slotErr = errors.ErrSlotNotFound
		}
	}()

	go func() {
		defer wg.Done()

		var socialGroup *model.Group

		socialGroup, socialGroupErr = s.storage.FindSocialGroupByID(ctx, socialGroupID)
		if socialGroupErr != nil {
			return
		}

		if socialGroup == nil {
			socialGroupErr = errors.ErrSocialGroupNotFound
		}
	}()

	wg.Wait()

	if slotErr != nil {
		return slotErr
	}

	if socialGroupErr != nil {
		return socialGroupErr
	}

	return nil
}

func (s *Service) SelectBanner(ctx context.Context, slotID, socialGroupID *uuid.UUID) (*model.Banner, error) {
	err := s.checkSlotAndSocialGroupExists(ctx, slotID, socialGroupID)
	if err != nil {
		return nil, err
	}

	stats, err := s.storage.FindStatsBySlotAndSocialGroup(ctx, slotID, socialGroupID)
	if err != nil {
		return nil, err
	}

	bannerIDs, err := s.storage.FindBannersInSlot(ctx, slotID)
	if err != nil {
		return nil, err
	}

	if len(bannerIDs) == 0 {
		return nil, errors.ErrNoOneBannerFoundForSlot
	}

	var (
		statsWithLink []*model.Stat
	)

nextBanner:
	for _, bannerID := range bannerIDs {
		for _, stat := range stats {
			if stat.BannerID == *bannerID {
				statsWithLink = append(statsWithLink, stat)
				continue nextBanner
			}
		}

		statsWithLink = append(statsWithLink, &model.Stat{
			BannerID: *bannerID,
			SlotID:   *slotID,
			GroupID:  *socialGroupID,
			Shows:    0,
			Clicks:   0,
		})
	}

	selectedStat := mab.UCB1(statsWithLink)

	selectedBanner, err := s.storage.FindBannerByID(ctx, &selectedStat.BannerID)
	if err != nil {
		return nil, err
	}

	bdStat, err := s.storage.FindStatByParams(ctx, &selectedStat.BannerID, &selectedStat.SlotID, &selectedStat.GroupID)
	if err != nil {
		return nil, err
	}
	if bdStat == nil {
		err = s.storage.CreateStat(ctx, selectedStat)
		if err != nil {
			return nil, err
		}
	}

	err = s.storage.AddShowToStat(ctx, selectedStat)
	if err != nil {
		return nil, err
	}

	return selectedBanner, nil
}

func (s *Service) checkBannerAndSlotAndSocialGroupExists(ctx context.Context, bannerID, slotID, socialGroupID *uuid.UUID) error {
	wg := sync.WaitGroup{}
	wg.Add(3)

	var (
		bannerErr      error
		slotErr        error
		socialGroupErr error
	)

	go func() {
		defer wg.Done()

		var banner *model.Banner

		banner, bannerErr = s.storage.FindBannerByID(ctx, bannerID)
		if bannerErr != nil {
			return
		}

		if banner == nil {
			bannerErr = errors.ErrBannerNotFound
		}
	}()

	go func() {
		defer wg.Done()

		var slot *model.Slot

		slot, slotErr = s.storage.FindSlotByID(ctx, slotID)
		if slotErr != nil {
			return
		}

		if slot == nil {
			slotErr = errors.ErrSlotNotFound
		}
	}()

	go func() {
		defer wg.Done()

		var socialGroup *model.Group

		socialGroup, socialGroupErr = s.storage.FindSocialGroupByID(ctx, socialGroupID)
		if socialGroupErr != nil {
			return
		}

		if socialGroup == nil {
			socialGroupErr = errors.ErrSocialGroupNotFound
		}
	}()

	wg.Wait()

	if bannerErr != nil {
		return bannerErr
	}

	if slotErr != nil {
		return slotErr
	}

	if socialGroupErr != nil {
		return socialGroupErr
	}

	return nil
}

func (s *Service) AddClick(ctx context.Context, bannerID, slotID, socialGroupID *uuid.UUID) error {

	err := s.checkBannerAndSlotAndSocialGroupExists(ctx, bannerID, slotID, socialGroupID)
	if err != nil {
		return err
	}

	stat, err := s.storage.FindStatByParams(ctx, bannerID, slotID, socialGroupID)
	if err != nil {
		return err
	}

	if stat == nil {
		err = s.storage.CreateStat(ctx, &model.Stat{
			BannerID: *bannerID,
			SlotID:   *slotID,
			GroupID:  *socialGroupID,
			Shows:    0,
			Clicks:   1,
		})
		if err != nil {
			return err
		}

		return nil
	}

	return s.storage.AddClickToStat(ctx, stat)
}
