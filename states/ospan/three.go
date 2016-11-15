package ospan

import (
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
	"github.com/haisum/focusedu/set"
)

const (
	demoThreeIntroTemplate = "demo3_intro.gohtml"
)

type DemoThreeState struct {
	s session.Session
}

func (ds *DemoThreeState) Render(w io.Writer, values url.Values) error {
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
	if ds.s.Get(session.ShowGridSession) != nil {
		return renderTemplate(w, letterGridTemplate, nil)
	}
	if ds.s.Get("result") != nil {
		results := ds.s.Get(session.ResultsSession).(map[int]set.SetResult)
		result := results[currentSetIndex]
		user := ds.s.Get(session.UserSession).(*models.User)
		return renderTemplate(w, resultTemplate, map[string]string{
			"Timeout":        strconv.FormatInt(user.QuestionTimeout, 10),
			"Total":          strconv.FormatInt(int64(result.Total), 10),
			"CorrectAnswers": strconv.FormatInt(int64(result.CorrectAnswers), 10),
			"CorrectLetters": strconv.FormatInt(int64(result.CorrectLetters), 10),
			"Percentage":     strconv.FormatFloat((float64(result.CorrectAnswers) / float64(result.Total) * 100.0), 'f', 1, 64),
		})
	}
	switch ds.s.Get(session.CurrentItemStateSession).(int) {
	case stateQuestion:
		user := ds.s.Get(session.UserSession).(*models.User)
		return renderTemplate(w, questionTemplate, map[string]string{
			"Question": sets[currentSetIndex][currentItemIndex].Question.Question,
			"Timeout":  strconv.FormatInt(user.QuestionTimeout, 10),
		})
	case stateAnswer:
		return renderTemplate(w, answerTemplate, map[string]string{
			"Option": sets[currentSetIndex][currentItemIndex].Question.Option,
		})
	case stateLetter:
		return renderTemplate(w, letterTemplate, map[string]string{
			"Letter": sets[currentSetIndex][currentItemIndex].Letter,
		})
	}
	return nil
}

func (ds *DemoThreeState) renderIntro(w io.Writer, values url.Values) error {
	user := ds.s.Get("user").(*models.User)
	return renderTemplate(w, demoThreeIntroTemplate, map[string]string{"Name": user.Name})
}

func (ds *DemoThreeState) Process(values url.Values) error {
	if ds.s.Get(session.SetsSession) == nil {
		log.Info("No sets defined, setting them")
		sets, err := set.GetSets(ds.s.Get(session.UserSession).(*models.User).ID, 3, true, true)
		if err != nil {
			log.WithError(err).Error("Couldn't get sets")
			return err
		}
		ds.s.Set(session.SetsSession, sets)
		ds.s.Set(session.CurrentItemSession, 0)
		ds.s.Set(session.CurrentSetSession, 0)
		ds.s.Set(session.CurrentItemStateSession, 0)
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
	if currentItemIndex == len(sets[currentSetIndex])-1 && currentItemState == stateLetter { //all items showed
		log.Info("All items have been shown")

		//let's show results
		if ds.s.Get(session.ShowResultSession) != nil {
			ds.s.Set(session.ShowResultSession, nil)
			log.Info("Moving on to next set")
			if currentSetIndex == len(sets)-1 { //all sets showed
				log.Info("All sets shown, moving to next state")
				user := ds.s.Get(session.UserSession).(*models.User)
				user.CurrentStep = models.StepOSPAN
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
			ds.s.Set(session.ShowResultSession, nil)
			err := ds.s.Save()
			return err
		}

		if ds.s.Get(session.ShowGridSession) != nil {
			log.Info("Record grid input and show result")
			// process grid here
			givenLetters := strings.Split(values.Get("Letters"), ",")
			correctCount := 0
			for i := 0; i < len(givenLetters) && i < len(sets[currentSetIndex]); i++ {
				if givenLetters[i] == sets[currentSetIndex][i].Letter {
					correctCount = correctCount + 1
				}
			}
			results := ds.s.Get(session.ResultsSession).(map[int]set.SetResult)
			result := set.SetResult{}
			if v, ok := results[currentSetIndex]; ok {
				result = v
			}
			result.CorrectLetters = correctCount
			results[currentSetIndex] = result
			log.Infof("New result %+v", result)
			ds.s.Set(session.ResultsSession, results)
			//set show grid to nil
			ds.s.Set(session.ShowGridSession, nil)
			ds.s.Set(session.ShowResultSession, 1)
			return ds.s.Save()
		}
		log.Info("Show grid")
		ds.s.Set(session.ShowGridSession, 1)
		err := ds.s.Save()
		return err
	}
	if currentItemState == stateLetter {
		//increase item
		ds.s.Set(session.CurrentItemStateSession, 0)
		ds.s.Set(session.CurrentItemSession, currentItemIndex+1)
		log.Infof("Incrementing item index to %d", currentItemIndex+1)
	} else {
		log.Info("Increasing item state")
		ds.s.Set(session.CurrentItemStateSession, currentItemState+1)
		log.Infof("Incrementing item state to %d", currentItemState+1)
	}
	return ds.s.Save()
}
func (ds *DemoThreeState) SetSession(s session.Session) {
	ds.s = s
}
