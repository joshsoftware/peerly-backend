package validation

import (
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func GetIntranetUserDataValidation(user dto.IntranetUserData) (err error) {

	if user.PublicProfile.FirstName == "" || user.PublicProfile.LastName == "" || user.EmpolyeeDetail.Designation.Name == "" || user.Email == "" || user.EmpolyeeDetail.Grade == "" {
		logger.Errorf("\ninvalid user data\nfirst name:%s\nlast name:%s\ndesignation: %s\nemail:%s\ngrade:%s\nprofile_image:%s",
        user.PublicProfile.FirstName, 
        user.PublicProfile.LastName,
				user.EmpolyeeDetail.Designation.Name,
		  	user.Email,
				user.EmpolyeeDetail.Grade,
				user.PublicProfile.ProfileImgUrl,
      )
		err = apperrors.InvalidIntranetData
	}
	return
}
