package db

type SpyCat struct {
	Id                uint16  `json:"id"`
	Name              string  `json:"name"`
	YearsOfExperience uint16  `json:"years_of_experience"`
	Breed             string  `json:"breed"`
	BreedValidation   bool    `json:"breed_validation"`
	Salary            float32 `json:"salary"`
}

type Target struct {
	Id            uint16 `json:"id"`
	MissionId     uint16 `json:"mission_id"`
	Name          string `json:"name"`
	Country       string `json:"country"`
	Notes         string `json:"notes"`
	CompleteState bool   `json:"complete_state"`
}

type Mission struct {
	Id            uint16   `json:"id"`
	Name          string   `json:"name"`
	CatId         uint16   `json:"cat_id"`
	Targets       []Target `json:"targets"`
	CompleteState bool     `json:"complete_state"`
}

type Breed struct {
	Name string `json:"name"`
}
