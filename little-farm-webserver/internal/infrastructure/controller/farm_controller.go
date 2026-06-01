package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"little-farm-webserver/internal/domain/littlefarm"
	"little-farm-webserver/pkg/apperr"
	"little-farm-webserver/pkg/domain"
	"little-farm-webserver/pkg/platform/web"
)

type FarmDAO interface {
	Save(farm littlefarm.Farm) error
	Find(id domain.ID) (littlefarm.Farm, bool, error)
}

type FarmController struct {
	farmDAO FarmDAO
}

func NewFarmController(farmDAO FarmDAO) *FarmController {
	return &FarmController{
		farmDAO: farmDAO,
	}
}

func (controller *FarmController) CreateFarm(response http.ResponseWriter, request *http.Request) {
	var createFarmInput CreateFarmInputDTO
	if err := json.NewDecoder(request.Body).Decode(&createFarmInput); err != nil {
		web.RespondWithError(response, apperr.ErrUnprocessableEntity)
		return
	}

	farm, err := littlefarm.NewFarm(littlefarm.CreateFarmCmd{
		FarmID:    domain.GenerateID(),
		MapWidth:  createFarmInput.MapWidth,
		MapHeight: createFarmInput.MapHeight,
		Timestamp: time.Now(),
	})
	if err != nil {
		web.RespondWithError(response, err)
		return
	}
	if err = controller.farmDAO.Save(farm); err != nil {
		web.RespondWithError(response, err)
		return
	}

	response.Header().Set("Location", "/farms/"+farm.ID().String())
	web.RespondWithJSON(response, http.StatusCreated, FarmToFarmOutputDTO(farm))
}

func (controller *FarmController) FindFarm(response http.ResponseWriter, request *http.Request) {
	farmID, err := domain.NewID(request.PathValue("farm_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidFarmID)
		return
	}

	farm, found, err := controller.farmDAO.Find(farmID)
	if err != nil {
		web.RespondWithError(response, err)
		return
	}
	if !found {
		web.RespondWithError(response, littlefarm.ErrFarmNotFound)
		return
	}

	web.RespondWithJSON(response, http.StatusOK, FarmToFarmOutputDTO(farm))
}

func (controller *FarmController) AdvanceOneSecond(response http.ResponseWriter, request *http.Request) {
	farmID, err := domain.NewID(request.PathValue("farm_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidFarmID)
		return
	}

	farm, found, err := controller.farmDAO.Find(farmID)
	if err != nil {
		web.RespondWithError(response, err)
		return
	}
	if !found {
		web.RespondWithError(response, littlefarm.ErrFarmNotFound)
		return
	}
	if err = farm.AdvanceOneSecond(littlefarm.AdvanceOneSecondCmd{
		FarmID:    farmID,
		Timestamp: time.Now(),
	}); err != nil {
		web.RespondWithError(response, err)
		return
	}
	if err = controller.farmDAO.Save(farm); err != nil {
		web.RespondWithError(response, err)
		return
	}

	web.RespondWithJSON(response, http.StatusCreated, FarmToFarmOutputDTO(farm))
}

func (controller *FarmController) HarvestCorn(response http.ResponseWriter, request *http.Request) {
	farmID, err := domain.NewID(request.PathValue("farm_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidFarmID)
		return
	}
	cornID, err := domain.NewID(request.PathValue("corn_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidCornID)
		return
	}

	farm, found, err := controller.farmDAO.Find(farmID)
	if err != nil {
		web.RespondWithError(response, err)
		return
	}
	if !found {
		web.RespondWithError(response, littlefarm.ErrFarmNotFound)
		return
	}
	if err = farm.HarvestCorn(littlefarm.HarvestCornCmd{
		FarmID:    farmID,
		CornID:    cornID,
		Timestamp: time.Now(),
	}); err != nil {
		web.RespondWithError(response, err)
		return
	}
	if err = controller.farmDAO.Save(farm); err != nil {
		web.RespondWithError(response, err)
		return
	}

	web.RespondWithJSON(response, http.StatusCreated, FarmToFarmOutputDTO(farm))
}

func (controller *FarmController) HarvestGrass(response http.ResponseWriter, request *http.Request) {
	farmID, err := domain.NewID(request.PathValue("farm_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidFarmID)
		return
	}
	grassID, err := domain.NewID(request.PathValue("grass_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidGrassID)
		return
	}

	farm, found, err := controller.farmDAO.Find(farmID)
	if err != nil {
		web.RespondWithError(response, err)
		return
	}
	if !found {
		web.RespondWithError(response, littlefarm.ErrFarmNotFound)
		return
	}
	if err = farm.HarvestGrass(littlefarm.HarvestGrassCmd{
		FarmID:    farmID,
		GrassID:   grassID,
		Timestamp: time.Now(),
	}); err != nil {
		web.RespondWithError(response, err)
		return
	}
	if err = controller.farmDAO.Save(farm); err != nil {
		web.RespondWithError(response, err)
		return
	}

	web.RespondWithJSON(response, http.StatusCreated, FarmToFarmOutputDTO(farm))
}

func (controller *FarmController) PlantCorn(response http.ResponseWriter, request *http.Request) {
	farmID, err := domain.NewID(request.PathValue("farm_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidFarmID)
		return
	}
	dirtID, err := domain.NewID(request.PathValue("dirt_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidDirtID)
		return
	}

	farm, found, err := controller.farmDAO.Find(farmID)
	if err != nil {
		web.RespondWithError(response, err)
		return
	}
	if !found {
		web.RespondWithError(response, littlefarm.ErrFarmNotFound)
		return
	}
	if err = farm.PlantCorn(littlefarm.PlantCornCmd{
		FarmID:    farmID,
		DirtID:    dirtID,
		Timestamp: time.Now(),
	}); err != nil {
		web.RespondWithError(response, err)
		return
	}
	if err = controller.farmDAO.Save(farm); err != nil {
		web.RespondWithError(response, err)
		return
	}

	web.RespondWithJSON(response, http.StatusOK, FarmToFarmOutputDTO(farm))
}

func (controller *FarmController) PlantWheat(response http.ResponseWriter, request *http.Request) {
	farmID, err := domain.NewID(request.PathValue("farm_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidFarmID)
		return
	}
	dirtID, err := domain.NewID(request.PathValue("dirt_id"))
	if err != nil {
		web.RespondWithError(response, littlefarm.ErrInvalidDirtID)
		return
	}

	farm, found, err := controller.farmDAO.Find(farmID)
	if err != nil {
		web.RespondWithError(response, err)
		return
	}
	if !found {
		web.RespondWithError(response, littlefarm.ErrFarmNotFound)
		return
	}
	if err = farm.PlantWheat(littlefarm.PlantWheatCmd{
		FarmID:    farmID,
		DirtID:    dirtID,
		Timestamp: time.Now(),
	}); err != nil {
		web.RespondWithError(response, err)
		return
	}
	if err = controller.farmDAO.Save(farm); err != nil {
		web.RespondWithError(response, err)
		return
	}

	web.RespondWithJSON(response, http.StatusOK, FarmToFarmOutputDTO(farm))
}
