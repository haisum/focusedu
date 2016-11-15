package states

import (
	"io"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
	"github.com/haisum/focusedu/states/info"
	"github.com/haisum/focusedu/states/login"
	"github.com/haisum/focusedu/states/module"
	"github.com/haisum/focusedu/states/ospan"
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
	log.Info("Getting state info.")
	if user == nil {
		log.Info("No state found, setting to login state")
		state = &login.LoginState{}
	} else {
		log.Infof("Current step set to %d", user.CurrentStep)
		switch user.CurrentStep {
		case models.StepInfo:
			state = &info.InfoState{}
		case models.StepDemoOne:
			state = &ospan.DemoOneState{}
		case models.StepDemoTwo:
			state = &ospan.DemoTwoState{}
		case models.StepDemoThree:
			state = &ospan.DemoThreeState{}
		case models.StepOSPAN:
			rs := &ospan.RealOSPANState{}
			rs.SetUser(user)
			state = rs
		case models.StepModuleOne:
			state = &module.ModuleOneState{}
		default:
			state = &info.FinishedState{}
		}
	}
	if state != nil {
		log.Info("Setting session.")
		state.SetSession(s)
	}
	return state, nil
}
