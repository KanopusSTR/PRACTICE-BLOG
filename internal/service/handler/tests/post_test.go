package tests

import (
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"net/http"
	"server/internal/entities"
	"server/internal/models"
	handler "server/internal/service/handler"
	mock "server/internal/service/handler/tests/minimock"
	"testing"
)

func TestDeletePost(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.DeletePost, error)
		mail     string
		postID   int
		code     int
		ans      models.ResultResponseBody
		error
	}{
		{"success",
			func() (models.DeletePost, error) {

				return models.DeletePost{PostId: 0, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			http.StatusOK,
			models.ResultResponseBody{Message: "success"},
			nil,
		},

		{"getPostError",
			func() (models.DeletePost, error) {

				return models.DeletePost{PostId: 0, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			http.StatusNotFound,
			models.ResultResponseBody{Message: "deletePost error: "},
			errors.New(""),
		},

		{"deletePostError",
			func() (models.DeletePost, error) {

				return models.DeletePost{PostId: 0, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			http.StatusBadRequest,
			models.ResultResponseBody{Message: "deletePost error: "},
			errors.New(""),
		},

		{"noAccessError",
			func() (models.DeletePost, error) {

				return models.DeletePost{PostId: 0, Mail: "b@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			http.StatusForbidden,
			models.ResultResponseBody{Message: "deletePost error: you do not have permission"},
			errors.New(""),
		},

		{"funError",
			func() (models.DeletePost, error) {

				return models.DeletePost{PostId: 0, Mail: "a@mail.ru"}, errors.New("fun error")
			},
			"a@mail.ru",
			0,
			http.StatusBadRequest,
			models.ResultResponseBody{Message: "postId is required and it must be non-negative"},
			errors.New(""),
		},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			users := mock.NewUsersMock(mc)
			h := handler.New(users, nil)
			if tc.testName != "funError" {
				if tc.testName == "getPostError" {
					users.GetPostMock.Expect(tc.postID).Return(&entities.Post{PostId: tc.postID, AuthorMail: tc.mail}, tc.error)
				} else {
					users.GetPostMock.Expect(tc.postID).Return(&entities.Post{PostId: tc.postID, AuthorMail: tc.mail}, nil)
					if tc.testName != "noAccessError" {
						users.DeletePostMock.Expect(tc.postID).Return(tc.error)
					}
				}
			}
			c, m := h.DeletePost(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}

func TestGetPosts(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.GetCommentsRequest, error)
		mail     string
		postID   int
		posts    []interface{}
		code     int
		ans      models.GetPostsResponse
		error
	}{
		{"success",
			func() (models.GetCommentsRequest, error) {

				return models.GetCommentsRequest{PostId: 0}, nil
			},
			"a@mail.ru",
			0,
			[]interface{}{},
			http.StatusOK,
			models.GetPostsResponse{Posts: []interface{}{}},
			nil,
		},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			users := mock.NewUsersMock(mc)
			h := handler.New(users, nil)
			if tc.testName != "funError" {
				users.GetPostsMock.Expect().Return(tc.posts)
			}
			c, m := h.GetPosts()
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}

func TestGetPost(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.GetPost, error)
		mail     string
		postID   int
		post     entities.Post
		code     int
		ans      models.GetPostResponse
		error
	}{
		{"success",
			func() (models.GetPost, error) {

				return models.GetPost{Id: 0}, nil
			},
			"a@mail.ru",
			0,
			entities.Post{PostId: 0, AuthorMail: "a@mail.ru"},
			http.StatusOK,
			models.GetPostResponse{Posts: entities.Post{PostId: 0, AuthorMail: "a@mail.ru"}, Message: "success"},
			nil,
		},

		{"getPostError",
			func() (models.GetPost, error) {

				return models.GetPost{Id: 1}, nil
			},
			"a@mail.ru",
			1,
			entities.Post{PostId: 0, AuthorMail: "a@mail.ru"},
			http.StatusNotFound,
			models.GetPostResponse{Message: "getPost error: "},
			errors.New(""),
		},

		{"funError",
			func() (models.GetPost, error) {

				return models.GetPost{Id: 0}, errors.New("fun error")
			},
			"a@mail.ru",
			0,
			entities.Post{PostId: 0, AuthorMail: "a@mail.ru"},
			http.StatusBadRequest,
			models.GetPostResponse{Message: "postId is required and it must be non-negative"},
			errors.New(""),
		},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			users := mock.NewUsersMock(mc)
			h := handler.New(users, nil)
			if tc.testName != "funError" {
				users.GetPostMock.Expect(tc.postID).Return(&tc.post, tc.error)
			}
			c, m := h.GetPost(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}

func TestWritePost(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.WritePost, error)
		text     string
		mail     string
		postID   int
		code     int
		ans      models.WritePostResponse
		error
	}{
		{"success",
			func() (models.WritePost, error) {

				return models.WritePost{Mail: "a@mail.ru"}, nil
			},
			"text",
			"a@mail.ru",
			0,
			http.StatusOK,
			models.WritePostResponse{Message: "success", PostId: 0},
			nil,
		},

		{"writePostError",
			func() (models.WritePost, error) {

				return models.WritePost{Mail: "a@mail.ru"}, nil
			},
			"text",
			"a@mail.ru",
			1,
			http.StatusBadRequest,
			models.WritePostResponse{Message: "writePost error: "},
			errors.New(""),
		},

		{"funError",
			func() (models.WritePost, error) {

				return models.WritePost{Mail: "a@mail.ru"}, errors.New("fun error")
			},
			"a@mail.ru",
			"text",
			0,
			http.StatusBadRequest,
			models.WritePostResponse{Message: "writePost error: fun error"},
			errors.New(""),
		},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			users := mock.NewUsersMock(mc)
			h := handler.New(users, nil)
			if tc.testName != "funError" {
				users.WritePostMock.ExpectAuthorMailParam4(tc.mail).Return(tc.postID, tc.error)
			}
			c, m := h.WritePost(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}

func TestEditPost(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.EditPost, error)
		mail     string
		postID   int
		code     int
		ans      models.ResultResponseBody
		error
	}{
		{"success",
			func() (models.EditPost, error) {

				return models.EditPost{PostId: 0, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			http.StatusOK,
			models.ResultResponseBody{Message: "success"},
			nil,
		},

		{"getPostError",
			func() (models.EditPost, error) {

				return models.EditPost{PostId: 0, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			http.StatusNotFound,
			models.ResultResponseBody{Message: "editPost error: "},
			errors.New(""),
		},

		{"editPostError",
			func() (models.EditPost, error) {

				return models.EditPost{PostId: 0, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			http.StatusNotFound,
			models.ResultResponseBody{Message: "editPost error: "},
			errors.New(""),
		},

		{"noAccessError",
			func() (models.EditPost, error) {

				return models.EditPost{PostId: 0, Mail: "b@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			http.StatusForbidden,
			models.ResultResponseBody{Message: "editPost error: you do not have permission"},
			errors.New(""),
		},

		{"funError",
			func() (models.EditPost, error) {

				return models.EditPost{PostId: 0, Mail: "a@mail.ru"}, errors.New("fun error")
			},
			"a@mail.ru",
			0,
			http.StatusBadRequest,
			models.ResultResponseBody{Message: "editPost error: fun error"},
			errors.New(""),
		},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			users := mock.NewUsersMock(mc)
			h := handler.New(users, nil)
			if tc.testName != "funError" {
				if tc.testName == "getPostError" {
					users.GetPostMock.Expect(tc.postID).Return(&entities.Post{PostId: tc.postID, AuthorMail: tc.mail}, tc.error)
				} else {
					users.GetPostMock.Expect(tc.postID).Return(&entities.Post{PostId: tc.postID, AuthorMail: tc.mail}, nil)
					if tc.testName != "noAccessError" {
						users.EditPostMock.ExpectIdParam1(tc.postID).Return(tc.error)
					}
				}
			}
			c, m := h.EditPost(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}
