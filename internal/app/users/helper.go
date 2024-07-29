package user

import (
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)


func MapActiveUserDbtoDto(activeUserDb repository.ActiveUser)dto.ActiveUser{
	badgeName := ""
	if activeUserDb.BadgeName.Valid {
		badgeName = activeUserDb.BadgeName.String
	}
	profileImageURL := ""
	if activeUserDb.ProfileImageURL.Valid {
		profileImageURL = activeUserDb.ProfileImageURL.String
	}
	return dto.ActiveUser{
		ID: activeUserDb.ID,
		FirstName: activeUserDb.FirstName,
		LastName: activeUserDb.LastName,
		ProfileImageURL: profileImageURL,
		BadgeName: badgeName,
		AppreciationPoints: activeUserDb.AppreciationPoints,
	}
}
