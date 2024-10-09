package parser

import "ethereum-tx-parser/internal/model"

// Parser defines the public interface for the Ethereum blockchain parser
// This interface can be implemented for use via command line or HTTP API.
type Parser interface {
	// GetCurrentBlock retrieves the last parsed block number
	GetCurrentBlock() (int, error)

	// Subscribe adds an address to be observed for transactions
	Subscribe(address string) (bool, error)

	// GetTransactions retrieves the list of transactions for a specific address
	GetTransactions(address string) ([]model.Transaction, error)
}

// RpcRequest represents the structure of an Ethereum JSON-RPC request
type RpcRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// RpcResponse represents the structure of an Ethereum JSON-RPC response
type RpcResponse struct {
	ID      int    `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}
