package appreciation

import (
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

func mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(info repository.AppreciationResponse) dto.AppreciationResponse {

	receiverImageURL := info.ReceiverImageURL.String
	senderImageURL := info.SenderImageURL.String

	dtoApprResp := dto.AppreciationResponse{
		ID:                  info.ID,
		CoreValueName:       info.CoreValueName,
		CoreValueDesc:       info.CoreValueDesc,
		Description:         info.Description,
		TotalRewardPoints:   info.TotalRewardPoints,
		Quarter:             info.Quarter,
		SenderID:            info.SenderID,
		SenderFirstName:     info.SenderFirstName,
		SenderLastName:      info.SenderLastName,
		SenderImageURL:      senderImageURL,
		SenderDesignation:   info.SenderDesignation,
		ReceiverID:          info.ReceiverID,
		ReceiverFirstName:   info.ReceiverFirstName,
		ReceiverLastName:    info.ReceiverLastName,
		ReceiverImageURL:    receiverImageURL,
		ReceiverDesignation: info.ReceiverDesignation,
		TotalRewards:        info.TotalRewards,
		GivenRewardPoint:    info.GivenRewardPoint,
		ReportedFlag:        info.ReportedFlag,
		CreatedAt:           info.CreatedAt,
		UpdatedAt:           info.UpdatedAt,
	}

	return dtoApprResp
}

// DtoPagination returns modified response pagination struct
func dtoPagination(pagination repository.Pagination) dto.Pagination {
	return dto.Pagination{
		CurrentPage:  pagination.CurrentPage,
		TotalPage:    pagination.TotalPage,
		PageSize:     pagination.RecordPerPage,
		TotalRecords: pagination.TotalRecords,
	}
}
