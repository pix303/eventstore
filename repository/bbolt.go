package repository

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

type BBoltRepository struct {
	db *bbolt.DB
}

// NewStore return new instance of Eventstore
func NewBBoltRepository(dbPath string) (*BBoltRepository, error) {
	db, err := bbolt.Open(dbPath, 0666, nil)

	if err != nil {
		return nil, err
	}
	
	bbr := BBoltRepository{
		db,
	}

	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			return fmt.Errorf("error on create bucket: %s", err.Error())
		}
		return nil
	})

	return &bbr, nil
}

// BuildKey return a key for a StoreEvent
func (bbr *BBoltRepository) BuildKey(event event.StoreEvent) string {
	return fmt.Sprintf("%s-%s-%s", event.AggregateName, event.AggregateID, event.CreatedAt)
}

// Close closes db
func (bbr *BBoltRepository) Close() error {
	return bbr.db.Close()
}

// Add adds an event in bucket with key
func (bbr *BBoltRepository) Add(event event.StoreEvent) error {
	err := bbr.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))

		key := bbr.BuildKey(event)
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

// Get returns a single event by key
func (bbr *BBoltRepository) Get(key string) (event.StoreEvent, error) {
	var e event.StoreEvent
	err := bbr.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		value := bucket.Get([]byte(key))
		err := json.Unmarshal(value, &e)
		if err != nil {
			return err
		}
		return nil
	})
if err != nil{
	return e, err
}
	return e, nil
}

// GetByPrefix returns an event group by partial key (typically and aggregate + ID)
func (bbr *BBoltRepository) GetByPrefix(prefix string) ([]event.StoreEvent, error) {

	var events []event.StoreEvent = make([]event.StoreEvent, 0)

	err := bbr.db.View(func(tx *bbolt.Tx) error {
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

	if err != nil{
		return nil, err
	}
	return events, nil
}
