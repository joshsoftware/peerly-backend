package badges

import (
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

func mapRepoBadgeToDTOBadge(dbBadge repository.Badge)dto.Badge{
	return dto.Badge{
		ID: dbBadge.ID,
		Name: dbBadge.Name,
		RewardPoints: dbBadge.RewardPoints,
	}
}