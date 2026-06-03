package littlefarm

import (
	"fmt"
	"slices"
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

	activeElements := generateInitialActiveElements(cmd.MapWidth, cmd.MapHeight)

	var farm Farm
	evt := FarmCreatedEvt{
		FarmID:          cmd.FarmID,
		MapWidth:        cmd.MapWidth,
		MapHeight:       cmd.MapHeight,
		Resources:       Resources{},
		ActiveElements:  activeElements,
		PassiveElements: []ActiveElement{},
		Timestamp:       cmd.Timestamp,
	}
	farm.applyFarmCreatedEvt(evt)
	farm.Record(evt)
	return farm, nil
}

func generateInitialActiveElements(mapWidth, mapHeight int) []ActiveElement {
	cornPositions := []string{
		fmt.Sprintf("%d_%d", 0, 0),
		fmt.Sprintf("%d_%d", 1, 0),
	}
	grassPositions := []string{
		fmt.Sprintf("%d_%d", 0, 1),
		fmt.Sprintf("%d_%d", 1, 1),
	}
	var activeElements []ActiveElement
	for x := range mapWidth {
		for y := range mapHeight {
			position := fmt.Sprintf("%d_%d", x, y)
			if slices.Contains(cornPositions, position) {
				c := NewCorn(domain.GenerateID(), x, y)
				activeElements = append(activeElements, &c)
				continue
			}
			if slices.Contains(grassPositions, position) {
				g := NewGrass(domain.GenerateID(), x, y)
				activeElements = append(activeElements, &g)
				continue
			}
			d := NewDirt(domain.GenerateID(), x, y)
			activeElements = append(activeElements, &d)
		}
	}
	return activeElements
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

func (farm *Farm) HarvestPumpkin(cmd HarvestPumpkinCmd) error {
	pumpkin, found := farm.findPumpkinByID(cmd.PumpkinID)
	if !found {
		return ErrPumpkinNotFound
	}
	produced, err := pumpkin.Harvest()
	if err != nil {
		return err
	}
	evt := PumpkinHarvestedEvt{
		FarmID:    cmd.FarmID,
		PumpkinID: cmd.PumpkinID,
		Produced:  produced,
		Timestamp: cmd.Timestamp,
	}
	farm.applyPumpkinHarvestedEvt(evt)
	farm.Record(evt)
	return nil
}

func (farm *Farm) PlantCorn(cmd PlantCornCmd) error {
	_, found := farm.findDirtByID(cmd.DirtID)
	if !found {
		return ErrDirtNotFound
	}
	if farm.resources.Corn < 1 {
		return ErrNotEnoughCornResources
	}

	evt := CornPlantedEvt{
		CornID:    domain.GenerateID(),
		FarmID:    cmd.FarmID,
		DirtID:    cmd.DirtID,
		Timestamp: cmd.Timestamp,
	}
	farm.Record(evt)
	farm.applyCornPlantedEvt(evt)
	return nil
}

func (farm *Farm) PlantWheat(cmd PlantWheatCmd) error {
	_, found := farm.findDirtByID(cmd.DirtID)
	if !found {
		return ErrDirtNotFound
	}
	if farm.resources.Seeds < 1 {
		return ErrNotEnoughSeedResources
	}

	evt := WheatPlantedEvt{
		WheatID:   domain.GenerateID(),
		FarmID:    cmd.FarmID,
		DirtID:    cmd.DirtID,
		Timestamp: cmd.Timestamp,
	}
	farm.Record(evt)
	farm.applyWheatPlantedEvt(evt)
	return nil
}

func (farm *Farm) PlantPumpkin(cmd PlantPumpkinCmd) error {
	_, found := farm.findDirtByID(cmd.DirtID)
	if !found {
		return ErrDirtNotFound
	}
	if farm.resources.Corn < 1 {
		return ErrNotEnoughCornResources
	}

	evt := PumpkinPlantedEvt{
		PumpkinID: domain.GenerateID(),
		FarmID:    cmd.FarmID,
		DirtID:    cmd.DirtID,
		Timestamp: cmd.Timestamp,
	}
	farm.Record(evt)
	farm.applyPumpkinPlantedEvt(evt)
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

func (farm *Farm) findDirtByID(id domain.ID) (Dirt, bool) {
	for _, e := range farm.activeElements {
		if e.ID().Equals(id) {
			d, isDirt := e.(*Dirt)
			if !isDirt {
				return Dirt{}, false
			}
			return *d, true
		}
	}
	return Dirt{}, false
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

func (farm *Farm) findPumpkinByID(id domain.ID) (Pumpkin, bool) {
	for _, e := range farm.activeElements {
		if e.ID().Equals(id) {
			p, isPumpkin := e.(*Pumpkin)
			if !isPumpkin {
				return Pumpkin{}, false
			}
			return *p, true
		}
	}
	return Pumpkin{}, false
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
			farm.activeElements = append(farm.activeElements[:i], farm.activeElements[i+1:]...)
			d := NewDirt(domain.GenerateID(), e.XPosition(), e.YPosition())
			farm.activeElements = append(farm.activeElements, &d)
			break
		}
	}
	farm.resources.Corn += evt.Produced
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) applyGrassHarvestedEvt(evt GrassHarvestedEvt) {
	for i, e := range farm.activeElements {
		if e.ID().Equals(evt.GrassID) {
			farm.activeElements = append(farm.activeElements[:i], farm.activeElements[i+1:]...)
			d := NewDirt(domain.GenerateID(), e.XPosition(), e.YPosition())
			farm.activeElements = append(farm.activeElements, &d)
			break
		}
	}
	farm.resources.Seeds += evt.Produced
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) applyPumpkinHarvestedEvt(evt PumpkinHarvestedEvt) {
	for i, e := range farm.activeElements {
		if e.ID().Equals(evt.PumpkinID) {
			farm.activeElements = append(farm.activeElements[:i], farm.activeElements[i+1:]...)
			d := NewDirt(domain.GenerateID(), e.XPosition(), e.YPosition())
			farm.activeElements = append(farm.activeElements, &d)
			break
		}
	}
	farm.resources.Corn += evt.Produced
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) applyCornPlantedEvt(evt CornPlantedEvt) {
	for i, e := range farm.activeElements {
		if e.ID().Equals(evt.DirtID) {
			d, isDirt := e.(*Dirt)
			if isDirt {
				farm.activeElements = append(farm.activeElements[:i], farm.activeElements[i+1:]...)
				c := NewCorn(evt.CornID, d.XPosition(), d.YPosition())
				farm.activeElements = append(farm.activeElements, &c)
				break
			}
		}
	}
	farm.resources.Corn--
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) applyWheatPlantedEvt(evt WheatPlantedEvt) {
	for i, e := range farm.activeElements {
		if e.ID().Equals(evt.DirtID) {
			d, isDirt := e.(*Dirt)
			if isDirt {
				farm.activeElements = append(farm.activeElements[:i], farm.activeElements[i+1:]...)
				w := NewWheat(evt.WheatID, d.XPosition(), d.YPosition())
				farm.activeElements = append(farm.activeElements, &w)
				break
			}
		}
	}
	farm.resources.Seeds--
	farm.updatedAt = evt.Timestamp
}

func (farm *Farm) applyPumpkinPlantedEvt(evt PumpkinPlantedEvt) {
	for i, e := range farm.activeElements {
		if e.ID().Equals(evt.DirtID) {
			d, isDirt := e.(*Dirt)
			if isDirt {
				farm.activeElements = append(farm.activeElements[:i], farm.activeElements[i+1:]...)
				p := NewPumpkin(evt.PumpkinID, d.XPosition(), d.YPosition())
				farm.activeElements = append(farm.activeElements, &p)
				break
			}
		}
	}
	farm.resources.Corn--
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
