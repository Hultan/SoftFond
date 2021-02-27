package data

import "time"


type Rates struct {
	List []Rate `json:"Rates"`
}

type Rate struct {
	Date time.Time
	Rate float64
}