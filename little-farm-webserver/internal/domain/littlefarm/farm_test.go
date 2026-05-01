package littlefarm

import (
	"testing"
	"time"

	"little-farm-webserver/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func newValidFarm(t *testing.T) Farm {
	t.Helper()
	farm, err := NewFarm(CreateFarmCmd{
		FarmID:    domain.GenerateID(),
		MapWidth:  10,
		MapHeight: 10,
		Timestamp: time.Now(),
	})
	assert.NoError(t, err)
	return farm
}

func cornIDFromFarm(farm Farm) domain.ID {
	evt := farm.UncommittedEvents()[0].(FarmCreatedEvt)
	for _, e := range evt.ActiveElements {
		if _, isCorn := e.(*Corn); isCorn {
			return e.ID()
		}
	}
	panic("no corn found in farm")
}

func grassIDFromFarm(farm Farm) domain.ID {
	evt := farm.UncommittedEvents()[0].(FarmCreatedEvt)
	for _, e := range evt.ActiveElements {
		if _, isGrass := e.(*Grass); isGrass {
			return e.ID()
		}
	}
	panic("no grass found in farm")
}

func TestNewFarm(t *testing.T) {
	t.Run("It should return a Farm when MapWidth and MapHeight are valid", func(t *testing.T) {
		farmID := domain.GenerateID()
		now := time.Now()
		farm, err := NewFarm(CreateFarmCmd{
			FarmID:    farmID,
			MapWidth:  4,
			MapHeight: 4,
			Timestamp: now,
		})
		assert.NoError(t, err)
		assert.Equal(t, farmID, farm.ID())
	})

	t.Run("It should record a FarmCreatedEvt when a farm is created", func(t *testing.T) {
		farm, err := NewFarm(CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  4,
			MapHeight: 4,
			Timestamp: time.Now(),
		})
		assert.NoError(t, err)
		events := farm.UncommittedEvents()
		assert.Equal(t, 1, len(events))
		_, isFarmCreatedEvt := events[0].(FarmCreatedEvt)
		assert.True(t, isFarmCreatedEvt)
	})

	t.Run("It should start with version 0 after creation", func(t *testing.T) {
		farm, err := NewFarm(CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  4,
			MapHeight: 4,
			Timestamp: time.Now(),
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, farm.Version())
	})

	t.Run("It should return ErrInvalidMapWidth when MapWidth is less than 4", func(t *testing.T) {
		_, err := NewFarm(CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  3,
			MapHeight: 10,
			Timestamp: time.Now(),
		})
		assert.ErrorIs(t, err, ErrInvalidMapWidth)
	})

	t.Run("It should return ErrInvalidMapHeight when MapHeight is less than 4", func(t *testing.T) {
		_, err := NewFarm(CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  10,
			MapHeight: 3,
			Timestamp: time.Now(),
		})
		assert.ErrorIs(t, err, ErrInvalidMapHeight)
	})

	t.Run("It should initialize with empty resources", func(t *testing.T) {
		farm := newValidFarm(t)
		evt := farm.UncommittedEvents()[0].(FarmCreatedEvt)
		assert.Equal(t, Resources{}, evt.Resources)
	})

	t.Run("It should initialize with 4 active elements", func(t *testing.T) {
		farm := newValidFarm(t)
		evt := farm.UncommittedEvents()[0].(FarmCreatedEvt)
		assert.Equal(t, 4, len(evt.ActiveElements))
	})
}

func TestFarm_AdvanceOneSecond(t *testing.T) {
	t.Run("It should return no error when advancing one second", func(t *testing.T) {
		farm := newValidFarm(t)
		err := farm.AdvanceOneSecond(AdvanceOneSecondCmd{
			FarmID:    farm.ID(),
			Timestamp: time.Now(),
		})
		assert.NoError(t, err)
	})

	t.Run("It should record a OneSecondAdvancedEvt when advancing one second", func(t *testing.T) {
		farm := newValidFarm(t)
		farm.Commit()
		err := farm.AdvanceOneSecond(AdvanceOneSecondCmd{
			FarmID:    farm.ID(),
			Timestamp: time.Now(),
		})
		assert.NoError(t, err)
		events := farm.UncommittedEvents()
		assert.Equal(t, 1, len(events))
		_, isOneSecondAdvancedEvt := events[0].(OneSecondAdvancedEvt)
		assert.True(t, isOneSecondAdvancedEvt)
	})

	t.Run("It should increment the version when advancing one second", func(t *testing.T) {
		farm := newValidFarm(t)
		versionBefore := farm.Version()
		farm.AdvanceOneSecond(AdvanceOneSecondCmd{
			FarmID:    farm.ID(),
			Timestamp: time.Now(),
		})
		assert.Equal(t, versionBefore+1, farm.Version())
	})
}

