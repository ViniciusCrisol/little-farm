package littlefarm

import (
	"errors"
	"little-farm-webserver/pkg/domain"
	"time"
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
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type PassiveElement interface {
	ID() domain.ID
	Kind() string
	AdvanceOneSecond()
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type Farm struct {
	domain.AggregateRoot
	mapWidth        int
	mapHeight       int
	resources       Resources
	activeElements  []ActiveElement
	passiveElements []ActiveElement
	createdAt       time.Time
	updatedAt       time.Time
}

func NewFarm(cmd CreateFarmCmd) (Farm, error) {
	if cmd.MapWidth < 4 {
		return Farm{}, errors.New("validation err: invalid map width (move it to apperr)")
	}
	if cmd.MapHeight < 4 {
		return Farm{}, errors.New("validation err: invalid map height (move it to apperr)")
	}

	corn0_0 := NewCorn(domain.GenerateID(), 0, 0, cmd.Timestamp)
	corn1_0 := NewCorn(domain.GenerateID(), 0, 0, cmd.Timestamp)
	grass0_1 := NewGrass(domain.GenerateID(), 0, 0, cmd.Timestamp)
	grass1_1 := NewGrass(domain.GenerateID(), 0, 0, cmd.Timestamp)

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
