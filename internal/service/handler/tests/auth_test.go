package tests

import (
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"net/http"
	"server/internal/models"
	handler "server/internal/service/handler"
	mock "server/internal/service/handler/tests/minimock"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Helper()
	testCasesAuth := []struct {
		testName string
		fun      func() (models.LoginRequest, error)
		mail     string
		code     int
		ans      models.Response
		error
	}{
		{"success",
			func() (models.LoginRequest, error) {

				return models.LoginRequest{Mail: "a@mail.ru", Password: ""}, nil
			},
			"a@mail.ru",
			http.StatusOK,
			models.Response{Message: "success", Data: ""},
			nil,
		},

		{"authError",
			func() (models.LoginRequest, error) {

				return models.LoginRequest{Mail: "a@mail.ru", Password: ""}, nil
			},
			"a@mail.ru",
			http.StatusUnprocessableEntity,
			models.Response{Message: "authorization error: "},
			errors.New(""),
		},

		{"tokenError",
			func() (models.LoginRequest, error) {

				return models.LoginRequest{Mail: "a@mail.ru", Password: ""}, nil
			},
			"a@mail.ru",
			http.StatusBadRequest,
			models.Response{Message: "server error"},
			errors.New(""),
		},

		{"funError",
			func() (models.LoginRequest, error) {

				return models.LoginRequest{Mail: "a@mail.ru", Password: ""}, errors.New("fun error")
			},
			"a@mail.ru",
			http.StatusBadRequest,
			models.Response{Message: "authorization error: fun error"},
			errors.New(""),
		},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesAuth {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			users := mock.NewUsersMock(mc)
			tkn := mock.NewTokenMock(mc)
			h := handler.New(users, tkn)
			if tc.testName != "funError" {
				if tc.testName == "authError" {
					users.AuthorizationMock.Expect(tc.mail, "").Return(tc.error)
				} else {
					users.AuthorizationMock.Expect(tc.mail, "").Return(nil)
					tkn.CreateTokenMock.Expect(tc.mail).Return("", tc.error)
				}
			}
			c, m := h.Login(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}
}

func TestRegister(t *testing.T) {
	t.Helper()
	testCasesRegister := []struct {
		testName string
		fun      func() (models.RegisterRequest, error)
		name     string
		mail     string
		code     int
		ans      models.Response
		error
	}{
		{"success",
			func() (models.RegisterRequest, error) {

				return models.RegisterRequest{Name: "a", Mail: "a@mail.ru", Password: ""}, nil
			},
			"a",
			"a@mail.ru",
			http.StatusOK,
			models.Response{Message: "success"},
			nil,
		},

		{"authError",
			func() (models.RegisterRequest, error) {
				return models.RegisterRequest{Name: "a", Mail: "a@mail.ru", Password: ""}, nil
			},
			"a",
			"a@mail.ru",
			http.StatusUnprocessableEntity,
			models.Response{Message: "register error: "},
			errors.New(""),
		},

		{"funError",
			func() (models.RegisterRequest, error) {

				return models.RegisterRequest{Name: "a", Mail: "a@mail.ru", Password: ""}, errors.New("fun error")
			},
			"a",
			"a@mail.ru",
			http.StatusBadRequest,
			models.Response{Message: "register error: fun error"},
			errors.New(""),
		},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesRegister {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			users := mock.NewUsersMock(mc)
			h := handler.New(users, nil)
			if tc.testName != "funError" {
				users.RegisterMock.Expect(tc.name, tc.mail, "").Return(tc.error)
			}
			c, m := h.Register(tc.fun)
			require.Equal(t, tc.code, c)
			require.Equal(t, tc.ans, m)
		})
	}

}
