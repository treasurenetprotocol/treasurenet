#!/bin/bash

# This script automates the setup of a Treasurenet blockchain network with multiple nodes.
# It performs the following tasks:
# 1. Copies genesis transactions (gentx) from all nodes to the primary node (node0)
# 2. Adds genesis accounts from a JSON configuration file
# 3. Collects all gentx files to create the final genesis file
# 4. Configures the network settings
# 5. Distributes the genesis file to all nodes
# 6. Records node IDs for network configuration

# Define the list of nodes in the network
nodes=("node0" "node1" "node2" "node3")

# Copy all gentx files from each node to the primary node (node0)
# -a: archive mode (preserve attributes)
# -u: update (only copy when source is newer)
# -v: verbose output
for node in "${nodes[@]}"; do
  cp -auv "/data/${node}/.treasurenetd/config/gentx/"*.json \
    /data/node0/.treasurenetd/config/gentx/
done

# Set HOME environment variable to node0's data directory
# This is required for the treasurenetd commands to work correctly
export HOME=/data/node0

# Path to the JSON file containing account information
json_file="/data/account.json"

# Add genesis accounts for each key in the JSON file
# Skip validator0 and orchestrator0 as they are special accounts
for key in $(jq -r 'keys_unsorted[]' "$json_file"); do
  if [[ "$key" != "validator0" && "$key" != "orchestrator0" ]]; then
    # Extract account address from JSON file
    ACCOUNT=$(jq -r ".${key}" "$json_file")
    echo "Adding genesis account for $key with address $ACCOUNT"
    
    # Add account to genesis with initial token allocations
    # --trace: show stack trace for errors
    # --keyring-backend test: use test keyring (no OS security)
    treasurenetd add-genesis-account --trace --keyring-backend test $ACCOUNT \
      10000000000000000000000aunit,10000000000stake,10000000000footoken,10000000000footoken2,10000000000ibc/nometadatatoken
  fi
done

# Change to gentx directory and collect all gentx files
# This creates the final genesis file by combining all validator transactions
cd /data/node0/.treasurenetd/config/gentx
treasurenetd collect-gentxs

# Backup and replace config.toml with a predefined version
# This ensures consistent network configuration across all nodes
cd /data/node0/.treasurenetd/config/
mv config.toml config.toml1  # Backup original config
cp /data/node1/.treasurenetd/config/config.toml ./  # Use config from node1

# Record each node's ID in the .env file for network configuration
# This is used by the Ansible deployment scripts
for node in "${nodes[@]}"; do
  export HOME="/data/$node"
  # Get the node's tendermint ID
  node_id=$(treasurenetd tendermint show-node-id)
  # Append to .env file in the format: nodeX_address=<node_id>
  echo "${node}_address=$node_id" >> /data/actions-runner/_work/treasurenet/treasurenet/.github/scripts/ansible/docker/.env
done

# Display the contents of the updated .env file
cat /data/actions-runner/_work/treasurenet/treasurenet/.github/scripts/ansible/docker/.env
echo "Node IDs appended to .env file."

# Distribute the final genesis.json file to all nodes
# This ensures all nodes start with the same genesis state
for node in "${nodes[@]}"; do
  cd /data/node0/.treasurenetd/config/
  cp -a genesis.json "/data/${node}/.treasurenetd/config/genesis.json"
done

# Reset HOME environment variable to default
export HOME=/home/ubuntu

# Copy all gentx files from the primary node to other nodes
# This ensures all nodes have a complete set of genesis transactions
for node in "${nodes[@]}"; do
  cd /data/node0/.treasurenetd/config/
  cp -auv "/data/node0/.treasurenetd/config/gentx/"*.json \
    /data/${node}/.treasurenetd/config/gentx/
done