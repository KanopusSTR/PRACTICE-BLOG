package app

import (
	"server/internal/api/handlers"
	"server/internal/repo"
	"server/internal/service"
)

type serviceProvider struct {
	commentRepo repo.Comment
	userRepo    repo.User
	postRepo    repo.Post

	service *service.I
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

func (s *serviceProvider) serviceImpl() *service.I {
	s.service = service.New(s.repo())
	return s.service
}

func (s *serviceProvider) handler() handlers.Handler {
	return handlers.New(s.serviceImpl())
}
