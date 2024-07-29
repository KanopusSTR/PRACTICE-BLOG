package tests

import (
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"server/internal/service/users"
	mock "server/internal/service/users/tests/minimock"
	"server/pkg/myErrors"
	"testing"
	"time"
)

func TestWriteComment(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName     string
		text         string
		date         time.Time
		mail         string
		postId       int
		errorMessage error
	}{
		{"simple", "hoho", time.Now(), "hohoho@hoho.com", 0, nil},
		{"empty", "", time.Now(), "hohoho@hoho.com", 0, myErrors.EmptyField},
		{"post doesn't exist", "hoho", time.Now(), "hohoho@hoho.com", 1, myErrors.PostNotFound},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			commentMock := mock.NewCommentMock(mc)
			postMock := mock.NewPostMock(mc)
			userMock := mock.NewUserMock(mc)
			userS := users.New(userMock, postMock, commentMock)

			if tc.postId == 0 {
				postMock.GetPostMock.Expect(0).Return(nil, nil)
			} else {
				postMock.GetPostMock.Expect(tc.postId).Return(nil, myErrors.PostNotFound)
			}
			if tc.errorMessage == nil {
				commentMock.AddMock.Expect(&tc.text, tc.date, tc.mail, tc.postId).Return(0)
			}
			err := userS.WriteComment(&tc.text, tc.date, tc.mail, tc.postId)

			require.Equal(t, tc.errorMessage, err)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName     string
		postId       int
		commentId    int
		errorMessage error
	}{
		{"simple", 0, 0, nil},
		{"post doesn't exist", 1, 0, myErrors.PostNotFound},
		{"comment doesn't exist", 0, 1, myErrors.CommentNotFound},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			commentMock := mock.NewCommentMock(mc)
			postMock := mock.NewPostMock(mc)
			userMock := mock.NewUserMock(mc)
			userS := users.New(userMock, postMock, commentMock)

			if tc.postId == 0 {
				postMock.GetPostMock.Expect(0).Return(nil, nil)
			} else {
				postMock.GetPostMock.Expect(tc.postId).Return(nil, myErrors.PostNotFound)
			}
			if !errors.Is(tc.errorMessage, myErrors.PostNotFound) {
				commentMock.RemoveMock.Expect(tc.postId, tc.commentId).Return(tc.errorMessage)
			}
			err := userS.DeleteComment(tc.postId, tc.commentId)

			require.Equal(t, tc.errorMessage, err)
		})
	}
}

func TestGetComments(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName     string
		postId       int
		errorMessage error
	}{
		{"simple", 0, nil},
		{"post doesn't exist", 1, myErrors.PostNotFound},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			commentMock := mock.NewCommentMock(mc)
			postMock := mock.NewPostMock(mc)
			userMock := mock.NewUserMock(mc)
			userS := users.New(userMock, postMock, commentMock)

			if tc.postId == 0 {
				postMock.GetPostMock.Expect(0).Return(nil, nil)
			} else {
				postMock.GetPostMock.Expect(tc.postId).Return(nil, myErrors.PostNotFound)
			}
			if !errors.Is(tc.errorMessage, myErrors.PostNotFound) {
				commentMock.GetPostCommentsMock.Expect(tc.postId).Return(nil)
			}
			_, err := userS.GetComments(tc.postId)
			require.Equal(t, tc.errorMessage, err)
		})
	}
}

func TestGetComment(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName     string
		postId       int
		commentId    int
		errorMessage error
	}{
		{"simple", 0, 0, nil},
		{"post doesn't exist", 1, 0, myErrors.PostNotFound},
		{"comment doesn't exist", 0, 1, myErrors.CommentNotFound},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			commentMock := mock.NewCommentMock(mc)
			postMock := mock.NewPostMock(mc)
			userMock := mock.NewUserMock(mc)
			userS := users.New(userMock, postMock, commentMock)

			if tc.postId == 0 {
				postMock.GetPostMock.Expect(0).Return(nil, nil)
			} else {
				postMock.GetPostMock.Expect(tc.postId).Return(nil, myErrors.PostNotFound)
			}
			if !errors.Is(tc.errorMessage, myErrors.PostNotFound) {
				commentMock.GetPostCommentMock.Expect(tc.postId, tc.commentId).Return(nil, tc.errorMessage)
			}
			_, err := userS.GetComment(tc.postId, tc.commentId)

			require.Equal(t, tc.errorMessage, err)
		})
	}
}
