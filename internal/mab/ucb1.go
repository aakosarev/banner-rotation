package mab

import (
	"github.com/aakosarev/banner-rotation/internal/model"
	"math"
)

func UCB1(stats []*model.Stat) *model.Stat {

	var (
		maxConfidence  float64
		rotationToShow *model.Stat
		totalShows     int64
	)

	for _, stat := range stats {
		if stat.Shows == 0 {
			return stat
		} else {
			totalShows += int64(stat.Shows)
			avgIncome := float64(stat.Clicks) / float64(stat.Shows)
			confidence := avgIncome + math.Sqrt(2*math.Log(float64(totalShows))/float64(stat.Shows))
			if confidence >= maxConfidence {
				maxConfidence = confidence
				rotationToShow = stat
			}
		}
	}
	return rotationToShow
}
