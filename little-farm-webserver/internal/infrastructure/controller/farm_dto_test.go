package controller

import (
	"testing"
	"time"

	"little-farm-webserver/internal/domain/littlefarm"
	"little-farm-webserver/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func newValidFarm(t *testing.T) littlefarm.Farm {
	t.Helper()
	farm, err := littlefarm.NewFarm(littlefarm.CreateFarmCmd{
		FarmID:    domain.GenerateID(),
		MapWidth:  10,
		MapHeight: 10,
		Timestamp: time.Now(),
	})
	assert.NoError(t, err)
	return farm
}

func TestFarmToFarmOutputDTO(t *testing.T) {
	t.Run("It should map the farm ID to the DTO", func(t *testing.T) {
		farm := newValidFarm(t)
		dto := FarmToFarmOutputDTO(farm)
		assert.Equal(t, farm.ID().String(), dto.FarmID)
	})

	t.Run("It should map MapWidth to the DTO", func(t *testing.T) {
		farm, err := littlefarm.NewFarm(littlefarm.CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  7,
			MapHeight: 5,
			Timestamp: time.Now(),
		})
		assert.NoError(t, err)
		dto := FarmToFarmOutputDTO(farm)
		assert.Equal(t, 7, dto.MapWidth)
	})

	t.Run("It should map MapHeight to the DTO", func(t *testing.T) {
		farm, err := littlefarm.NewFarm(littlefarm.CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  5,
			MapHeight: 8,
			Timestamp: time.Now(),
		})
		assert.NoError(t, err)
		dto := FarmToFarmOutputDTO(farm)
		assert.Equal(t, 8, dto.MapHeight)
	})

	t.Run("It should map ActiveFor to the DTO", func(t *testing.T) {
		farm := newValidFarm(t)
		now := time.Now()
		farm.AdvanceOneSecond(littlefarm.AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		farm.AdvanceOneSecond(littlefarm.AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		dto := FarmToFarmOutputDTO(farm)
		assert.Equal(t, 2, dto.ActiveFor)
	})

	t.Run("It should map Resources to the DTO", func(t *testing.T) {
		farm := newValidFarm(t)
		dto := FarmToFarmOutputDTO(farm)
		assert.Equal(t, ResourcesOutputDTO{Corn: 0, Seeds: 0, MoneyInCents: 0}, dto.Resources)
	})

	t.Run("It should map ActiveElements with correct ID, Kind, XPosition and YPosition to the DTO", func(t *testing.T) {
		farm := newValidFarm(t)
		dto := FarmToFarmOutputDTO(farm)
		elements := farm.ActiveElements()
		assert.Equal(t, len(elements), len(dto.ActiveElements))
		for i, e := range elements {
			assert.Equal(t, e.ID().String(), dto.ActiveElements[i].ID)
			assert.Equal(t, e.Kind(), dto.ActiveElements[i].Kind)
			assert.Equal(t, e.XPosition(), dto.ActiveElements[i].XPosition)
			assert.Equal(t, e.YPosition(), dto.ActiveElements[i].YPosition)
		}
	})

	t.Run("It should map PassiveElements to an empty slice in the DTO", func(t *testing.T) {
		farm := newValidFarm(t)
		dto := FarmToFarmOutputDTO(farm)
		assert.Equal(t, 0, len(dto.PassiveElements))
	})

	t.Run("It should map CreatedAt to the DTO", func(t *testing.T) {
		now := time.Now()
		farm, err := littlefarm.NewFarm(littlefarm.CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  4,
			MapHeight: 4,
			Timestamp: now,
		})
		assert.NoError(t, err)
		dto := FarmToFarmOutputDTO(farm)
		assert.Equal(t, now, dto.CreatedAt)
	})

	t.Run("It should map UpdatedAt to the DTO", func(t *testing.T) {
		now := time.Now()
		farm, err := littlefarm.NewFarm(littlefarm.CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  4,
			MapHeight: 4,
			Timestamp: now,
		})
		assert.NoError(t, err)
		later := now.Add(5 * time.Second)
		farm.AdvanceOneSecond(littlefarm.AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: later})
		dto := FarmToFarmOutputDTO(farm)
		assert.Equal(t, later, dto.UpdatedAt)
	})
}
