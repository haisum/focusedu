package set

import "github.com/haisum/focusedu/db/models"

type SetItem struct {
	Question models.OSPANQuestion
	Letter   string
}

type Set struct {
	Items        []SetItem
	HasLetters   bool
	HasQuestions bool
}

func GetSets(total int, haveLetters, haveQuestions bool) ([]Set, error) {
	var sets []Set

	return sets, nil
}
