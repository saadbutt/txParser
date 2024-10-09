// internal/model/storage.go
package model

import "sync"

// EthereumRPC holds the RPC URL for connecting to the Ethereum network
type EthereumRPC struct {
	rpcURL string
}

// JSON-RPC response for getting a block
type JsonRPCResponse struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Block  `json:"result"`
}

// Transaction represents an Ethereum transaction
type Transaction struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

// Block represents the structure of an Ethereum block
type Block struct {
	Number       string
	Transactions []Transaction
}

// BlockStorage is an in-memory storage for block-related data
type BlockStorage struct {
	mu           sync.RWMutex
	currentBlock string          // Latest block number processed by the listener
	subscribers  map[string]bool // Subscribed addresses
	transactions map[string][]Transaction
}

// NewBlockStorage initializes an in-memory block storage
func NewBlockStorage() *BlockStorage {
	return &BlockStorage{
		transactions: make(map[string][]Transaction),
		subscribers:  make(map[string]bool),
		currentBlock: "0x0", // Initialize to 0 or any appropriate value
	}
}

// NewEthereumRPC initializes an EthereumRPC client with the given URL
func NewEthereumRPC(rpcURL string) *EthereumRPC {
	return &EthereumRPC{
		rpcURL: rpcURL,
	}
}
