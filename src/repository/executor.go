package repository

import (
	"spy-cat/src/db"
)

func InsertSpyCat(cat db.SpyCat) (err error) {
	_, err = db.Conn().Exec("INSERT INTO spy_cats (name, years_of_experience, breed, salary) VALUES ($1, $2, $3, $4)",
		cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary)
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
	err = db.Conn().QueryRow("SELECT id, name, years_of_experience, breed, salary FROM spy_cats WHERE id = $1", id).Scan(
		&cat.Id, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary)
	return err
}

func GetSpyCatList() (err error, listOfCats []db.SpyCat) {
	rows, err := db.Conn().Query("SELECT id, name, years_of_experience, breed, salary FROM spy_cats")
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	var cats []db.SpyCat
	for rows.Next() {
		var cat db.SpyCat
		if err := rows.Scan(&cat.Id, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary); err != nil {
			return err, nil
		}
		cats = append(cats, cat)
	}
	return err, listOfCats
}

func InsertMission(mission db.Mission) (err error) {
	tx, err := db.Conn().Begin()
	if err != nil {
		return err
	}

	var missionID int
	err = tx.QueryRow("INSERT INTO missions (cat_id, complete_state) VALUES ($1, $2) RETURNING id",
		mission.CatId, mission.CompleteState).Scan(&missionID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, target := range mission.Targets {
		_, err = tx.Exec("INSERT INTO targets (mission_id, name, country, notes, complete_state) VALUES ($1, $2, $3, $4, $5)",
			missionID, target.Name, target.Country, target.Notes, target.CompleteState)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}

func DeleteMission(id int) (err error, catID int) {
	err = db.Conn().QueryRow("SELECT cat_id FROM missions WHERE id = $1", id).Scan(&catID)
	if err != nil {
		return err, catID
	}

	if catID != 0 {
		return nil, catID
	}

	_, err = db.Conn().Exec("DELETE FROM missions WHERE id = $1", id)
	if err != nil {

		return err, catID
	}

	return nil, catID
}

func CompleteMission(mission db.Mission, id int) (err error) {
	_, err = db.Conn().Exec("UPDATE missions SET complete_state = $1 WHERE id = $2", mission.CompleteState, id)
	return err
}

func GetSingleMission(mission *db.Mission, id int) (err error) {
	err = db.Conn().QueryRow("SELECT id, cat_id, complete_state FROM missions WHERE id = $1", id).Scan(
		&mission.Id, &mission.CatId, &mission.CompleteState)
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
	rows, err := db.Conn().Query("SELECT id, cat_id, complete_state FROM missions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var mission db.Mission
		if err := rows.Scan(&mission.Id, &mission.CatId, &mission.CompleteState); err != nil {
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

func InsertTarget(missionID int, target *db.Target) (err error, missionComplete bool) {
	err = db.Conn().QueryRow("SELECT complete_state FROM missions WHERE id = $1", missionID).Scan(&missionComplete)
	if err != nil {
		return err, missionComplete
	}

	if missionComplete {
		return nil, missionComplete
	}

	_, err = db.Conn().Exec("INSERT INTO targets (mission_id, name, country, notes, complete_state) VALUES ($1, $2, $3, $4, $5)",
		missionID, target.Name, target.Country, target.Notes, target.CompleteState)
	if err != nil {
		return err, missionComplete
	}
	return err, missionComplete
}

func DeleteTarget(targetID, missionID int) (err error, targetComplete bool) {
	err = db.Conn().QueryRow("SELECT complete_state FROM targets WHERE id = $1", targetID).Scan(&targetComplete)
	if err != nil {
		return err, targetComplete
	}

	if targetComplete {
		return nil, targetComplete
	}

	_, err = db.Conn().Exec("DELETE FROM targets WHERE id = $1", targetID)
	if err != nil {
		return err, targetComplete
	}

	return nil, targetComplete
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