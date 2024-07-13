package ethparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	// endpoint = "https://cloudflare-eth.com"
	endpoint = "https://public-en-cypress.klaytn.net"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type EthereumParser struct {
	currentBlock int
	addresses    map[string]bool
	mu           sync.Mutex
}

func NewEthereumParser() *EthereumParser {
	return &EthereumParser{
		addresses: make(map[string]bool),
	}
}

func (p *EthereumParser) GetCurrentBlock() int {
	response, err := p.callRPC("eth_blockNumber", []interface{}{})
	if err != nil {
		log.Printf("Error getting current block: %v", err)
		return 0
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Printf("Error unmarshaling current block response: %v", err)
		return 0
	}

	blockNumberHex := result["result"].(string)
	blockNumber := HexToInt(blockNumberHex)
	p.currentBlock = blockNumber

	return blockNumber
}

func (p *EthereumParser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.addresses[address]; exists {
		return false
	}
	p.addresses[address] = true
	return true
}

func (p *EthereumParser) GetTransactions(address string) []Transaction {
	var transactions []Transaction

	latestBlock := p.GetCurrentBlock()
	lastGetBlock := latestBlock - 5

	for blockNumber := latestBlock; blockNumber > lastGetBlock; blockNumber-- {
		response, err := p.callRPC("eth_getBlockByNumber", []interface{}{IntToHex(blockNumber), true})
		if err != nil {
			log.Printf("Error getting block by number: %v", err)
			continue
		}

		var blockResult map[string]interface{}
		if err := json.Unmarshal(response, &blockResult); err != nil {
			log.Printf("Error unmarshaling block response: %v", err)
			continue
		}

		block := blockResult["result"].(map[string]interface{})
		txs := block["transactions"].([]interface{})
		for _, tx := range txs {
			txMap := tx.(map[string]interface{})
			if txMap["from"].(string) == address || txMap["to"].(string) == address {
				transactions = append(transactions, Transaction{
					From:        txMap["from"].(string),
					To:          txMap["to"].(string),
					Value:       txMap["value"].(string),
					BlockNumber: HexToInt(block["number"].(string)),
					Gas:         HexToInt(txMap["gas"].(string)),
					GasPrice:    txMap["gasPrice"].(string),
					Hash:        txMap["hash"].(string),
					Nonce:       HexToInt(txMap["nonce"].(string)),
					Timestamp:   int64(HexToInt(block["timestamp"].(string))),
				})
			}
		}
	}

	return transactions
}

func (p *EthereumParser) callRPC(method string, params []interface{}) ([]byte, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(endpoint, "application/json", strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func HexToInt(hexStr string) int {
	var i int
	fmt.Sscanf(hexStr, "0x%x", &i)
	return i
}

func IntToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}
