package info

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
)

type InfoState struct {
	s session.Session
}

type templateData struct {
	Errors map[string]string
}

const infoTemplate = "info.gohtml"

func (is *InfoState) Render(w io.Writer, values url.Values) error {
	b, err := ioutil.ReadFile("templates/" + infoTemplate)
	if err != nil {
		log.WithError(err).Errorf("Couldn't read file templates/%s", infoTemplate)
		return err
	}
	tmpl, err := template.New("test").Parse(string(b[:]))
	if err != nil {
		log.WithError(err).Errorf("Couldn't parse template.", infoTemplate)
		return err
	}
	var tplData templateData
	if v := is.s.Get("errors"); v != nil {
		tplData.Errors = v.(map[string]string)
		is.s.Set("errors", nil)
		is.s.Save()
	}
	err = tmpl.Execute(w, tplData)
	if err != nil {
		log.WithError(err).Error("Couldn't execute template")
		return err
	}
	return nil
}
func (is *InfoState) Process(values url.Values) error {
	user := is.s.Get("user").(*models.User)
	user.Age = getIntOrZero(values.Get("Age"))
	user.Gender = models.Gender(getIntOrZero(values.Get("Gender")))
	user.MidtermScore = getIntOrZero(values.Get("MidtermScore"))
	user.Name = values.Get("Name")
	errors := user.Validate()
	if len(errors) > 0 {
		log.WithField("errors", errors).Error("Errors occurred")
		is.s.Set("errors", errors)
		is.s.Save()
		return nil
	}
	user.CurrentStep = models.StepDemoOne
	err := user.Update()
	if err != nil {
		log.WithError(err).Error("Error saving user.")
		return err
	}
	is.s.Set("user", user)
	err = is.s.Save()
	return err
}
func (is *InfoState) SetSession(s session.Session) {
	is.s = s
}

func getIntOrZero(num string) int {
	n, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		return 0
	}
	return int(n)
}
