package dto

import (
	"strings"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
)

type Badge struct {
	ID           int8  `json:"id"`
	Name         string `json:"name"`
	RewardPoints int16  `json:"reward_points"`
}

func (bdg *Badge) ValidateCreateBadge()error {

	bdg.Name = strings.TrimSpace(bdg.Name)

	if bdg.Name == ""{
		return apperrors.InvalidBadgeName
	}

	if bdg.RewardPoints <= 0 {
		return apperrors.InvalidRewardPoints
	}
	return nil
}

func (bdg *Badge) ValidateUpdateBadge()error {

	bdg.Name = strings.TrimSpace(bdg.Name)

	if bdg.ID <= 0 {
		return apperrors.InvalidId
	}
	if bdg.RewardPoints < 0 {
		return apperrors.InvalidRewardPoints
	}
	return nil
}


