<p align="center">
  <a href="https://treasurenet.io">
    <img alt="Treasurenet Logo" src="https://raw.githubusercontent.com/treasurenetprotocol/docs/feature/1.0.3/static/img/logo_tn_github.png" width="350" />
  </a>
</p>

<h1 align="center">Treasurenet Protocol</h1>

<p align="center">
  <strong>A Layer 1 blockchain protocol bridging real-world economic value with distributed ledger technology</strong>
</p>

<p align="center">
  <a href="https://github.com/treasurenetprotocol/treasurenet/blob/main/LICENSE">
    <img alt="License: LGPL v3" src="https://img.shields.io/badge/License-LGPL%20v3-blue.svg" />
  </a>
  <a href="https://golang.org/dl/">
    <img alt="Go Version" src="https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go" />
  </a>
  <a href="https://github.com/treasurenetprotocol/treasurenet/releases">
    <img alt="Version" src="https://img.shields.io/github/v/release/treasurenetprotocol/treasurenet?include_prereleases&logo=github" />
  </a>
  <a href="https://discord.com/invite/treasurenet">
    <img alt="Discord" src="https://img.shields.io/discord/treasurenet?logo=discord&logoColor=white" />
  </a>
  <a href="https://twitter.com/treasurenet_io">
    <img alt="X (formerly Twitter) Follow" src="https://img.shields.io/twitter/follow/treasurenet_io?logo=x&style=social" />
  </a>
</p>

<p align="center">
  <a href="#-overview">Overview</a> •
  <a href="#-key-features">Features</a> •
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-documentation">Documentation</a> •
  <a href="#-community">Community</a> •
  <a href="#-contributing">Contributing</a>
</p>

---

## 📋 Overview

Treasurenet Protocol is a groundbreaking Layer 1 blockchain solution designed to address the significant shortfall of enduring and substantial value within the cryptocurrency realm. By merging real-world economic forces with the scalability of distributed ledger technology, Treasurenet aims to set a new precedent for maintaining value across both fiat and digital domains.

### 🎯 Mission

Our mission is to create a sustainable blockchain ecosystem that:
- Bridges traditional finance with decentralized technology
- Provides lasting value preservation mechanisms
- Enables seamless cross-chain interoperability
- Supports real-world asset tokenization

## ✨ Key Features

### 🏗️ Technical Architecture

- **Consensus Mechanism**: Advanced proof-of-stake (PoS) consensus for energy efficiency and security
- **EVM Compatibility**: Full Ethereum Virtual Machine support for smart contract deployment
- **Cross-Chain Communication**: Built-in IBC (Inter-Blockchain Communication) protocol support
- **High Performance**: Capable of processing thousands of transactions per second
- **Low Latency**: Sub-second block finality for optimal user experience

### 💡 Core Capabilities

- **🔐 Security First**: Enterprise-grade security with regular audits and formal verification
- **🌐 Interoperability**: Native support for cross-chain asset transfers and communication
- **📊 Scalability**: Horizontal scaling capabilities to meet growing demand
- **🛠️ Developer Friendly**: Comprehensive SDKs and tooling for rapid development
- **💰 Economic Sustainability**: Novel tokenomics model ensuring long-term value stability

## 🚀 Quick Start

### Prerequisites

Before you begin, ensure your system meets the following requirements:

- **Operating System**: Linux, macOS, or Windows (WSL2)
- **Go**: Version 1.22 or higher
- **Git**: Latest version
- **Hardware**:
  - CPU: 4+ cores
  - RAM: 8GB minimum (16GB recommended)
  - Storage: 100GB+ SSD

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/treasurenetprotocol/treasurenet.git
   cd treasurenet
   ```

2. **Install dependencies**
   ```bash
   make install-deps
   ```

3. **Build from source**
   ```bash
   make build
   ```

4. **Verify installation**
   ```bash
   ./build/treasurenetd version
   ```

### 🐳 Docker Installation

For a containerized setup:

```bash
docker pull treasurenet/node:latest
docker run -d --name treasurenet-node treasurenet/node:latest
```

### 🏃 Running a Node

#### Local Development Node

```bash
# Initialize the node
treasurenetd init my-node --chain-id treasurenet-testnet-1

# Start the node
treasurenetd start
```

#### Joining Mainnet

```bash
# Download genesis file
curl -o ~/.treasurenet/config/genesis.json https://raw.githubusercontent.com/treasurenetprotocol/mainnet/main/genesis.json

