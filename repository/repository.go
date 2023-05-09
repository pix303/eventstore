package repository

import "github.com/pix303/eventstore/event"

type EventRepository interface {
	Add(event event.StoreEvent) error
	GetAggregateEvents(name string, id string) ([]event.StoreEvent, error)
	GetAllAggregates(name string) ([]event.StoreEvent, error)
	Close() error
}
