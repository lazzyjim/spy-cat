package repository

import (
	"spy-cat/src/db"
)

func InsertSpyCat(cat db.SpyCat) (err error) {
	_, err = db.Conn().Exec("INSERT INTO spy_cats (name, years_of_experience, breed, breed_validation, salary) VALUES ($1, $2, $3, $4, $5)",
		cat.Name, cat.YearsOfExperience, cat.Breed, cat.BreedValidation, cat.Salary)
	return err
}

func DeleteSpyCat(id int) (err error) {
	_, err = db.Conn().Exec("DELETE FROM spy_cats WHERE id = $1", id)
	return err
}

func UpdateSpyCat(cat db.SpyCat, id int) (err error) {
	_, err = db.Conn().Exec("UPDATE spy_cats SET salary = $1 WHERE id = $2", cat.Salary, id)
	return err
}

func GetSingleSpyCat(cat *db.SpyCat, id int) (err error) {
	err = db.Conn().QueryRow("SELECT id, name, years_of_experience, breed, breed_validation, salary FROM spy_cats WHERE id = $1", id).Scan(
		&cat.Id, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.BreedValidation, &cat.Salary)
	return err
}

func GetSpyCatList() (err error, cats []db.SpyCat) {
	rows, err := db.Conn().Query("SELECT id, name, years_of_experience, breed, breed_validation, salary FROM spy_cats")
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var cat db.SpyCat
		if err := rows.Scan(&cat.Id, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.BreedValidation, &cat.Salary); err != nil {
			return err, nil
		}
		cats = append(cats, cat)
	}
	return err, cats
}

func InsertMission(mission db.Mission) (err error) {
	var missionID int
	err = db.Conn().QueryRow("INSERT INTO missions (name) VALUES ($1) RETURNING id",
		mission.Name).Scan(&missionID)
	if err != nil {
		return err
	}

	for _, target := range mission.Targets {
		_, err = db.Conn().Exec("INSERT INTO targets (mission_id, name, country, notes, complete_state) VALUES ($1, $2, $3, $4, $5)",
			missionID, target.Name, target.Country, target.Notes, target.CompleteState)
		if err != nil {

			return err
		}
	}

	return err
}

func DeleteMission(id int) (err error) {
	_, err = db.Conn().Exec("DELETE FROM missions WHERE id = $1 AND cat_id IS NULL", id)
	if err != nil {
		return err
	}
	_, err = db.Conn().Exec("DELETE FROM targets WHERE mission_id = $1", id)
	return nil
}

func CompleteMission(mission db.Mission, id int) (err error) {
	_, err = db.Conn().Exec("UPDATE missions SET complete_state = $1 WHERE id = $2", mission.CompleteState, id)
	return err
}

func CompleteTarget(targetID, missionID int, target db.Target) (err error) {
	_, err = db.Conn().Exec("UPDATE targets SET complete_state = $1 WHERE mission_id = $2 AND id = $3", target.CompleteState, missionID, targetID)
	return err
}

func GetSingleMission(mission *db.Mission, id int) (err error) {
	err = db.Conn().QueryRow("SELECT id, name, cat_id, complete_state FROM missions WHERE id = $1", id).Scan(
		&mission.Id, &mission.Name, &mission.CatId, &mission.CompleteState)
	if err != nil {
		return err
	}

	rows, err := db.Conn().Query("SELECT id, mission_id, name, country, notes, complete_state FROM targets WHERE mission_id = $1", id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var target db.Target
		if err := rows.Scan(&target.Id, &target.MissionId, &target.Name, &target.Country, &target.Notes, &target.CompleteState); err != nil {
			return err
		}
		mission.Targets = append(mission.Targets, target)
	}
	return err
}

func GetMissionsList(missions *[]db.Mission) (err error) {
	rows, err := db.Conn().Query("SELECT id, name, cat_id, complete_state FROM missions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var mission db.Mission
		if err := rows.Scan(&mission.Id, &mission.Name, &mission.CatId, &mission.CompleteState); err != nil {
			return err
		}

		targetRows, err := db.Conn().Query("SELECT id, mission_id, name, country, notes, complete_state FROM targets WHERE mission_id = $1", mission.Id)
		if err != nil {
			return err
		}
		defer targetRows.Close()

		for targetRows.Next() {
			var target db.Target
			if err := targetRows.Scan(&target.Id, &target.MissionId, &target.Name, &target.Country, &target.Notes, &target.CompleteState); err != nil {
				return err
			}
			mission.Targets = append(mission.Targets, target)
		}

		*missions = append(*missions, mission)
	}
	return err
}

func InsertTarget(missionID int, target db.Target) (err error) {
	_, err = db.Conn().Exec("INSERT INTO targets (mission_id, name, country, notes, complete_state) VALUES ($1, $2, $3, $4, $5)",
		missionID, target.Name, target.Country, target.Notes, target.CompleteState)
	if err != nil {
		return err
	}
	return err
}

func DeleteTarget(targetID, missionID int) (err error) {
	_, err = db.Conn().Exec("DELETE FROM targets WHERE id = $1 AND mission_id = $2", targetID, missionID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTargetNotes(missionID, targetID int, target *db.Target) (err error, missionComplete, targetComplete bool) {
	err = db.Conn().QueryRow("SELECT complete_state FROM missions WHERE id = $1", missionID).Scan(&missionComplete)
	if err != nil {
		return err, missionComplete, targetComplete
	}

	err = db.Conn().QueryRow("SELECT complete_state FROM targets WHERE id = $1", targetID).Scan(&targetComplete)
	if err != nil {
		return err, missionComplete, targetComplete
	}

	if missionComplete || targetComplete {
		return nil, missionComplete, targetComplete
	}

	_, err = db.Conn().Exec("UPDATE targets SET notes = $1 WHERE id = $2", target.Notes, targetID)
	if err != nil {
		return err, missionComplete, targetComplete
	}

	return nil, missionComplete, targetComplete
}

func AssignCatToMission(mission *db.Mission, missionID int) (err error) {
	_, err = db.Conn().Exec("UPDATE missions SET cat_id = $1 WHERE id = $2", mission.CatId, missionID)
	if err != nil {
		return err
	}

	return nil
}