func TestFarm_HarvestCorn(t *testing.T) {
	t.Run("It should return ErrCornNotFound when the corn ID does not exist", func(t *testing.T) {
		farm := newValidFarm(t)
		unknownID := domain.GenerateID()
		err := farm.HarvestCorn(HarvestCornCmd{
			FarmID:    farm.ID(),
			CornID:    unknownID,
			Timestamp: time.Now(),
		})
		assert.ErrorIs(t, err, ErrCornNotFound)
	})

	t.Run("It should return ErrCornNotReady when the corn has not grown enough", func(t *testing.T) {
		farm := newValidFarm(t)
		cornID := cornIDFromFarm(farm)
		err := farm.HarvestCorn(HarvestCornCmd{
			FarmID:    farm.ID(),
			CornID:    cornID,
			Timestamp: time.Now(),
		})
		assert.ErrorIs(t, err, ErrCornNotReady)
	})

	t.Run("It should return no error and increase corn resources when corn is ready", func(t *testing.T) {
		farm := newValidFarm(t)
		cornID := cornIDFromFarm(farm)
		now := time.Now()
		for i := 0; i < 5; i++ {
			farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		}
		err := farm.HarvestCorn(HarvestCornCmd{
			FarmID:    farm.ID(),
			CornID:    cornID,
			Timestamp: now,
		})
		assert.NoError(t, err)
	})

	t.Run("It should record a CornHarvestedEvt when corn is successfully harvested", func(t *testing.T) {
		farm := newValidFarm(t)
		cornID := cornIDFromFarm(farm)
		now := time.Now()
		for i := 0; i < 5; i++ {
			farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		}
		farm.Commit()
		farm.HarvestCorn(HarvestCornCmd{
			FarmID:    farm.ID(),
			CornID:    cornID,
			Timestamp: now,
		})
		events := farm.UncommittedEvents()
		assert.Equal(t, 1, len(events))
		_, isCornHarvestedEvt := events[0].(CornHarvestedEvt)
		assert.True(t, isCornHarvestedEvt)
	})
}

func TestFarm_HarvestGrass(t *testing.T) {
	t.Run("It should return ErrGrassNotFound when the grass ID does not exist", func(t *testing.T) {
		farm := newValidFarm(t)
		unknownID := domain.GenerateID()
		err := farm.HarvestGrass(HarvestGrassCmd{
			FarmID:    farm.ID(),
			GrassID:   unknownID,
			Timestamp: time.Now(),
		})
		assert.ErrorIs(t, err, ErrGrassNotFound)
	})

	t.Run("It should return ErrGrassNotReady when the grass has not grown enough", func(t *testing.T) {
		farm := newValidFarm(t)
		grassID := grassIDFromFarm(farm)
		err := farm.HarvestGrass(HarvestGrassCmd{
			FarmID:    farm.ID(),
			GrassID:   grassID,
			Timestamp: time.Now(),
		})
		assert.ErrorIs(t, err, ErrGrassNotReady)
	})

	t.Run("It should return no error and increase seed resources when grass is ready", func(t *testing.T) {
		farm := newValidFarm(t)
		grassID := grassIDFromFarm(farm)
		now := time.Now()
		for i := 0; i < 3; i++ {
			farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		}
		err := farm.HarvestGrass(HarvestGrassCmd{
			FarmID:    farm.ID(),
			GrassID:   grassID,
			Timestamp: now,
		})
		assert.NoError(t, err)
	})

	t.Run("It should record a GrassHarvestedEvt when grass is successfully harvested", func(t *testing.T) {
		farm := newValidFarm(t)
		grassID := grassIDFromFarm(farm)
		now := time.Now()
		for i := 0; i < 3; i++ {
			farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		}
		farm.Commit()
		farm.HarvestGrass(HarvestGrassCmd{
			FarmID:    farm.ID(),
			GrassID:   grassID,
			Timestamp: now,
		})
		events := farm.UncommittedEvents()
		assert.Equal(t, 1, len(events))
		_, isGrassHarvestedEvt := events[0].(GrassHarvestedEvt)
		assert.True(t, isGrassHarvestedEvt)
	})
}
