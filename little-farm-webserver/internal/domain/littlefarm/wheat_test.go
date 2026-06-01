package littlefarm

import (
	"fmt"
	"testing"

	"little-farm-webserver/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewWheat(t *testing.T) {
	t.Run("It should return a Wheat with the provided ID and positions", func(t *testing.T) {
		id := domain.GenerateID()
		wheat := NewWheat(id, 3, 7)
		assert.Equal(t, id, wheat.ID())
		assert.Equal(t, 3, wheat.XPosition())
		assert.Equal(t, 7, wheat.YPosition())
	})

	t.Run("It should return a Wheat with size 0 when first created", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		assert.Equal(t, 0, wheat.Size())
	})
}

func TestWheat_ID(t *testing.T) {
	t.Run("It should return the ID passed at construction", func(t *testing.T) {
		id := domain.GenerateID()
		wheat := NewWheat(id, 1, 2)
		assert.Equal(t, id, wheat.ID())
	})
}

func TestWheat_XPosition(t *testing.T) {
	t.Run("It should return the x position passed at construction", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 5, 9)
		assert.Equal(t, 5, wheat.XPosition())
	})
}

func TestWheat_YPosition(t *testing.T) {
	t.Run("It should return the y position passed at construction", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 5, 9)
		assert.Equal(t, 9, wheat.YPosition())
	})
}

func TestWheat_Size(t *testing.T) {
	t.Run("It should return 0 when activeFor is less than 5", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 4; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, 0, wheat.Size())
	})

	t.Run("It should return 1 when activeFor is exactly 5", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, 1, wheat.Size())
	})

	t.Run("It should return 1 when activeFor is between 5 and 9", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 9; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, 1, wheat.Size())
	})

	t.Run("It should return 2 when activeFor is exactly 10", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, 2, wheat.Size())
	})

	t.Run("It should return 2 when activeFor is greater than 10", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 20; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, 2, wheat.Size())
	})
}

func TestWheat_Kind(t *testing.T) {
	t.Run("It should return wheat:0 when size is 0", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		assert.Equal(t, fmt.Sprintf("wheat:%d", 0), wheat.Kind())
	})

	t.Run("It should return wheat:1 when size is 1", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("wheat:%d", 1), wheat.Kind())
	})

	t.Run("It should return wheat:2 when size is 2", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, fmt.Sprintf("wheat:%d", 2), wheat.Kind())
	})
}

func TestWheat_AdvanceOneSecond(t *testing.T) {
	t.Run("It should increase the wheat size after 5 calls", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		assert.Equal(t, 0, wheat.Size())
		for i := 0; i < 5; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, 1, wheat.Size())
	})

	t.Run("It should increase the wheat size again after 10 calls", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			wheat.AdvanceOneSecond()
		}
		assert.Equal(t, 2, wheat.Size())
	})
}

func TestWheat_Harvest(t *testing.T) {
	t.Run("It should return an error when wheat is not ready (size 0)", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		produced, err := wheat.Harvest()
		assert.ErrorIs(t, err, ErrWheatNotReady)
		assert.Equal(t, 0, produced)
	})

	t.Run("It should return 2 when wheat size is 1", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 5; i++ {
			wheat.AdvanceOneSecond()
		}
		// size == 1: produced = 1 + rand.IntN(1) + 1 = 2 (rand.IntN(1) is always 0)
		produced, err := wheat.Harvest()
		assert.NoError(t, err)
		assert.Equal(t, 2, produced)
	})

	t.Run("It should return a value between 3 and 4 when wheat size is 2", func(t *testing.T) {
		wheat := NewWheat(domain.GenerateID(), 0, 0)
		for i := 0; i < 10; i++ {
			wheat.AdvanceOneSecond()
		}
		// size == 2: produced = 2 + rand.IntN(2) + 1 = 3 or 4
		produced, err := wheat.Harvest()
		assert.NoError(t, err)
		assert.True(t, produced >= 3 && produced <= 4)
	})
}
