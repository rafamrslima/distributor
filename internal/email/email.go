package email

import "net/mail"

func IsValidEmail(addr string) bool {
	_, err := mail.ParseAddress(addr)
	return err == nil
}
