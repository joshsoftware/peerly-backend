package repository

import (
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
