package validation

import (
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

func GetIntranetUserDataValidation(user dto.IntranetUserData) (err error) {

	if user.EmpolyeeDetail.Designation.Name == "" || user.Email == "" || user.EmpolyeeDetail.Grade == "" {
		fmt.Println("Invalid user data")
		fmt.Println("First Name: ", user.FirstName)
		fmt.Println("Last Name: ", user.LastName)
		fmt.Println("Designation: ", user.EmpolyeeDetail.Designation.Name)
		fmt.Println("Email: ", user.Email)
		fmt.Println("Grade: ", user.EmpolyeeDetail.Grade)
		fmt.Println("Profile image: ", user.PublicProfile.ProfileImgUrl)

		err = apperrors.InvalidIntranetData
	}
	return
}
