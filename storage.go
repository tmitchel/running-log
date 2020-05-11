package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Storage wraps the contents of the file used for saving
// runs.
type Storage struct {
	// Name of storage file
	Name string

	// List of saved runs
	Entries []Run
}

// Open the file (create it if it doesn't exist) and read
// the contents into the Storage struct.
func Open(file string) (*Storage, error) {
	s := &Storage{Name: file}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		s.Entries = make([]Run, 0)
		return s, nil
	}

	// read the file
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// parse the content
	err = json.Unmarshal(content, &s.Entries)
	if err != nil {
		return nil, err
	}

	return s, err
}

// AddRun persists a run to the database.
func (s *Storage) AddRun(run *Run) error {
	s.Entries = append(s.Entries, *run)

	// write file in the background
	go func() {
		file, err := os.Create(s.Name)
		if err != nil {
			logrus.Fatal(err)
		}
		defer file.Close()

		content, err := json.MarshalIndent(s.Entries, "", "  ")
		if err != nil {
			logrus.Fatal(err)
		}

		_, err = file.Write(content)
	}()
	return nil
}

// GetRuns reformats the Runs read from the file to add
// more information.
func (s *Storage) GetRuns() ([]LoggedRun, error) {
	runs := make([]LoggedRun, len(s.Entries))
	for i, run := range s.Entries {
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

		runs[i] = logged
	}
	return runs, nil
}

// Reset deletes everything.
func (s *Storage) Reset() error {
	s.Entries = make([]Run, 0)
	return os.Remove(s.Name)
}
