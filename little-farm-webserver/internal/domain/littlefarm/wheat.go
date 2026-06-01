package littlefarm

import (
	"fmt"
	"math/rand/v2"

	"little-farm-webserver/pkg/domain"
)

type Wheat struct {
	id        domain.ID
	xPosition int
	yPosition int
	activeFor int
}

func NewWheat(
	id domain.ID,
	xPosition int,
	yPosition int,
) Wheat {
	return Wheat{
		id:        id,
		xPosition: xPosition,
		yPosition: yPosition,
	}
}

func (wheat *Wheat) ID() domain.ID {
	return wheat.id
}

func (wheat *Wheat) Kind() string {
	return fmt.Sprintf("wheat:%d", wheat.Size())
}

func (wheat *Wheat) Size() int {
	if wheat.activeFor < 5 {
		return 0
	}
	if wheat.activeFor < 10 {
		return 1
	}
	return 2
}

func (wheat *Wheat) XPosition() int {
	return wheat.xPosition
}

func (wheat *Wheat) YPosition() int {
	return wheat.yPosition
}

func (wheat *Wheat) AdvanceOneSecond() {
	wheat.activeFor++
}

func (wheat *Wheat) Harvest() (int, error) {
	size := wheat.Size()
	if size == 0 {
		return 0, ErrWheatNotReady
	}
	produced := size
	produced += rand.IntN(size) + 1
	return produced, nil
}
