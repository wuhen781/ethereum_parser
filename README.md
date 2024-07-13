# Ethereum Blockchain Parser

This project implements an Ethereum blockchain parser in Go. The parser interacts with the Ethereum blockchain using the JSON-RPC API and provides functionality to:
- Get the current block number.
- Subscribe to Ethereum addresses.
- Get transactions for a subscribed address.

## Project Structure
├── ethparser/
│ ├── ethparser.go
│ ├── ethparser_test.go
│ └── transaction.go
└── app/cli/main.go
└── app/api/main.go


### `app/cli/main.go`
Contains the HTTP server implementation that exposes the parser functionality through REST endpoints.

### `app/cli/main.go`
Provides a command-line interface for interacting with the parser.

### `ethparser/ethparser.go`
Contains the core logic for interacting with the Ethereum blockchain.

### `ethparser/ethparser_test.go`
Contains tests for the core logic in `ethparser.go`.

### `ethparser/transaction.go`
Defines the `Transaction` struct used to represent Ethereum transactions.


## Getting Started

### Prerequisites

- Go 1.16 or later
- Access to the Ethereum JSON-RPC endpoint (e.g., https://cloudflare-eth.com)

### Installation

1. Clone the repository:
```sh
git clone https://github.com/wuhen781/ethereum_parser.git
cd ethereum_parser

2. Build main.go:
go build -o ethparser_cli app/cli/main.go
go build -o ethparser_api app/api/main.go

3. Run ethparse_cli
# To get the current block number
./ethparser_cli -currentBlock

# To subscribe to an address
./ethparser_cli -subscribe=0xYourAddress

# To get transactions for an address
./ethparser_cli -transactions=0xYourAddress

4. Run ethparser_api 
# start the api server
./ethparser_api  

# To get the current block number
curl 127.0.0.1:8080/currentBlock

# To subscribe to an address
curl 127.0.0.1:8080/subscribe?address=0xyouraddress

# To get transactions for an address
curl 127.0.0.1:8080/transactions?address=0xyouraddress
