package littlefarm

import (
	"testing"

	"little-farm-webserver/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewDirt(t *testing.T) {
	t.Run("It should return a Dirt with the provided ID and positions", func(t *testing.T) {
		id := domain.GenerateID()
		dirt := NewDirt(id, 2, 5)
		assert.Equal(t, id, dirt.ID())
		assert.Equal(t, 2, dirt.XPosition())
		assert.Equal(t, 5, dirt.YPosition())
	})
}

func TestDirt_ID(t *testing.T) {
	t.Run("It should return the ID passed at construction", func(t *testing.T) {
		id := domain.GenerateID()
		dirt := NewDirt(id, 0, 0)
		assert.Equal(t, id, dirt.ID())
	})
}

func TestDirt_XPosition(t *testing.T) {
	t.Run("It should return the x position passed at construction", func(t *testing.T) {
		dirt := NewDirt(domain.GenerateID(), 7, 3)
		assert.Equal(t, 7, dirt.XPosition())
	})
}

func TestDirt_YPosition(t *testing.T) {
	t.Run("It should return the y position passed at construction", func(t *testing.T) {
		dirt := NewDirt(domain.GenerateID(), 7, 3)
		assert.Equal(t, 3, dirt.YPosition())
	})
}

func TestDirt_Kind(t *testing.T) {
	t.Run("It should return dirt:0", func(t *testing.T) {
		dirt := NewDirt(domain.GenerateID(), 0, 0)
		assert.Equal(t, "dirt:0", dirt.Kind())
	})
}
