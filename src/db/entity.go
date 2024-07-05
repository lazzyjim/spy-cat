package db

import "time"

type SpyCat struct {
	Id                uint16
	Name              string
	YearsOfExperience uint16
	Breed             string
	BreedValidation   time.Time
	Salary            uint16
}

type Target struct {
	Id            uint16
	MissionId     uint16
	Name          string
	Country       string
	Notes         string
	CompleteState time.Time
}

type Mission struct {
	Id            uint16
	CatId         uint16
	Targets       []Target
	CompleteState time.Time
}
