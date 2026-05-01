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
	corn1_0 := NewCorn(domain.GenerateID(), 0, 0)
	grass0_1 := NewGrass(domain.GenerateID(), 0, 0)
	grass1_1 := NewGrass(domain.GenerateID(), 0, 0)

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

func (farm *Farm) AdvanceOneSecond(cmd AdvanceOneSecondCmd) error {
	evt := OneSecondAdvancedEvt{
		FarmID:    cmd.FarmID,
		Timestamp: cmd.Timestamp,
	}
	farm.applyOneSecondAdvancedEvt(evt)
	farm.Record(evt)
	return nil
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
