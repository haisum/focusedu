package states

import (
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
)

type State interface {
	Render() error
}

func getState(s session.Session) (State, error) {
	var user *models.User
	if val := s.Get("user"); val != nil {
		user = val.(*models.User)
	}
	if user == nil {
		//render login
	} else {
		switch user.CurrentStep {

		}
	}
	return nil, nil
}
