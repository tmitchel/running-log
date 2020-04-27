package app

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*mux.Router
	DB    *Database
	Index *template.Template
}

func NewServer(db *Database) (*Server, error) {
	server := &Server{
		Router: mux.NewRouter().StrictSlash(true),
		DB:     db,
		Index:  template.Must(template.ParseFiles("views/index.html")),
	}

	server.HandleFunc("/api_v1/run", server.AddRun()).Methods("POST")
	server.HandleFunc("/api_v1/reset", server.Reset())
	server.Handle("/", server.ServeIndex())

	return server, nil
}

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
		logrus.Infof("%+v\n", Runs{runs})
		if err := s.Index.Execute(w, Runs{runs}); err != nil {
			logrus.Error(err)
		}
	}
}

func (s *Server) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Resetting database.")
		s.DB.Reset()
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *Server) AddRun() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		distanceInMiles, err := strconv.ParseFloat(r.FormValue("distance_in_miles"), 64)
		if err != nil {
			logrus.Error(err)
			return
		}

		duration, err := time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			logrus.Error(err)
			return
		}

		startTime, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", r.FormValue("start_time"))
		if err != nil {
			logrus.Error(err)
			return
		}

		run := Run{
			ID:              uuid.New(),
			DistanceInMiles: distanceInMiles,
			Duration:        duration,
			StartTime:       startTime,
			Podcast:         r.FormValue("podcast"),
			Episode:         r.FormValue("episode"),
			Quality:         r.FormValue("quality"),
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
