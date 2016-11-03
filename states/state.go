package states

import (
	"io"
	"net/url"

	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
	"github.com/haisum/focusedu/states/login"
)

type State interface {
	SetSession(s session.Session)
	Render(w io.Writer, values url.Values) error
	Process(values url.Values) error
}

func getState(s session.Session) (State, error) {
	var user *models.User
	if val := s.Get("user"); val != nil {
		user = val.(*models.User)
	}
	var state State
	if user == nil {
		state = login.LoginState{}
	} else {
		switch user.CurrentStep {

		}
	}
	return state, nil
}
