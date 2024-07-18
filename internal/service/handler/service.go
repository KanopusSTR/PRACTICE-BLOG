package handlers

import (
	"server/internal/models"
	"server/internal/service/token"
	"server/internal/service/users"
)

type Service interface {
	WritePost(fun func() (models.WritePost, error)) (int, models.WritePostResponse)
	GetPost(fun func() (models.GetPost, error)) (int, models.GetPostResponse)
	GetPosts() (int, models.GetPostsResponse)
	EditPost(fun func() (models.EditPost, error)) (int, models.ResultResponseBody)
	DeletePost(fun func() (models.DeletePost, error)) (int, models.ResultResponseBody)

	GetUser(fun func() (models.GetUser, error)) (int, models.GetProfileResponse)

	LoginMiddleware(fun func() (models.LoginMiddleware, error)) (int, models.ResultResponseBody, string)

	WriteComment(fun func() (models.WriteComment, error)) (int, models.ResultResponseBody)
	GetComments(fun func() (models.GetCommentsRequest, error)) (int, models.GetCommentsResponse)
	DeleteComment(fun func() (models.DeleteComment, error)) (int, models.ResultResponseBody)

	Login(fun func() (models.LoginRequest, error)) (int, models.LoginResponse)
	Register(fun func() (models.RegisterRequest, error)) (int, models.ResultResponseBody)
}

type handlerService struct {
	users users.Service
	token token.Service
}

func New(s users.Service, t token.Service) Service {
	return &handlerService{s, t}
}
