package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubEvent struct {
	name string
}

func TestAggregateRoot_NewAggregateRoot(t *testing.T) {
	t.Run("It should return an aggregate root with the given ID when a valid ID is provided", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)
		assert.Equal(t, id, aggregateRoot.ID())
	})

	t.Run("It should initialize version as -1 when a new aggregate root is created", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)
		assert.Equal(t, -1, aggregateRoot.Version())
	})

	t.Run("It should initialize with no uncommitted events when a new aggregate root is created", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)
		assert.Empty(t, aggregateRoot.UncommittedEvents())
	})
}

func TestAggregateRoot_Record(t *testing.T) {
	t.Run("It should append the event to uncommitted events when an event is recorded", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)
		event := stubEvent{name: "created"}

		aggregateRoot.Record(event)

		assert.Equal(t, []Event{event}, aggregateRoot.UncommittedEvents())
	})

	t.Run("It should increment the version by one when an event is recorded", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)

		aggregateRoot.Record(stubEvent{name: "created"})

		assert.Equal(t, 0, aggregateRoot.Version())
	})

	t.Run("It should increment version for each event recorded when multiple events are recorded", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)

		aggregateRoot.Record(stubEvent{name: "created"})
		aggregateRoot.Record(stubEvent{name: "updated"})
		aggregateRoot.Record(stubEvent{name: "deleted"})

		assert.Equal(t, 2, aggregateRoot.Version())
	})

	t.Run("It should accumulate all events in order when multiple events are recorded", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)
		event1 := stubEvent{name: "created"}
		event2 := stubEvent{name: "updated"}

		aggregateRoot.Record(event1)
		aggregateRoot.Record(event2)

		assert.Equal(t, []Event{event1, event2}, aggregateRoot.UncommittedEvents())
	})
}

func TestAggregateRoot_Commit(t *testing.T) {
	t.Run("It should clear all uncommitted events when commit is called", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)
		aggregateRoot.Record(stubEvent{name: "created"})
		aggregateRoot.Record(stubEvent{name: "updated"})

		aggregateRoot.Commit()

		assert.Empty(t, aggregateRoot.UncommittedEvents())
	})

	t.Run("It should preserve the version after commit when events were previously recorded", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)
		aggregateRoot.Record(stubEvent{name: "created"})
		aggregateRoot.Record(stubEvent{name: "updated"})

		aggregateRoot.Commit()

		assert.Equal(t, 1, aggregateRoot.Version())
	})

	t.Run("It should have no effect when commit is called on a fresh aggregate root", func(t *testing.T) {
		id := GenerateID()
		aggregateRoot := NewAggregateRoot(id)

		aggregateRoot.Commit()

		assert.Equal(t, -1, aggregateRoot.Version())
		assert.Empty(t, aggregateRoot.UncommittedEvents())
	})
}
