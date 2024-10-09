package parser

import "ethereum-tx-parser/internal/model"

// Transaction represents a basic transaction structure

// Parser is the main interface for the Ethereum blockchain parser
type Parser interface {
	GetCurrentBlock() int                               // Retrieve the last parsed block
	Subscribe(address string) bool                      // Add an address to be observed
	GetTransactions(address string) []model.Transaction // Get the list of transactions for an address
}

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RpcResponse struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}
