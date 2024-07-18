package users

import (
	"bytes"
	"server/internal/entities"
	"server/pkg"
	"server/pkg/myErrors"
)

func (s *service) Authorization(mail, password string) error {
	user, err := s.users.Get(mail)
	if err != nil {
		return err
	}

	sum := pkg.GetStringHash(password)

	if !bytes.Equal(user.Password, sum) {
		return myErrors.IncorrectPassword
	}
	return nil
}

func (s *service) Register(name, mail, password string) error {
	if !pkg.IsMailValid(mail) {
		return myErrors.InvalidMail
	}
	if name == "" || mail == "" || password == "" {
		return myErrors.EmptyRegisterData
	}
	sum := pkg.GetStringHash(password)
	user := entities.User{Name: name, Mail: mail, Password: sum}
	return s.users.Add(user)
}

func (s *service) GetProfile(mail string) (*entities.User, error) {
	user, err := s.users.Get(mail)
	if err == nil {
		return &entities.User{Name: user.Name, Mail: mail}, nil
	}
	return nil, err
}
