package littlefarm

import (
	"fmt"
	"math/rand/v2"

	"little-farm-webserver/pkg/domain"
)

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
	return fmt.Sprintf("grass:%d", grass.Size())
}

func (grass *Grass) Size() int {
	if grass.activeFor < 3 {
		return 0
	}
	if grass.activeFor < 6 {
		return 1
	}
	if grass.activeFor < 9 {
		return 2
	}
	if grass.activeFor < 12 {
		return 3
	}
	return 4
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

func (grass *Grass) Harvest() (int, error) {
	size := grass.Size()
	if size == 0 {
		return 0, ErrGrassNotReady
	}
	produced := size
	produced += rand.IntN(size) + 1
	return produced, nil
}
