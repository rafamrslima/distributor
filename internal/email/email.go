package email

import (
	"net/mail"

	"github.com/rafamrslima/distributor/internal/domain"
)

func SendEmail(emailInfo domain.Message) error {
	// todo
	return nil
}

func IsValidEmail(addr string) bool {
	_, err := mail.ParseAddress(addr)
	return err == nil
}
