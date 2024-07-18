package handlers

import (
	"net/http"
	"server/internal/models"
)

func (hs *handlerService) GetUser(fun func() (models.GetUser, error)) (int, models.GetProfileResponse) {
	req, err := fun()
	if err != nil {
		return http.StatusInternalServerError, models.GetProfileResponse{Message: "getProfile error: " + err.Error()}
	}
	profile, err := hs.users.GetProfile(req.Mail)
	if err != nil {
		return http.StatusNotFound, models.GetProfileResponse{Message: "getProfile error: " + err.Error()}
	}
	return http.StatusOK, models.GetProfileResponse{
		Message: "success",
		Profile: models.ProfileResponse{Name: profile.Name, Mail: profile.Mail},
	}

}
