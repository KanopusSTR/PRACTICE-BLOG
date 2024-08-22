package tests

import (
	"github.com/stretchr/testify/require"
	"server/internal/entities"
	"server/internal/repo/internaldb"
	"server/pkg/myErrors"
	"testing"
)

func TestAddUser(t *testing.T) {
	t.Helper()
	testCasesUsers := []struct {
		testName     string
		name         string
		mail         string
		password     []byte
		errorMessage error
	}{
		{"simple", "hoho", "hohoho@hoho.com", []byte("hohoho"), nil},
		{"chinese", "尽快。", "chinese@hoho.com", []byte("hohoho"), nil},
		{"arabian", "على الطاولة", "arabian@hoho.com", []byte("hohoho"), nil},
		{"double", "hoho", "hohoho@hoho.com", []byte("hohoho"), myErrors.UserAlreadyExists},
	}

	userRepo := internaldb.NewUser()

	for _, tc := range testCasesUsers {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			err := userRepo.Add(entities.User{Name: tc.name, Mail: tc.mail, Password: tc.password})
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				user, err := userRepo.Get(tc.mail)
				require.Nil(t, err)
				require.Equal(t, &entities.User{
					Name:     tc.name,
					Mail:     tc.mail,
					Password: tc.password,
				}, user)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Helper()
	testCasesUsers := []struct {
		testName     string
		mail         string
		userId       int
		users        []*entities.User
		errorMessage error
	}{
		{"simple", "hohoho@hoho.com", 0, []*entities.User{
			{"hoho", "hohoho@hoho.com", []byte("hohoho")},
		}, nil},
		{"several users", "hohoho2@hoho.com", 1, []*entities.User{
			{"hoho", "hohoho@hoho.com", []byte("hohoho")},
			{"hoho2", "hohoho2@hoho.com", []byte("hohoho2")},
		}, nil},
		{"user not exist", "hohoho3@hoho.com", -1, []*entities.User{
			{"hoho", "hohoho@hoho.com", []byte("hohoho")},
			{"hoho2", "hohoho2@hoho.com", []byte("hohoho2")},
		}, myErrors.UserNotFound},
	}

	userRepo := internaldb.NewUser()

	for _, tc := range testCasesUsers {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			for _, user := range tc.users {
				_ = userRepo.Add(*user)
			}
			user, err := userRepo.Get(tc.mail)
			require.Equal(t, err, tc.errorMessage)
			if err == nil {
				require.Equal(t, tc.users[tc.userId], user)
			}
		})
	}
}
