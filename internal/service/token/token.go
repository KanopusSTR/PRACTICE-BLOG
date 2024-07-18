package token

import (
	"github.com/golang-jwt/jwt/v5"
	"server/internal/config"
	"server/internal/models"
	"time"
)

func (s *service) CreateToken(mail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		models.Token{
			Mail: mail,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour).UTC()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
	return token.SignedString([]byte(config.GetJWTKey()))
}

func (s *service) ParseToken(stringToken string) (models.Token, *jwt.Token, error) {
	t := &models.Token{}
	tkn, err := jwt.ParseWithClaims(stringToken, t, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTKey()), nil
	})
	return *t, tkn, err
}
