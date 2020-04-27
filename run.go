package app

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Run contains all information about a single run.
type Run struct {
	ID              uuid.UUID
	DistanceInMiles float64
	Duration        time.Duration
	StartTime       time.Time
	Podcast         string
	Episode         string
	Quality         string
}

func (r *Run) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Run) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

type LoggedRun struct {
	ID              uuid.UUID
	DistanceInMiles float64
	Duration        time.Duration
	StartDate       string
	StartTime       string
	Podcast         string
	Episode         string
	Quality         string
}
