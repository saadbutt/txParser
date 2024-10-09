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

type storeInterface interface {
	GetCurrentBlock() string
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
	SaveBlock(blockNumber string) error
	GetAllSubscriptions() map[string]bool
	SaveTransaction(string, Transaction) error
}
