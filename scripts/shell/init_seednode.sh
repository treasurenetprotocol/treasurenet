#!/bin/bash

#############################################################################
# Treasurenet Bootnode Initialization Script
#
# Purpose: Sets up bootnodes for Treasurenet network by:
# 1. Initializing node directories
# 2. Copying genesis configuration
# 3. Recording node IDs for network configuration
#
# Targets: bootnode-1 and bootnode-2
#############################################################################

# Initialize first bootnode
echo "Initializing bootnode-1..."
sudo rm -rf /home/ubuntu/.treasurenetd
treasurenetd init bootnode-1 --chain-id "treasurenet_5005-1"

# Set up data directory for bootnode-1
sudo mkdir -p /data/bootnode-1/.treasurenetd
sudo rm -rf /data/bootnode-1/.treasurenetd
sudo mv /home/ubuntu/.treasurenetd /data/bootnode-1/.treasurenetd

# Copy genesis configuration from validator
cd /data/genesis-validator-1/.treasurenetd/config/
cp -a genesis.json /data/bootnode-1/.treasurenetd/config/genesis.json

# Initialize second bootnode
echo "Initializing bootnode-2..."
sudo rm -rf /home/ubuntu/.treasurenetd
treasurenetd init bootnode-2 --chain-id "treasurenet_5005-1"

# Set up data directory for bootnode-2
sudo mkdir -p /data/bootnode-2/.treasurenetd
sudo rm -rf /data/bootnode-2/.treasurenetd
sudo mv /home/ubuntu/.treasurenetd /data/bootnode-2/.treasurenetd

# Copy genesis configuration from validator
cd /data/genesis-validator-1/.treasurenetd/config/
cp -a genesis.json /data/bootnode-2/.treasurenetd/config/genesis.json

# Define bootnodes array
nodes=("bootnode-1" "bootnode-2")

# Record node IDs to .env file for Ansible configuration
echo "Recording node IDs..."
for node in "${nodes[@]}"; do
  export HOME="/data/$node"
  node_id=$(treasurenetd tendermint show-node-id)
  node_name=$(echo "$node" | tr '-' '_')  # Convert hyphens to underscores
  echo "${node_name}_address=$node_id" >> /data/actions-runner/_work/treasurenet/treasurenet/.github/scripts/ansible/docker/.env
done

# Restore default HOME environment variable
export HOME=/home/ubuntu
echo "Bootnode initialization completed successfully."