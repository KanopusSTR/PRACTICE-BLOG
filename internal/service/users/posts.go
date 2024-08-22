package users

import (
	"server/internal/entities"
	"server/pkg/myErrors"
	"time"
)

func (s *service) WritePost(header, body *string, date time.Time, authorMail string) error {
	if *header == "" || *body == "" {
		return myErrors.EmptyPost
	}
	return s.posts.Add(header, body, date, authorMail)
}

func (s *service) EditPost(id int, header, body *string) error {
	if *header == "" || *body == "" {
		return myErrors.EmptyPost
	}
	return s.posts.Update(id, header, body)
}

func (s *service) DeletePost(postId int) error {
	return s.posts.Remove(postId)
}

func (s *service) GetPosts() ([]interface{}, error) {
	return s.posts.GetPosts()
}

func (s *service) GetPost(postId int) (*entities.Post, error) {
	return s.posts.GetPost(postId)
}
