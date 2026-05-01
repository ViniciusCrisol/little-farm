package littlefarm

import (
	"fmt"
	"testing"

	"little-farm-webserver/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewGrass(t *testing.T) {
	t.Run("It should return a Grass with the provided ID and positions", func(t *testing.T) {
		id := domain.GenerateID()
		grass := NewGrass(id, 4, 8)
		assert.Equal(t, id, grass.ID())
		assert.Equal(t, 4, grass.XPosition())
		assert.Equal(t, 8, grass.YPosition())
	})

	t.Run("It should return a Grass with size 0 when first created", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		assert.Equal(t, 0, grass.Size())
	})
}

func TestGrass_ID(t *testing.T) {
	t.Run("It should return the ID passed at construction", func(t *testing.T) {
		id := domain.GenerateID()
		grass := NewGrass(id, 1, 2)
		assert.Equal(t, id, grass.ID())
	})
}

func TestGrass_XPosition(t *testing.T) {
	t.Run("It should return the x position passed at construction", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 6, 9)
		assert.Equal(t, 6, grass.XPosition())
	})
}

func TestGrass_YPosition(t *testing.T) {
	t.Run("It should return the y position passed at construction", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 6, 9)
		assert.Equal(t, 9, grass.YPosition())
	})
}

func TestGrass_Size(t *testing.T) {
	t.Run("It should return 0 when activeFor is less than 3", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 2; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 0, grass.Size())
	})

	t.Run("It should return 1 when activeFor is exactly 3", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 3; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 1, grass.Size())
	})

	t.Run("It should return 1 when activeFor is between 3 and 5", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 1, grass.Size())
	})

	t.Run("It should return 2 when activeFor is exactly 6", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 6; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 2, grass.Size())
	})

	t.Run("It should return 2 when activeFor is between 6 and 8", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 8; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 2, grass.Size())
	})

	t.Run("It should return 3 when activeFor is exactly 9", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 9; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 3, grass.Size())
	})

	t.Run("It should return 3 when activeFor is between 9 and 11", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 11; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 3, grass.Size())
	})

	t.Run("It should return 4 when activeFor is exactly 12", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 12; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 4, grass.Size())
	})

	t.Run("It should return 4 when activeFor is greater than 12", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 20; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 4, grass.Size())
	})
}

func TestGrass_Kind(t *testing.T) {
	t.Run("It should return grass:0 when size is 0", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		assert.Equal(t, fmt.Sprintf("grass:%d", 0), grass.Kind())
	})

	t.Run("It should return grass:1 when size is 1", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 3; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("grass:%d", 1), grass.Kind())
	})

	t.Run("It should return grass:2 when size is 2", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 6; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("grass:%d", 2), grass.Kind())
	})

	t.Run("It should return grass:3 when size is 3", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 9; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("grass:%d", 3), grass.Kind())
	})

	t.Run("It should return grass:4 when size is 4", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 12; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("grass:%d", 4), grass.Kind())
	})
}

func TestGrass_AdvanceOneSecond(t *testing.T) {
	t.Run("It should increase the grass size after 3 calls", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		assert.Equal(t, 0, grass.Size())
		for i := 0; i < 3; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 1, grass.Size())
	})

	t.Run("It should reach size 4 after 12 calls", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 12; i++ {
			grass.AdvanceOneSecond()
		}
		assert.Equal(t, 4, grass.Size())
	})
}

func TestGrass_Harvest(t *testing.T) {
	t.Run("It should return an error when grass is not ready (size 0)", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		produced, err := grass.Harvest()
		assert.ErrorIs(t, err, ErrGrassNotReady)
		assert.Equal(t, 0, produced)
	})

	t.Run("It should return 2 when grass size is 1", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 3; i++ {
			grass.AdvanceOneSecond()
		}
		// size == 1: produced = 1 + rand.IntN(1) + 1 = 2 (rand.IntN(1) is always 0)
		produced, err := grass.Harvest()
		assert.NoError(t, err)
		assert.Equal(t, 2, produced)
	})

	t.Run("It should return a value between 3 and 4 when grass size is 2", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 6; i++ {
			grass.AdvanceOneSecond()
		}
		// size == 2: produced = 2 + rand.IntN(2) + 1 = 3 or 4
		produced, err := grass.Harvest()
		assert.NoError(t, err)
		assert.True(t, produced >= 3 && produced <= 4)
	})

	t.Run("It should return a value between 4 and 6 when grass size is 3", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 9; i++ {
			grass.AdvanceOneSecond()
		}
		// size == 3: produced = 3 + rand.IntN(3) + 1 = 4, 5, or 6
		produced, err := grass.Harvest()
		assert.NoError(t, err)
		assert.True(t, produced >= 4 && produced <= 6)
	})

	t.Run("It should return a value between 5 and 8 when grass size is 4", func(t *testing.T) {
		grass := NewGrass(domain.GenerateID(), 0, 0)
		for i := 0; i < 12; i++ {
			grass.AdvanceOneSecond()
		}
		// size == 4: produced = 4 + rand.IntN(4) + 1 = 5, 6, 7, or 8
		produced, err := grass.Harvest()
		assert.NoError(t, err)
		assert.True(t, produced >= 5 && produced <= 8)
	})
}
