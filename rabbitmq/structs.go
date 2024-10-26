package rabbitmq

// Queue
type Queue struct {
	// Name
	Name string
	// Exchange
	Exchange string
	// Bindings
	Bindings []string
	// Connections
	Connections int
	// Channels
	Channels int
	// Pre Fetch
	PreFetch int
}
