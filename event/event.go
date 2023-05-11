package event

// StoreEvent represents single event that changes aggregates
type StoreEvent struct {
	Type    string
	Payload string
	Metadata      string
	UserID        string
	CreatedAt     string
	AggregateName string
	AggregateID   string
}
