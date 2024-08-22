package tests

import (
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
	"regexp"
	"server/internal/entities"
	"server/internal/repo/psql"
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

	for _, tc := range testCasesUsers {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {

			t.Parallel()
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)
			//mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_user($1, $2, $3)")).
				WithArgs(tc.name, tc.mail, tc.password).
				WillReturnRows().WillReturnError(tc.errorMessage)

			userRepo := psql.NewUser(mock)
			err = userRepo.Add(entities.User{Name: tc.name, Mail: tc.mail, Password: tc.password})
			require.Equal(t, tc.errorMessage, err)
			if err == nil {
				rows := mock.NewRows([]string{"name", "mail", "password"}).
					AddRow(tc.name, tc.mail, tc.password)
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_user($1)")).
					WithArgs(tc.mail).
					WillReturnRows(rows)
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

	for _, tc := range testCasesUsers {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			mock, err := pgxmock.NewPool()
			require.Nil(t, err)
			userRepo := psql.NewUser(mock)

			for _, user := range tc.users {
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM create_user($1, $2, $3)")).WithArgs(user.Name, user.Mail, user.Password).WillReturnRows()
				err := userRepo.Add(*user)
				require.Nil(t, err)
			}
			u := &entities.User{}
			if tc.testName != "user not exist" {
				u = tc.users[tc.userId]
				rows := mock.NewRows([]string{"name", "mail", "password"}).
					AddRow(u.Name, u.Mail, u.Password)
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_user($1)")).
					WithArgs(tc.mail).
					WillReturnRows(rows)
			} else {
				rows := mock.NewRows([]string{"name", "mail", "password"}).
					AddRow(u.Name, u.Mail, u.Password)
				//mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM read_user($1)")).
					WithArgs(tc.mail).
					WillReturnRows(rows).WillReturnError(tc.errorMessage)
			}
			user, err := userRepo.Get(tc.mail)
			require.Equal(t, err, tc.errorMessage)
			if err == nil {
				require.Equal(t, u, user)
			}
		})
	}
}
