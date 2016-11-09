package models

import (
	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db"
)

/*
 *ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                         Name TEXT, Age INTEGER, RollNo TEXT, Gender INTEGER, RegisteredAt INTEGER, MidtermScore INTEGER, CurrentScore INTEGER
*/

type User struct {
	ID           int
	Name         string
	Age          int
	RollNo       string
	Gender       int
	RegisteredAt int
	MidtermScore int
	CurrentScore int
	CurrentStep  Step
}

type Step int

const (
	Info Step = iota
	Demo
	OSPAN
	Module
	ModuleTest
	Feedback
)

func GetUser(RollNo string) (User, error) {
	db := db.Get()
	user := User{}
	stmt, err := db.Preparex("SELECT * FROM User WHERE RollNo=?")
	err = stmt.Get(&user, RollNo)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error":  err,
			"RollNo": RollNo,
		}).Error("Error getting user.")
		return user, err
	}
	return user, nil
}

func MakeUser(RollNo string) error {
	db := db.Get()
	stmt, err := db.Preparex("INSERT INTO USER (RollNo) VALUES(?)")
	_, err = stmt.Exec(RollNo)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error":  err,
			"RollNo": RollNo,
		}).Error("Error inserting user.")
		return err
	}
	return nil
}
