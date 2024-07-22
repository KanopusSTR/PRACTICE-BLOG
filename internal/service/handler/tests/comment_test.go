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

func TestDeleteComment(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName  string
		fun       func() (models.DeleteComment, error)
		mail      string
		postID    int
		commentId int
		code      int
		ans       models.Response
		error
	}{
		{"success",
			func() (models.DeleteComment, error) {

				return models.DeleteComment{PostId: 0, CommentId: 0, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			0,
			http.StatusOK,
			models.Response{Message: "success"},
			nil,
		},

		{"getCommentError",
			func() (models.DeleteComment, error) {

				return models.DeleteComment{PostId: 0, CommentId: 1, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			1,
			http.StatusBadRequest,
			models.Response{Message: "deleteComment error: "},
			errors.New(""),
		},

		{"deleteCommentError",
			func() (models.DeleteComment, error) {

				return models.DeleteComment{PostId: 0, CommentId: 1, Mail: "a@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			1,
			http.StatusBadRequest,
			models.Response{Message: "deleteComment error: "},
			errors.New(""),
		},

		{"noAccessError",
			func() (models.DeleteComment, error) {

				return models.DeleteComment{PostId: 0, CommentId: 0, Mail: "b@mail.ru"}, nil
			},
			"a@mail.ru",
			0,
			0,
			http.StatusForbidden,
			models.Response{Message: "deleteComment error: you do not have permission"},
			errors.New(""),
		},

		{"funError",
			func() (models.DeleteComment, error) {

				return models.DeleteComment{PostId: 0, CommentId: 0, Mail: "a@mail.ru"}, errors.New("fun error")
			},
			"a@mail.ru",
			0,
			0,
			http.StatusBadRequest,
			models.Response{Message: "deleteComments error: fun error"},
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
				if tc.testName == "getCommentError" {
					users.GetCommentMock.Expect(tc.postID, tc.commentId).Return(&entities.Comment{AuthorMail: "a@mail.ru"}, tc.error)
				} else {
					users.GetCommentMock.Expect(tc.postID, tc.commentId).Return(&entities.Comment{AuthorMail: "a@mail.ru"}, nil)
					if tc.testName != "noAccessError" {
						users.DeleteCommentMock.Expect(tc.postID, tc.commentId).Return(tc.error)
					}

				}
			}
			c, m := h.DeleteComment(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}

func TestGetComments(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.GetCommentsRequest, error)
		mail     string
		postID   int
		comments []interface{}
		code     int
		ans      models.Response
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
			models.Response{Message: "success", Data: []interface{}{}},
			nil,
		},

		{"getCommentsError",
			func() (models.GetCommentsRequest, error) {

				return models.GetCommentsRequest{PostId: 1}, nil
			},
			"a@mail.ru",
			1,
			[]interface{}{},
			http.StatusBadRequest,
			models.Response{Message: "getComments error: "},
			errors.New(""),
		},

		{"funError",
			func() (models.GetCommentsRequest, error) {

				return models.GetCommentsRequest{PostId: 0}, errors.New("fun error")
			},
			"a@mail.ru",
			0,
			[]interface{}{},
			http.StatusBadRequest,
			models.Response{Message: "getComments error: fun error"},
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
				users.GetCommentsMock.Expect(tc.postID).Return(tc.comments, tc.error)
			}
			c, m := h.GetComments(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}

func TestWriteComment(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.WriteComment, error)
		text     string
		mail     string
		postID   int
		code     int
		ans      models.Response
		error
	}{
		{"success",
			func() (models.WriteComment, error) {

				return models.WriteComment{Text: "text", PostId: 0, Mail: "a@mail.ru"}, nil
			},
			"text",
			"a@mail.ru",
			0,
			http.StatusOK,
			models.Response{Message: "success"},
			nil,
		},

		{"writeCommentError",
			func() (models.WriteComment, error) {

				return models.WriteComment{Text: "text", PostId: 1, Mail: "a@mail.ru"}, nil
			},
			"text",
			"a@mail.ru",
			1,
			http.StatusBadRequest,
			models.Response{Message: "writeComment error: "},
			errors.New(""),
		},

		{"funError",
			func() (models.WriteComment, error) {

				return models.WriteComment{Text: "text", PostId: 0, Mail: "a@mail.ru"}, errors.New("fun error")
			},
			"text",
			"a@mail.ru",
			0,
			http.StatusBadRequest,
			models.Response{Message: "writeComment error: fun error"},
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
				users.WriteCommentMock.ExpectTextParam1(&tc.text).ExpectAuthorMailParam3(tc.mail).ExpectPostIdParam4(tc.postID).Return(tc.error)
			}
			c, m := h.WriteComment(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}
