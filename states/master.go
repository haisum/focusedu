package states

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/session"
)

var MasterHandler = func(w http.ResponseWriter, r *http.Request) {
	s, err := session.GetHTTPSession(w, r)
	if err != nil {
		log.WithError(err).Errorf("Can't get session.")
	}
	state, err := getState(s)
	if err != nil || state == nil {
		log.WithError(err).Error("Error in getting state.")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	if r.Method == http.MethodGet {
		err = state.Render(w, url.Values(r.Form))
		if err != nil {
			log.WithError(err).Error("Error in rendering state.")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		err = state.Process(url.Values(r.PostForm))
		if err != nil {
			log.WithError(err).Error("Error in processing state.")
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
