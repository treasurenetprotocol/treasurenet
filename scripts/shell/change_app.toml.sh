#!/bin/bash

#############################################################################
# Treasurenet Configuration Update Script
#
# Purpose: Modifies app.toml configuration for multiple validator and service nodes
#          to enable telemetry and API features, then distributes the updated config.
#
# Changes:
# 1. Enables telemetry with hostname and service labels
# 2. Sets Prometheus metrics retention to 100 blocks
# 3. Activates API, Swagger UI, and CORS
#
# Targets:
# - Validators 1-6
# - RPC nodes 1-2
# - Bootnodes 1-2
#############################################################################

# Navigate to primary validator's config directory
cd /data/genesis-validator-1/.treasurenetd/config/

# Define configuration filenames
input_file="app.toml"            # Original configuration
output_file="modified_app.toml"  # Temporary modified file

# Modify configuration using sed with extended regex (-E)
sed -E '

# Telemetry section modifications
/^\[telemetry\]/,/^\[/ {
    s/^(enabled) = false/\1 = true/                   # Enable telemetry collection
    s/^(enable-hostname) = false/\1 = true/           # Include hostnames
    s/^(enable-hostname-label) = false/\1 = true/     # Add hostname labels
    s/^(enable-service-label) = false/\1 = true/      # Add service labels
    s/^(prometheus-retention-time) = 0/\1 = 100/      # Keep metrics for 100 blocks
}

# API section modifications
/^\[api\]/,/^\[/ {
    s/^(enable) = false/\1 = true/                    # Enable API endpoints
    s/^(swagger) = false/\1 = true/                   # Enable Swagger UI
    s/^(enabled-unsafe-cors) = false/\1 = true/       # Allow Cross-Origin requests
}' "$input_file" > "$output_file"

# Replace original config with modified version
mv "$output_file" "$input_file"

# List all target nodes for config distribution
nodes=(
    "/data/genesis-validator-2/.treasurenetd/config"
    "/data/genesis-validator-3/.treasurenetd/config"
    "/data/genesis-validator-4/.treasurenetd/config"
    "/data/genesis-validator-5/.treasurenetd/config"
    "/data/genesis-validator-6/.treasurenetd/config"
    "/data/rpc-1/.treasurenetd/config"
    "/data/rpc-2/.treasurenetd/config"
    "/data/bootnode-1/.treasurenetd/config"
    "/data/bootnode-2/.treasurenetd/config"
)

# Distribute modified config to all nodes
for node_dir in "${nodes[@]}"; do
    cp "$input_file" "$node_dir/$input_file"
    echo "Updated config copied to $node_dir"
done

echo "Configuration update completed for all nodes"