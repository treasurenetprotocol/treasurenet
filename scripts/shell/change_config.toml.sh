#!/bin/bash

#############################################################################
# Node Configuration Modifier Script
#
# Purpose: Updates config.toml for multiple Treasurenet nodes to:
#          1. Enable Prometheus metrics
#          2. Open RPC service to all interfaces
#
# Targets: 
# - Genesis validators 1-6
# - Bootnodes 1-2 
# - RPC nodes 1-2
#
# Changes:
# - prometheus = false → true
# - laddr from 127.0.0.1:26657 → 0.0.0.0:26657
#############################################################################

# Configuration filenames
input_file="config.toml"
output_file="modified_config.toml"


# Process all target nodes
for node in genesis-validator-{1..6} bootnode-{1..2} rpc-{1..3}; do
    # Node configuration directory
    target_dir="/data/$node/.treasurenetd/config/"
    
    # Change directory with error checking
    if ! cd "$target_dir"; then
        echo "Error: Failed to access directory $target_dir" >&2
        exit 1
    fi

    # Apply modifications using sed:
    # 1. Enable Prometheus monitoring
    # 2. Change RPC laddr from localhost to all interfaces
    sed -E \
        -e '/^prometheus[[:space:]]*=[[:space:]]*false$/s/false/true/' \
        -e 's|^(laddr[[:space:]]*=[[:space:]]*"tcp://)127.0.0.1:26657"|\10.0.0.0:26657"|' \
        "$input_file" > "$output_file"
    if [[ $node == bootnode* ]]; then
        sed -E -e  '/^addr_book_strict[[:space:]]*=[[:space:]]*true/s/true/false/' "$output_file" 
        echo " $node Turn off address book strict mode"
    fi
    # Safely replace original config file
    mv -- "$output_file" "$input_file"
    
    echo "Updated configuration for $node"
done

echo "All node configurations updated successfully"
