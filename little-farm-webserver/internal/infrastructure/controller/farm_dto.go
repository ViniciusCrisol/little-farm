package controller

import (
	"little-farm-webserver/internal/domain/littlefarm"
	"time"
)

type CreateFarmInputDTO struct {
	MapWidth  int `json:"map_width"`
	MapHeight int `json:"map_height"`
}

type FarmOutputDTO struct {
	FarmID          string                    `json:"farm_id"`
	MapWidth        int                       `json:"map_width"`
	MapHeight       int                       `json:"map_height"`
	ActiveFor       int                       `json:"active_for"`
	Resources       ResourcesOutputDTO        `json:"resources"`
	ActiveElements  []ActiveElementOutputDTO  `json:"active_elements"`
	PassiveElements []PassiveElementOutputDTO `json:"passive_elements"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

type ResourcesOutputDTO struct {
	Corn         int `json:"corn"`
	Seeds        int `json:"seeds"`
	MoneyInCents int `json:"money_in_cents"`
}

type ActiveElementOutputDTO struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	XPosition int    `json:"x_position"`
	YPosition int    `json:"y_position"`
}

type PassiveElementOutputDTO struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

func FarmToFarmOutputDTO(farm littlefarm.Farm) FarmOutputDTO {
	resources := farm.Resources()

	activeElements := make([]ActiveElementOutputDTO, len(farm.ActiveElements()))
	for i, e := range farm.ActiveElements() {
		activeElements[i] = ActiveElementOutputDTO{
			ID:        e.ID().String(),
			Kind:      e.Kind(),
			XPosition: e.XPosition(),
			YPosition: e.YPosition(),
		}
	}

	passiveElements := make([]PassiveElementOutputDTO, len(farm.PassiveElements()))
	for i, e := range farm.PassiveElements() {
		passiveElements[i] = PassiveElementOutputDTO{
			ID:   e.ID().String(),
			Kind: e.Kind(),
		}
	}

	return FarmOutputDTO{
		FarmID:    farm.ID().String(),
		MapWidth:  farm.MapWidth(),
		MapHeight: farm.MapHeight(),
		ActiveFor: farm.ActiveFor(),
		Resources: ResourcesOutputDTO{
			Corn:         resources.Corn,
			Seeds:        resources.Seeds,
			MoneyInCents: resources.MoneyInCents,
		},
		ActiveElements:  activeElements,
		PassiveElements: passiveElements,
		CreatedAt:       farm.CreatedAt(),
		UpdatedAt:       farm.UpdatedAt(),
	}
}
