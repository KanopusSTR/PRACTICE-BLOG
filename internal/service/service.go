package service

import (
	"server/internal/repo"
	handlers "server/internal/service/handler"
	"server/internal/service/token"
	"server/internal/service/users"
)

type I struct {
	Users   users.Service
	Token   token.Service
	Handler handlers.Service
}

func New(user repo.User, posts repo.Post, comments repo.Comment) *I {
	u := users.New(user, posts, comments)
	t := token.New()
	return &I{Users: u, Token: t, Handler: handlers.New(u, t)}
}
