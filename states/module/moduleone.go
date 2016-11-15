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
	moduleOneTestTemplate  = "moduleone_test.gohtml"
)

var (
	moduleOneCorrectAnswers = []int{0, 2, 0, 3, 3, 2, 1, 3, 0, 0}
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
			"Timeout":  strconv.FormatInt(user.QuestionTimeout, 10),
			"AtOnce":   AtOnce,
			"UserType": strconv.FormatInt(int64(user.Type), 10),
		})
	}
	// render test
	return renderTemplate(w, moduleOneTestTemplate, nil)
}
func (ms *ModuleOneState) Process(values url.Values) error {
	if ms.s.Get(session.ShowModuleOneSession) == nil {
		ms.s.Set(session.ShowModuleOneSession, 1)
		return ms.s.Save()
	}
	if ms.s.Get(session.TestModuleOneSession) == nil {
		log.Info("Processing module stats")
		//process distraction etc
		user := ms.s.Get(session.UserSession).(*models.User)
		distraction := getIntOrZero(values.Get("distractionCount"))
		example := getIntOrZero(values.Get("exampleCount"))
		grasping := getIntOrZero(values.Get("graspingCount"))
		log.WithFields(log.Fields{
			"distraction": distraction,
			"example":     example,
			"grasping":    grasping,
		}).Info("New stats")
		err := user.UpdateModuleOneStats(distraction, example, grasping)
		if err != nil {
			log.WithError(err).Error("Error in updating module stats")
			return err
		}
		ms.s.Set(session.TestModuleOneSession, 1)
		return ms.s.Save()
	}
	log.Info("Processing test scores")
	correctAnswers := 0
	for k, v := range moduleOneCorrectAnswers {
		given := values.Get("answer-" + strconv.FormatInt(int64(k), 10))
		if given == strconv.FormatInt(int64(v), 10) {
			correctAnswers += 1
		}
	}
	log.WithField("Score", correctAnswers).Info("Calculated score for module one.")
	user := ms.s.Get(session.UserSession).(*models.User)
	err := user.UpdateModuleOneScore(correctAnswers)
	if err != nil {
		log.WithError(err).Error("Couldn't update module one score.")
		return err
	}
	ms.s.Set(session.ShowModuleOneSession, nil)
	ms.s.Set(session.TestModuleOneSession, nil)
	user.CurrentStep = models.StepModuleTwo
	err = user.Update()
	if err != nil {
		log.WithError(err).Error("Couldn't update user in module one.")
		return err
	}
	ms.s.Set(session.UserSession, user)
	return ms.s.Save()
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
