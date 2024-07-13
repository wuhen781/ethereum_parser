package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/wuhen781/ethereum_parser/ethparser"
)

var parser = ethparser.NewEthereumParser()

func main() {
	initLogging()
	defer closeLogging()

	http.HandleFunc("/currentBlock", getCurrentBlockHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/transactions", getTransactionsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var logFile *os.File

func initLogging() {
	var err error
	logFile, err = os.OpenFile("errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func closeLogging() {
	if logFile != nil {
		logFile.Close()
	}
}

func getCurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	block := parser.GetCurrentBlock()
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]int{"currentBlock": block}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}
	subscribed := parser.Subscribe(address)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"subscribed": subscribed}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}
	transactions := parser.GetTransactions(address)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
