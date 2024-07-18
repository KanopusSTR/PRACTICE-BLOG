package token

import (
	"github.com/golang-jwt/jwt/v5"
	"server/internal/models"
)

type Service interface {
	CreateToken(mail string) (string, error)

	ParseToken(stringToken string) (models.Token, *jwt.Token, error)
}

type service struct{}

func New() Service {
	return &service{}
}
