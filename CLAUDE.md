# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Treasurenet is a Layer 1 blockchain protocol built on the Cosmos SDK that combines:
- Tendermint consensus with Ethereum Virtual Machine (EVM) compatibility
- IBC (Inter-Blockchain Communication) support
- Cross-chain bridge functionality (Gravity)
- Custom fee market implementation

## Essential Commands

### Build
```bash
make build                   # Build main binary (treasurenetd)
make build-linux            # Build for Linux
make docker-build           # Build Docker image
```

### Testing
```bash
make test                   # Run unit tests (excludes simulation)
make run-integration-tests  # Run integration tests (requires nix-shell)
```

### Linting & Formatting
```bash
make lint                   # Run golangci-lint
make lint-fix              # Run golangci-lint with auto-fix
make format                # Check code formatting
make format-fix            # Auto-format code
```

### Protobuf
```bash
make proto-gen             # Generate protobuf files
make proto-format          # Format proto files
make proto-lint            # Lint proto files
```

### Local Development
```bash
make localnet-start        # Start 4-node local testnet
make localnet-stop         # Stop local testnet
make localnet-clean        # Clean testnet data
```

## Architecture Overview

### Core Modules (`x/`)
- **`x/evm/`**: Ethereum Virtual Machine implementation for smart contract compatibility
- **`x/feemarket/`**: EIP-1559 style dynamic fee market
- **`x/gravity/`**: Cross-chain bridge for asset transfers

### Key Components
- **`app/`**: Application initialization and module wiring
- **`cmd/treasurenetd/`**: CLI commands for node operations
- **`rpc/`**: Ethereum JSON-RPC server for Web3 compatibility
- **`crypto/ethsecp256k1/`**: Ethereum-compatible key management

### Custom Dependencies
The project uses forked versions of core dependencies (see go.mod replace directives):
- Custom Tendermint fork in `tendermint_pack/`
- Modified Cosmos SDK and IBC-Go versions

## Development Notes

- Binary name: `treasurenetd`
- Main module: `github.com/treasurenetprotocol/treasurenet`
- Go version: 1.18+
- Supports multiple database backends: cleveldb, badgerdb, rocksdb, boltdb
- Ledger hardware wallet support enabled by default