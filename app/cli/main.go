package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wuhen781/ethereum_parser/ethparser"
)

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

func main() {
	initLogging()
	defer closeLogging()

	var (
		currentBlock bool
		subscribe    string
		transactions string
	)

	flag.BoolVar(&currentBlock, "currentBlock", false, "Get the current block number")
	flag.StringVar(&subscribe, "subscribe", "", "Subscribe to an address")
	flag.StringVar(&transactions, "transactions", "", "Get transactions for an address")
	flag.Parse()

	parser := ethparser.NewEthereumParser()

	if currentBlock {
		block := parser.GetCurrentBlock()
		fmt.Printf("Current Block: %d\n", block)
	}

	if subscribe != "" {
		if parser.Subscribe(subscribe) {
			fmt.Printf("Subscribed to address: %s\n", subscribe)
		} else {
			fmt.Printf("Already subscribed to address: %s\n", subscribe)
		}
	}

	if transactions != "" {
		txs := parser.GetTransactions(transactions)
		output, err := json.MarshalIndent(txs, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal transactions: %v", err)
		}
		fmt.Printf("Transactions for address %s:\n%s\n", transactions, string(output))
	}

	if !currentBlock && subscribe == "" && transactions == "" {
		flag.Usage()
	}
}
