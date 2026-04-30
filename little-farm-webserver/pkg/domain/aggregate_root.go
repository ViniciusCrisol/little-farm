package domain

type Event any

type AggregateRoot struct {
	id                ID
	version           int
	uncommittedEvents []Event
}

func NewAggregateRoot(id ID) AggregateRoot {
	return AggregateRoot{
		id:                id,
		version:           -1,
		uncommittedEvents: []Event{},
	}
}

func (aggregate *AggregateRoot) Commit() {
	aggregate.uncommittedEvents = []Event{}
}

func (aggregate *AggregateRoot) Record(event Event) {
	aggregate.version++
	aggregate.uncommittedEvents = append(aggregate.uncommittedEvents, event)
}

func (aggregate *AggregateRoot) IncrementVersion() {
	aggregate.version++
}

func (aggregate *AggregateRoot) ID() ID {
	return aggregate.id
}

func (aggregate *AggregateRoot) Version() int {
	return aggregate.version
}

func (aggregate *AggregateRoot) UncommittedEvents() []Event {
	return aggregate.uncommittedEvents
}
