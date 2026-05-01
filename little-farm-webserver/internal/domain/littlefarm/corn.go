package littlefarm

import (
	"fmt"
	"math/rand/v2"

	"little-farm-webserver/pkg/domain"
)

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
	return fmt.Sprintf("corn:%d", corn.Size())
}

func (corn *Corn) Size() int {
	if corn.activeFor < 5 {
		return 0
	}
	if corn.activeFor < 10 {
		return 1
	}
	return 2
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

func (corn *Corn) Harvest() (int, error) {
	size := corn.Size()
	if size == 0 {
		return 0, ErrCornNotReady
	}
	produced := size
	produced += rand.IntN(size) + 1
	return produced, nil
}
