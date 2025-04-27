#!/bin/bash

#############################################################################
# Treasurenet Network Initialization Script
#
# Purpose: Configures the genesis setup for a Treasurenet blockchain network by:
# 1. Collecting genesis transactions (gentx) from all nodes
# 2. Adding genesis accounts from a JSON file
# 3. Creating the final genesis file
# 4. Distributing configuration to all nodes
#
# Targets:
# - Genesis validators 1-6
# - RPC nodes 1-2
#############################################################################

# Define all network nodes
nodes=(
    "genesis-validator-1" "genesis-validator-2" "genesis-validator-3"
    "genesis-validator-4" "genesis-validator-5" "genesis-validator-6"
    "rpc-1" "rpc-2"
)

# Phase 1: Collect genesis transactions from all nodes
echo "Collecting genesis transactions from all nodes..."
for node in "${nodes[@]}"; do
    echo "Copying gentx from $node..."
    cp -auv "/data/${node}/.treasurenetd/config/gentx/"*.json \
        /data/genesis-validator-1/.treasurenetd/config/gentx/
done

# Phase 2: Set up genesis accounts
export HOME=/data/genesis-validator-1

treasurenetd keys add Airdrop_administrator --keyring-backend file
printf "$KEYRING_SECRET\n"
address=$(printf "%s\n" "$KEYRING_SECRET" | treasurenetd keys show Airdrop_administrator -a --keyring-backend file 2>/dev/null)
printf "$KEYRING_SECRET\n" | treasurenetd add-genesis-account \
            --trace \
            --keyring-backend file \
            "$address" \
            000000000000000000aunit
treasurenetd keys add contract_deployer --keyring-backend file 
printf "$KEYRING_SECRET\n"
address=$(printf "%s\n" "$KEYRING_SECRET" | treasurenetd keys show contract_deployer -a --keyring-backend file 2>/dev/null)
printf "$KEYRING_SECRET\n" | treasurenetd add-genesis-account \
            --trace \
            --keyring-backend file \
            "$address" \
            50000000000000000000aunit
json_file="/data/account.json"

echo "Adding genesis accounts from $json_file..."
for key in $(jq -r 'keys_unsorted[]' "$json_file"); do
    # Skip validator1 and orchestrator1 accounts
    if [[ "$key" != "validator1" && "$key" != "orchestrator1" ]]; then
        ACCOUNT=$(jq -r ".${key}" "$json_file" 2>/dev/null)
        
        # Validate account address
        if [[ -z "$ACCOUNT" ]]; then
            echo "Warning: Skipping invalid account for $key"
            continue
        fi

        echo "Adding account: $key ($ACCOUNT)"
        printf "$KEYRING_SECRET\n" | treasurenetd add-genesis-account \
            --trace \
            --keyring-backend file \
            "$ACCOUNT" \
            200000000000000000000aunit
    fi
done

# Phase 3: Finalize genesis file
echo "Creating final genesis file..."
cd /data/genesis-validator-1/.treasurenetd/config/gentx
treasurenetd collect-gentxs

# Phase 4: Update node configuration
echo "Updating node configurations..."
cd /data/genesis-validator-1/.treasurenetd/config/
mv config.toml config.toml1  # Backup original config
cp /data/node1/.treasurenetd/config/config.toml ./  # Use standard config

# Phase 5: Record node IDs for network configuration
echo "Recording node IDs..."
for node in "${nodes[@]}"; do
    export HOME="/data/$node"
    node_id=$(treasurenetd tendermint show-node-id)
    node_name=$(echo "$node" | tr '-' '_')  # Convert to env-friendly name
    echo "${node_name}_address=$node_id" >> /data/actions-runner/_work/treasurenet/treasurenet/scripts/shell/.env
done

# Verify recorded node IDs
echo "Node IDs recorded:"
cat /data/actions-runner/_work/treasurenet/treasurenet/scripts/shell/.env

# Phase 6: Distribute genesis configuration
echo "Distributing genesis file to all nodes..."
for node in "${nodes[@]}"; do
    cd /data/genesis-validator-1/.treasurenetd/config/
    cp -a genesis.json "/data/${node}/.treasurenetd/config/genesis.json"
    echo "Distributed to $node"
done

# Phase 7: Distribute gentx files (if needed)
echo "Distributing gentx files..."
for node in "${nodes[@]}"; do
    cd /data/genesis-validator-1/.treasurenetd/config/
    cp -auv "/data/genesis-validator-1/.treasurenetd/config/gentx/"*.json \
        "/data/${node}/.treasurenetd/config/gentx/"
done

# Clean up
export HOME=/home/ubuntu
echo "Network initialization completed successfully!"