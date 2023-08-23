package models

import "time"

type URLs struct {
	Url       string
	Slug      string
	Active    bool
	CreatedAt time.Time
}
