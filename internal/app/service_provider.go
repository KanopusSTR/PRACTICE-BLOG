package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"server/internal/api/handlers"
	"server/internal/repo"
	"server/internal/repo/internaldb"
	"server/internal/repo/psql"
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

	db *pgxpool.Pool
}

func newServiceProvider(db *pgxpool.Pool) *serviceProvider {
	return &serviceProvider{db: db}
}

func (s *serviceProvider) postgresRepo() (repo.User, repo.Post, repo.Comment) {
	s.postRepo = psql.NewPost(s.db)
	s.commentRepo = psql.NewComment(s.db)
	s.userRepo = psql.NewUser(s.db)
	return s.userRepo, s.postRepo, s.commentRepo
}

func (s *serviceProvider) internalRepo() (repo.User, repo.Post, repo.Comment) {
	s.postRepo = internaldb.NewPost()
	s.commentRepo = internaldb.NewComment()
	s.userRepo = internaldb.NewUser()
	return s.userRepo, s.postRepo, s.commentRepo
}

func (s *serviceProvider) serviceImpl() (users.Service, token.Service, hs.Service) {
	s.usersS = users.New(s.postgresRepo())
	s.tokenS = token.New()
	s.handlerS = hs.New(s.usersS, s.tokenS)
	return s.usersS, s.tokenS, s.handlerS
}

func (s *serviceProvider) handler() handlers.Handler {
	s.serviceImpl()
	return handlers.New(s.handlerS)
}
