# Ethereum Transaction Parser

Implement Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

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

## Installation

To get started, clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/txParser.git
cd txParser
```
# Ethereum Transaction Parser

## Prerequisites

- **Go** (version 1.16 or later)

## Usage

### Running the Application

1. Navigate to the root directory of the project.
2. Run the application:

 ```bash
go run cmd/server/main.go
run ./txParser
  ```

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
  ```
