package repository

import "github.com/pix303/eventstore/event"

type EventRepository interface {
	Add(event event.StoreEvent) error
	BuildKey(event event.StoreEvent) string
	Get(key string) (event.StoreEvent,error)
	GetByPrefix(prefix string) ( []event.StoreEvent, error )
	Close() error
}
