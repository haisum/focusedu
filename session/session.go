package session

import (
	"encoding/gob"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/set"
)

type Session interface {
	Get(key string) interface{}
	Set(key string, val interface{})
	Save() error
	Destroy() error
}

const (
	CurrentLetterSession    = "currentletter"
	CurrentSetSession       = "currentset"
	SetsSession             = "sets"
	CurrentItemSession      = "currentitem"
	CurrentItemStateSession = "currentitemstate"
	ResultsSession          = "results"
	ShowResultSession       = "result"
	UserSession             = "user"
	ShowGridSession         = "showgrid"
	TotalTimeSession        = "totaltime"
	StartTimeSession        = "starttime"
	ShowModuleOneSession    = "showmoduleone"
	TestModuleOneSession    = "testmoduleone"
	ShowModuleTwoSession    = "showmoduletwo"
	TestModuleTwoSession    = "testmoduletwo"
)

type httpsession struct {
	r       *http.Request
	w       http.ResponseWriter
	store   sessions.Store
	session *sessions.Session
}

const (
	authKey     = "mKjMrUnLQ+IKMOaPFrLiLg8gsfgdQI/Zcx3MJbQZ1R0="
	encKey      = "y6cPoxR3e2UceU0aEFQd8w=="
	sessionName = "focusedu-session"
)

var s *httpsession

func GetHTTPSession(w http.ResponseWriter, r *http.Request) (Session, error) {
	if s == nil {
		registerGobTypes()
		s = &httpsession{}
		s.store = sessions.NewFilesystemStore("", []byte(authKey), []byte(encKey))
	}
	s.r = r
	s.w = w
	var err error
	s.session, err = s.store.Get(r, sessionName)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s *httpsession) Get(key string) interface{} {
	if val, ok := s.session.Values[key]; ok {
		return val
	}
	return nil
}

func (s *httpsession) Set(key string, val interface{}) {
	s.session.Values[key] = val
}

func (s *httpsession) Save() error {
	err := s.session.Save(s.r, s.w)
	if err != nil {
		log.WithError(err).Error("Error in saving session.")
	}
	return err
}

func (s *httpsession) Destroy() error {
	for k, _ := range s.session.Values {
		s.session.Values[k] = nil
	}
	err := s.session.Save(s.r, s.w)
	if err != nil {
		log.WithError(err).Error("Error in saving session.")
	}
	return err
}
func registerGobTypes() {
	gob.Register(&models.User{})
	gob.Register(map[string]string{})
	gob.Register(set.Sets{})
	gob.Register(map[int]set.SetResult{})
	gob.Register(map[int]int64{})
}
