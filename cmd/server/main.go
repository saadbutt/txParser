package main

import (
	"ethereum-tx-parser/internal/api"
	"ethereum-tx-parser/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize logging
	logger := log.New(os.Stdout, "ethereum-parser: ", log.LstdFlags|log.Lshortfile)

	// Initialize in-memory storage
	service.InitializeModelLayer()

	// Start the block processing service in the background
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go startBlockProcessingService(logger)

	// Start the HTTP server
	if err := api.StartServer(); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}

	// Wait for termination signal
	<-stopChan
	logger.Println("shutting down gracefully...")

}

// startBlockProcessingService runs the block fetching and transaction filtering loop
func startBlockProcessingService(logger *log.Logger) {
	for {
		currentBlockNum := service.GetBlockNumber()
		if currentBlockNum == "0x0" {
			logger.Printf("failed to get current block number")
			latestBlock, _ := service.GetLatestETHBlock()
			service.SaveLatestBlock(latestBlock)

			continue
		}

		latestBlock, err := service.GetLatestETHBlock()
		if err != nil {
			logger.Printf("failed to get latest Ethereum block: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if latestBlock >= currentBlockNum {
			block, err := service.GetEthBlockByNumber(currentBlockNum)
			if err != nil {
				logger.Printf("failed to fetch block %d: %v", currentBlockNum, err)
				time.Sleep(2 * time.Second)
				continue
			}

			if err := service.FilterTransactionsByAddress(block.Transactions); err != nil {
				logger.Printf("failed to filter transactions for block %d: %v", currentBlockNum, err)
			}

			if err := service.IncrementBlockNumber(currentBlockNum); err != nil {
				logger.Printf("failed to increment block number: %v", err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
