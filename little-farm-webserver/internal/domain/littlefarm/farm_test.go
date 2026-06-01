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

func cornIDFromFarm(t *testing.T, farm Farm) domain.ID {
	t.Helper()
	evt := farm.UncommittedEvents()[0].(FarmCreatedEvt)
	for _, e := range evt.ActiveElements {
		if _, isCorn := e.(*Corn); isCorn {
			return e.ID()
		}
	}
	t.Fatal("no corn found in farm")
	return domain.ID{}
}

func grassIDFromFarm(t *testing.T, farm Farm) domain.ID {
	t.Helper()
	evt := farm.UncommittedEvents()[0].(FarmCreatedEvt)
	for _, e := range evt.ActiveElements {
		if _, isGrass := e.(*Grass); isGrass {
			return e.ID()
		}
	}
	t.Fatal("no grass found in farm")
	return domain.ID{}
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

	t.Run("It should initialize with 100 active elements", func(t *testing.T) {
		farm := newValidFarm(t)
		evt := farm.UncommittedEvents()[0].(FarmCreatedEvt)
		assert.Equal(t, 100, len(evt.ActiveElements))
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
		cornID := cornIDFromFarm(t, farm)
		err := farm.HarvestCorn(HarvestCornCmd{
			FarmID:    farm.ID(),
			CornID:    cornID,
			Timestamp: time.Now(),
		})
		assert.ErrorIs(t, err, ErrCornNotReady)
	})

	t.Run("It should return no error and increase corn resources when corn is ready", func(t *testing.T) {
		farm := newValidFarm(t)
		cornID := cornIDFromFarm(t, farm)
		var xPosition, yPosition int
		for _, e := range farm.ActiveElements() {
			if e.ID().Equals(cornID) {
				xPosition = e.XPosition()
				yPosition = e.YPosition()
				break
			}
		}
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
		assert.Equal(t, 100, len(farm.ActiveElements()))
		var foundDirt bool
		for _, e := range farm.ActiveElements() {
			if e.ID().Equals(cornID) {
				assert.Fail(t, "harvested corn should no longer be present in active elements")
			}
			if d, isDirt := e.(*Dirt); isDirt && d.XPosition() == xPosition && d.YPosition() == yPosition {
				foundDirt = true
			}
		}
		assert.True(t, foundDirt)
	})

	t.Run("It should record a CornHarvestedEvt when corn is successfully harvested", func(t *testing.T) {
		farm := newValidFarm(t)
		cornID := cornIDFromFarm(t, farm)
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
		grassID := grassIDFromFarm(t, farm)
		err := farm.HarvestGrass(HarvestGrassCmd{
			FarmID:    farm.ID(),
			GrassID:   grassID,
			Timestamp: time.Now(),
		})
		assert.ErrorIs(t, err, ErrGrassNotReady)
	})

	t.Run("It should return no error and increase seed resources when grass is ready", func(t *testing.T) {
		farm := newValidFarm(t)
		grassID := grassIDFromFarm(t, farm)
		var xPosition, yPosition int
		for _, e := range farm.ActiveElements() {
			if e.ID().Equals(grassID) {
				xPosition = e.XPosition()
				yPosition = e.YPosition()
				break
			}
		}
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
		assert.Equal(t, 100, len(farm.ActiveElements()))
		var foundDirt bool
		for _, e := range farm.ActiveElements() {
			if e.ID().Equals(grassID) {
				assert.Fail(t, "harvested grass should no longer be present in active elements")
			}
			if d, isDirt := e.(*Dirt); isDirt && d.XPosition() == xPosition && d.YPosition() == yPosition {
				foundDirt = true
			}
		}
		assert.True(t, foundDirt)
	})

	t.Run("It should record a GrassHarvestedEvt when grass is successfully harvested", func(t *testing.T) {
		farm := newValidFarm(t)
		grassID := grassIDFromFarm(t, farm)
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

func TestFarm_MapWidth(t *testing.T) {
	t.Run("It should return the map width set during creation", func(t *testing.T) {
		farm, err := NewFarm(CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  7,
			MapHeight: 4,
			Timestamp: time.Now(),
		})
		assert.NoError(t, err)
		assert.Equal(t, 7, farm.MapWidth())
	})
}

func TestFarm_MapHeight(t *testing.T) {
	t.Run("It should return the map height set during creation", func(t *testing.T) {
		farm, err := NewFarm(CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  4,
			MapHeight: 9,
			Timestamp: time.Now(),
		})
		assert.NoError(t, err)
		assert.Equal(t, 9, farm.MapHeight())
	})
}

func TestFarm_ActiveFor(t *testing.T) {
	t.Run("It should return 0 when the farm has just been created", func(t *testing.T) {
		farm := newValidFarm(t)
		assert.Equal(t, 0, farm.ActiveFor())
	})

	t.Run("It should return the number of seconds advanced", func(t *testing.T) {
		farm := newValidFarm(t)
		now := time.Now()
		farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		assert.Equal(t, 3, farm.ActiveFor())
	})
}

func TestFarm_Resources(t *testing.T) {
	t.Run("It should return empty resources when the farm has just been created", func(t *testing.T) {
		farm := newValidFarm(t)
		assert.Equal(t, Resources{}, farm.Resources())
	})

	t.Run("It should return updated corn resources after harvesting corn", func(t *testing.T) {
		farm := newValidFarm(t)
		cornID := cornIDFromFarm(t, farm)
		now := time.Now()
		for i := 0; i < 5; i++ {
			farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: now})
		}
		farm.HarvestCorn(HarvestCornCmd{FarmID: farm.ID(), CornID: cornID, Timestamp: now})
		assert.True(t, farm.Resources().Corn > 0)
	})
}

