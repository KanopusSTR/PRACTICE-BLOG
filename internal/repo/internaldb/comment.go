package internaldb

import (
	"github.com/emirpasic/gods/maps/treemap"
	"server/internal/entities"
	"server/internal/repo"
	"server/pkg/myErrors"
	"time"
)

type comment struct {
	comments  map[int]*treemap.Map
	commentId int
}

func NewComment() repo.Comment {
	return &comment{comments: make(map[int]*treemap.Map), commentId: -1}
}

func (repo *comment) Add(text *string, date time.Time, authorMail string, postId int) error {
	if repo.comments[postId] == nil {
		repo.comments[postId] = treemap.NewWithIntComparator()
	}
	repo.commentId++
	repo.comments[postId].Put(repo.commentId, &entities.Comment{
		CommentId:  repo.commentId,
		Text:       *text,
		Date:       date,
		AuthorMail: authorMail,
		PostId:     postId})
	return nil
}

func (repo *comment) Remove(postId, commentId int) error {
	if _, found := repo.comments[postId].Get(commentId); found {
		repo.comments[postId].Remove(commentId)
		return nil
	}
	return myErrors.CommentNotFound
}

func (repo *comment) GetPostComments(postId int) ([]interface{}, error) {
	if repo.comments[postId] == nil {
		repo.comments[postId] = treemap.NewWithIntComparator()
	}
	return repo.comments[postId].Values(), nil
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
