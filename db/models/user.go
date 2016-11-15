package models

import (
	"math"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/haisum/focusedu/db"
)

/*
 *ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                         Name TEXT, Age INTEGER, RollNo TEXT, Gender INTEGER, RegisteredAt INTEGER, MidtermScore INTEGER, CurrentScore INTEGER
*/

type User struct {
	ID                        int
	Name                      string
	Age                       int
	RollNo                    string
	Gender                    Gender
	RegisteredAt              int64
	MidtermScore              int
	CurrentScore              int
	CurrentStep               Step
	QuestionTimeout           int64
	UsedQuestions             string
	OSPANScore                int
	TotalCorrect              int
	SpeedErrors               int
	AccuracyErrors            int
	MathErrors                int
	Type                      int
	ModuleOneDistractionCount int
	ModuleOneExampleCount     int
	ModuleOneGraspingCount    int
	ModuleOneCorrect          int
	ModuleTwoDistractionCount int
	ModuleTwoExampleCount     int
	ModuleTwoGraspingCount    int
	ModuleTwoCorrect          int
}

type Step int

const (
	StepInfo Step = iota
	StepDemoOne
	StepDemoTwo
	StepDemoThree
	StepOSPAN
	StepModuleOne
	StepModuleOneTest
	StepModuleTwo
	StepModuleTwoTest
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
	userType := 1
	if strings.HasPrefix(RollNo, "2-") {
		userType = 2
	}
	stmt, err := db.Preparex("INSERT INTO USER (RollNo, CurrentStep, RegisteredAt, Type) VALUES(?, ?, ?, ?)")
	_, err = stmt.Exec(RollNo, StepInfo, time.Now().Unix(), userType)
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
	stmt, err := db.Preparex(`UPDATE USER SET Age = ? , Name= ?, Gender=?, MidtermScore=?, CurrentStep=?,  
								OSPANScore=?, TotalCorrect=?, SpeedErrors=?, AccuracyErrors=?, MathErrors=?
								WHERE ID=?`)
	_, err = stmt.Exec(u.Age, u.Name, u.Gender, u.MidtermScore, u.CurrentStep, u.OSPANScore, u.TotalCorrect, u.SpeedErrors, u.AccuracyErrors, u.MathErrors, u.ID)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"User":  u,
		}).Error("Error updating user.")
	}
	return err
}

func (u *User) UpdateModuleOneStats(userid, distraction, example, grasping int) error {
	db := db.Get()
	stmt, err := db.Preparex(`UPDATE USER SET ModuleOneDistractionCount = ? , ModuleOneExampleCount= ?, ModuleOneGraspingCount=?
								WHERE ID=?`)
	_, err = stmt.Exec(userid, distraction, example, grasping)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error updating user.")
	}
	return err
}

func (u *User) UpdateModuleOneScore(userid, score int) error {
	db := db.Get()
	stmt, err := db.Preparex(`UPDATE USER SET ModuleOneCorrect = ?
								WHERE ID=?`)
	_, err = stmt.Exec(userid, score)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error updating user.")
	}
	return err
}

func (u *User) UpdateModuleTwoScore(userid, score int) error {
	db := db.Get()
	stmt, err := db.Preparex(`UPDATE USER SET ModuleTwoCorrect = ?
								WHERE ID=?`)
	_, err = stmt.Exec(userid, score)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error updating user.")
	}
	return err
}

func (u *User) UpdateModuleTwoStats(userid, distraction, example, grasping int) error {
	db := db.Get()
	stmt, err := db.Preparex(`UPDATE USER SET ModuleTwoDistractionCount = ? , ModuleTwoExampleCount= ?, ModuleTwoGraspingCount=?
								WHERE ID=?`)
	_, err = stmt.Exec(userid, distraction, example, grasping)
	defer stmt.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error updating user.")
	}
	return err
}

func (u *User) SetTimeout(totalTimes map[int]int64) {
	var totalTime int64
	for _, v := range totalTimes {
		totalTime += v
	}
	mean := float64(totalTime) / float64(len(totalTimes))
	var variance float64
	for _, v := range totalTimes {
		variance += math.Pow(float64(v)-mean, 2)
	}
	variance = variance / float64(len(totalTimes))
	stDev := math.Sqrt(variance)
	u.QuestionTimeout = int64(math.Ceil(mean + (2.5 * stDev)))
	log.WithFields(log.Fields{
		"totalTime":       totalTime,
		"mean":            mean,
		"variance":        variance,
		"stDev":           stDev,
		"QuestionTimeout": u.QuestionTimeout,
		"IndividualTimes": totalTimes,
	}).Info("Setting user timeout")
}
