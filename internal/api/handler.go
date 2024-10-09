package api

import (
	"encoding/json"
	"ethereum-tx-parser/internal/model"
	"ethereum-tx-parser/internal/service"
	"log"
	"net/http"
	"strings"
)

// StartServer initializes and starts the HTTP server
func StartServer() error {
	http.HandleFunc("/currentBlock", CurrentBlockHandler)
	http.HandleFunc("/subscribe", SaveSubscriptionHandler)
	http.HandleFunc("/transactions", ListTransactionsHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}
	return nil
}

// CurrentBlockHandler returns the current block number
func CurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	setJSONResponseHeaders(w)
	response, _ := model.SharedStore().GetCurrentBlock()

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("error encoding current block response: %v", err)
	}
}

// SaveSubscriptionHandler handles subscription requests for an Ethereum address
func SaveSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	setJSONResponseHeaders(w)

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	subscribed, err := service.Subscribe(address)

	status := "Already Subscribed"
	if subscribed {
		status = "Subscribed"
	}
	var response map[string]string
	if err != nil {
		response = map[string]string{"status": err.Error(), "address": address}

	} else {
		response = map[string]string{"status": status, "address": address}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("error encoding subscription response for address %s: %v", address, err)
	}
}

// ListTransactionsHandler returns a list of transactions for a given Ethereum address
func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	setJSONResponseHeaders(w)

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	transactions := model.SharedStore().GetTransactions(strings.ToLower(address))
	// Check if transactions are empty
	if len(transactions) == 0 {
		// Return an empty array or a message, depending on your preference
		response := map[string]interface{}{
			"address":      address,
			"transactions": []string{}, // Returning an empty list to indicate no transactions
			"message":      "No transactions found for this address",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("error encoding empty transactions response for address %s: %v", address, err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("error encoding transactions for address %s: %v", address, err)
	}
}

// setJSONResponseHeaders sets common headers for JSON responses
func setJSONResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allows cross-origin requests
}
