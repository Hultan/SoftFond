package data

import "time"

func NewFunds() *Funds {
	f := new(Funds)
	return f
}

func NewRate(date time.Time, rate float64) *Rate {
	return new(Rate)
}