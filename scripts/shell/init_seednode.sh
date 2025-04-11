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

# Create data directory (with parent directories if needed)
sudo mkdir -p /data/seednode0/.treasurenetd

# Remove old configuration in data directory if exists
sudo rm -rf /data/seednode0/.treasurenetd

# Move newly initialized node config to data directory
sudo mv /home/ubuntu/.treasurenetd /data/seednode0/.treasurenetd

# Change to main node's config directory
cd /data/node0/.treasurenetd/config/

# Copy genesis file to seed node (preserving all attributes)
cp -a genesis.json /data/seednode0/.treasurenetd/config/genesis.json