package restful

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"spy-cat/src/common"
	"spy-cat/src/db"
	"spy-cat/src/models"
	"spy-cat/src/repository"
	"strconv"
)

func Handler() func(ctx *fasthttp.RequestCtx) {
	router := routing.New()
	router.Use(common.LoggingMiddleware)

	spyCat := router.Group("/spy_cat")
	spyCat.Post("/create", CreateCat)
	spyCat.Patch("/update/<id>", UpdateCat)
	spyCat.Delete("/delete/<id>", DeleteCat)
	spyCat.Get("/single/<id>", GetSingleSpyCat)
	spyCat.Get("/list", GetSpyCatList)

	mission := router.Group("/mission")
	mission.Post("/create", CreateMission)
	mission.Post("/<mission_id>/assign_to_cat", AssignCatToMission)
	mission.Patch("/<mission_id>/update/add_target", AddTarget)
	mission.Delete("/<mission_id>/update/delete_target/<target_id>", DeleteTarget)
	mission.Patch("/<mission_id>/update/target/<target_id>/notes", UpdateTargetNotes)
	mission.Patch("/<mission_id>/update/target/<target_id>/complete", CompleteTarget)
	mission.Delete("/delete/<mission_id>", DeleteMission)
	mission.Get("/list", GetMissionsList)
	mission.Get("/single/<mission_id>", GetSingleMission)
	mission.Patch("/complete/<mission_id>", CompleteMission)
	return router.HandleRequest
}

func CreateCat(ctx *routing.Context) error {
	var cat db.SpyCat
	if err := json.Unmarshal(ctx.PostBody(), &cat); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	v, err := models.ValidateBreed(cat.Breed)
	if err != nil {
		return err
	}
	cat.BreedValidation = v
	err = repository.InsertSpyCat(cat)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	return nil
}

func DeleteCat(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	err = repository.DeleteSpyCat(id)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}

func UpdateCat(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var cat db.SpyCat
	if err := json.Unmarshal(ctx.PostBody(), &cat); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	err = repository.UpdateSpyCat(cat, id)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

func GetSingleSpyCat(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var cat db.SpyCat
	err = repository.GetSingleSpyCat(&cat, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		return err
	}

	response, err := json.Marshal(cat)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetContentType("application/json")
	ctx.SetBody(response)
	return nil
}

func GetSpyCatList(ctx *routing.Context) error {
	err, listOfCats := repository.GetSpyCatList()

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	response, err := json.Marshal(listOfCats)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetContentType("application/json")
	ctx.SetBody(response)
	return nil
}

func CreateMission(ctx *routing.Context) error {
	var mission db.Mission
	if err := json.Unmarshal(ctx.PostBody(), &mission); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	if mission.Name == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("Mission name is required"))
		return nil
	}

	if len(mission.Targets) < 1 || len(mission.Targets) > 3 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("A mission must have between 1 and 3 targets"))
		return nil
	}

	err := repository.InsertMission(mission)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	return nil
}

func GetSingleMission(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("mission_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var mission db.Mission
	err = repository.GetSingleMission(&mission, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		return err
	}

	response, err := json.Marshal(mission)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetContentType("application/json")
	ctx.SetBody(response)
	return nil
}

func CompleteMission(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("mission_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var mission db.Mission
	if err := json.Unmarshal(ctx.PostBody(), &mission); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}
	err = repository.CompleteMission(mission, id)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

func DeleteMission(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("mission_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var catID sql.NullInt64
	err = db.Conn().QueryRow("SELECT cat_id FROM missions WHERE id = $1", id).Scan(&catID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		return err
	}

	if catID.Valid {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return nil
	}

	err = repository.DeleteMission(id)

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}

func GetMissionsList(ctx *routing.Context) error {
	var missions []db.Mission

	err := repository.GetMissionsList(&missions)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	response, err := json.Marshal(missions)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetContentType("application/json")
	ctx.SetBody(response)
	return nil
}

func UpdateTargetNotes(ctx *routing.Context) error {
	missionID, err := strconv.Atoi(ctx.Param("mission_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	targetID, err := strconv.Atoi(ctx.Param("target_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var target db.Target
	if err := json.Unmarshal(ctx.PostBody(), &target); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	err, missionComplete, targetComplete := repository.UpdateTargetNotes(missionID, targetID, &target)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	if missionComplete || targetComplete {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

func AddTarget(ctx *routing.Context) error {
	missionID, err := strconv.Atoi(ctx.Param("mission_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var target db.Target
	if err := json.Unmarshal(ctx.PostBody(), &target); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var missionComplete bool

	err = db.Conn().QueryRow("SELECT complete_state FROM missions WHERE id = $1", missionID).Scan(&missionComplete)
	if err != nil {
		return err
	}

	if missionComplete {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return nil
	}

	var targetCount int
	err = db.Conn().QueryRow("SELECT COUNT(*) FROM targets WHERE mission_id = $1", missionID).Scan(&targetCount)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	if targetCount >= 3 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("A mission cannot have more than 3 targets"))
		return nil
	}

	var existingTargetID int
	err = db.Conn().QueryRow("SELECT id FROM targets WHERE mission_id = $1 AND name = $2 AND country = $3 AND notes = $4", missionID, target.Name, target.Country, target.Notes).Scan(&existingTargetID)
	if err == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("Target with this name already exists in the mission"))
		return nil
	}

	err = repository.InsertTarget(missionID, target)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	return nil
}

func DeleteTarget(ctx *routing.Context) error {
	missionID, err := strconv.Atoi(ctx.Param("mission_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	targetID, err := strconv.Atoi(ctx.Param("target_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var targetComplete bool

	err = db.Conn().QueryRow("SELECT complete_state FROM targets WHERE id = $1 AND mission_id = $2", targetID, missionID).Scan(&targetComplete)
	if err != nil {
		return err
	}

	if targetComplete {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return nil
	}

	err = repository.DeleteTarget(targetID, missionID)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}

func CompleteTarget(ctx *routing.Context) error {
	missionID, err := strconv.Atoi(ctx.Param("mission_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	targetID, err := strconv.Atoi(ctx.Param("target_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var target db.Target
	if err := json.Unmarshal(ctx.PostBody(), &target); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	err = repository.CompleteTarget(targetID, missionID, target)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}

func AssignCatToMission(ctx *routing.Context) error {
	missionID, err := strconv.Atoi(ctx.Param("mission_id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var mission db.Mission
	if err := json.Unmarshal(ctx.PostBody(), &mission); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	var existingMissionID int
	err = db.Conn().QueryRow("SELECT id FROM missions WHERE cat_id = $1 AND complete_state = false", mission.CatId).Scan(&existingMissionID)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return err
		}
		ctx.SetStatusCode(fasthttp.StatusConflict)
		ctx.SetBody([]byte("Cat is already assigned to another mission"))
		return nil
	}

	err = repository.AssignCatToMission(&mission, missionID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}
