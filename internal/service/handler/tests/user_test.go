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

func TestGetProfile(t *testing.T) {
	t.Helper()
	testCases := []struct {
		testName string
		fun      func() (models.GetUser, error)
		name     string
		mail     string
		code     int
		ans      models.Response
		error
	}{
		{"success",
			func() (models.GetUser, error) {

				return models.GetUser{Mail: "a@mail.ru"}, nil
			},
			"a",
			"a@mail.ru",
			http.StatusOK,
			models.Response{Message: "success", Data: models.ProfileResponse{Name: "a", Mail: "a@mail.ru"}},
			nil,
		},

		{"getProfileError",
			func() (models.GetUser, error) {
				return models.GetUser{Mail: "a@mail.ru"}, nil
			},
			"a",
			"a@mail.ru",
			http.StatusNotFound,
			models.Response{Message: "getProfile error: "},
			errors.New(""),
		},

		{"funError",
			func() (models.GetUser, error) {
				return models.GetUser{Mail: "a@mail.ru"}, errors.New("fun error")
			},
			"a",
			"a@mail.ru",
			http.StatusInternalServerError,
			models.Response{Message: "getProfile error: fun error"},
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
				users.GetProfileMock.Expect(tc.mail).Return(&entities.User{Name: tc.name, Mail: tc.mail}, tc.error)
			}
			code, mod := h.GetUser(tc.fun)
			require.Equal(t, tc.code, code)
			require.Equal(t, tc.ans, mod)
		})
	}
}
