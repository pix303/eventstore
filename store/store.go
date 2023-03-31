package store

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/pix303/eventstore/event"
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

func withRepository(repo repository.EventRepository) EventStoreConfiguration {

	return func(store *EventStore) error {
		if repo == nil {
			return errors.New("no repo")
		}

		store.repo = repo
		return nil
	}
}

func WithBBoltRepository() EventStoreConfiguration {
	repo, err := repository.NewBBoltRepository("store.db")
	if err != nil {
		log.Fatalln(err)
	}
	return withRepository(repo)
}

func (es *EventStore) Close() {
	es.repo.Close()
}

func (es *EventStore) Add( eventType string, aggregate string, id string, payload any ) error {
	e := event.StoreEvent{
		Type:          eventType,
		Payload:       payload,
		CreatedAt:    time.Now().Format("20060102_150405.00000"),
		AggregateName: aggregate,
		AggregateID:   id,
	}

	err := es.repo.Add(e)
	if err != nil {
		return err
	}
	return nil
}

func (es *EventStore) GetAggregate(aggregate string, id string ) ([]event.StoreEvent, error) {
	return es.repo.GetByPrefix(fmt.Sprintf("%s-%s",aggregate,id))
} 

