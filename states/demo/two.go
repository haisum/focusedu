package demo

import (
	"errors"
	"io"
	"net/url"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
	"github.com/haisum/focusedu/set"
)

const (
	demoTwoIntroTemplate    = "demo2_intro.gohtml"
	setsSession             = "demotwo_sets"
	currentItemSession      = "demotwo_currentitem"
	currentItemStateSession = "demotwo_currentitemstate"
	questionTemplate        = "ospan_question.gohtml"
	answerTemplate          = "ospan_answer.gohtml"
	resultSession           = "demotwo_results"
	stateQuestion           = 0
	stateAnswer             = 1
	stateLetter             = 2
)

type DemoTwoState struct {
	s session.Session
}

func (ds *DemoTwoState) Render(w io.Writer, values url.Values) error {
	if ds.s.Get(setsSession) == nil {
		return ds.renderIntro(w, values)
	}
	sets := ds.s.Get(setsSession).(set.Sets)
	//show letter
	currentSetIndex := ds.s.Get(currentSetSession).(int)
	currentItemIndex := ds.s.Get(currentItemSession).(int)
	if currentSetIndex >= len(sets) || currentItemIndex >= len(sets[currentSetIndex]) {
		log.Error("This shouldn't happen!. Something's wrong in process function of demo two.")
		return errors.New("Error in process function of demo two")
	}
	if ds.s.Get("result") != nil {
		results := ds.s.Get(resultSession).(map[int]set.SetResult)
		result := results[currentSetIndex]
		return renderTemplate(w, resultTemplate, map[string]string{
			"Total":          strconv.FormatInt(int64(result.Total), 10),
			"CorrectAnswers": strconv.FormatInt(int64(result.CorrectAnswers), 10),
			"Percentage":     strconv.FormatFloat((float64(result.CorrectAnswers) / float64(result.Total) * 100.0), 'f', 1, 64),
		})
	}
	switch ds.s.Get(currentItemStateSession).(int) {
	case stateQuestion:
		return renderTemplate(w, questionTemplate, map[string]string{
			"Question": sets[currentSetIndex][currentItemIndex].Question.Question,
		})
	case stateAnswer:
		return renderTemplate(w, answerTemplate, map[string]string{
			"Option": sets[currentSetIndex][currentItemIndex].Question.Option,
		})
	}
	return nil
}

func (ds *DemoTwoState) renderIntro(w io.Writer, values url.Values) error {
	user := ds.s.Get("user").(*models.User)
	return renderTemplate(w, demoTwoIntroTemplate, map[string]string{"Name": user.Name})
}

func (ds *DemoTwoState) Process(values url.Values) error {
	if ds.s.Get(setsSession) == nil {
		log.Info("No sets defined, setting them")
		sets, err := set.GetSets(3, false, true)
		if err != nil {
			log.WithError(err).Error("Couldn't get sets")
			return err
		}
		ds.s.Set(setsSession, sets)
		ds.s.Set(currentItemSession, 0)
		ds.s.Set(currentSetSession, 0)
		ds.s.Set(currentItemStateSession, 0)
		err = ds.s.Save()
		return err
	}
	sets := ds.s.Get(setsSession).(set.Sets)
	currentSetIndex := ds.s.Get(currentSetSession).(int)
	currentItemIndex := ds.s.Get(currentItemSession).(int)

	currentItemState := ds.s.Get(currentItemStateSession).(int)
	if currentItemState == stateAnswer {
		log.Info("Recording answer.")
		var results = make(map[int]set.SetResult)
		if ds.s.Get(resultSession) == nil {
			ds.s.Set(resultSession, results)
		}
		results = ds.s.Get(resultSession).(map[int]set.SetResult)
		result := set.SetResult{}
		if v, ok := results[currentSetIndex]; ok {
			result = v
		}
		// record result
		result.Total = result.Total + 1
		if values.Get("IsTrue") == strconv.FormatInt(int64(sets[currentSetIndex][currentItemIndex].Question.IsTrue), 10) {
			result.CorrectAnswers = result.CorrectAnswers + 1
		}
		results[currentSetIndex] = result
		log.WithField("result", result).Info("New result")
		ds.s.Set(resultSession, results)
	}
	if currentItemIndex == len(sets[currentSetIndex])-1 && currentItemState == stateAnswer { //all items showed
		log.Info("All items have been shown")
		//let's show results
		if ds.s.Get("result") != nil {
			ds.s.Set("result", nil)
			log.Info("Moving on to next set")
			if currentSetIndex == len(sets)-1 { //all sets showed
				log.Info("All sets shown, moving to next state")
				user := ds.s.Get("user").(*models.User)
				user.CurrentStep = models.StepDemoThree
				err := user.Update()
				if err != nil {
					return err
				}
				ds.s.Set("user", user)
				ds.s.Set("result", nil)
				ds.s.Set(setsSession, nil)
				ds.s.Set(currentItemSession, nil)
				ds.s.Set(currentSetSession, nil)
				ds.s.Set(currentItemStateSession, nil)
				err = ds.s.Save()
				return err
			}
			ds.s.Set(currentSetSession, currentSetIndex+1)
			ds.s.Set(currentItemSession, 0)
			ds.s.Set(currentItemStateSession, 0)
			err := ds.s.Save()
			return err
		}
		ds.s.Set("result", 1)
		err := ds.s.Save()
		return err
	}
	if currentItemState == stateAnswer {
		//increase item
		ds.s.Set(currentItemStateSession, 0)
		ds.s.Set(currentItemSession, currentItemIndex+1)
		log.Infof("Incrementing item index to %d", currentItemIndex+1)
	} else {
		log.Info("Increasing item state to show answer in next run")
		ds.s.Set(currentItemStateSession, currentItemState+1)
		log.Infof("Incrementing item state to %d", currentItemState+1)
	}
	return ds.s.Save()
}
func (ds *DemoTwoState) SetSession(s session.Session) {
	ds.s = s
}
