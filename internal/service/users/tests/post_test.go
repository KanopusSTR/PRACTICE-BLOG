package tests

import (
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"server/internal/service/users"
	minimock2 "server/internal/service/users/tests/minimock"
	"server/pkg/myErrors"
	"testing"
	"time"
)

func TestWritePost(t *testing.T) {
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
		{"empty body", "hoho", "", time.Now(), "hohoho@hoho.com", myErrors.EmptyPost},
		{"empty header", "", "hoho", time.Now(), "hohoho@hoho.com", myErrors.EmptyPost},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesPosts {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			postMock := minimock2.NewPostMock(mc)
			userS := users.New(nil, postMock, nil)

			if tc.errorMessage == nil {
				postMock.AddMock.Expect(&tc.header, &tc.body, tc.date, tc.mail).Return(nil)
			}
			err := userS.WritePost(&tc.header, &tc.body, tc.date, tc.mail)

			require.Equal(t, tc.errorMessage, err)
		})
	}
}

func TestEditPost(t *testing.T) {
	t.Helper()
	testCasesPosts := []struct {
		testName     string
		postId       int
		header       string
		body         string
		errorMessage error
	}{
		{"simple", 0, "hoho", "hoho", nil},
		{"empty body", 0, "hoho", "", myErrors.EmptyPost},
		{"empty header", 0, "", "hoho", myErrors.EmptyPost},
		{"post doesn't exist", 1, "hoho", "hoho", myErrors.PostNotFound},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesPosts {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			postMock := minimock2.NewPostMock(mc)
			userS := users.New(nil, postMock, nil)

			if tc.body == "" || tc.header == "" {
				require.Equal(t, tc.errorMessage, myErrors.EmptyPost)
			} else {

				if !errors.Is(tc.errorMessage, myErrors.EmptyPost) {
					postMock.UpdateMock.Expect(tc.postId, &tc.header, &tc.body).Return(tc.errorMessage)
				}
				err := userS.EditPost(tc.postId, &tc.header, &tc.body)

				require.Equal(t, tc.errorMessage, err)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	t.Helper()
	testCasesPosts := []struct {
		testName     string
		postId       int
		errorMessage error
	}{
		{"simple", 0, nil},
		{"post doesn't exist", 1, myErrors.PostNotFound},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesPosts {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			postMock := minimock2.NewPostMock(mc)
			userS := users.New(nil, postMock, nil)

			postMock.RemoveMock.Expect(tc.postId).Return(tc.errorMessage)
			err := userS.DeletePost(tc.postId)

			require.Equal(t, tc.errorMessage, err)
		})
	}
}

func TestGetPosts(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName     string
		errorMessage error
	}{
		{"simple", nil},
	}

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			postMock := minimock2.NewPostMock(mc)
			userS := users.New(nil, postMock, nil)

			postMock.GetPostsMock.Expect().Return(nil, nil)

			_, err := userS.GetPosts()
			require.Nil(t, err)
			require.Equal(t, tc.errorMessage, nil)
		})
	}
}

func TestGetPost(t *testing.T) {
	t.Helper()
	testCasesComments := []struct {
		testName     string
		postId       int
		errorMessage error
	}{
		{"simple", 0, nil},
		{"post doesn't exist", 1, myErrors.PostNotFound},
	}

	for _, tc := range testCasesComments {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			postMock := minimock2.NewPostMock(mc)
			userS := users.New(nil, postMock, nil)

			postMock.GetPostMock.Expect(tc.postId).Return(nil, tc.errorMessage)

			_, err := userS.GetPost(tc.postId)
			require.Equal(t, tc.errorMessage, err)
		})
	}
}
