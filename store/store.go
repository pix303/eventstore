package store

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pix303/eventstore/event"
	"go.etcd.io/bbolt"
)

const (
	BUCKET_NAME = "event-store"
)

type EventStore struct {
	db *bbolt.DB
}

// NewStore return new instance of Eventstore
func NewStore(dbPath string) (*EventStore, error) {
	db, err := bbolt.Open(dbPath, 0666, nil)

	es := EventStore{
		db,
	}
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			return fmt.Errorf("error on create bucket: %s", err.Error())
		}
		return nil
	})

	return &es, nil
}

// BuildKey return a key for a StoreEvent
func (es *EventStore) BuildKey(event event.StoreEvent) string {
	return fmt.Sprintf("%s-%s-%s", event.AggregateName, event.AggregateID, event.CreatedAt)
}

// Close set db to close
func (es *EventStore) Close() error {
	return es.db.Close()
}

// Add adds a store event in bucket with key
func (es *EventStore) Add(event event.StoreEvent) error {
	err := es.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))

		key := es.BuildKey(event)
		value, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("error on convert in json event: %s", err.Error())
		}

		err = bucket.Put([]byte(key), value)
		if err != nil {
			return fmt.Errorf("error on store event: %s", err.Error())
		}

		return nil
	})

	return err
}

// Get return single event by key
func (es *EventStore) Get(key string) event.StoreEvent {
	var e event.StoreEvent
	_ = es.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		value := bucket.Get([]byte(key))
		err := json.Unmarshal(value, &e)
		if err != nil {
			return err
		}
		return nil
	})

	return e
}

// GetByPrefix return event group by partial key
func (es *EventStore) GetByPrefix(prefix string) []event.StoreEvent {

	var events []event.StoreEvent = make([]event.StoreEvent, 0)

	_ = es.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket([]byte(BUCKET_NAME)).Cursor()
		prefixInput := []byte(prefix)

		for k, v := cursor.Seek(prefixInput); k != nil && bytes.HasPrefix(k, prefixInput); k, v = cursor.Next() {
			var e event.StoreEvent
			err := json.Unmarshal(v, &e)
			if err != nil {
				return err
			}
			events = append(events, e)
		}

		return nil
	})

	return events
}
