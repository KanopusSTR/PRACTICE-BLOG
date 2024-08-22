package repo

import (
	"server/internal/entities"
	"time"
)

type Post interface {
	Add(header, body *string, date time.Time, authorMail string) error
	Remove(postId int) error
	Update(postId int, header, body *string) error
	GetPost(id int) (*entities.Post, error)
	GetPosts() ([]interface{}, error)
}
