package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pix303/eventstore/event"
	"github.com/pix303/localemgmtES-commons/model"
)

// PgsqlRepository wraps a postgresql db client
type PgsqlRepository struct {
	db *pgx.Conn
}

// NewPgsqlRepository create a new pgsql db client
func NewPgsqlRepository(dbURL string) (*PgsqlRepository, error) {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS eventstore (
    event_id serial NOT NULL PRIMARY KEY,
    event_type varchar(128),
    aggregate_id varchar(128),
    aggregate_name varchar(1024),
    created_at timestamp DEFAULT clock_timestamp(),
    payload text,
	metadata text,
	user_id varchar(128), 
	)
	`)

	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	db := PgsqlRepository{
		db: conn,
	}
	return &db, nil
}

// Add implements repository interface for adding an event in store
func (pgr *PgsqlRepository) Add(event event.StoreEvent) error {
	tx, err := pgr.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var payload any
	err = json.Unmarshal([]byte(event.Payload), &payload)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO public.eventstore(event_type, aggregate_id, aggregate_name, payload) VALUES ($1,$2,$3,$4);",
		event.Type, event.AggregateID, event.AggregateName, payload,
	)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// GetAggregateEvents return all events about an aggregate
func (pgr *PgsqlRepository) GetAggregateEvents(name string, id string) ([]event.StoreEvent, error) {
	return getAggregate(pgr.db, name, &id)
}

// GetAllAggregates return all events about an aggregate type
func (pgr *PgsqlRepository) GetAllAggregates(name string) ([]event.StoreEvent, error) {
	return getAggregate(pgr.db, name, nil)
}

func getAggregate(db *pgx.Conn, name string, id *string) ([]event.StoreEvent, error) {

	var events []event.StoreEvent

	var typez string
	var aggID string
	var aggName string
	var payloadValue string
	var metadataValue string
	var userID string
	var createdAt pgtype.Timestamp

	var rows pgx.Rows
	var err error

	baseSQLStatement := "SELECT event_type, aggregate_id, aggregate_name, payload, metadata, user_id, created_at FROM eventstore WHERE aggregate_id = $1"

	if id != nil {
		rows, err = db.Query(
			context.Background(),
			fmt.Sprintf("%s AND aggregate_name = $2", baseSQLStatement),
			*id, name,
		)
	} else if id == nil {
		rows, err = db.Query(
			context.Background(),
			baseSQLStatement,
			name,
		)
	}

	if err != nil {
		return []event.StoreEvent{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&typez,
			&aggID,
			&aggName,
			&payloadValue,
			&metadataValue,
			&userID,
			&createdAt,
		)
		var event = event.StoreEvent{
			Type:          typez,
			AggregateName: aggName,
			AggregateID:   aggID,
			Payload:       payloadValue,
			Metadata:      metadataValue,
			UserID:        userID,
			CreatedAt:     createdAt.Time.Format(model.DateTimeFormat),
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return []event.StoreEvent{}, err
	}

	return events, nil
}

// Close permits close db connection
func (pgr PgsqlRepository) Close() error {
	return pgr.db.Close(context.Background())
}
