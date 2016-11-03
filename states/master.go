package states

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/session"
)

//var store = sessions.NewCookieStore([]byte("mKjMrUnLQ+IKMOaPFrLiLg8gsfgdQI/Zcx3MJbQZ1R0="))
var MasterHandler = func(w http.ResponseWriter, r *http.Request) {
	s, err := session.GetHTTPSession(w, r)
	if err != nil {
		log.WithError(err).Errorf("Can't get session.")
	}
	if s.Get("username") == nil {
		s.Set("username", "haisum")
		err = s.Save()
		if err != nil {
			log.WithError(err).Error("Error saving session.")
		}
	} else {
		log.Info("username is %s", s.Get("username").(string))
	}
	fmt.Fprintf(w, "Hello World")
}
