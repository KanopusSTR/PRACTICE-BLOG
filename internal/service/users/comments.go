package users

import (
	"server/internal/entities"
	"server/pkg/myErrors"
	"time"
)

func (s *service) WriteComment(text *string, date time.Time, authorMail string, postId int) error {
	if _, e := s.posts.GetPost(postId); e != nil {
		return e
	}
	if text == nil || *text == "" {
		return myErrors.EmptyField
	}
	s.comments.Add(text, date, authorMail, postId)
	return nil
}

func (s *service) DeleteComment(postId, commentId int) error {
	if _, e := s.posts.GetPost(postId); e != nil {
		return e
	}
	return s.comments.Remove(postId, commentId)
}

func (s *service) GetComments(postId int) ([]interface{}, error) {
	if _, e := s.posts.GetPost(postId); e != nil {
		return nil, e
	}
	return s.comments.GetPostComments(postId), nil
}

func (s *service) GetComment(postId, commentId int) (*entities.Comment, error) {
	if _, e := s.posts.GetPost(postId); e != nil {
		return nil, e
	}
	return s.comments.GetPostComment(postId, commentId)
}
