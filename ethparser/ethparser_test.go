package ethparser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewEthereumParser(t *testing.T) {
	parser := NewEthereumParser()
	if parser == nil {
		t.Fatal("Expected parser to be non-nil")
	}
	if parser.addresses == nil {
		t.Fatal("Expected addresses map to be initialized")
	}
}

func TestEthereumParser_GetCurrentBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  "0x4b7",
		})
	}))
	defer server.Close()

	parser := NewEthereumParser()
	endpoint = server.URL
	block := parser.GetCurrentBlock()
	expectedBlock := 1207
	if block != expectedBlock {
		t.Errorf("Expected block %d, got %d", expectedBlock, block)
	}
}

func TestEthereumParser_Subscribe(t *testing.T) {
	parser := NewEthereumParser()

	address := "0x123"
	subscribed := parser.Subscribe(address)
	if !subscribed {
		t.Errorf("Expected subscribe to return true for a new address")
	}

	subscribed = parser.Subscribe(address)
	if subscribed {
		t.Errorf("Expected subscribe to return false for an already subscribed address")
	}
}

func TestEthereumParser_GetTransactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 67 { // GetCurrentBlock()
			json.NewEncoder(w).Encode(map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"result":  "0x4b7",
			})
		} else { //GetTransactions()
			json.NewEncoder(w).Encode(map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"result": map[string]interface{}{
					"number":    "0x4b7",
					"timestamp": "0x5c0",
					"transactions": []interface{}{
						map[string]interface{}{
							"from":     "0x123",
							"to":       "0x456",
							"value":    "0x1",
							"gas":      "0x5208",
							"gasPrice": "0x3b9aca00",
							"hash":     "0x789",
							"nonce":    "0x0",
						},
					},
				},
			})
		}
	}))
	defer server.Close()

	parser := NewEthereumParser()
	endpoint = server.URL
	transactions := parser.GetTransactions("0x123")

	if len(transactions) < 1 {
		t.Fatalf("Expected at least 1 transaction, got %d", len(transactions))
	}

	tx := transactions[0]
	if tx.From != "0x123" {
		t.Errorf("Expected From to be '0x123', got %s", tx.From)
	}
	if tx.To != "0x456" {
		t.Errorf("Expected To to be '0x456', got %s", tx.To)
	}
	if tx.Value != "0x1" {
		t.Errorf("Expected Value to be '0x1', got %s", tx.Value)
	}
	if tx.Gas != 21000 {
		t.Errorf("Expected Gas to be 21000, got %d", tx.Gas)
	}
	if tx.GasPrice != "0x3b9aca00" {
		t.Errorf("Expected GasPrice to be '0x3b9aca00', got %s", tx.GasPrice)
	}
	if tx.Hash != "0x789" {
		t.Errorf("Expected Hash to be '0x789', got %s", tx.Hash)
	}
	if tx.Nonce != 0 {
		t.Errorf("Expected Nonce to be 0, got %d", tx.Nonce)
	}
	if tx.BlockNumber != 1207 {
		t.Errorf("Expected BlockNumber to be 1207, got %d", tx.BlockNumber)
	}
	if tx.Timestamp != 1472 {
		t.Errorf("Expected Timestamp to be 1472, got %d", tx.Timestamp)
	}
}

func TestHexToInt(t *testing.T) {
	hexStr := "0x4b7"
	expected := 1207
	result := HexToInt(hexStr)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestIntToHex(t *testing.T) {
	num := 1207
	expected := "0x4b7"
	result := IntToHex(num)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
