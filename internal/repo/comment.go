package repo

import (
	"github.com/emirpasic/gods/maps/treemap"
	"server/internal/entities"
	"server/pkg/myErrors"
	"time"
)

type Comment interface {
	Add(text *string, date time.Time, authorMail string, postId int) int
	Remove(postId, commentId int) error
	GetPostComments(postId int) []interface{}
	GetPostComment(postId, commentId int) (*entities.Comment, error)
}

type comment struct {
	comments map[int]*treemap.Map
}

func NewComment() Comment {
	return &comment{comments: make(map[int]*treemap.Map)}
}

func (repo *comment) Add(text *string, date time.Time, authorMail string, postId int) int {
	if repo.comments[postId] == nil {
		repo.comments[postId] = treemap.NewWithIntComparator()
	}
	commentId := repo.comments[postId].Size()
	repo.comments[postId].Put(commentId, &entities.Comment{
		CommentId:  commentId,
		Text:       *text,
		Date:       date,
		AuthorMail: authorMail,
		PostId:     postId})
	return commentId
}

func (repo *comment) Remove(postId, commentId int) error {
	if _, found := repo.comments[postId].Get(commentId); found {
		repo.comments[postId].Remove(commentId)
		return nil
	}
	return myErrors.CommentNotFound
}

func (repo *comment) GetPostComments(postId int) []interface{} {
	if repo.comments[postId] == nil {
		repo.comments[postId] = treemap.NewWithIntComparator()
	}
	return repo.comments[postId].Values()
}

func (repo *comment) GetPostComment(postId, commentId int) (*entities.Comment, error) {
	if repo.comments[postId] == nil {
		return nil, myErrors.CommentNotFound
	}
	if comment, found := repo.comments[postId].Get(commentId); found {
		return comment.(*entities.Comment), nil
	}
	return nil, myErrors.CommentNotFound
}
