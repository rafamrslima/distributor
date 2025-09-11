package domain

import "time"

type Message struct {
	ClientName        string
	ReportName        string
	Email             string
	MessageReceivedAt time.Time
}
