package module

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

type ModuleOneState struct {
	s session.Session
}

const (
	moduleOneIntroTemplate = "moduleone_intro.gohtml"
	moduleOneTemplate      = "moduleone.gohtml"
)

func (ms *ModuleOneState) Render(w io.Writer, values url.Values) error {
	if ms.s.Get(session.ShowModuleOneSession) == nil {
		return renderTemplate(w, moduleOneIntroTemplate, nil)
	}
	if ms.s.Get(session.TestModuleOneSession) == nil {
		user := ms.s.Get(session.UserSession).(*models.User)
		AtOnce := "2"
		if user.OSPANScore > 8 {
			AtOnce = "4"
		} else if user.OSPANScore > 18 {
			AtOnce = "8"
		}
		return renderTemplate(w, moduleOneTemplate, map[string]string{
			"Timeout": strconv.FormatInt(user.QuestionTimeout, 10),
			"AtOnce":  AtOnce,
		})
	}
	// render test

	return nil
}
func (ms *ModuleOneState) Process(values url.Values) error {
	if ms.s.Get(session.ShowModuleOneSession) == nil {
		ms.s.Set(session.ShowModuleOneSession, 1)
		return ms.s.Save()
	}
	if ms.s.Get(session.TestModuleOneSession) == nil {
		//process distraction etc
		user := ms.s.Get(session.UserSession).(*models.User)
		distraction := getIntOrZero(values.Get("distractionCount"))
		example := getIntOrZero(values.Get("exampleCount"))
		grasping := getIntOrZero(values.Get("graspingCount"))
		user.UpdateModuleOneStats(user.ID, distraction, example, grasping)
		ms.s.Set(session.TestModuleOneSession, 1)
		return ms.s.Save()
	}
	//save module test results here
	//make sessions nil
	//proceed to next state
	return nil
}
func (ms *ModuleOneState) SetSession(s session.Session) {
	ms.s = s
}

func getIntOrZero(num string) int {
	n, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		return 0
	}
	return int(n)
}

func renderTemplate(w io.Writer, tpl string, data map[string]string) error {
	b, err := ioutil.ReadFile("templates/" + tpl)
	if err != nil {
		log.WithError(err).Errorf("Couldn't read file templates/%s", tpl)
		return err
	}
	tmpl, err := template.New("test").Parse(string(b[:]))
	if err != nil {
		log.WithError(err).Errorf("Couldn't parse template.", tpl)
		return err
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.WithError(err).Error("Couldn't execute template")
		return err
	}
	return nil
}
