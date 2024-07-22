package handlers

import (
	"net/http"
	"server/internal/models"
)

func (hs *handlerService) GetUser(fun func() (models.GetUser, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusInternalServerError, models.Response{Message: "getProfile error: " + err.Error()}
	}
	profile, err := hs.users.GetProfile(req.Mail)
	if err != nil {
		return http.StatusNotFound, models.Response{Message: "getProfile error: " + err.Error()}
	}
	return http.StatusOK, models.Response{
		Message: "success",
		Data:    models.ProfileResponse{Name: profile.Name, Mail: profile.Mail},
	}

}