# Start with state sync
treasurenetd start --p2p.seeds="seed1@treasurenet.io:26656,seed2@treasurenet.io:26656"
```

## 📚 Documentation

### Core Documentation

- 📖 [Official Documentation](https://wiki.treasurenet.io) - Comprehensive guides and tutorials
- 🔧 [Installation Guide](https://wiki.treasurenet.io/docs/For-Validators/quickStart/) - Detailed setup instructions
- 👨‍💻 [Developer Portal](https://wiki.treasurenet.io/docs/For-Developers/) - Smart contract and dApp development
- 🔍 [API Reference](http://124.70.23.119:8282/#/) - REST and gRPC Gateway documentation

### Technical Specifications

- **Block Time**: ~5 seconds
- **Transaction Throughput**: 2,000+ TPS
- **Smart Contract Support**: Solidity, Vyper (EVM compatible)
- **Token Standard**: TNT (native), ERC-20/721/1155 compatible
- **Network ID**: 
  - Mainnet: `treasurenet-mainnet-1`
  - Testnet: `treasurenet-testnet-1`

## 🛠️ Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/treasurenetprotocol/treasurenet.git
cd treasurenet

# Checkout a specific version (optional)
git checkout v1.5.0

# Build the binary
make build

# Run tests
make test

# Run linter
make lint
```

### Development Tools

- **Treasurenet CLI**: Command-line interface for node operations
- **Web3 SDK**: JavaScript/TypeScript library for dApp development
- **Block Explorer**: [explorer.treasurenet.io](https://explorer.treasurenet.io)
- **Faucet**: [faucet.treasurenet.io](https://faucet.treasurenet.io) (Testnet only)

## 🔌 API Reference

### JSON-RPC Endpoints

```javascript
// Mainnet
const mainnetRPC = "https://rpc.treasurenet.io";

// Testnet
const testnetRPC = "https://testnet-rpc.treasurenet.io";

// WebSocket
const wsEndpoint = "wss://ws.treasurenet.io";
```

### Example: Query Balance

```bash
curl -X POST https://rpc.treasurenet.io \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0x742d35Cc6634C0532925a3b844Bc9e7595f7F8e", "latest"],
    "id": 1
  }'
```

### SDK Usage

```javascript
import { TreasurenetSDK } from '@treasurenet/sdk';

const sdk = new TreasurenetSDK({
  rpcUrl: 'https://rpc.treasurenet.io',
  chainId: 'treasurenet-mainnet-1'
});

// Query account balance
const balance = await sdk.bank.balance('treasurenet1...');
console.log(`Balance: ${balance.amount} TNT`);
```

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

### Local devnet credentials

The local devnet and integration tests require `scripts/.env`, but that file is
intentionally not tracked. Generate fresh local-only credentials before running
them:

```bash
./scripts/generate-test-env.sh
```

The generated file has owner-only permissions. Never commit, share, or reuse it
outside local development; use the environment's approved key-management
process for all production credentials.

### Getting Started

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit your changes**: `git commit -m 'Add amazing feature'`
4. **Push to the branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Contribution Guidelines

- 📝 [Contributing Guide](CONTRIBUTING.md) - Detailed contribution guidelines
- 🐛 [Issue Templates](https://github.com/treasurenetprotocol/treasurenet/issues/new/choose) - Report bugs or request features
- 💻 [Code of Conduct](CODE_OF_CONDUCT.md) - Community standards and expectations
- 🏷️ [Good First Issues](https://github.com/treasurenetprotocol/treasurenet/labels/good%20first%20issue) - Perfect for newcomers

### Development Workflow

1. **Code Style**: Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
2. **Testing**: Write unit tests for new features (minimum 80% coverage)
3. **Documentation**: Update relevant documentation with your changes
4. **Commit Messages**: Use [Conventional Commits](https://www.conventionalcommits.org/)

## 👥 Community

Join our vibrant community and stay updated with the latest developments:

### 🌐 Official Channels

- 💬 [Discord](https://discord.com/invite/treasurenet) - Real-time discussions and support
- 📱 [Telegram](https://t.me/treasurenet) - Announcements and community chat
- 🐦 [X (Twitter)](https://twitter.com/treasurenet_io) - Latest news and updates
- 📺 [YouTube](https://youtube.com/@treasurenet) - Tutorials and webinars
- 📝 [Medium](https://medium.com/@treasurenet) - Technical articles and insights

### 🗓️ Community Calls

- **Developer Calls**: Every Tuesday at 15:00 UTC
- **Governance Calls**: First Thursday of each month at 16:00 UTC
- **Office Hours**: Fridays at 14:00 UTC

## 📄 License

Treasurenet Protocol is licensed under the [GNU Lesser General Public License v3.0](LICENSE).

### Third-Party Software

Treasurenet Protocol includes third-party open-source code. In general, a source subtree with a LICENSE or COPYRIGHT file is from a third party, and our modifications thereto are licensed under the same third-party open source license.

## 🙏 Acknowledgments

We would like to thank:

- The Cosmos SDK team for their foundational work
- The Ethereum community for EVM innovations
- All our contributors and community members
- Our partners and validators supporting the network

---

<p align="center">
  <strong>Build the future of value with Treasurenet Protocol</strong>
  <br>
  <sub>Copyright © 2024 Treasurenet Protocol. All rights reserved.</sub>
</p>
