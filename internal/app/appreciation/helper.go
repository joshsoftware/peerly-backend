package appreciation

import (
	"fmt"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

// Function to map AppreciationDB to AppreciationDTO
func MapAppreciationDBToDTO(dbAppreciation repository.Appreciation) dto.Appreciation {

	fmt.Println("db: ",dbAppreciation)
	// fmt.Println("dto: ",)
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

func mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(info repository.AppreciationInfo) dto.ResponseAppreciation {

	receiverImageURL := ""
	if info.ReceiverImageURL.Valid {
		receiverImageURL = info.ReceiverImageURL.String
	}

	senderImageURL := ""
	if info.SenderImageURL.Valid {
		senderImageURL = info.SenderImageURL.String
	}

	var dtoApprResp dto.ResponseAppreciation

	dtoApprResp.ID = info.ID
	dtoApprResp.CoreValueName = info.CoreValueName
	dtoApprResp.Description = info.Description
	dtoApprResp.TotalRewards = info.TotalRewards
	dtoApprResp.Quarter = info.Quarter
	dtoApprResp.SenderFirstName = info.SenderFirstName
	dtoApprResp.SenderLastName = info.SenderLastName
	dtoApprResp.SenderImageURL = senderImageURL
	dtoApprResp.SenderDesignation = info.SenderDesignation
	dtoApprResp.ReceiverFirstName = info.ReceiverFirstName
	dtoApprResp.ReceiverLastName = info.ReceiverLastName
	dtoApprResp.ReceiverImageURL = receiverImageURL
	dtoApprResp.ReceiverDesignation = info.ReceiverDesignation
	dtoApprResp.CreatedAt = info.CreatedAt
	dtoApprResp.UpdatedAt = info.UpdatedAt
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

func DtoPagination (pagination repository.Pagination)dto.Pagination {
	var pagenationResp dto.Pagination
	pagenationResp.CurrentPage = pagination.CurrentPage
	pagenationResp.Next = pagination.Next
	pagenationResp.Previous = pagination.Previous
	pagenationResp.RecordPerPage = pagination.RecordPerPage
	pagenationResp.TotalPage =pagination.TotalPage
	pagenationResp.TotalRecords = pagination.TotalRecords
	return pagenationResp
}