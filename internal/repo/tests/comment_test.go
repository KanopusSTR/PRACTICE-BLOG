package tests

import (
	"errors"
	"github.com/stretchr/testify/require"
	"server/internal/entities"
	"server/internal/repo"
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

	commentRepo := repo.NewComment()

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			id := commentRepo.Add(&tc.text, tc.date, tc.mail, tc.postId)
			comment, err := commentRepo.GetPostComment(tc.postId, id)
			require.Nil(t, err)
			require.Equal(t, &entities.Comment{
				CommentId:  id,
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
			commentRepo := repo.NewComment()
			_ = commentRepo.Add(&base.text, base.date, base.mail, tc.postId)
			err := commentRepo.Remove(tc.postId, tc.commentId)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
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
			commentRepo := repo.NewComment()
			for _, comment := range tc.comments {
				_ = commentRepo.Add(&comment.Text, comment.Date, comment.AuthorMail, comment.PostId)
			}
			comments := commentRepo.GetPostComments(tc.postId)
			require.Equal(t, len(tc.comments), len(comments))
			for i, comment := range tc.comments {
				require.Equal(t, &comment, comments[i].(*entities.Comment))
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
			commentRepo := repo.NewComment()
			for _, comment := range tc.comments {
				_ = commentRepo.Add(&comment.Text, comment.Date, comment.AuthorMail, comment.PostId)
			}
			comment, err := commentRepo.GetPostComment(tc.postId, tc.commentId)
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				require.Equal(t, &tc.comments[tc.commentId], comment)
			}
		})
	}
}
