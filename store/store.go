package store

import (
	"errors"

	"github.com/pix303/eventstore/repository"
)

type EventStore struct {
	repo repository.EventRepository
}

type EventStoreConfiguration func(store *EventStore) error

// NewStore return new instance of Eventstore
func NewStore(config EventStoreConfiguration) (*EventStore, error) {
	es := &EventStore{}
	err := config(es)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func withRepository(repo repository.EventRepository, err error) EventStoreConfiguration {

	return func(store *EventStore) error {
		if err != nil {
			return err
		}

		if repo == nil {
			return errors.New("no repo")
		}

		store.repo = repo
		return nil
	}
}

func WithBBoltRepository() EventStoreConfiguration {
	return withRepository(repository.NewBBoltRepository("db"))
}

func (es *EventStore) Close() {
	es.repo.Closed()
}
