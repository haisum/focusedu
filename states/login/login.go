package login

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
)

type LoginState struct {
	s session.Session
}

const loginTemplate = "login.gohtml"

func (ls *LoginState) Render(w io.Writer, values url.Values) error {
	b, err := ioutil.ReadFile("templates/" + loginTemplate)
	if err != nil {
		log.WithError(err).Error("Couldn't read file templates/login.gohtml")
		return err
	}
	w.Write(b)
	return nil
}
func (ls *LoginState) Process(values url.Values) error {
	if _, ok := values["RollNo"]; !ok || len(values["RollNo"]) != 1 {
		return fmt.Errorf("RollNo is required.")
	}
	rollno := values["RollNo"][0]
	user, err := models.GetUser(rollno)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Infof("Creating new user with rollno %s", rollno)
			err = models.MakeUser(rollno)
			if err != nil {
				log.WithError(err).Error("Couldn't create user.")
				return err
			}
			user, err = models.GetUser(rollno)
			if err != nil {
				return err
			}
		} else {
			log.WithError(err).Error("Couldn't get user from db.")
			return fmt.Errorf("Couldn't log in. Check with your system admin.")
		}
	}
	ls.s.Set("user", user)
	err = ls.s.Save()
	if err != nil {
		log.WithError(err).Error("Error saving session.")
		return err
	}
	log.WithField("user", user).Info("Just logged in")
	return nil
}
func (ls *LoginState) SetSession(s session.Session) {
	ls.s = s
}
