package tests

import (
	"errors"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
	"regexp"
	"server/internal/entities"
	"server/internal/repo/psql"
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
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			postRepo := psql.NewPost(mock)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_post($1, $2, $3, $4)")).
				WithArgs(tc.header, tc.body, tc.date, tc.mail).
				WillReturnRows().WillReturnError(tc.errorMessage)
			err = postRepo.Add(&tc.header, &tc.body, tc.date, tc.mail)
			require.Nil(t, err)
			rows := mock.NewRows([]string{"id", "header", "body", "date", "mail"}).
				AddRow(0, tc.header, tc.body, tc.date, tc.mail)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_post($1)")).WithArgs(0).WillReturnRows(rows)
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
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)
			postRepo := psql.NewPost(mock)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_post($1, $2, $3, $4)")).
				WithArgs(base.header, base.body, base.date, base.mail).WillReturnRows()
			_ = postRepo.Add(&base.header, &base.body, base.date, base.mail)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM delete_post($1)")).
				WithArgs(tc.postId).WillReturnRows().WillReturnError(tc.errorMessage)
			err = postRepo.Remove(tc.postId)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_post($1)")).
					WithArgs(tc.postId).WillReturnRows().WillReturnError(myErrors.PostNotFound)
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
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)
			postRepo := psql.NewPost(mock)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_post($1, $2, $3, $4)")).
				WithArgs(base.header, base.body, base.date, base.mail).WillReturnRows()
			_ = postRepo.Add(&base.header, &base.body, base.date, base.mail)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM update_post($1, $2, $3)")).
				WithArgs(tc.postId, newBase.header, newBase.body).WillReturnRows().WillReturnError(tc.errorMessage)
			err = postRepo.Update(tc.postId, &newBase.header, &newBase.body)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				rows := mock.NewRows([]string{"id", "header", "body", "date", "mail"}).
					AddRow(0, newBase.header, newBase.body, nil, nil)
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_post($1)")).
					WithArgs(tc.postId).WillReturnRows(rows).WillReturnError(nil)
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
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)
			postRepo := psql.NewPost(mock)
			rows := mock.NewRows([]string{"id", "header", "body", "date", "mail"})
			for id, post := range tc.posts {
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_post($1, $2, $3, $4)")).
					WithArgs(post.Header, post.Body, post.Date, post.AuthorMail).WillReturnRows()
				rows.AddRow(id, post.Header, post.Body, post.Date, post.AuthorMail)
				_ = postRepo.Add(&post.Header, &post.Body, post.Date, post.AuthorMail)
			}
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_posts()")).WillReturnRows(rows)
			posts, err := postRepo.GetPosts()
			require.Nil(t, err)
			require.Equal(t, len(tc.posts), len(posts))
			for i, comment := range tc.posts {
				require.Equal(t, comment, posts[i].(entities.Post))
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
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)
			postRepo := psql.NewPost(mock)
			for _, post := range tc.posts {
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_post($1, $2, $3, $4)")).
					WithArgs(post.Header, post.Body, post.Date, post.AuthorMail).WillReturnRows()
				_ = postRepo.Add(&post.Header, &post.Body, post.Date, post.AuthorMail)
			}
			rows := mock.NewRows([]string{"id", "header", "body", "date", "mail"})
			if tc.testName != "not exist post" {
				rows.AddRow(tc.postId, tc.posts[tc.postId].Header, tc.posts[tc.postId].Body, tc.posts[tc.postId].Date, tc.posts[tc.postId].AuthorMail)
			}
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_post($1)")).WithArgs(tc.postId).WillReturnRows(rows).WillReturnError(tc.errorMessage)
			post, err := postRepo.GetPost(tc.postId)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				require.Equal(t, &tc.posts[tc.postId], post)
			}
		})
	}
}
