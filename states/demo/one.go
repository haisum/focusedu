package demo

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
)

const (
	TotalSets int = 10
)

type DemoOneState struct {
	s session.Session
}

type demoOneData struct {
	Name string
}

type DemoOneSession struct {
	CurrentSet int
	Letters    []string
}

const demoOneIntroTemplate = "demo1_intro.gohtml"

func (ds *DemoOneState) Render(w io.Writer, values url.Values) error {
	b, err := ioutil.ReadFile("templates/" + demoOneIntroTemplate)
	if err != nil {
		log.WithError(err).Errorf("Couldn't read file templates/%s", demoOneIntroTemplate)
		return err
	}
	tmpl, err := template.New("test").Parse(string(b[:]))
	if err != nil {
		log.WithError(err).Errorf("Couldn't parse template.", demoOneIntroTemplate)
		return err
	}
	var tplData demoOneData
	user := ds.s.Get("user").(*models.User)
	tplData.Name = user.Name
	err = tmpl.Execute(w, tplData)
	if err != nil {
		log.WithError(err).Error("Couldn't execute template")
		return err
	}
	return nil
}
func (ds *DemoOneState) Process(values url.Values) error {
	return nil
}
func (ds *DemoOneState) SetSession(s session.Session) {
	ds.s = s
}
