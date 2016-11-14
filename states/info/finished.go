package info

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/session"
)

type FinishedState struct {
	s session.Session
}

const finishedTpl = "finished.gohtml"

func (is *FinishedState) Render(w io.Writer, values url.Values) error {
	b, err := ioutil.ReadFile("templates/" + finishedTpl)
	if err != nil {
		log.WithError(err).Errorf("Couldn't read file templates/%s", finishedTpl)
		return err
	}
	tmpl, err := template.New("test").Parse(string(b[:]))
	if err != nil {
		log.WithError(err).Errorf("Couldn't parse template.", finishedTpl)
		return err
	}
	var tplData templateData
	if v := is.s.Get("errors"); v != nil {
		tplData.Errors = v.(map[string]string)
		is.s.Set("errors", nil)
		is.s.Save()
	}
	err = tmpl.Execute(w, tplData)
	if err != nil {
		log.WithError(err).Error("Couldn't execute template")
		return err
	}
	return nil
}
func (is *FinishedState) Process(values url.Values) error {
	return nil
}
func (is *FinishedState) SetSession(s session.Session) {
	is.s = s
}
