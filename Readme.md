# Ethereum Transaction Parser

This project is an Ethereum transaction parser designed to monitor and analyze Ethereum blockchain transactions. It interacts with the Ethereum network to fetch block data and process transactions based on specific criteria.

## Features

- Fetch the latest Ethereum block number.
- Retrieve detailed information about a specific Ethereum block.
- Filter transactions by address.
- Save and manage blockchain data using a mock store for testing purposes.
- Support for subscription management for Ethereum addresses.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Running Tests](#running-tests)
- [Contributing](#contributing)
- [License](#license)

## Installation

To get started, clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/txParser.git
cd txParser

# Ethereum Transaction Parser

## Prerequisites

- **Go** (version 1.16 or later)
- **Access to an Ethereum node** (can be local or through services like Infura or Alchemy)

## Configuration

1. Set up your Ethereum RPC URL in the code or as an environment variable. You can find various Ethereum nodes to connect to [here](https://eth.wiki/en/Nodes).
2. Update the `EthereumRPCURL` variable in the code with your node's endpoint.

## Usage

### Running the Application

1. Navigate to the root directory of the project.
2. Run the application:

   ```bash
   go run cmd/server/main.go or run ./txParser
## Available Functions

- **GetCurrentBlock**: Retrieves the last parsed block number.
- **Subscribe**: Subscribes an Ethereum address for transaction updates.
- **GetTransactions**: Filters transactions based on a specified Ethereum address.



## Running Tests

To ensure everything is working correctly, run the tests included in the project:

1. Reset the Go module cache (optional):

   ```bash
   go clean -modcache
   go test ./...
