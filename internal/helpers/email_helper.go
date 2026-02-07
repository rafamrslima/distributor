package helpers

import "net/mail"

func IsEmailValid(addr string) bool {
	_, err := mail.ParseAddress(addr)
	return err == nil
}
