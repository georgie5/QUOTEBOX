package models

import (
	"time"
)

// A struct to hold a quote
type Quote struct {
	Quote_id    int
	Created_at  time.Time
	Author_name string
	Category    string
	Body        string
}
