package appreciation

import (
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

// Function to map AppreciationDB to AppreciationDTO
func MapAppreciationDBToDTO(dbAppreciation repository.Appreciation) dto.Appreciation {

	return dto.Appreciation{
		ID:           dbAppreciation.ID,
		CoreValueID:  dbAppreciation.CoreValueID,
		Description:  dbAppreciation.Description,
		TotalRewards: dbAppreciation.TotalRewards,
		Quarter:      dbAppreciation.Quarter,
		Sender:       dbAppreciation.Sender,
		Receiver:     dbAppreciation.Receiver,
		CreatedAt:    dbAppreciation.CreatedAt,
		UpdatedAt:    dbAppreciation.UpdatedAt,
	}
}

func mapAppreciationInfoToResponse(info repository.AppreciationInfo) dto.ResponseAppreciation {
	return dto.ResponseAppreciation{
		ID:                  info.ID,
		CoreValueName:       info.CoreValueName,
		Description:         info.Description,
		IsValid:             info.IsValid,
		TotalRewards:        info.TotalRewards,
		Quarter:             info.Quarter,
		SenderFirstName:     info.SenderFirstName,
		SenderLastName:      info.SenderLastName,
		SenderImageURL:      info.SenderImageURL,
		SenderDesignation:   info.SenderDesignation,
		ReceiverFirstName:   info.ReceiverFirstName,
		ReceiverLastName:    info.ReceiverLastName,
		ReceiverImageURL:    info.ReceiverImageURL,
		ReceiverDesignation: info.ReceiverDesignation,
		CreatedAt:           info.CreatedAt,
		UpdatedAt:           info.UpdatedAt,
	}
}

func GetQuarter() int {
	month := int(time.Now().Month())
	if month >= 1 && month <= 3 {
		return 1
	} else if month >= 4 && month <= 6 {
		return 2
	} else if month >= 7 && month <= 9 {
		return 3
	} else if month >= 10 && month <= 12 {
		return 4
	}
	return -1
}
