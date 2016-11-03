package models

type User struct {
	CurrentStep Step
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
