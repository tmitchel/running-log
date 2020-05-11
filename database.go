package app

import (
	"strings"

	"github.com/go-redis/redis/v7"
)

// Database wraps the redis client.
type Database struct {
	*redis.Client
}

// OpenStorage connects the redis client and returns a pointer to
// the Database wrapper.
func OpenStorage(addr, password string) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &Database{client}, nil
}

func (d *Database) Reset() error {
	d.Client.FlushAll()
	return nil
}

// AddRun persists a run to the database.
func (d *Database) AddRun(run *Run) error {
	return d.Client.Set(run.ID.String(), run, 0).Err()
}

// GetRuns returns all runs from the database.
func (d *Database) GetRuns() ([]LoggedRun, error) {
	var runs []LoggedRun
	keys := d.Client.Keys("*").Val()
	for _, key := range keys {
		var run Run
		err := d.Client.Get(key).Scan(&run)
		if err != nil {
			return nil, err
		}

		logged := LoggedRun{
			ID:              run.ID,
			DistanceInMiles: run.DistanceInMiles,
			Duration:        run.Duration,
			Podcast:         run.Podcast,
			Episode:         run.Episode,
			Quality:         run.Quality,
			Temperature:     run.Temperature,
			HeartRate:       run.HeartRate,
			Walk:            run.Walk,
		}

		stime := strings.Split(run.StartTime.String(), " ")
		logged.StartDate = stime[0]
		logged.StartTime = stime[1]
		runs = append(runs, logged)
	}
	return runs, nil
}