func TestFarm_ActiveElements(t *testing.T) {
	t.Run("It should return 100 active elements when the farm has just been created", func(t *testing.T) {
		farm := newValidFarm(t)
		assert.Equal(t, 100, len(farm.ActiveElements()))
	})
}

func TestFarm_PassiveElements(t *testing.T) {
	t.Run("It should return an empty slice when the farm has just been created", func(t *testing.T) {
		farm := newValidFarm(t)
		assert.Equal(t, 0, len(farm.PassiveElements()))
	})
}

func TestFarm_CreatedAt(t *testing.T) {
	t.Run("It should return the timestamp provided during creation", func(t *testing.T) {
		now := time.Now()
		farm, err := NewFarm(CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  4,
			MapHeight: 4,
			Timestamp: now,
		})
		assert.NoError(t, err)
		assert.Equal(t, now, farm.CreatedAt())
	})
}

func TestFarm_UpdatedAt(t *testing.T) {
	t.Run("It should equal CreatedAt when the farm has just been created", func(t *testing.T) {
		now := time.Now()
		farm, err := NewFarm(CreateFarmCmd{
			FarmID:    domain.GenerateID(),
			MapWidth:  4,
			MapHeight: 4,
			Timestamp: now,
		})
		assert.NoError(t, err)
		assert.Equal(t, farm.CreatedAt(), farm.UpdatedAt())
	})

	t.Run("It should return the timestamp of the last advance when time is advanced", func(t *testing.T) {
		farm := newValidFarm(t)
		later := time.Now().Add(10 * time.Second)
		farm.AdvanceOneSecond(AdvanceOneSecondCmd{FarmID: farm.ID(), Timestamp: later})
		assert.Equal(t, later, farm.UpdatedAt())
	})
}

func TestGenerateInitialActiveElements(t *testing.T) {
	t.Run("It should return the correct total number of elements", func(t *testing.T) {
		mapWidth := 10
		mapHeight := 10
		assert.Equal(t, mapWidth*mapHeight, len(generateInitialActiveElements(mapWidth, mapHeight)))
	})

	t.Run("It should return exactly 2 Corn elements at positions (0,0) and (1,0)", func(t *testing.T) {
		cornCount := 0
		for _, e := range generateInitialActiveElements(4, 4) {
			if _, isCorn := e.(*Corn); isCorn {
				cornCount++
			}
		}
		assert.Equal(t, 2, cornCount)
	})

	t.Run("It should return exactly 2 Grass elements at positions (0,1) and (1,1)", func(t *testing.T) {
		grassCount := 0
		for _, e := range generateInitialActiveElements(4, 4) {
			if _, isGrass := e.(*Grass); isGrass {
				grassCount++
			}
		}
		assert.Equal(t, 2, grassCount)
	})

	t.Run("It should fill remaining positions with Dirt elements", func(t *testing.T) {
		mapWidth := 5
		mapHeight := 5
		dirtCount := 0
		for _, e := range generateInitialActiveElements(mapWidth, mapHeight) {
			if _, isDirt := e.(*Dirt); isDirt {
				dirtCount++
			}
		}
		expectedDirtCount := (mapWidth * mapHeight) - 4
		assert.Equal(t, expectedDirtCount, dirtCount)
	})
}
