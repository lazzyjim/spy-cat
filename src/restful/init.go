package restful

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"spy-cat/src/db"
	"spy-cat/src/repository"
	"strconv"
)

func Handler() func(ctx *fasthttp.RequestCtx) {
	router := routing.New()

	spyCat := router.Group("/spy_cat")
	spyCat.Post("/create", CreateCat)
	spyCat.Patch("/update", UpdateCat)
	spyCat.Delete("/delete", DeleteCat)
	spyCat.Get("/single", GetSingleSpyCat)
	spyCat.Get("/list", GetSpyCatList)

	mission := router.Group("/mission")
	mission.Post("/create", CreateMission)
	mission.Post("/assign_to_cat", AssignCatToMission)
	mission.Patch("/update/add_target", AddTarget)
	mission.Patch("/update/delete_target", DeleteTarget)
	mission.Patch("/update/target/notes", UpdateTargetNotes)
	mission.Delete("/delete", DeleteMission)
	mission.Get("/list", GetMissionsList)
	mission.Get("/single", GetSingleMission)
	mission.Patch("/complete", CompleteMission)
	return router.HandleRequest
}

func CreateCat(ctx *routing.Context) error {
	var cat db.SpyCat
	if err := json.Unmarshal(ctx.PostBody(), &cat); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	err := repository.InsertSpyCat(cat)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	return nil
}

// DeleteCat handles deleting a cat by ID
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

// UpdateCat handles updating a cat's salary
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

// GetSingleSpyCat handles fetching a single cat by ID
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

// GetSpyCatList handles listing all cats
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

	err := repository.InsertMission(mission)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	return nil
}

func GetSingleMission(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
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

// CompleteMission handles updating a mission's completion status
func CompleteMission(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
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

// DeleteMission handles deleting a mission by ID
func DeleteMission(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	err, catID := repository.DeleteMission(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		return err
	}

	if catID != 0 {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}

// GetMissionsList handles listing all missions
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

// UpdateTargetNotes handles updating notes of a target
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

// AddTarget handles adding a new target to an existing mission
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

	err, missionComplete := repository.InsertTarget(missionID, &target)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	if missionComplete {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	return nil
}

// DeleteTarget handles deleting a target from an existing mission
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

	err, targetComplete := repository.DeleteTarget(targetID, missionID)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	if targetComplete {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}

// AssignCatToMission handles assigning a cat to a mission
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

	err = repository.AssignCatToMission(&mission, missionID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}
