package user

import (
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

func syncData(intranetUserData dto.IntranetUserData, peerlyUserData dto.GetUserResp) (syncNeeded bool, dataToBeUpdated dto.UpdateUserData) {
	syncNeeded = false
	if intranetUserData.PublicProfile.FirstName != peerlyUserData.FirstName || intranetUserData.PublicProfile.LastName != peerlyUserData.LastName || intranetUserData.PublicProfile.ProfileImgUrl != peerlyUserData.ProfileImgUrl || intranetUserData.EmpolyeeDetail.Designation.Name != peerlyUserData.Designation || intranetUserData.EmpolyeeDetail.Grade != peerlyUserData.Grade {
		syncNeeded = true
		dataToBeUpdated.FirstName = intranetUserData.PublicProfile.FirstName
		dataToBeUpdated.LastName = intranetUserData.PublicProfile.LastName
		dataToBeUpdated.ProfileImgUrl = intranetUserData.PublicProfile.ProfileImgUrl
		dataToBeUpdated.Designation = intranetUserData.EmpolyeeDetail.Designation.Name
		dataToBeUpdated.Grade = intranetUserData.EmpolyeeDetail.Grade
		dataToBeUpdated.Email = intranetUserData.Email
	}
	return
}