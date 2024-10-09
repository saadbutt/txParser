package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockStorage(t *testing.T) {
	storage := NewBlockStorage()

	resp, _ := storage.GetCurrentBlock()
	// Test initial current block
	assert.Equal(t, "0x0", resp, "Initial block should be 0x0")

	// Test saving a block
	err := storage.SaveBlock("0x1")
	assert.NoError(t, err, "Error should be nil when saving block")
	resp, _ = storage.GetCurrentBlock()
	assert.Equal(t, "0x1", resp, "Current block should be updated to 0x1")

	// Test subscribing to an address
	address := "0x1234567890abcdef1234567890abcdef12345678"
	subscription, _ := storage.Subscribe(address)
	assert.True(t, subscription, "Should return true for new subscription")
	assert.Equal(t, 1, len(storage.GetAllSubscriptions()), "Should have one subscription")

	// Test subscribing to the same address again
	subscription, _ = storage.Subscribe(address)
	assert.False(t, subscription, "Should return false for existing subscription")
	assert.Equal(t, 1, len(storage.GetAllSubscriptions()), "Should still have one subscription")

	// Test retrieving transactions for an address that has no transactions
	transactions := storage.GetTransactions(address)
	assert.Empty(t, transactions, "Should have no transactions for a new address")

	// Test saving a transaction
	tx := Transaction{
		Hash:  "0xabc",
		From:  address,
		To:    "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		Value: "1000000000000000000", // 1 ETH in Wei
	}
	err = storage.SaveTransaction(address, tx)
	assert.NoError(t, err, "Error should be nil when saving transaction")

	// Test retrieving transactions for an address with transactions
	transactions = storage.GetTransactions(address)
	assert.Equal(t, 1, len(transactions), "Should have one transaction")
	assert.Equal(t, tx, transactions[0], "The retrieved transaction should match the saved transaction")
}

func TestMultipleTransactions(t *testing.T) {
	storage := NewBlockStorage()
	address := "0x1234567890abcdef1234567890abcdef12345678"

	// Save multiple transactions
	for i := 0; i < 5; i++ {
		tx := Transaction{
			Hash:  "0xabc" + fmt.Sprint(i),
			From:  address,
			To:    "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			Value: "1000000000000000000", // 1 ETH in Wei
		}
		err := storage.SaveTransaction(address, tx)
		assert.NoError(t, err, "Error should be nil when saving transaction")
	}

	// Test retrieving transactions for the address
	transactions := storage.GetTransactions(address)
	assert.Equal(t, 5, len(transactions), "Should have five transactions")
}
