package event

type StoreEvent struct {
	Type          string
	Payload       any
	CreatedAt     string
	AggregateName string
	AggregateID   string
}
