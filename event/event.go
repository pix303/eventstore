package event

type StoreEvent struct {
	Type          string
	Payload       string
	CreatedAt     string
	AggregateName string
	AggregateID   string
}
