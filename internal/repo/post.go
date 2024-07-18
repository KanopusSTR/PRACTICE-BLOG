package repo

import (
	"server/internal/entities"
	"server/pkg/myErrors"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

type Post interface {
	Add(header, body *string, date time.Time, authorMail string) int
	Remove(postId int) error
	Update(postId int, header, body *string) error
	GetPost(id int) (*entities.Post, error)
	GetPosts() []interface{}
}

type post struct {
	posts  *treemap.Map
	postId int
}

func NewPost() Post {
	return &post{posts: treemap.NewWithIntComparator(), postId: -1}
}

func (repo *post) Add(header, body *string, date time.Time, authorMail string) int {
	repo.postId++
	repo.posts.Put(repo.postId, &entities.Post{
		PostId:     repo.postId,
		Header:     *header,
		Body:       *body,
		Date:       date,
		AuthorMail: authorMail})
	return repo.postId
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

func (repo *post) GetPosts() []interface{} {
	return repo.posts.Values()
}

func (repo *post) GetPost(id int) (*entities.Post, error) {
	if post, found := repo.posts.Get(id); found {
		return post.(*entities.Post), nil
	}
	return nil, myErrors.PostNotFound
}
