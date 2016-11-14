package ospan

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
	"github.com/haisum/focusedu/set"
)

const (
	demoTwoIntroTemplate = "demo2_intro.gohtml"
	totalQuestions       = 8
	questionTemplate     = "ospan_question.gohtml"
	answerTemplate       = "ospan_answer.gohtml"
	stateQuestion        = 0
	stateAnswer          = 1
	stateLetter          = 2
)

type DemoTwoState struct {
	s session.Session
}

func (ds *DemoTwoState) Render(w io.Writer, values url.Values) error {
	if ds.s.Get(session.SetsSession) == nil {
		return ds.renderIntro(w, values)
	}
	sets := ds.s.Get(session.SetsSession).(set.Sets)
	//show letter
	currentSetIndex := ds.s.Get(session.CurrentSetSession).(int)
	currentItemIndex := ds.s.Get(session.CurrentItemSession).(int)
	if currentSetIndex >= len(sets) || currentItemIndex >= len(sets[currentSetIndex]) {
		log.Error("This shouldn't happen!. Something's wrong in process function of demo two.")
		return errors.New("Error in process function of demo two")
	}
	if ds.s.Get("result") != nil {
		results := ds.s.Get(session.ResultsSession).(map[int]set.SetResult)
		result := results[currentSetIndex]
		return renderTemplate(w, resultTemplate, map[string]string{
			"Total":          strconv.FormatInt(int64(result.Total), 10),
			"CorrectAnswers": strconv.FormatInt(int64(result.CorrectAnswers), 10),
			"Percentage":     strconv.FormatFloat((float64(result.CorrectAnswers) / float64(result.Total) * 100.0), 'f', 1, 64),
			"Time":           fmt.Sprintf("%+v", ds.s.Get(session.TotalTimeSession).(map[int]int64)),
		})
	}
	switch ds.s.Get(session.CurrentItemStateSession).(int) {
	case stateQuestion:
		ds.s.Set(session.StartTimeSession, time.Now().UnixNano())
		ds.s.Save()
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
	if ds.s.Get(session.SetsSession) == nil {
		log.Info("No sets defined, setting them")
		sets, err := set.GetQuestionsSet(ds.s.Get(session.UserSession).(*models.User).ID, totalQuestions)
		if err != nil {
			log.WithError(err).Error("Couldn't get sets")
			return err
		}
		ds.s.Set(session.SetsSession, sets)
		ds.s.Set(session.CurrentItemSession, 0)
		ds.s.Set(session.CurrentSetSession, 0)
		ds.s.Set(session.CurrentItemStateSession, 0)
		ds.s.Set(session.TotalTimeSession, map[int]int64{})
		err = ds.s.Save()
		return err
	}
	sets := ds.s.Get(session.SetsSession).(set.Sets)
	currentSetIndex := ds.s.Get(session.CurrentSetSession).(int)
	currentItemIndex := ds.s.Get(session.CurrentItemSession).(int)

	currentItemState := ds.s.Get(session.CurrentItemStateSession).(int)
	if currentItemState == stateAnswer {
		log.Info("Recording answer.")
		var results = make(map[int]set.SetResult)
		if ds.s.Get(session.ResultsSession) == nil {
			ds.s.Set(session.ResultsSession, results)
		}
		results = ds.s.Get(session.ResultsSession).(map[int]set.SetResult)
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
		ds.s.Set(session.ResultsSession, results)
	}
	if currentItemIndex == len(sets[currentSetIndex])-1 && currentItemState == stateAnswer { //all items showed
		log.Info("All items have been shown")
		//let's show results
		if ds.s.Get(session.ShowResultSession) != nil {
			ds.s.Set(session.ShowResultSession, nil)
			log.Info("Moving on to next set")
			if currentSetIndex == len(sets)-1 { //all sets showed
				log.Info("All sets shown, moving to next state")
				user := ds.s.Get(session.UserSession).(*models.User)
				user.CurrentStep = models.StepDemoThree
				user.SetTimeout(ds.s.Get(session.TotalTimeSession).(map[int]int64))
				err := user.Update()
				if err != nil {
					return err
				}
				ds.s.Set(session.UserSession, user)
				ds.s.Set(session.ShowResultSession, nil)
				ds.s.Set(session.SetsSession, nil)
				ds.s.Set(session.CurrentItemSession, nil)
				ds.s.Set(session.CurrentSetSession, nil)
				ds.s.Set(session.CurrentItemStateSession, nil)
				ds.s.Set(session.ResultsSession, nil)
				err = ds.s.Save()
				return err
			}
			ds.s.Set(session.CurrentSetSession, currentSetIndex+1)
			ds.s.Set(session.CurrentItemSession, 0)
			ds.s.Set(session.CurrentItemStateSession, 0)
			err := ds.s.Save()
			return err
		}
		ds.s.Set(session.ShowResultSession, 1)
		err := ds.s.Save()
		return err
	}
	if currentItemState == stateAnswer {
		//increase item
		ds.s.Set(session.CurrentItemStateSession, 0)
		ds.s.Set(session.CurrentItemSession, currentItemIndex+1)
		log.Infof("Incrementing item index to %d", currentItemIndex+1)
	} else {
		log.Info("Increasing item state to show answer in next run")
		totalTime := ds.s.Get(session.TotalTimeSession).(map[int]int64)
		startTime := ds.s.Get(session.StartTimeSession).(int64)
		//save time diff in milli seconds
		totalTime[currentItemIndex] = (time.Now().UnixNano() - startTime) / int64(time.Millisecond)
		ds.s.Set(session.TotalTimeSession, totalTime)
		ds.s.Save()
		ds.s.Set(session.CurrentItemStateSession, currentItemState+1)
		log.Infof("Incrementing item state to %d", currentItemState+1)
	}
	return ds.s.Save()
}
func (ds *DemoTwoState) SetSession(s session.Session) {
	ds.s = s
}
