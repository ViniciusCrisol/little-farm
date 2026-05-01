package littlefarm

import (
	"fmt"
	"testing"

	"little-farm-webserver/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewCorn(t *testing.T) {
	t.Run("It should return a Corn with the provided ID and positions", func(t *testing.T) {
		id := domain.GenerateID()
		corn := NewCorn(id, 3, 7)
		assert.Equal(t, id, corn.ID())
		assert.Equal(t, 3, corn.XPosition())
		assert.Equal(t, 7, corn.YPosition())
	})

	t.Run("It should return a Corn with size 0 when first created", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		assert.Equal(t, 0, corn.Size())
	})
}

func TestCorn_ID(t *testing.T) {
	t.Run("It should return the ID passed at construction", func(t *testing.T) {
		id := domain.GenerateID()
		corn := NewCorn(id, 1, 2)
		assert.Equal(t, id, corn.ID())
	})
}

func TestCorn_XPosition(t *testing.T) {
	t.Run("It should return the x position passed at construction", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 5, 9)
		assert.Equal(t, 5, corn.XPosition())
	})
}

func TestCorn_YPosition(t *testing.T) {
	t.Run("It should return the y position passed at construction", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 5, 9)
		assert.Equal(t, 9, corn.YPosition())
	})
}

func TestCorn_Size(t *testing.T) {
	t.Run("It should return 0 when activeFor is less than 5", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 4; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, 0, corn.Size())
	})

	t.Run("It should return 1 when activeFor is exactly 5", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, 1, corn.Size())
	})

	t.Run("It should return 1 when activeFor is between 5 and 9", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 9; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, 1, corn.Size())
	})

	t.Run("It should return 2 when activeFor is exactly 10", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, 2, corn.Size())
	})

	t.Run("It should return 2 when activeFor is greater than 10", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 20; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, 2, corn.Size())
	})
}

func TestCorn_Kind(t *testing.T) {
	t.Run("It should return corn:0 when size is 0", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		assert.Equal(t, fmt.Sprintf("corn:%d", 0), corn.Kind())
	})

	t.Run("It should return corn:1 when size is 1", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("corn:%d", 1), corn.Kind())
	})

	t.Run("It should return corn:2 when size is 2", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("corn:%d", 2), corn.Kind())
	})
}

func TestCorn_AdvanceOneSecond(t *testing.T) {
	t.Run("It should increase the corn size after 5 calls", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		assert.Equal(t, 0, corn.Size())
		for i := 0; i < 5; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, 1, corn.Size())
	})

	t.Run("It should increase the corn size again after 10 calls", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			corn.AdvanceOneSecond()
		}
		assert.Equal(t, 2, corn.Size())
	})
}

func TestCorn_Harvest(t *testing.T) {
	t.Run("It should return an error when corn is not ready (size 0)", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		produced, err := corn.Harvest()
		assert.ErrorIs(t, err, ErrCornNotReady)
		assert.Equal(t, 0, produced)
	})

	t.Run("It should return 2 when corn size is 1", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			corn.AdvanceOneSecond()
		}
		// size == 1: produced = 1 + rand.IntN(1) + 1 = 2 (rand.IntN(1) is always 0)
		produced, err := corn.Harvest()
		assert.NoError(t, err)
		assert.Equal(t, 2, produced)
	})

	t.Run("It should return a value between 3 and 4 when corn size is 2", func(t *testing.T) {
		corn := NewCorn(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			corn.AdvanceOneSecond()
		}
		// size == 2: produced = 2 + rand.IntN(2) + 1 = 3 or 4
		produced, err := corn.Harvest()
		assert.NoError(t, err)
		assert.True(t, produced >= 3 && produced <= 4)
	})
}
