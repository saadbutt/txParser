// internal/model/storage.go
package model

import (
	"strings"
)

// SaveBlock stores the block number for a given address
func (s *BlockStorage) SaveBlock(blockNum string) error {
	s.mu.Lock() // Use lock to protect write access
	defer s.mu.Unlock()
	s.currentBlock = blockNum
	return nil
}

// GetCurrentBlock retrieves the latest block number
func (s *BlockStorage) GetCurrentBlock() (string, error) {
	s.mu.RLock() // Use read lock for thread-safe reading
	defer s.mu.RUnlock()
	return s.currentBlock, nil
}

// SaveTransaction stores a transaction for an address
func (s *BlockStorage) SaveTransaction(address string, tx Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[address] = append(s.transactions[address], tx)
	return nil
}

// Get All Subscription
func (s *BlockStorage) GetAllSubscriptions() map[string]bool {
	return s.subscribers
}

// GetTransactions retrieves all transactions for a given address
func (s *BlockStorage) GetTransactions(address string) []Transaction {
	s.mu.RLock() // Read lock for concurrent reads
	defer s.mu.RUnlock()
	return s.transactions[address]
}

// Subscribe adds an address to the list of observed addresses
func (s *BlockStorage) Subscribe(address string) (bool, error) {
	if _, exists := s.subscribers[strings.ToLower(address)]; !exists {
		s.subscribers[strings.ToLower(address)] = true // Mark the address as subscribed
		return true, nil
	}

	return false, nil
}
