package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tmitchel/running-log/views"
)

// Server wraps the router, storage, and the template
// for displaying the web page.
type Server struct {
	*mux.Router
	DB    *Storage
	Index *views.View
}

// NewServer takes the Storage dependency, sets up the routes
// and returns a pointer to itself.
func NewServer(db *Storage) (*Server, error) {
	server := &Server{
		Router: mux.NewRouter().StrictSlash(true),
		DB:     db,
		Index:  views.NewView("bootstrap.html", "views/index.html"),
	}

	server.HandleFunc("/api_v1/run", server.AddRun()).Methods("POST")
	server.HandleFunc("/api_v1/reset", server.Reset())
	server.Handle("/", server.ServeIndex())

	return server, nil
}

// ServeIndex handles serving the web page.
func (s *Server) ServeIndex() http.HandlerFunc {
	type Runs struct {
		Runs []LoggedRun
	}
	return func(w http.ResponseWriter, r *http.Request) {
		runs, err := s.DB.GetRuns()
		if err != nil {
			logrus.Error(err)
			return
		}

		if err := s.Index.Render(w, Runs{runs}); err != nil {
			logrus.Error(err)
		}
	}
}

// Reset is used to delete the database.
func (s *Server) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Resetting database.")
		s.DB.Reset()
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// AddRun handles parsing information from the form and saving
// the new run.
func (s *Server) AddRun() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		distanceInMiles, err := strconv.ParseFloat(r.FormValue("distance_in_miles"), 64)
		if err != nil {
			logrus.Error(err)
			return
		}

		_, err = time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			logrus.Error(err)
			return
		}

		startTime, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", r.FormValue("start_time"))
		if err != nil {
			logrus.Error(err)
			return
		}

		temp, err := strconv.ParseInt(r.FormValue("temperature"), 10, 64)
		if err != nil {
			logrus.Error(err)
			return
		}

		hr, err := strconv.ParseInt(r.FormValue("heart_rate"), 10, 64)
		if err != nil {
			logrus.Error(err)
			return
		}

		walk := false
		if r.FormValue("walk") == "walked" {
			walk = true
		}

		run := Run{
			ID:              uuid.New(),
			DistanceInMiles: distanceInMiles,
			Duration:        r.FormValue("duration"),
			StartTime:       startTime,
			Podcast:         r.FormValue("podcast"),
			Episode:         r.FormValue("episode"),
			Quality:         r.FormValue("quality"),
			Temperature:     int(temp),
			HeartRate:       int(hr),
			Walk:            walk,
		}

		logrus.Infof("Adding run: %+v\n", run)

		err = s.DB.AddRun(&run)
		if err != nil {
			logrus.Error(err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
