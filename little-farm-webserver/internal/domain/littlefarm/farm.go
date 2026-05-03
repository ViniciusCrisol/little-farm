package littlefarm

import (
	"time"

	"little-farm-webserver/pkg/domain"
)

type Resources struct {
	Corn         int
	Seeds        int
	MoneyInCents int
}

type ActiveElement interface {
	ID() domain.ID
	Kind() string
	XPosition() int
	YPosition() int
	AdvanceOneSecond()
}

type PassiveElement interface {
	ID() domain.ID
	Kind() string
	AdvanceOneSecond()
}

type Farm struct {
	domain.AggregateRoot
	mapWidth        int
	mapHeight       int
	activeFor       int
	resources       Resources
	activeElements  []ActiveElement
	passiveElements []ActiveElement
	createdAt       time.Time
	updatedAt       time.Time
}

func NewFarm(cmd CreateFarmCmd) (Farm, error) {
	if cmd.MapWidth < 4 {
		return Farm{}, ErrInvalidMapWidth
	}
	if cmd.MapHeight < 4 {
		return Farm{}, ErrInvalidMapHeight
	}

	corn0_0 := NewCorn(domain.GenerateID(), 0, 0)
	corn1_0 := NewCorn(domain.GenerateID(), 1, 0)
	grass0_1 := NewGrass(domain.GenerateID(), 0, 1)
	grass1_1 := NewGrass(domain.GenerateID(), 1, 1)

	var farm Farm
	evt := FarmCreatedEvt{
		FarmID:    cmd.FarmID,
		MapWidth:  cmd.MapWidth,
		MapHeight: cmd.MapHeight,
		Resources: Resources{},
		ActiveElements: []ActiveElement{
			&corn0_0, &corn1_0,
			&grass0_1, &grass1_1,
		},
		PassiveElements: []ActiveElement{},
		Timestamp:       cmd.Timestamp,
	}
	farm.applyFarmCreatedEvt(evt)
	farm.Record(evt)
	return farm, nil
}

func (farm *Farm) HarvestCorn(cmd HarvestCornCmd) error {
	corn, found := farm.findCornByID(cmd.CornID)
	if !found {
		return ErrCornNotFound
	}
	produced, err := corn.Harvest()
	if err != nil {
		return err
	}
	evt := CornHarvestedEvt{
		FarmID:    cmd.FarmID,
		CornID:    cmd.CornID,
		Produced:  produced,
		Timestamp: cmd.Timestamp,
	}
	farm.applyCornHarvestedEvt(evt)
	farm.Record(evt)
	return nil
}

func (farm *Farm) HarvestGrass(cmd HarvestGrassCmd) error {
	grass, found := farm.findGrassByID(cmd.GrassID)
	if !found {
		return ErrGrassNotFound
	}
	produced, err := grass.Harvest()
	if err != nil {
		return err
	}
	evt := GrassHarvestedEvt{
		FarmID:    cmd.FarmID,
		GrassID:   cmd.GrassID,
		Produced:  produced,
		Timestamp: cmd.Timestamp,
	}
	farm.applyGrassHarvestedEvt(evt)
	farm.Record(evt)
	return nil
}

func (farm *Farm) AdvanceOneSecond(cmd AdvanceOneSecondCmd) error {
	evt := OneSecondAdvancedEvt{
		FarmID:    cmd.FarmID,
		Timestamp: cmd.Timestamp,
	}
	farm.applyOneSecondAdvancedEvt(evt)
	farm.Record(evt)
	return nil
}

func (farm *Farm) findCornByID(id domain.ID) (Corn, bool) {
	for _, e := range farm.activeElements {
		if e.ID().Equals(id) {
			c, isCorn := e.(*Corn)
			if !isCorn {
				return Corn{}, false
			}
			return *c, true
		}
	}
	return Corn{}, false
}

func (farm *Farm) findGrassByID(id domain.ID) (Grass, bool) {
	for _, e := range farm.activeElements {
		if e.ID().Equals(id) {
			g, isGrass := e.(*Grass)
			if !isGrass {
				return Grass{}, false
			}
			return *g, true
		}
	}
	return Grass{}, false
}

func (farm *Farm) applyFarmCreatedEvt(evt FarmCreatedEvt) {
	farm.AggregateRoot = domain.NewAggregateRoot(evt.FarmID)
	farm.mapWidth = evt.MapWidth
	farm.mapHeight = evt.MapHeight
	farm.resources = evt.Resources
	farm.activeElements = evt.ActiveElements
	farm.passiveElements = evt.PassiveElements
	farm.createdAt = evt.Timestamp
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) applyCornHarvestedEvt(evt CornHarvestedEvt) {
	for i, e := range farm.activeElements {
		if e.ID().Equals(evt.CornID) {
			farm.activeElements = append(
				farm.activeElements[:i],
				farm.activeElements[i+1:]...,
			)
			break
		}
	}
	farm.resources.Corn += evt.Produced
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) applyGrassHarvestedEvt(evt GrassHarvestedEvt) {
	for i, e := range farm.activeElements {
		if e.ID().Equals(evt.GrassID) {
			farm.activeElements = append(
				farm.activeElements[:i],
				farm.activeElements[i+1:]...,
			)
			break
		}
	}
	farm.resources.Seeds += evt.Produced
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) applyOneSecondAdvancedEvt(evt OneSecondAdvancedEvt) {
	farm.activeFor++
	for i, e := range farm.activeElements {
		e.AdvanceOneSecond()
		farm.activeElements[i] = e
	}
	for i, e := range farm.passiveElements {
		e.AdvanceOneSecond()
		farm.passiveElements[i] = e
	}
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) MapWidth() int {
	return farm.mapWidth
}

func (farm *Farm) MapHeight() int {
	return farm.mapHeight
}

func (farm *Farm) ActiveFor() int {
	return farm.activeFor
}

func (farm *Farm) Resources() Resources {
	return farm.resources
}

func (farm *Farm) ActiveElements() []ActiveElement {
	return farm.activeElements
}

func (farm *Farm) PassiveElements() []ActiveElement {
	return farm.passiveElements
}

func (farm *Farm) CreatedAt() time.Time {
	return farm.createdAt
}

func (farm *Farm) UpdatedAt() time.Time {
	return farm.updatedAt
}
