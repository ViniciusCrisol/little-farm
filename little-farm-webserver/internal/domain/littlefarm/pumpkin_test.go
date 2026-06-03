package littlefarm

import (
	"fmt"
	"testing"

	"little-farm-webserver/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewPumpkin(t *testing.T) {
	t.Run("It should return a Pumpkin with the provided ID and positions", func(t *testing.T) {
		id := domain.GenerateID()
		pumpkin := NewPumpkin(id, 3, 7)
		assert.Equal(t, id, pumpkin.ID())
		assert.Equal(t, 3, pumpkin.XPosition())
		assert.Equal(t, 7, pumpkin.YPosition())
	})

	t.Run("It should return a Pumpkin with size 0 when first created", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		assert.Equal(t, 0, pumpkin.Size())
	})
}

func TestPumpkin_ID(t *testing.T) {
	t.Run("It should return the ID passed at construction", func(t *testing.T) {
		id := domain.GenerateID()
		pumpkin := NewPumpkin(id, 1, 2)
		assert.Equal(t, id, pumpkin.ID())
	})
}

func TestPumpkin_XPosition(t *testing.T) {
	t.Run("It should return the x position passed at construction", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 5, 9)
		assert.Equal(t, 5, pumpkin.XPosition())
	})
}

func TestPumpkin_YPosition(t *testing.T) {
	t.Run("It should return the y position passed at construction", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 5, 9)
		assert.Equal(t, 9, pumpkin.YPosition())
	})
}

func TestPumpkin_Size(t *testing.T) {
	t.Run("It should return 0 when activeFor is less than 5", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 4; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, 0, pumpkin.Size())
	})

	t.Run("It should return 1 when activeFor is exactly 5", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, 1, pumpkin.Size())
	})

	t.Run("It should return 1 when activeFor is between 5 and 9", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 9; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, 1, pumpkin.Size())
	})

	t.Run("It should return 2 when activeFor is exactly 10", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, 2, pumpkin.Size())
	})

	t.Run("It should return 2 when activeFor is greater than 10", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 20; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, 2, pumpkin.Size())
	})
}

func TestPumpkin_Kind(t *testing.T) {
	t.Run("It should return pumpkin:0 when size is 0", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		assert.Equal(t, fmt.Sprintf("pumpkin:%d", 0), pumpkin.Kind())
	})

	t.Run("It should return pumpkin:1 when size is 1", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("pumpkin:%d", 1), pumpkin.Kind())
	})

	t.Run("It should return pumpkin:2 when size is 2", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("pumpkin:%d", 2), pumpkin.Kind())
	})
}

func TestPumpkin_AdvanceOneSecond(t *testing.T) {
	t.Run("It should increase the pumpkin size after 5 calls", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		assert.Equal(t, 0, pumpkin.Size())
		for i := 0; i < 5; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, 1, pumpkin.Size())
	})

	t.Run("It should increase the pumpkin size again after 10 calls", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			pumpkin.AdvanceOneSecond()
		}
		assert.Equal(t, 2, pumpkin.Size())
	})
}

func TestPumpkin_Harvest(t *testing.T) {
	t.Run("It should return an error when pumpkin is not ready (size 0)", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		produced, err := pumpkin.Harvest()
		assert.ErrorIs(t, err, ErrPumpkinNotReady)
		assert.Equal(t, 0, produced)
	})

	t.Run("It should return 2 when pumpkin size is 1", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			pumpkin.AdvanceOneSecond()
		}
		// size == 1: produced = 1 + rand.IntN(1) + 1 = 2 (rand.IntN(1) is always 0)
		produced, err := pumpkin.Harvest()
		assert.NoError(t, err)
		assert.Equal(t, 2, produced)
	})

	t.Run("It should return a value between 3 and 4 when pumpkin size is 2", func(t *testing.T) {
		pumpkin := NewPumpkin(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			pumpkin.AdvanceOneSecond()
		}
		// size == 2: produced = 2 + rand.IntN(2) + 1 = 3 or 4
		produced, err := pumpkin.Harvest()
		assert.NoError(t, err)
		assert.True(t, produced >= 3 && produced <= 4)
	})
}
