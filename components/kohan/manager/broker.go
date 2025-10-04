package manager

import "github.com/amanhigh/go-fun/models/tax"

// Broker represents a brokerage service that can parse transaction data.
// All broker implementations (DriveWealth, Interactive Brokers, etc.) must implement this interface.
type Broker interface {
	Parse() (tax.BrokerageInfo, error)
}
