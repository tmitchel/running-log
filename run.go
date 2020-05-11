package app

import (
	"time"

	"github.com/google/uuid"
)

// Run contains all information about a single run.
type Run struct {
	ID              uuid.UUID `json:"id"`
	DistanceInMiles float64   `json:"distance_in_miles"`
	Duration        string    `json:"duration"`
	StartTime       time.Time `json:"start_time"`
	Podcast         string    `json:"podcast"`
	Episode         string    `json:"episode"`
	Quality         string    `json:"quality"`
	Temperature     int       `json:"temperature"`
	HeartRate       int       `json:"heart_rate"`
	Walk            bool      `json:"walk"`
}

// LoggedRun represents a run with some of the information
// combined into more useful things.
type LoggedRun struct {
	ID              uuid.UUID
	DistanceInMiles float64
	Duration        string
	StartDate       string
	StartTime       string
	Podcast         string
	Episode         string
	Quality         string
	Temperature     int
	HeartRate       int
	Walk            bool
}
