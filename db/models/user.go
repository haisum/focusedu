package models

import (
	"time"

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
	Gender       Gender
	RegisteredAt int64
	MidtermScore int
	CurrentScore int
	CurrentStep  Step
}

type Step int

const (
	StepInfo Step = iota
	StepDemoOne
	StepDemoTwo
	StepDemoThree
	StepOSPAN
	StepModule
	StepModuleTest
	StepFeedback
)

type Gender int

const (
	Male Gender = iota
	Female
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
	stmt, err := db.Preparex("INSERT INTO USER (RollNo, CurrentStep, RegisteredAt) VALUES(?, ?, ?)")
	_, err = stmt.Exec(RollNo, StepInfo, time.Now().Unix())
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error":  err,
			"RollNo": RollNo,
		}).Error("Error inserting user.")
	}
	return err
}

func (u *User) Validate() map[string]string {
	errors := make(map[string]string)
	if u.Age < 5 || u.Age > 80 {
		errors["Age"] = "Please provide valid age."
	}
	if u.Name == "" {
		errors["Name"] = "Name is required"
	}
	if u.Gender != Male && u.Gender != Female {
		errors["Gender"] = "Unknown Gender"
	}
	if u.MidtermScore < 1 {
		errors["MidtermScore"] = "Invalid midterm score."
	}
	return errors
}

func (u *User) Update() error {
	db := db.Get()
	stmt, err := db.Preparex("UPDATE USER SET AGE = ? , Name= ?, Gender=?, MidtermScore=?, CurrentStep=? WHERE ID=?")
	_, err = stmt.Exec(u.Age, u.Name, u.Gender, u.MidtermScore, u.CurrentStep, u.ID)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"User":  u,
		}).Error("Error updating user.")
	}
	return err
}
