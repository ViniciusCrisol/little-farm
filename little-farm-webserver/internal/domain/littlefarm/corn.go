package littlefarm

import "little-farm-webserver/pkg/domain"

type Corn struct {
	id        domain.ID
	xPosition int
	yPosition int
	activeFor int
}

func NewCorn(
	id domain.ID,
	xPosition int,
	yPosition int,
) Corn {
	return Corn{
		id:        id,
		xPosition: xPosition,
		yPosition: yPosition,
	}
}

func (corn *Corn) ID() domain.ID {
	return corn.id
}

func (corn *Corn) Kind() string {
	if corn.activeFor < 5 {
		return "corn:0"
	}
	if corn.activeFor < 10 {
		return "corn:1"
	}
	return "corn:2"
}

func (corn *Corn) XPosition() int {
	return corn.xPosition
}

func (corn *Corn) YPosition() int {
	return corn.yPosition
}

func (corn *Corn) AdvanceOneSecond() {
	corn.activeFor++
}
