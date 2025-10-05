package manager

import "github.com/amanhigh/go-fun/models/tax"

// Broker defines the interface for all broker parsers.
type Broker interface {
	Parse(year int) (tax.BrokerageInfo, error)
	GetName() string
}
