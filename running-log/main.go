package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	app "github.com/tmitchel/running-log"
)

func main() {
	db, err := app.OpenStorage("localhost:6379", "")
	if err != nil {
		logrus.Fatal(err)
	}

	server, err := app.NewServer(db)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Fatal(http.ListenAndServe(":8000", server))
}
