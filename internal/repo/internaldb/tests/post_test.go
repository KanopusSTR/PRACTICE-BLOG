package tests

import (
	"errors"
	"github.com/stretchr/testify/require"
	"server/internal/entities"
	"server/internal/repo/internaldb"
	"server/pkg/myErrors"
	"testing"
	"time"
)

func TestAddPost(t *testing.T) {
	t.Helper()
	testCasesPosts := []struct {
		testName     string
		header       string
		body         string
		date         time.Time
		mail         string
		errorMessage error
	}{
		{"simple", "hoho", "hoho", time.Now(), "hohoho@hoho.com", nil},
		{"chinese", "尽快。", "尽快。", time.Now(), "hohoho@hoho.com", nil},
		{"arabian", "على الطاولة", "الكتاب على الطاولة", time.Now(), "hohoho@hoho.com", nil},
	}

	for _, tc := range testCasesPosts {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			postRepo := internaldb.NewPost()
			err := postRepo.Add(&tc.header, &tc.body, tc.date, tc.mail)
			require.Nil(t, err)
			post, err := postRepo.GetPost(0)
			require.Nil(t, err)
			require.Equal(t, &entities.Post{
				Id:         0,
				Header:     tc.header,
				Body:       tc.body,
				Date:       tc.date,
				AuthorMail: tc.mail,
			}, post)
		})
	}
}

func TestRemovePost(t *testing.T) {
	t.Helper()
	base := struct {
		header string
		body   string
		date   time.Time
		mail   string
	}{"hoho", "hoho", time.Now(), "hohoho@hoho.com"}
	testCasesPosts := []struct {
		testName     string
		postId       int
		errorMessage error
	}{
		{"post exist", 0, nil},
		{"post doesn't exist", 1, myErrors.PostNotFound},
	}

	for _, tc := range testCasesPosts {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			postRepo := internaldb.NewPost()
			_ = postRepo.Add(&base.header, &base.body, base.date, base.mail)
			err := postRepo.Remove(tc.postId)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				_, err := postRepo.GetPost(tc.postId)
				if !errors.Is(err, myErrors.PostNotFound) {
					t.Error("expected error, but got post")
				}
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	t.Helper()
	base := struct {
		header string
		body   string
		date   time.Time
		mail   string
	}{"hoho", "hoho", time.Now(), "hohoho@hoho.com"}
	newBase := struct {
		header string
		body   string
	}{"new hoho", "new hoho"}
	testCasesPosts := []struct {
		testName     string
		postId       int
		errorMessage error
	}{
		{"post exist", 0, nil},
		{"post doesn't exist", 1, myErrors.PostNotFound},
	}

	for _, tc := range testCasesPosts {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			postRepo := internaldb.NewPost()
			_ = postRepo.Add(&base.header, &base.body, base.date, base.mail)
			err := postRepo.Update(tc.postId, &newBase.body, &newBase.body)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				post, _ := postRepo.GetPost(0)
				require.Equal(t, newBase.header, post.Header)
				require.Equal(t, newBase.body, post.Body)
			}
		})
	}
}

func TestGetPosts(t *testing.T) {
	t.Helper()
	testCasesPosts := []struct {
		testName string
		posts    []entities.Post
	}{
		{
			testName: "two posts",
			posts: []entities.Post{
				{0, "hoho1", "hoho1", time.Now(), "hohoho@hoho.com"},
				{1, "hoho2", "hoho2", time.Now(), "hohoho@hoho.com"},
			},
		},
		{
			testName: "one comment",
			posts: []entities.Post{
				{0, "hoho1", "hoho1", time.Now(), "hohoho@hoho.com"},
			},
		},
		{
			testName: "tree posts",
			posts: []entities.Post{
				{0, "hoho1", "hoho1", time.Now(), "hohoho@hoho.com"},
				{1, "hoho2", "hoho2", time.Now(), "hohoho@hoho.com"},
				{2, "hoho3", "hoho3", time.Now(), "hohoho@hoho.com"},
			},
		},
		{
			testName: "zero posts",
			posts:    []entities.Post{},
		},
	}

	for _, tc := range testCasesPosts {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			postRepo := internaldb.NewPost()
			for _, post := range tc.posts {
				_ = postRepo.Add(&post.Header, &post.Body, post.Date, post.AuthorMail)
			}
			posts, err := postRepo.GetPosts()
			require.Nil(t, err)
			require.Equal(t, len(tc.posts), len(posts))
			for i, comment := range tc.posts {
				require.Equal(t, &comment, posts[i].(*entities.Post))
			}
		})
	}
}

func TestGetPost(t *testing.T) {
	t.Helper()
	testCasesPosts := []struct {
		testName     string
		postId       int
		errorMessage error
		posts        []entities.Post
	}{
		{
			testName:     "simple",
			postId:       0,
			errorMessage: nil,
			posts: []entities.Post{
				{0, "hoho", "hoho", time.Now(), "hohoho@hoho.com"},
			},
		},
		{
			testName:     "not exist post",
			postId:       1,
			errorMessage: myErrors.PostNotFound,
			posts: []entities.Post{
				{0, "hoho", "hoho", time.Now(), "hohoho@hoho.com"},
			},
		},
		{
			testName:     "tree posts",
			postId:       1,
			errorMessage: nil,
			posts: []entities.Post{
				{0, "hoho1", "hoho1", time.Now(), "hohoho@hoho.com"},
				{1, "hoho2", "hoho2", time.Now(), "hohoho@hoho.com"},
				{2, "hoho3", "hoho3", time.Now(), "hohoho@hoho.com"},
			},
		},
	}

	for _, tc := range testCasesPosts {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			postRepo := internaldb.NewPost()
			for _, post := range tc.posts {
				_ = postRepo.Add(&post.Header, &post.Body, post.Date, post.AuthorMail)
			}
			post, err := postRepo.GetPost(tc.postId)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				require.Equal(t, &tc.posts[tc.postId], post)
			}
		})
	}
}
