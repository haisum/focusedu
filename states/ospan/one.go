package ospan

import (
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/session"
)

const (
	totalDemoOneSets     = 5
	demoOneIntroTemplate = "demo1_intro.gohtml"
	letterGridTemplate   = "ospan_letter_grid.gohtml"
	letterTemplate       = "ospan_letter.gohtml"
	resultTemplate       = "ospan_result.gohtml"
)

type DemoOneState struct {
	s session.Session
}

var (
	letters = [][]string{
		{
			"H", "L", "Q",
		},
		{
			"R", "N",
		},
		{
			"Y", "T",
		},
	}
)

func (ds *DemoOneState) Render(w io.Writer, values url.Values) error {
	if ds.s.Get(session.CurrentSetSession) == nil {
		return ds.renderIntro(w, values)
	}
	if ds.s.Get(session.ShowGridSession) != nil {
		return renderTemplate(w, letterGridTemplate, nil)
	}
	//show letter
	currentSetIndex := ds.s.Get(session.CurrentSetSession).(int)
	currentLetterIndex := ds.s.Get(session.CurrentLetterSession).(int)
	if currentSetIndex > len(letters) || currentLetterIndex > len(letters[currentSetIndex]) {
		log.Error("This shouldn't happen!. Something's wrong in process function of demo one.")
		return errors.New("Error in process function of demo one")
	}
	if ds.s.Get(session.ShowResultSession) != nil {
		return renderTemplate(w, resultTemplate, map[string]string{
			"Total":          strconv.FormatInt(int64(len(letters[currentSetIndex])), 10),
			"CorrectLetters": strconv.FormatInt(int64(ds.s.Get(session.ShowResultSession).(int)), 10),
			"Percentage":     strconv.FormatFloat((float64(ds.s.Get(session.ShowResultSession).(int)) / float64(len(letters[currentSetIndex])) * 100.0), 'f', 1, 64),
		})
	}
	return renderTemplate(w, letterTemplate, map[string]string{
		"Letter": letters[currentSetIndex][currentLetterIndex],
	})
}

func (ds *DemoOneState) renderIntro(w io.Writer, values url.Values) error {
	user := ds.s.Get(session.UserSession).(*models.User)
	return renderTemplate(w, demoOneIntroTemplate, map[string]string{"Name": user.Name})
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

func (ds *DemoOneState) Process(values url.Values) error {
	if ds.s.Get(session.CurrentSetSession) == nil {
		log.Info("No sets defined, setting them")
		ds.s.Set(session.CurrentSetSession, 0)
		ds.s.Set(session.CurrentLetterSession, 0)
		err := ds.s.Save()
		return err
	}
	currentSetIndex := ds.s.Get(session.CurrentSetSession).(int)
	currentLetterIndex := ds.s.Get(session.CurrentLetterSession).(int)
	if currentLetterIndex == len(letters[currentSetIndex])-1 { //all letters showed
		log.Info("All letters have been shown, setting show grid to true")
		ds.s.Set("showgrid", true)
		ds.s.Set(session.CurrentLetterSession, currentLetterIndex+1)
		err := ds.s.Save()
		return err
	}
	if currentLetterIndex == len(letters[currentSetIndex]) { //process grid
		//process grid
		ds.s.Set("showgrid", nil)
		if ds.s.Get(session.ShowResultSession) != nil {
			ds.s.Set(session.ShowResultSession, nil)
			log.Info("Moving on to next set")
			if currentSetIndex == len(letters)-1 { //all sets showed
				log.Info("All sets shown, moving to next state")
				user := ds.s.Get(session.UserSession).(*models.User)
				user.CurrentStep = models.StepDemoTwo
				err := user.Update()
				if err != nil {
					return err
				}
				ds.s.Set(session.ShowResultSession, nil)
				ds.s.Set(session.CurrentSetSession, nil)
				ds.s.Set(session.CurrentLetterSession, nil)
				ds.s.Set(session.UserSession, user)
				err = ds.s.Save()
				return err
			}
			ds.s.Set(session.CurrentSetSession, currentSetIndex+1)
			ds.s.Set(session.CurrentLetterSession, 0)
			err := ds.s.Save()
			return err
		}
		log.Info("processing grid, saving results")
		givenLetters := strings.Split(values.Get("Letters"), ",")
		correctCount := 0
		for i := 0; i < len(givenLetters) && i < len(letters[currentSetIndex]); i++ {
			if givenLetters[i] == letters[currentSetIndex][i] {
				correctCount = correctCount + 1
			}
		}
		ds.s.Set(session.ShowResultSession, correctCount)
		err := ds.s.Save()
		return err
	}
	log.Infof("Incrementing letter index to %d", currentLetterIndex+1)
	ds.s.Set(session.CurrentLetterSession, currentLetterIndex+1)
	return ds.s.Save()
}
func (ds *DemoOneState) SetSession(s session.Session) {
	ds.s = s
}
