package littlefarm

import (
	"fmt"
	"math/rand/v2"

	"little-farm-webserver/pkg/domain"
)

type Pumpkin struct {
	id        domain.ID
	xPosition int
	yPosition int
	activeFor int
}

func NewPumpkin(
	id domain.ID,
	xPosition int,
	yPosition int,
) Pumpkin {
	return Pumpkin{
		id:        id,
		xPosition: xPosition,
		yPosition: yPosition,
	}
}

func (pumpkin *Pumpkin) ID() domain.ID {
	return pumpkin.id
}

func (pumpkin *Pumpkin) Kind() string {
	return fmt.Sprintf("pumpkin:%d", pumpkin.Size())
}

func (pumpkin *Pumpkin) Size() int {
	if pumpkin.activeFor < 5 {
		return 0
	}
	if pumpkin.activeFor < 10 {
		return 1
	}
	return 2
}

func (pumpkin *Pumpkin) XPosition() int {
	return pumpkin.xPosition
}

func (pumpkin *Pumpkin) YPosition() int {
	return pumpkin.yPosition
}

func (pumpkin *Pumpkin) AdvanceOneSecond() {
	pumpkin.activeFor++
}

func (pumpkin *Pumpkin) Harvest() (int, error) {
	size := pumpkin.Size()
	if size == 0 {
		return 0, ErrPumpkinNotReady
	}
	produced := size
	produced += rand.IntN(size) + 1
	return produced, nil
}
