package handlers

import (
	"net/http"
	"server/internal/models"
)

func (hs *handlerService) LoginMiddleware(fun func() (models.LoginMiddleware, error)) (int, models.ResultResponseBody, string) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.ResultResponseBody{Message: "middleware error: " + err.Error()}, ""
	}
	t, tkn, err := hs.token.ParseToken(req.Token)
	if err != nil {
		return http.StatusUnauthorized, models.ResultResponseBody{Message: "accessibility error: you do not have permission"}, ""
	}
	if !tkn.Valid {
		return http.StatusForbidden, models.ResultResponseBody{Message: "accessibility error: you do not have permission"}, ""
	}
	return http.StatusOK, models.ResultResponseBody{Message: "success"}, t.Mail
}
