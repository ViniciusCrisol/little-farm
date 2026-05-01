package littlefarm

import "little-farm-webserver/pkg/domain"

type Grass struct {
	id        domain.ID
	xPosition int
	yPosition int
	activeFor int
}

func NewGrass(
	id domain.ID,
	xPosition int,
	yPosition int,
) Grass {
	return Grass{
		id:        id,
		xPosition: xPosition,
		yPosition: yPosition,
	}
}

func (grass *Grass) ID() domain.ID {
	return grass.id
}

func (grass *Grass) Kind() string {
	if grass.activeFor < 3 {
		return "grass:0"
	}
	if grass.activeFor < 6 {
		return "grass:1"
	}
	if grass.activeFor < 9 {
		return "grass:2"
	}
	if grass.activeFor < 12 {
		return "grass:3"
	}
	return "grass:4"
}

func (grass *Grass) XPosition() int {
	return grass.xPosition
}

func (grass *Grass) YPosition() int {
	return grass.yPosition
}

func (grass *Grass) AdvanceOneSecond() {
	grass.activeFor++
}
