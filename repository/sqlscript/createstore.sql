CREATE TABLE IF NOT EXISTS eventstore (
    event_id serial NOT NULL PRIMARY KEY,
    event_type varchar(128),
    aggregate_id UUID,
    aggregate_name varchar(1024),
    created_at timestamp DEFAULT clock_timestamp(),
    payload json
)
