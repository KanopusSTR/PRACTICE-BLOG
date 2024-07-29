package app

import (
	"server/internal/api/handlers"
	"server/internal/repo"
	hs "server/internal/service/handler"
	"server/internal/service/token"
	"server/internal/service/users"
)

type serviceProvider struct {
	commentRepo repo.Comment
	userRepo    repo.User
	postRepo    repo.Post

	usersS   users.Service
	tokenS   token.Service
	handlerS hs.Service
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) repo() (repo.User, repo.Post, repo.Comment) {
	s.postRepo = repo.NewPost()
	s.commentRepo = repo.NewComment()
	s.userRepo = repo.NewUser()
	return s.userRepo, s.postRepo, s.commentRepo
}

func (s *serviceProvider) serviceImpl() (users.Service, token.Service, hs.Service) {
	s.usersS = users.New(s.repo())
	s.tokenS = token.New()
	s.handlerS = hs.New(s.usersS, s.tokenS)
	return s.usersS, s.tokenS, s.handlerS
}

func (s *serviceProvider) handler() handlers.Handler {
	s.serviceImpl()
	return handlers.New(s.handlerS)
}
