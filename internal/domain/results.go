package domain

import "time"

type Results struct {
	ClientEmail string
	ReportName  string
	Gains       float64
	Losses      float64
	InfoDate    time.Time
}
