package set

import (
	"math/rand"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/random"
)

type SetItem struct {
	Question *models.OSPANQuestion
	Letter   string
}

type Set []SetItem

type SetResult struct {
	CorrectLetters int
	CorrectAnswers int
	Total          int
}

type Sets []Set

var randLetters = []string{"F", "H", "J", "K", "L", "N", "P", "Q", "R", "S", "T", "Y"}

func GetSets(userid, total int, haveLetters, haveQuestions bool) (Sets, error) {
	var sets Sets
	for i := 0; i < total; i++ {
		var s Set
		letterRand := random.NewUniqueRand(len(randLetters))
		//return 0 or 1 randomly
		totalSetElems := rand.Intn(2) + 2
		var (
			qs  []*models.OSPANQuestion
			err error
		)
		log.Infof("Total elems: %d", totalSetElems)
		if haveQuestions {
			qs, err = models.GetQuestions(userid, totalSetElems)
		}
		if err != nil {
			log.WithError(err).Error("Error in fetching questions from db")
			return sets, err
		}
		// 2 or 3 elements
		for j := 0; j < totalSetElems; j++ {
			setItem := SetItem{}
			if haveLetters {
				setItem.Letter = randLetters[letterRand.Int()]
			}
			if haveQuestions {
				setItem.Question = qs[j]
			}
			s = append(s, setItem)
		}
		sets = append(sets, s)
	}

	return sets, nil
}
func GetQuestionsSet(userid, total int) (Sets, error) {
	var sets Sets
	var s Set
	qs, err := models.GetQuestions(userid, total)
	if err != nil {
		log.WithError(err).Error("Error in fetching questions from db")
		return sets, err
	}
	for i := 0; i < total; i++ {
		setItem := SetItem{}
		setItem.Question = qs[i]
		s = append(s, setItem)
	}
	sets = append(sets, s)
	return sets, nil
}

func GetRealSets(userid int) (Sets, error) {
	var sets Sets
	setRand := random.NewUniqueRand(5)
	for i := 0; i < 5; i++ {
		var s Set
		letterRand := random.NewUniqueRand(len(randLetters))
		totalSetElems := setRand.Int() + 3
		var (
			qs  []*models.OSPANQuestion
			err error
		)
		log.Infof("Total elems: %d", totalSetElems)
		qs, err = models.GetQuestions(userid, totalSetElems)
		if err != nil {
			log.WithError(err).Error("Error in fetching questions from db")
			return sets, err
		}
		for j := 0; j < totalSetElems; j++ {
			setItem := SetItem{}
			setItem.Letter = randLetters[letterRand.Int()]
			setItem.Question = qs[j]
			s = append(s, setItem)
		}
		sets = append(sets, s)
	}
	return sets, nil
}
