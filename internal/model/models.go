// internal/model/storage.go
package model

import "sync"

type EthereumRPC struct {
	rpcURL string
}

// JSON-RPC response for getting a block
type JsonRPCResponse struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Block  `json:"result"`
}

type Transaction struct {
	Hash  string
	From  string
	To    string
	Value string
}

// Struct for the block result from eth_getBlockByNumber
type Block struct {
	Number       string
	Transactions []Transaction
}

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
