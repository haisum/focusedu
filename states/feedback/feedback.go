package feedback

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/session"
)

type FeedbackState struct {
	s session.Session
}

const feedbackTpl = "feedback.gohtml"

func (is *FeedbackState) Render(w io.Writer, values url.Values) error {
	b, err := ioutil.ReadFile("templates/" + feedbackTpl)
	if err != nil {
		log.WithError(err).Errorf("Couldn't read file templates/%s", feedbackTpl)
		return err
	}
	tmpl, err := template.New("test").Parse(string(b[:]))
	if err != nil {
		log.WithError(err).Errorf("Couldn't parse template.", feedbackTpl)
		return err
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.WithError(err).Error("Couldn't execute template")
		return err
	}
	return nil
}
func (is *FeedbackState) Process(values url.Values) error {
	return nil
}
func (is *FeedbackState) SetSession(s session.Session) {
	is.s = s
}
