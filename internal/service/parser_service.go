package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"ethereum-tx-parser/internal/model"
	"ethereum-tx-parser/internal/parser"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	EthereumRPCURL = "https://ethereum-rpc.publicnode.com"
)

// InitializeModelLayer sets up the block storage model
func InitializeModelLayer() {
	blockStorage := model.NewBlockStorage()
	model.InitializeStore(blockStorage)
}

// GetBlockNumber retrieves the current block number from the store
func GetBlockNumber() (string, error) {
	return model.SharedStore().GetCurrentBlock()
}

// SaveLatestBlock saves the latest block number to the store with error handling
func SaveLatestBlock(currentBlockNum string) error {
	if err := model.SharedStore().SaveBlock(currentBlockNum); err != nil {
		log.Printf("Error saving block: %v", err)
		return errors.New("error updating blocks")
	}
	return nil
}

// GetLatestETHBlock retrieves the latest Ethereum block number via RPC
func GetLatestETHBlock() (string, error) {
	requestPayload := parser.RpcRequest{
		JsonRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      1,
	}

	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshaling request payload: %v", err)
		return "0x0", err
	}

	resp, err := http.Post(EthereumRPCURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error making RPC request: %v", err)
		return "0x0", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading RPC response: %v", err)
		return "0x0", err
	}

	var rpcResp parser.RpcResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		log.Printf("Error unmarshaling RPC response: %v", err)
		return "0x0", err
	}

	return rpcResp.Result, nil
}

// GetEthBlockByNumber retrieves a block by its number via RPC
func GetEthBlockByNumber(blockNumber string) (model.Block, error) {
	if !strings.HasPrefix(blockNumber, "0x") {
		return model.Block{}, fmt.Errorf("block number must be in hex format, e.g., '0x4b7'")
	}

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockNumber, true},
		"id":      1,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request payload: %v", err)
		return model.Block{}, err
	}

	resp, err := http.Post(EthereumRPCURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error making RPC request: %v", err)
		return model.Block{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return model.Block{}, err
	}

	var rpcResp model.JsonRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		log.Printf("Error unmarshaling block response: %v", err)
		return model.Block{}, err
	}

	return rpcResp.Result, nil
}

// FilterTransactionsByAddress filters transactions for subscribed addresses
func FilterTransactionsByAddress(transactions []model.Transaction) error {
	addressMap := model.SharedStore().GetAllSubscriptions()
	if len(addressMap) == 0 {
		// Map is empty, handle the empty case here
		fmt.Println("No subscriptions found.")
		return nil
	}

	for _, tx := range transactions {
		lowerFrom := strings.ToLower(tx.From)
		lowerTo := strings.ToLower(tx.To)

		// Check if From or To address is subscribed
		if addressMap[lowerFrom] {
			if err := model.SharedStore().SaveTransaction(lowerFrom, tx); err != nil {
				log.Printf("Error saving transaction for address %s: %v", lowerFrom, err)
				return err
			}
		}

		if addressMap[lowerTo] && lowerFrom != lowerTo {
			if err := model.SharedStore().SaveTransaction(lowerTo, tx); err != nil {
				log.Printf("Error saving transaction for address %s: %v", lowerTo, err)
				return err
			}
		}
	}

	return nil
}

// IncrementBlockNumber increments and stores the block number
func IncrementBlockNumber(blockHex string) error {
	blockHex = strings.TrimPrefix(blockHex, "0x")
	blockNumber, err := strconv.ParseInt(blockHex, 16, 64)
	if err != nil {
		log.Printf("Error parsing hex block number: %v", err)
		return errors.New("error parsing block number")
	}

	blockNumber++
	newBlockHex := fmt.Sprintf("0x%x", blockNumber)
	if err := SaveLatestBlock(newBlockHex); err != nil {
		log.Printf("Error saving incremented block number: %v", err)
		return err
	}

	log.Printf("Successfully incremented block number to: %d %s", blockNumber, newBlockHex)
	return nil
}

func isValidEthereumAddress(address string) bool {
	return strings.HasPrefix(address, "0x") && len(address) == 42
}

// Subscribe adds an address to the list of observed addresses
func Subscribe(address string) (bool, error) {
	if !isValidEthereumAddress(address) {
		fmt.Println("Invalid Ethereum address")
		return false, errors.New("invalid Ethereum address")
	}
	return model.SharedStore().Subscribe(address)
}
