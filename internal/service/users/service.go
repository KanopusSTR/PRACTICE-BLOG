package users

import (
	"server/internal/entities"
	"server/internal/repo"
	"time"
)

type Service interface {
	Authorization(mail, password string) error
	Register(name, mail, password string) error
	GetProfile(mail string) (*entities.User, error)

	WritePost(header, body *string, date time.Time, authorMail string) error
	EditPost(id int, header, body *string) error
	DeletePost(postId int) error
	GetPosts() ([]interface{}, error)
	GetPost(postId int) (*entities.Post, error)

	WriteComment(text *string, date time.Time, authorMail string, postId int) error
	DeleteComment(postId, commentId int) error
	GetComments(postId int) ([]interface{}, error)
	GetComment(postId, commentId int) (*entities.Comment, error)
}

type service struct {
	users    repo.User
	posts    repo.Post
	comments repo.Comment
}

func New(users repo.User, posts repo.Post, comments repo.Comment) Service {
	return &service{
		users:    users,
		posts:    posts,
		comments: comments,
	}
}
