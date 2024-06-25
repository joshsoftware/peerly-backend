package validation

import (
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

func GetIntranetUserDataValidation(user dto.IntranetApiResp) (err error) {

	if user.FirstName == "" || user.Designation == "" || user.Email == "" || user.Grade == "" || user.LastName == "" || user.ProfileImgUrl == "" {
		fmt.Println("Invalid user data")
		fmt.Println("First Name: ", user.FirstName)
		fmt.Println("Last Name: ", user.LastName)
		fmt.Println("Designation: ", user.Designation)
		fmt.Println("Email: ", user.Email)
		fmt.Println("Grade: ", user.Grade)
		fmt.Println("Profile image: ", user.ProfileImgUrl)

		err = apperrors.InvalidIntranetData
	}
	return
}
