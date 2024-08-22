package internaldb

import (
	"server/internal/entities"
	"server/internal/repo"
	"server/pkg/myErrors"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

type post struct {
	posts  *treemap.Map
	postId int
}

func NewPost() repo.Post {
	return &post{posts: treemap.NewWithIntComparator(), postId: -1}
}

func (repo *post) Add(header, body *string, date time.Time, authorMail string) error {
	repo.postId++
	repo.posts.Put(repo.postId, &entities.Post{
		Id:         repo.postId,
		Header:     *header,
		Body:       *body,
		Date:       date,
		AuthorMail: authorMail})
	return nil
}

func (repo *post) Remove(postId int) error {
	if _, found := repo.posts.Get(postId); found {
		repo.posts.Remove(postId)
		return nil
	}
	return myErrors.PostNotFound
}

func (repo *post) Update(postId int, header, body *string) error {
	if post, found := repo.posts.Get(postId); found {
		post.(*entities.Post).Body = *body
		post.(*entities.Post).Header = *header
		return nil
	}
	return myErrors.PostNotFound
}

func (repo *post) GetPosts() ([]interface{}, error) {
	return repo.posts.Values(), nil
}

func (repo *post) GetPost(id int) (*entities.Post, error) {
	if post, found := repo.posts.Get(id); found {
		return post.(*entities.Post), nil
	}
	return nil, myErrors.PostNotFound
}
