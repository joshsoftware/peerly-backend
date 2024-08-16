package repository

import (
	"time"

	"github.com/joshsoftware/peerly-backend/internal/repository"
)

func getPaginationMetaData(page int16, limit int16, totalRecords int32) repository.Pagination {

	// Calculate pagination details
	totalPages := (totalRecords + int32(limit) - 1) / int32(limit)
	var pagination repository.Pagination

	pagination.TotalPage = int16(totalPages)
	pagination.CurrentPage = page
	pagination.RecordPerPage = limit
	pagination.TotalRecords = totalRecords
	return pagination

}

// func GetQuarterStartUnixTime() int64 {
//     month := int(time.Now().Month())
//     quarterStartMonth := ((month - 1) / 3) * 3 + 1
//     return time.Date(time.Now().Year(), time.Month(quarterStartMonth), 1, 0, 0, 0, 0, time.Local).UnixMilli()
// }

func GetQuarterStartUnixTime() int64 {
	// Example function to get the Unix timestamp of the start of the quarter
	now := time.Now()
	quarterStart := time.Date(now.Year(), (now.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.UTC)
	return quarterStart.Unix() * 1000 // convert to milliseconds
}
