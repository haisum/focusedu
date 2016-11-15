package module

import (
	"io"
	"net/url"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
)

type ModuleTwoState struct {
	s session.Session
}

const (
	moduleTwoIntroTemplate = "moduletwo_intro.gohtml"
	moduleTwoTemplate      = "moduletwo.gohtml"
	moduleTwoTestTemplate  = "moduletwo_test.gohtml"
)

var (
	moduleTwoCorrectAnswers = []int{3, 2, 3, 3, 3, 3, 3, 0, 1, 0}
)

func (ms *ModuleTwoState) Render(w io.Writer, values url.Values) error {
	if ms.s.Get(session.ShowModuleTwoSession) == nil {
		return renderTemplate(w, moduleTwoIntroTemplate, nil)
	}
	if ms.s.Get(session.TestModuleTwoSession) == nil {
		user := ms.s.Get(session.UserSession).(*models.User)
		AtOnce := "2"
		if user.OSPANScore > 8 {
			AtOnce = "4"
		} else if user.OSPANScore > 18 {
			AtOnce = "8"
		}
		return renderTemplate(w, moduleTwoTemplate, map[string]string{
			"Timeout":  strconv.FormatInt(user.QuestionTimeout, 10),
			"AtOnce":   AtOnce,
			"UserType": strconv.FormatInt(int64(user.Type), 10),
		})
	}
	// render test
	return renderTemplate(w, moduleTwoTestTemplate, nil)
}
func (ms *ModuleTwoState) Process(values url.Values) error {
	if ms.s.Get(session.ShowModuleTwoSession) == nil {
		ms.s.Set(session.ShowModuleTwoSession, 1)
		return ms.s.Save()
	}
	if ms.s.Get(session.TestModuleTwoSession) == nil {
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
		err := user.UpdateModuleTwoStats(distraction, example, grasping)
		if err != nil {
			log.WithError(err).Error("Error in updating module stats")
			return err
		}
		ms.s.Set(session.TestModuleTwoSession, 1)
		return ms.s.Save()
	}
	log.Info("Processing test scores")
	correctAnswers := 0
	for k, v := range moduleTwoCorrectAnswers {
		given := values.Get("answer-" + strconv.FormatInt(int64(k), 10))
		if given == strconv.FormatInt(int64(v), 10) {
			correctAnswers += 1
		}
	}
	log.WithField("Score", correctAnswers).Info("Calculated score for module two.")
	user := ms.s.Get(session.UserSession).(*models.User)
	err := user.UpdateModuleTwoScore(correctAnswers)
	if err != nil {
		log.WithError(err).Error("Couldn't update module two score.")
		return err
	}
	ms.s.Set(session.ShowModuleTwoSession, nil)
	ms.s.Set(session.TestModuleTwoSession, nil)
	user.CurrentStep = models.StepFeedback
	err = user.Update()
	if err != nil {
		log.WithError(err).Error("Couldn't update user in module two.")
		return err
	}
	ms.s.Set(session.UserSession, user)
	return ms.s.Save()
}
func (ms *ModuleTwoState) SetSession(s session.Session) {
	ms.s = s
}
