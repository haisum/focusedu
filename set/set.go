package set

import (
	"errors"
	"math/rand"

	"github.com/haisum/focusedu/db/models"
	"github.com/haisum/focusedu/random"
)

type SetItem struct {
	Question *models.OSPANQuestion
	Letter   string
}

type Set struct {
	items        []SetItem
	HasLetters   bool
	HasQuestions bool
}

type SetResult struct {
	CorrectLetters int
	CorrectAnswers int
	Total          int
}

func (s *Set) Pop() (SetItem, error) {
	var i SetItem
	total := len(s.items)
	if total == 0 {
		return i, errors.New("No items left")
	}
	i, s.items = s.items[total-1], s.items[:(total-1)]
	return i, nil
}

func (s *Set) Push(i SetItem) {
	s.items = append(s.items, i)
}

func (s *Set) Size() int {
	return len(s.items)
}

type Sets struct {
	items []Set
}

func (s *Sets) Pop() (Set, error) {
	var i Set
	total := len(s.items)
	if total == 0 {
		return i, errors.New("No items left")
	}
	i, s.items = s.items[total-1], s.items[:(total-1)]
	return i, nil
}

func (s *Sets) Push(i Set) {
	s.items = append(s.items, i)
}

func (s *Sets) Size() int {
	return len(s.items)
}

var randLetters = []string{"F", "H", "J", "K", "L", "N", "P", "Q", "R", "S", "T", "Y"}

func GetSets(total int, haveLetters, haveQuestions bool) (Sets, error) {
	var sets Sets
	letterRand := random.NewUniqueRand(len(randLetters))
	for i := 0; i < total; i++ {
		var s Set
		//return 0 or 1 randomly
		totalSetElems := rand.Intn(2)
		// 2 or 3 elements
		for j := 0; j <= totalSetElems+2; j++ {
			s.Push(SetItem{
				Letter: randLetters[letterRand.Int()],
			})
		}
		sets.Push(s)
	}

	return sets, nil
}
