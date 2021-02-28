package data

import "time"


type Rates struct {
	List []Rate `json:"Rates"`
}

type Rate struct {
	Date time.Time
	Rate float64
}

// RateNew : Create a new rate struct
func RateNew(date time.Time, rate float64) *Rate {
	return new(Rate)
}