package session

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/haisum/focusedu/db/models"
)

type Session interface {
	Get(key string) interface{}
	Set(key string, val interface{})
	Save() error
}

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
	return s.session.Save(s.r, s.w)
}
func registerGobTypes() {
	gob.Register(&models.User{})
}
