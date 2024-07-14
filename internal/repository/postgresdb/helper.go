package repository

import (
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/repository"
)

func GetPaginationMetaData(page int64, limit int64, totalRecords int64) repository.Pagination {
	// Calculate pagination details
	fmt.Println("totalRecords: ",totalRecords)
	fmt.Println("limit: ",limit)
	fmt.Println("page: ",page)
	totalPages := (totalRecords + limit - 1) / limit
	fmt.Println("rem: ",((totalRecords % limit) &1 ))
	fmt.Println("totalpages: ", totalPages)
	fmt.Println("total records: ", totalRecords)
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
