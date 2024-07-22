package handlers

import (
	"net/http"
	"server/internal/models"
	"time"
)

func (hs *handlerService) DeleteComment(fun func() (models.DeleteComment, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "deleteComments error: " + err.Error()}
	}
	mail := req.Mail
	if comment, err := hs.users.GetComment(req.PostId, req.CommentId); err != nil {
		return http.StatusBadRequest, models.Response{Message: "deleteComment error: " + err.Error()}
	} else if comment.AuthorMail != mail {
		return http.StatusForbidden, models.Response{Message: "deleteComment error: you do not have permission"}
	}
	if err := hs.users.DeleteComment(req.PostId, req.CommentId); err != nil {
		return http.StatusBadRequest, models.Response{Message: "deleteComment error: " + err.Error()}
	}
	return http.StatusOK, models.Response{Message: "success"}
}

func (hs *handlerService) GetComments(fun func() (models.GetCommentsRequest, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "getComments error: " + err.Error()}
	}
	comments, err := hs.users.GetComments(req.PostId)
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "getComments error: " + err.Error()}
	}
	return http.StatusOK, models.Response{Message: "success", Data: comments}
}

func (hs *handlerService) WriteComment(fun func() (models.WriteComment, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "writeComment error: " + err.Error()}
	}
	if err := hs.users.WriteComment(&req.Text, time.Now(), req.Mail, req.PostId); err != nil {
		return http.StatusBadRequest, models.Response{Message: "writeComment error: " + err.Error()}
	}
	return http.StatusOK, models.Response{Message: "success"}
}
