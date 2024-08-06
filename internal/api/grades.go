package api

import (
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/app/grades"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

func listGradesHandler(gradeSvc grades.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		resp, err := gradeSvc.ListGrades(ctx)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, 200, "grades list fetched successfully", resp)
	})
}
