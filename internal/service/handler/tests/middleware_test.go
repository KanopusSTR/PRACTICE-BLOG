package tests

import (
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"server/internal/models"
	handler "server/internal/service/handler"
	mock "server/internal/service/handler/tests/minimock"
	"testing"
)

func TestLoginMiddleware(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.LoginMiddleware, error)
		mail     string
		code     int
		ans      models.ResultResponseBody
		error
	}{
		{"success",
			func() (models.LoginMiddleware, error) {

				return models.LoginMiddleware{Token: "token"}, nil
			},
			"a@mail.ru",
			http.StatusOK,
			models.ResultResponseBody{Message: "success"},
			nil,
		},

		{"wrongTokenError",
			func() (models.LoginMiddleware, error) {
				return models.LoginMiddleware{Token: "token"}, nil
			},
			"",
			http.StatusForbidden,
			models.ResultResponseBody{Message: "accessibility error: you do not have permission"},
			nil,
		},

		{"parseError",
			func() (models.LoginMiddleware, error) {

				return models.LoginMiddleware{Token: "token"}, nil
			},
			"",
			http.StatusUnauthorized,
			models.ResultResponseBody{Message: "accessibility error: you do not have permission"},
			errors.New(""),
		},

		{"funError",
			func() (models.LoginMiddleware, error) {

				return models.LoginMiddleware{Token: "token"}, errors.New("fun error")
			},
			"",
			http.StatusBadRequest,
			models.ResultResponseBody{Message: "middleware error: fun error"},
			errors.New(""),
		},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			users := mock.NewUsersMock(mc)
			token := mock.NewTokenMock(mc)
			h := handler.New(users, token)

			if tc.testName != "funError" {
				tkn := jwt.New(jwt.SigningMethodHS256)
				if tc.testName == "wrongTokenError" {
					tkn.Valid = false
				} else {
					tkn.Valid = true
				}
				token.ParseTokenMock.Expect("token").Return(models.Token{Mail: "a@mail.ru"}, tkn, tc.error)
			}
			code, mod, mail := h.LoginMiddleware(tc.fun)
			require.Equal(t, tc.code, code)
			require.Equal(t, tc.ans, mod)
			require.Equal(t, tc.mail, mail)
		})
	}
}
