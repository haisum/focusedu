package login

import (
	"io"
	"io/ioutil"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/session"
)

type LoginState struct {
	s session.Session
}

const loginTemplate = "login.gohtml"

func (ls LoginState) Render(w io.Writer, values url.Values) error {
	b, err := ioutil.ReadFile("templates/" + loginTemplate)
	if err != nil {
		log.WithError(err).Error("Couldn't read file templates/login.gohtml")
		return err
	}
	w.Write(b)
	return nil
}
func (ls LoginState) Process(values url.Values) error {
	return nil
}
func (ls LoginState) SetSession(s session.Session) {
	ls.s = s
}
