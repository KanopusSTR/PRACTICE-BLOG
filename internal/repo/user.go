package repo

import (
	"server/internal/entities"
	"server/pkg/myErrors"
)

type User interface {
	Add(user entities.User) error
	Get(mail string) (*entities.User, error)
}

type user struct {
	users map[string]*entities.User
}

func NewUser() User {
	storage := &user{users: make(map[string]*entities.User)}
	return storage
}

func (repo *user) Add(user entities.User) error {
	if _, found := repo.users[user.Mail]; found {
		return myErrors.UserAlreadyExists
	}
	repo.users[user.Mail] = &user
	return nil
}

func (repo *user) Get(mail string) (*entities.User, error) {
	if user, found := repo.users[mail]; found {
		return user, nil
	}
	return nil, myErrors.UserNotFound
}
