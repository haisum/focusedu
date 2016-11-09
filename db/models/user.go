package models

import "github.com/haisum/focusedu/db"

type User struct {
	CurrentStep Step
	RollNo      string
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

func GetUser(RollNo string) User {
	db := db.Get()
}
