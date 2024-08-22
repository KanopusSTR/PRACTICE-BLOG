package repo

import (
	"server/internal/entities"
	"time"
)

type Comment interface {
	Add(text *string, date time.Time, authorMail string, postId int) error
	Remove(postId, commentId int) error
	GetPostComments(postId int) ([]interface{}, error)
	GetPostComment(postId, commentId int) (*entities.Comment, error)
}
