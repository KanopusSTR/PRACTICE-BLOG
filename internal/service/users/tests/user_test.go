package tests

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"server/internal/entities"
	"server/internal/service/users"
	mock "server/internal/service/users/tests/minimock"
	"server/pkg"
	"server/pkg/myErrors"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Helper()
	testCasesUsers := []struct {
		testName     string
		name         string
		mail         string
		password     string
		errorMessage error
	}{
		{"simple", "artur", "hohoho@hoho.com", "hoho", nil},
		{"empty name", "", "hohoho@hoho.com", "hoho", myErrors.EmptyRegisterData},
		{"empty password", "artur", "hohoho@hoho.com", "", myErrors.EmptyRegisterData},
		{"already exist", "artur", "alreadyexist@hoho.com", "hoho", myErrors.UserAlreadyExists},
		{"incorrect mail", "artur", "incorrect", "hoho", myErrors.InvalidMail},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesUsers {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			userRepo := mock.NewUserMock(mc)

			userS := users.New(userRepo, nil, nil)

			if tc.testName != "incorrect mail" && tc.testName != "empty name" && tc.testName != "empty password" {
				userRepo.AddUserMock.Expect(entities.User{
					Name:     tc.name,
					Mail:     tc.mail,
					Password: pkg.GetStringHash(tc.password),
				}).Return(tc.errorMessage)
			}

			err := userS.Register(tc.name, tc.mail, tc.password)
			require.Equal(t, tc.errorMessage, err)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Helper()
	testCasesUsers := []struct {
		testName     string
		mail         string
		password     string
		errorMessage error
	}{
		{"simple", "hohoho@hoho.com", "hoho", nil},
		{"incorrect password", "hohoho@hoho.com", "incorrect", myErrors.IncorrectPassword},
		{"user doesn't exist", "empty@hoho.com", "hoho", myErrors.UserNotFound},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesUsers {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			userRepo := mock.NewUserMock(mc)

			userS := users.New(userRepo, nil, nil)

			if tc.testName == "user doesn't exist" {
				userRepo.GetUserMock.Expect(tc.mail).Return(nil, myErrors.UserNotFound)
			} else {
				userRepo.GetUserMock.Expect(tc.mail).Return(&entities.User{Password: pkg.GetStringHash(tc.password)}, tc.errorMessage)
			}

			err := userS.Authorization(tc.mail, tc.password)
			require.Equal(t, tc.errorMessage, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Helper()
	testCasesUsers := []struct {
		testName     string
		mail         string
		errorMessage error
	}{
		{"simple", "hohoho@hoho.com", nil},
		{"user doesn't exist", "empty@hoho.com", myErrors.UserNotFound},
	}

	mc := minimock.NewController(t)

	for _, tc := range testCasesUsers {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			userRepo := mock.NewUserMock(mc)

			userS := users.New(userRepo, nil, nil)

			if tc.testName == "user doesn't exist" {
				userRepo.GetUserMock.Expect(tc.mail).Return(nil, myErrors.UserNotFound)
			} else {
				userRepo.GetUserMock.Expect(tc.mail).Return(&entities.User{Mail: tc.mail}, tc.errorMessage)
			}

			_, err := userS.GetProfile(tc.mail)
			require.Equal(t, tc.errorMessage, err)
		})
	}
}
