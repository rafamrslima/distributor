package domain

import "time"

type Message struct {
	Name              string
	Email             string
	EmailCc           string
	Content           []byte
	MessageReceivedAt time.Time
}
