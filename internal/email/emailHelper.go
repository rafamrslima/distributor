package email

import "net/mail"

func IsValid(addr string) bool {
	_, err := mail.ParseAddress(addr)
	return err == nil
}
