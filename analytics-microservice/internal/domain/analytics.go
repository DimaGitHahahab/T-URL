package domain

import "time"

type Stats struct {
	ShortURL       string
	NumberOfUsage  int
	LastAccessedAt time.Time
}
