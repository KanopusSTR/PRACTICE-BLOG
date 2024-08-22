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

func TestAddComment(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName string
		text     string
		date     time.Time
		mail     string
		postId   int
	}{
		{"simple", "hoho", time.Now(), "hohoho@hoho.com", 0},
		{"chinese", "尽快。", time.Now(), "hohoho@hoho.com", 0},
		{"arabian", "الكتاب على الطاولة", time.Now(), "hohoho@hoho.com", 0},
		{"empty", "", time.Now(), "hohoho@hoho.com", 0},
	}

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)

			commentRepo := psql.NewComment(mock)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_comment($1, $2, $3, $4)")).
				WithArgs(tc.text, tc.date, tc.mail, tc.postId).WillReturnRows()

			err = commentRepo.Add(&tc.text, tc.date, tc.mail, tc.postId)
			require.Nil(t, err)
			rows := mock.NewRows([]string{"id", "text", "date", "mail", "postId"}).
				AddRow(0, tc.text, tc.date, tc.mail, tc.postId)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_comment($1)")).
				WithArgs(0).WillReturnRows(rows)

			comment, err := commentRepo.GetPostComment(tc.postId, 0)
			require.Nil(t, err)
			require.Equal(t, &entities.Comment{
				CommentId:  0,
				Text:       tc.text,
				Date:       tc.date,
				AuthorMail: tc.mail,
				PostId:     tc.postId,
			}, comment)
		})
	}
}

func TestRemoveComment(t *testing.T) {
	t.Helper()
	base := struct {
		text string
		date time.Time
		mail string
	}{"hoho", time.Now(), "hohoho@hoho.com"}
	testCasesComments := []struct {
		testName     string
		postId       int
		commentId    int
		errorMessage error
	}{
		{"comment exist", 0, 0, nil},
		{"comment doesn't exist", 0, 1, myErrors.CommentNotFound},
	}

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)

			commentRepo := psql.NewComment(mock)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_comment($1, $2, $3, $4)")).
				WithArgs(base.text, base.date, base.mail, tc.postId).WillReturnRows()

			err = commentRepo.Add(&base.text, base.date, base.mail, tc.postId)
			require.Nil(t, err)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM delete_comment($1)")).
				WithArgs(tc.commentId).WillReturnRows().WillReturnError(tc.errorMessage)

			err = commentRepo.Remove(tc.postId, tc.commentId)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {

				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_comment($1)")).WithArgs(0).
					WillReturnRows().WillReturnError(myErrors.CommentNotFound)
				_, err := commentRepo.GetPostComment(tc.postId, 0)
				if !errors.Is(err, myErrors.CommentNotFound) {
					t.Error("expected error, but got comment")
				}
			}
		})
	}
}

func TestGetPostComments(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName string
		postId   int
		comments []entities.Comment
	}{
		{
			testName: "two comments",
			comments: []entities.Comment{
				{0, "hoho1", time.Now(), "hohoho@hoho.com", 0},
				{1, "hoho2", time.Now(), "hohoho@hoho.com", 0}},
		},
		{
			testName: "one comment",
			comments: []entities.Comment{
				{0, "hoho", time.Now(), "hohoho@hoho.com", 0}},
		},
		{
			testName: "tree comments",
			comments: []entities.Comment{
				{0, "hoho1", time.Now(), "hohoho@hoho.com", 0},
				{1, "hoho2", time.Now(), "hohoho@hoho.com", 0},
				{2, "hoho3", time.Now(), "hohoho@hoho.com", 0}},
		},
		{
			testName: "zero comments",
			comments: []entities.Comment{},
		},
	}

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)

			commentRepo := psql.NewComment(mock)
			rows := mock.NewRows([]string{"id", "text", "date", "mail", "postId"})
			for _, comment := range tc.comments {

				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_comment($1, $2, $3, $4)")).
					WithArgs(comment.Text, comment.Date, comment.AuthorMail, comment.PostId).WillReturnRows()
				rows.AddRow(comment.CommentId, comment.Text, comment.Date, comment.AuthorMail, comment.PostId)
				_ = commentRepo.Add(&comment.Text, comment.Date, comment.AuthorMail, comment.PostId)
			}
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_comments($1)")).WithArgs(tc.postId).WillReturnRows(rows)
			comments, err := commentRepo.GetPostComments(tc.postId)
			require.Nil(t, err)
			require.Equal(t, len(tc.comments), len(comments))
			for i, comment := range tc.comments {
				require.Equal(t, comment, comments[i].(entities.Comment))
			}
		})
	}
}

func TestGetPostComment(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName     string
		postId       int
		commentId    int
		errorMessage error
		comments     []entities.Comment
	}{
		{
			testName:     "simple",
			commentId:    0,
			errorMessage: nil,
			comments: []entities.Comment{
				{0, "hoho", time.Now(), "hohoho@hoho.com", 0}},
		},
		{
			testName:     "not exist comment",
			commentId:    1,
			errorMessage: myErrors.CommentNotFound,
			comments: []entities.Comment{
				{0, "hoho", time.Now(), "hohoho@hoho.com", 0}},
		},
		{
			testName:     "tree comments",
			commentId:    1,
			errorMessage: nil,
			comments: []entities.Comment{
				{0, "hoho1", time.Now(), "hohoho@hoho.com", 0},
				{1, "hoho2", time.Now(), "hohoho@hoho.com", 0},
				{2, "hoho3", time.Now(), "hohoho@hoho.com", 0}},
		},
	}

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)

			commentRepo := psql.NewComment(mock)
			for _, comment := range tc.comments {
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_comment($1, $2, $3, $4)")).
					WithArgs(comment.Text, comment.Date, comment.AuthorMail, comment.PostId).WillReturnRows()
				_ = commentRepo.Add(&comment.Text, comment.Date, comment.AuthorMail, comment.PostId)
			}
			rows := pgxmock.NewRows([]string{"id", "text", "date", "mail", "postId"})
			if tc.testName != "not exist comment" {
				com := tc.comments[tc.commentId]
				rows.AddRow(tc.commentId, com.Text, com.Date, com.AuthorMail, com.PostId)
			}
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_comment($1)")).WithArgs(tc.commentId).WillReturnRows(rows).WillReturnError(tc.errorMessage)
			comment, err := commentRepo.GetPostComment(tc.postId, tc.commentId)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				require.Equal(t, &tc.comments[tc.commentId], comment)
			}
		})
	}
}
