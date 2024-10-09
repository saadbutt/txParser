package model

// Store is a global variable to hold the instance of the storeInterface.
var store storeInterface

// SharedStore returns the current instance of the storeInterface.
func SharedStore() storeInterface {
	return store
}

// InitializeStore initializes the global store with a concrete implementation.
func InitializeStore(s storeInterface) {
	store = s
}

// storeInterface defines the behavior required for interacting with Ethereum-related data storage.
type storeInterface interface {
	// GetCurrentBlock retrieves the number of the current Ethereum block being tracked.
	GetCurrentBlock() (string, error)

	// Subscribe adds an Ethereum address to be tracked. Returns true if the subscription is new, false if the address is already subscribed.
	Subscribe(address string) (bool, error)

	// GetTransactions retrieves the list of transactions associated with a specific address.
	GetTransactions(address string) []Transaction

	// SaveBlock persists the latest block number.
	SaveBlock(blockNumber string) error

	// GetAllSubscriptions retrieves all Ethereum addresses that are currently being tracked.
	GetAllSubscriptions() map[string]bool

	// SaveTransaction saves a transaction associated with an Ethereum address.
	SaveTransaction(address string, tx Transaction) error
}
