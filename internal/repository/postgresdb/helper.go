package repository

import (
	"time"

	"github.com/joshsoftware/peerly-backend/internal/repository"
)

func GetPaginationMetaData(page int64, limit int64, totalRecords int64) repository.Pagination {
	// Calculate pagination details

	var totalPages int64
	if limit == 0 {
		totalPages = 1
	} else {
		totalPages = (totalRecords / limit) + ((totalRecords % limit) & 1)
	}

	var pagination repository.Pagination

	// Handle next and pre
	// next
	if page < totalPages {
		next := int64(page + 1)
		pagination.Next = &next
	}

	// pre
	if page > 1 {
		previous := min(int64(page-1), int64(totalPages))
		pagination.Previous = &previous
	}

	pagination.TotalPage = totalPages
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
