package states

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/session"
)

var MasterHandler = func(w http.ResponseWriter, r *http.Request) {
	s, err := session.GetHTTPSession(w, r)
	if err != nil {
		log.WithError(err).Errorf("Can't get session.")
	}
	state, err := getState(s)
	if err != nil {
		log.WithError(err).Error("Error in getting state")
	}
	err = state.Render()
	if err != nil {
		log.WithError(err).Error("Error in rendering state")
	}
}
