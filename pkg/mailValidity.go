package pkg

import (
	"net/mail"
)

func IsMailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
