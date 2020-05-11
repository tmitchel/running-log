package views

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type View struct {
	Template *template.Template
	Layout   string
}

func NewView(layout string, files ...string) *View {
	htmls, err := filepath.Glob("views/layouts/*.html")
	if err != nil {
		logrus.Fatal(err)
	}

	files = append(htmls, files...)
	temp, err := template.ParseFiles(files...)
	if err != nil {
		logrus.Fatal(err)
	}

	return &View{
		Template: temp,
		Layout:   layout,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}
