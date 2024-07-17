package appreciation

import (
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

// Function to map AppreciationDB to AppreciationDTO
func mapAppreciationDBToDTO(dbAppreciation repository.Appreciation) dto.Appreciation {
	return dto.Appreciation{
		ID:                dbAppreciation.ID,
		CoreValueID:       dbAppreciation.CoreValueID,
		Description:       dbAppreciation.Description,
		TotalRewardPoints: dbAppreciation.TotalRewardPoints,
		Quarter:           dbAppreciation.Quarter,
		Sender:            dbAppreciation.Sender,
		Receiver:          dbAppreciation.Receiver,
		CreatedAt:         dbAppreciation.CreatedAt,
		UpdatedAt:         dbAppreciation.UpdatedAt,
	}
}

func mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(info repository.AppreciationInfo) dto.ResponseAppreciation {

	receiverImageURL := ""
	if info.ReceiverImageURL.Valid {
		receiverImageURL = info.ReceiverImageURL.String
	}

	senderImageURL := ""
	if info.SenderImageURL.Valid {
		senderImageURL = info.SenderImageURL.String
	}

	dtoApprResp := dto.ResponseAppreciation{
		ID:                  info.ID,
		CoreValueName:       info.CoreValueName,
		CoreValueDesc:       info.CoreValueDesc,
		Description:         info.Description,
		TotalRewardPoints:   info.TotalRewardPoints,
		Quarter:             info.Quarter,
		SenderFirstName:     info.SenderFirstName,
		SenderLastName:      info.SenderLastName,
		SenderImageURL:      senderImageURL,
		SenderDesignation:   info.SenderDesignation,
		ReceiverFirstName:   info.ReceiverFirstName,
		ReceiverLastName:    info.ReceiverLastName,
		ReceiverImageURL:    receiverImageURL,
		ReceiverDesignation: info.ReceiverDesignation,
		TotalRewards:        info.TotalRewards,
		GivenRewardPoint:    info.GivenRewardPoint,
		CreatedAt:           info.CreatedAt,
		UpdatedAt:           info.UpdatedAt,
	}

	return dtoApprResp
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

func DtoPagination(pagination repository.Pagination) dto.Pagination {
	var pagenationResp dto.Pagination
	pagenationResp.CurrentPage = pagination.CurrentPage
	// pagenationResp.Next = pagination.Next
	// pagenationResp.Previous = pagination.Previous
	// pagenationResp.RecordPerPage = pagination.RecordPerPage
	pagenationResp.TotalPage = pagination.TotalPage
	pagenationResp.TotalRecords = pagination.TotalRecords
	return pagenationResp
}
