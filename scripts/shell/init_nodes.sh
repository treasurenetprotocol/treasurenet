#!/bin/bash

#############################################################################
# Node Configuration Generator Script
#
# Purpose: Automates node setup by:
# 1. Reading configuration from JSON file
# 2. Generating environment variables
# 3. Creating customized node initialization scripts
#
# Features:
# - Supports multiple node configurations
# - Template-based script generation
# - Secure environment variable handling
# - JSON configuration input
#############################################################################

# Configuration files
JSON_FILE="node_config.json"          # Node configuration file
TEMPLATE_FILE="init_node_template.sh" # Template script for node initialization

# Verify jq is installed for JSON processing
if ! command -v jq >/dev/null 2>&1; then
    echo >&2 "Error: jq not installed. Required for JSON processing."
    echo >&2 "Download from: https://stedolan.github.io/jq/download/"
    exit 1
fi

# Validate configuration files exist
for file in "$JSON_FILE" "$TEMPLATE_FILE"; do
    if [ ! -f "$file" ]; then
        echo >&2 "Error: Required file not found: $file"
        exit 1
    fi
done

# Process each node configuration in JSON array
node_count=$(jq length "$JSON_FILE")
for ((i = 0; i < node_count; i++)); do
    echo "Processing node configuration $((i + 1)) of $node_count"
    
    # Extract configuration parameters
    config_vars=(
        DATA_PATH HOME_PATH PROJECT_NAME BINARY_NAME
        CHAIN_ID ALLOCATION VALIDATOR_KEY ORCHESTRATOR_KEY
        MONIKER KEYRING KEYALGO DENOM NODE_NAME
    )
    
    declare -A config
    for var in "${config_vars[@]}"; do
        config[$var]=$(jq -r ".[$i].$var" "$JSON_FILE")
    done

    # Create temporary environment file
    ENV_FILE=$(mktemp)
    cat > "$ENV_FILE" <<EOF
export DATA_PATH=${config[DATA_PATH]}
export HOME_PATH=${config[HOME_PATH]}
export PROJECT_NAME=${config[PROJECT_NAME]}
export BINARY_NAME=${config[BINARY_NAME]}
export BIN=${config[BINARY_NAME]}
export ARGS="--home ${config[HOME_PATH]} --keyring-backend test"
export CHAIN_ID=${config[CHAIN_ID]}
export ALLOCATION=${config[ALLOCATION]}
export VALIDATOR_KEY=${config[VALIDATOR_KEY]}
export ORCHESTRATOR_KEY=${config[ORCHESTRATOR_KEY]}
export MONIKER=${config[MONIKER]}
export KEYRING=${config[KEYRING]}
export KEYALGO=${config[KEYALGO]}
export DENOM=${config[DENOM]}
export NODE_NAME=${config[NODE_NAME]}
export KEYRING_SECRET=${KEYRING_SECRET}
EOF

    # Load environment variables
    source "$ENV_FILE"

    # Debug output
    echo "--- Node Configuration ---"
    for var in "${config_vars[@]}"; do
        echo "$var=${config[$var]}"
    done
    echo "KEYRING_SECRET=[redacted]"

    # Generate node initialization script
    output_script="run_${config[NODE_NAME]}.sh"
    envsubst '${DATA_PATH} ${HOME_PATH} ${PROJECT_NAME} ${BINARY_NAME} 
              ${CHAIN_ID} ${ALLOCATION} ${VALIDATOR_KEY} ${ORCHESTRATOR_KEY}
              ${MONIKER} ${KEYRING} ${KEYALGO} ${DENOM} ${NODE_NAME} ${KEYRING_SECRET}' \
              < "$TEMPLATE_FILE" > "$output_script"

    # Make script executable
    chmod +x "$output_script"
    echo "Generated initialization script: $output_script"

    # Optional: Execute the script
    # echo "Executing initialization script..."
    ./"$output_script"

    # Clean up
    rm "$ENV_FILE"
done

echo "Node configuration completed successfully"