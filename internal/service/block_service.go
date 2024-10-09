package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"ethereum-tx-parser/internal/model"
	"ethereum-tx-parser/internal/parser"
	"strconv"
	"strings"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func InitializeModelLayer() {
	blockStorage := model.NewBlockStorage()
	model.InitializeStore(blockStorage)
}

func GetBlockNumber() string {
	currentBlock := model.SharedStore().GetCurrentBlock()

	return currentBlock
}

func SaveLatestBlock(currentBlockNum string) error {
	if err := model.SharedStore().SaveBlock(currentBlockNum); err != nil {
		fmt.Println(err)
		return errors.New("error updating blocks")
	}
	return nil
}

func GetLatestETHBlock() (string, error) {
	// Create the request payload
	requestPayload := parser.RpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      1,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshaling request payload: %v", err)
		return "0x0", err
	}

	// Send the request to the Ethereum RPC
	resp, err := http.Post("https://ethereum-rpc.publicnode.com", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error making RPC request: %v", err)
		return "0x0", err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading RPC response: %v", err)
		return "0x0", err
	}

	// Parse the response
	var rpcResp parser.RpcResponse
	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		log.Printf("Error unmarshaling RPC response: %v", err)
		return "0x0", err
	}

	return rpcResp.Result, nil
}

// Function to make a JSON-RPC request to an Ethereum node
func GetEthBlockByNumber(blockNumber string) (model.Block, error) {
	// Ensure blockNumber is in hex format
	if !strings.HasPrefix(blockNumber, "0x") {
		return model.Block{}, fmt.Errorf("block number must be in hex format, e.g., '0x4b7'")
	}

	// Create the JSON-RPC request
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockNumber, true}, // "true" to get full transaction objects
		"id":      1,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return model.Block{}, err
	}

	// Send the request to the Ethereum node
	resp, err := http.Post("https://ethereum-rpc.publicnode.com", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.Block{}, err
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Block{}, err
	}

	var rpcResp model.JsonRPCResponse
	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		return model.Block{}, err
	}

	// Return the block with transactions
	return rpcResp.Result, nil
}

// Function to filter transactions by address
func FilterTransactionsByAddress(transactions []model.Transaction) error {

	// Create a map to store lowercase versions of the addresses for fast lookup
	addressMap := model.SharedStore().GetAllSubscriptions()

	for _, tx := range transactions {

		lowerFrom := strings.ToLower(tx.From)
		lowerTo := strings.ToLower(tx.To)

		// Check if the From or To address exists in the address map
		if addressMap[lowerFrom] {
			fmt.Println("Save Tx:", lowerFrom, tx)
			return model.SharedStore().SaveTransaction(lowerFrom, tx)
		}

		if addressMap[lowerTo] && lowerFrom != lowerTo {
			fmt.Println("Save Tx:", lowerTo, tx)
			return model.SharedStore().SaveTransaction(lowerTo, tx)
		}
	}

	return nil
}

func IncrementBlockNumber(blockHex string) error {
	// Remove the "0x" prefix if it exists
	blockHex = strings.TrimPrefix(blockHex, "0x")

	// Convert the hex string to an integer (base 16)
	blockNumber, err := strconv.ParseInt(blockHex, 16, 64)
	if err != nil {
		fmt.Println("Error parsing hex:", err)
		return errors.New("error parsing")
	}

	// Increment the block number
	blockNumber++

	fmt.Println("BLOCK number:", blockNumber)

	// Convert the new block number back to hex, and add the "0x" prefix
	newBlockHex := fmt.Sprintf("0x%x", blockNumber)

	return SaveLatestBlock(newBlockHex)

}

// Subscribe adds an address to the list of observed addresses
func Subscribe(address string) bool {
	return model.SharedStore().Subscribe(address) // Use BlockStorage's Subscribe method
}
