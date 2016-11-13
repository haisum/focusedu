package set

import (
	"errors"
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
var questions = []*models.OSPANQuestion{
	{Question: "1/1 + 1", Option: "2", IsTrue: 1},
	{Question: "2/1 - 1", Option: "2", IsTrue: 0},
	{Question: "(2 * 3) + 9 - 2", Option: "10", IsTrue: 0},
	{Question: "(1/1 + 1) - 1", Option: "0", IsTrue: 1},
	{Question: "1/2 + 1", Option: "0.5", IsTrue: 0},
	{Question: "3/4 + 1", Option: "1.3", IsTrue: 1},
	{Question: "2-354*1 + 1", Option: "2", IsTrue: 0},
	{Question: "6/12 + 0.5", Option: "1", IsTrue: 1},
	{Question: "1/2 + 1", Option: "1.5", IsTrue: 1},
}

func GetSets(total int, haveLetters, haveQuestions bool) (Sets, error) {
	var sets Sets
	questionRand := random.NewUniqueRand(len(questions))
	for i := 0; i < total; i++ {
		var s Set
		letterRand := random.NewUniqueRand(len(randLetters))
		//return 0 or 1 randomly
		totalSetElems := rand.Intn(2)
		log.Infof("Rand val: %d", totalSetElems)
		// 2 or 3 elements
		for j := 0; j < totalSetElems+2; j++ {
			setItem := SetItem{}
			if haveLetters {
				setItem.Letter = randLetters[letterRand.Int()]
			}
			if haveQuestions {
				qIndex := questionRand.Int()
				log.Infof("qindex %d", qIndex)
				if qIndex == -1 {
					return sets, errors.New("Not enough questions in database!")
				}
				setItem.Question = questions[qIndex]
			}
			s = append(s, setItem)
		}
		sets = append(sets, s)
	}

	return sets, nil
}
