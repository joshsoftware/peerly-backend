package validation

import (
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

func GetIntranetUserDataValidation(user dto.IntranetUserData) (err error) {

	if user.PublicProfile.FirstName == "" || user.PublicProfile.LastName == "" || user.EmpolyeeDetail.Designation.Name == "" || user.Email == "" || user.EmpolyeeDetail.Grade == "" {
		fmt.Println("Invalid user data")
		fmt.Println("First Name: ", user.PublicProfile.FirstName)
		fmt.Println("Last Name: ", user.PublicProfile.LastName)
		fmt.Println("Designation: ", user.EmpolyeeDetail.Designation.Name)
		fmt.Println("Email: ", user.Email)
		fmt.Println("Grade: ", user.EmpolyeeDetail.Grade)
		fmt.Println("Profile image: ", user.PublicProfile.ProfileImgUrl)

		err = apperrors.InvalidIntranetData
	}
	return
}
