package model

import "time"

type ToperFlags struct {
	Period  time.Duration
	DueDate string
	Acronym string
}
