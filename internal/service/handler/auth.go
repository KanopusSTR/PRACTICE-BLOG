package handlers

import (
	"net/http"
	"server/internal/models"
)

func (hs *handlerService) Login(fun func() (models.LoginRequest, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "authorization error: " + err.Error()}
	}
	if err := hs.users.Authorization(req.Mail, req.Password); err != nil {
		return http.StatusUnprocessableEntity, models.Response{Message: "authorization error: " + err.Error()}
	}
	token, err := hs.token.CreateToken(req.Mail)
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "server error"}
	}
	return http.StatusOK, models.Response{Message: "success", Data: token}
}

func (hs *handlerService) Register(fun func() (models.RegisterRequest, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "register error: " + err.Error()}
	}
	if err := hs.users.Register(req.Name, req.Mail, req.Password); err != nil {
		return http.StatusUnprocessableEntity, models.Response{Message: "register error: " + err.Error()}
	}
	return http.StatusOK, models.Response{Message: "success"}
}
