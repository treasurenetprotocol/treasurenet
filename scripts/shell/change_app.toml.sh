#!/bin/bash

##############################################################################
# Treasurenet Configuration Modifier Script
#
# Purpose: Modifies app.toml configuration to:
#          1. Enable telemetry features
#          2. Activate API endpoints
# 
#
# Targets: node0-node3 and seednode0
# Changes:
# - Enables telemetry collection
# - Configures Prometheus retention
# 
#
# Note: Run from a node with write access to all node directories
##############################################################################

# Navigate to primary node's config directory
cd /data/node0/.treasurenetd/config/

# Define configuration filenames
input_file="app.toml"            # Original configuration file
output_file="modified_app.toml"  # Temporary modified file

# Modify configuration using sed with extended regex (-E)
sed -E '

# Telemetry section modifications (between [telemetry] and next section)
/^\[telemetry\]/,/^\[/ {
    s/^(enabled) = false/\1 = true/                   # Enable telemetry
    s/^(enable-hostname) = false/\1 = true/           # Record hostnames
    s/^(enable-hostname-label) = false/\1 = true/     # Add hostname labels
    s/^(enable-service-label) = false/\1 = true/      # Add service labels
    s/^(prometheus-retention-time) = 0/\1 = 100/      # Set metrics retention
}

# API section modifications (between [api] and next section)
/^\[api\]/,/^\[/ {
    s/^(enable) = false/\1 = true/                    # Enable API
    s/^(swagger) = false/\1 = true/                   # Enable Swagger UI
    s/^(enabled-unsafe-cors) = false/\1 = true/       # Allow CORS
}' "$input_file" > "$output_file"

# Replace original config with modified version
mv ./$output_file ./$input_file

# Distribute modified config to all nodes
cp ./$input_file /data/node1/.treasurenetd/config/$input_file  # Node1
cp ./$input_file /data/node2/.treasurenetd/config/$input_file  # Node2
cp ./$input_file /data/node3/.treasurenetd/config/$input_file  # Node3
cp ./$input_file /data/seednode0/.treasurenetd/config/$input_file  # Seednode