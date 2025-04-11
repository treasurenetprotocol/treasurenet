#!/bin/bash

# This script sets up a seed node for Treasurenet blockchain
# Main workflow:
# 1. Remove old node configuration
# 2. Initialize new seed node
# 3. Set up data directory
# 4. Copy genesis file from main node

# Remove old node configuration if exists
sudo rm -rf /home/ubuntu/.treasurenetd

# Initialize new seed node with specified chain ID
treasurenetd init seednode0 --chain-id "treasurenet_5005-1"

# Remove old configuration and create data directory
sudo rm -rf /data/seednode0/.treasurenetd
sudo mkdir -p /data/seednode0/.treasurenetd

# Move newly initialized node config to data directory
sudo mv /home/ubuntu/.treasurenetd /data/seednode0/.treasurenetd

# Change to main node's config directory
cd /data/node0/.treasurenetd/config/ || { echo "Failed to change directory: /data/node0/.treasurenetd/config/ not found"; exit 1; }

# Copy genesis file to seed node (preserving all attributes)
cp -a genesis.json /data/seednode0/.treasurenetd/config/genesis.json