package validations

import "github.com/joshsoftware/peerly-backend/internal/pkg/dto"

func CreateAppreciation(appr dto.Appreciation) (errorResponse dto.ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	if appr.CoreValueID <= 0 {
		fieldErrors["core_value_id"] = "enter valid core value id"
	}

	if appr.Description == "" {
		fieldErrors["description"] = "enter description"
	}

	if appr.Receiver <= 0 {
		fieldErrors["receiver"] = "enter valid receiver id"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = dto.ErrorResponse{
		Error: dto.ErrorObject{
			Code:          "invalid_data",
			MessageObject: dto.MessageObject{Message: "Please provide valid appreciation data"},
			Fields:        fieldErrors,
		},
	}

	return
}
