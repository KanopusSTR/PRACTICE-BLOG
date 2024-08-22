package repo

import "server/internal/entities"

type User interface {
	Add(user entities.User) error
	Get(mail string) (*entities.User, error)
}
