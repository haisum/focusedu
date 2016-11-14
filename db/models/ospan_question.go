package models

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db"
)

/*
CREATE TABLE OSPANQuestion(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                                  Question TEXT, OPTION TEXT, isTrue INTEGER);
*/
type OSPANQuestion struct {
	ID       int
	Question string
	Option   string
	IsTrue   int
}

func GetQuestions(userid, total int) ([]*OSPANQuestion, error) {
	db := db.Get()
	var questions []*OSPANQuestion
	usedQuestions := ""
	stmt, err := db.Preparex("SELECT UsedQuestions FROM USER WHERE ID=?")
	if err != nil {
		return questions, err
	}
	err = stmt.Get(&usedQuestions, userid)
	if err != nil {
		return questions, err
	}
	if usedQuestions == "" {
		usedQuestions = "-1"
	}
	log.Info("Executing " + "SELECT * FROM OSPANQuestion WHERE ID NOT IN (" + usedQuestions + ") ORDER BY RANDOM() LIMIT ?")
	stmt, err = db.Preparex("SELECT * FROM OSPANQuestion WHERE ID NOT IN (" + usedQuestions + ") ORDER BY RANDOM() LIMIT ?")
	if err != nil {
		log.WithError(err).Error("Couldn't prepare statement")
		return questions, err
	}
	err = stmt.Select(&questions, total)
	stmt.Close()
	if err != nil {
		log.WithError(err).Error("Couldn't get questions")
		return questions, err
	}
	if len(questions) != total {
		return questions, fmt.Errorf("Not enough questions in db. Expected: %d, Found: %d", total, len(questions))
	}
	usedQuestionsList := strings.Split(usedQuestions, ",")
	for _, v := range questions {
		usedQuestionsList = append(usedQuestionsList, strconv.FormatInt(int64(v.ID), 10))
	}
	stmt, err = db.Preparex("UPDATE USER SET UsedQuestions = ? WHERE ID = ?")
	if err != nil {
		return questions, err
	}
	_, err = stmt.Exec(strings.Join(usedQuestionsList, ","), userid)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error updating used questions list.")
		return questions, err
	}
	return questions, err
}
