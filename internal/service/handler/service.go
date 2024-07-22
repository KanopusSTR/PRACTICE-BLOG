package handlers

import (
	"server/internal/models"
	"server/internal/service/token"
	"server/internal/service/users"
)

type Service interface {
	WritePost(fun func() (models.WritePost, error)) (int, models.Response)
	GetPost(fun func() (models.GetPost, error)) (int, models.Response)
	GetPosts() (int, models.Response)
	EditPost(fun func() (models.EditPost, error)) (int, models.Response)
	DeletePost(fun func() (models.DeletePost, error)) (int, models.Response)

	GetUser(fun func() (models.GetUser, error)) (int, models.Response)

	LoginMiddleware(fun func() (models.LoginMiddleware, error)) (int, models.Response, string)

	WriteComment(fun func() (models.WriteComment, error)) (int, models.Response)
	GetComments(fun func() (models.GetCommentsRequest, error)) (int, models.Response)
	DeleteComment(fun func() (models.DeleteComment, error)) (int, models.Response)

	Login(fun func() (models.LoginRequest, error)) (int, models.Response)
	Register(fun func() (models.RegisterRequest, error)) (int, models.Response)
}

type handlerService struct {
	users users.Service
	token token.Service
}

func New(s users.Service, t token.Service) Service {
	return &handlerService{s, t}
}
