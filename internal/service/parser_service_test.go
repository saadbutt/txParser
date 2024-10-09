package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ethereum-tx-parser/internal/model"
	"ethereum-tx-parser/internal/parser"
)

// Mocking the SharedStore for testing
type MockStore struct {
	currentBlock  string
	subscriptions map[string]bool
	transactions  []model.Transaction
}

func (m *MockStore) GetCurrentBlock() (string, error) {
	return m.currentBlock, nil
}

func (m *MockStore) SaveBlock(blockNumber string) error {
	m.currentBlock = blockNumber
	return nil // Default behavior
}

func (m *MockStore) Subscribe(address string) (bool, error) {
	m.subscriptions[address] = true
	return true, nil
}

func (m *MockStore) GetAllSubscriptions() map[string]bool {
	return m.subscriptions
}

func (m *MockStore) SaveTransaction(address string, tx model.Transaction) error {
	return nil
}

func (m *MockStore) GetTransactions(address string) []model.Transaction {
	return m.transactions
}

// Mocking HTTP Client for RPC calls
func mockEthBlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	response := parser.RpcResponse{
		JsonRPC: "2.0",
		ID:      1,
		Result:  "0x10d4f", // Example block number in hex
	}
	json.NewEncoder(w).Encode(response)
}

func mockGetBlockByNumberHandler(w http.ResponseWriter, r *http.Request) {
	response := model.JsonRPCResponse{
		JsonRPC: "2.0",
		ID:      1,
		Result: model.Block{
			Number: "0x10d4f",
			// Add additional block data if necessary
		},
	}
	json.NewEncoder(w).Encode(response)
}

func TestGetEthBlockByNumber(t *testing.T) {
	// Setting up mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(mockGetBlockByNumberHandler))
	defer server.Close()

	block, err := GetEthBlockByNumber("0x10d4f")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if block.Number != "0x10d4f" {
		t.Errorf("expected block number: 0x10d4f, got: %s", block.Number)
	}
}

func TestFilterTransactionsByAddress(t *testing.T) {
	mockStore := &MockStore{
		subscriptions: make(map[string]bool),
	}
	model.InitializeStore(mockStore)

	// Subscribe to an address
	mockStore.Subscribe("0x1234567890abcdef1234567890abcdef12345678")

	transactions := []model.Transaction{
		{From: "0x1234567890abcdef1234567890abcdef12345678", To: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
		{From: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", To: "0x1234567890abcdef1234567890abcdef12345678"},
	}

	err := FilterTransactionsByAddress(transactions)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(mockStore.GetAllSubscriptions()) != 1 {
		t.Errorf("expected 1 subscription, got: %d", len(mockStore.GetAllSubscriptions()))
	}
}

func TestIncrementBlockNumber(t *testing.T) {
	mockStore := &MockStore{
		subscriptions: make(map[string]bool),
	}
	model.InitializeStore(mockStore)

	mockStore.SaveBlock("0x10d4f")

	err := IncrementBlockNumber("0x10d4f")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	blockNum, _ := mockStore.GetCurrentBlock()
	if blockNum != "0x10d50" {
		t.Errorf("expected incremented block number: 0x10d50, got: %s", blockNum)
	}
}

func TestSubscribe(t *testing.T) {
	mockStore := &MockStore{
		subscriptions: make(map[string]bool),
	}
	model.InitializeStore(mockStore)

	tests := []struct {
		name          string
		address       string
		expectedError bool
	}{
		{"Valid Ethereum Address", "0x1234567890abcdef1234567890abcdef12345678", false},
		{"Invalid Ethereum Address", "invalid_address", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Subscribe(tt.address)

			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
