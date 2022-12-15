package mab

import (
	"github.com/aakosarev/banner-rotation/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUCB1(t *testing.T) {
	stats := make([]*model.Stat, 0, 10)

	for i := 1; i <= cap(stats); i++ {
		stats = append(stats, &model.Stat{
			BannerID: uuid.New(),
			SlotID:   uuid.New(),
			GroupID:  uuid.New(),
		})
	}

	t.Run("all banners shown", func(t *testing.T) {
		for i := 1; i <= len(stats); i++ {
			stat := UCB1(stats)
			stat.Shows++
		}

		for _, stat := range stats {
			require.NotEqual(t, 0, stat.Shows)
		}
	})

	t.Run("the popular banner was shown more often than the others", func(t *testing.T) {
		popularBannerID := uuid.New()
		stats = append(stats, &model.Stat{
			BannerID: popularBannerID,
			SlotID:   uuid.New(),
			GroupID:  uuid.New(),
			Clicks:   100, //imitation of a large number of clicks
		})

		var (
			maxShows              int
			resultPopularBannerID uuid.UUID
		)

		for i := 1; i <= 200; i++ {
			stat := UCB1(stats)
			stat.Shows++

			if stat.Shows > maxShows {
				maxShows = stat.Shows
				resultPopularBannerID = stat.BannerID
			}
		}
		require.Equal(t, popularBannerID, resultPopularBannerID)
	})
}
