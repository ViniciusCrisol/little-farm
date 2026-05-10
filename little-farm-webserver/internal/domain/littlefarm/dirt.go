package littlefarm

import "little-farm-webserver/pkg/domain"

type Dirt struct {
	id        domain.ID
	xPosition int
	yPosition int
	activeFor int
}

func NewDirt(
	id domain.ID,
	xPosition int,
	yPosition int,
) Dirt {
	return Dirt{
		id:        id,
		xPosition: xPosition,
		yPosition: yPosition,
	}
}

func (dirt *Dirt) ID() domain.ID {
	return dirt.id
}

func (dirt *Dirt) Kind() string {
	return "dirt:0"
}

func (dirt *Dirt) XPosition() int {
	return dirt.xPosition
}

func (dirt *Dirt) YPosition() int {
	return dirt.yPosition
}

func (dirt *Dirt) AdvanceOneSecond() {
	dirt.activeFor++
}
