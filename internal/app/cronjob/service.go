package cronjob

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type service struct {
	rewardRepo       repository.RewardStorer
	appreciationRepo repository.AppreciationStorer

}

type Service interface {
	UpdateAppreciation(ctx context.Context,tx rune
	)

}

func NewService(rewardRepo repository.RewardStorer, appreciationRepo repository.AppreciationStorer) Service {
	return &service{
		rewardRepo:       rewardRepo,
		appreciationRepo: appreciationRepo,
	}
}
