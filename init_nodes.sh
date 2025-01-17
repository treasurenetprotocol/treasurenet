#!/bin/bash

# Load environment variables from .env file
source .env

# Ensure node variable is set in .env
if [ -z "$node" ]; then
  echo "Error: 'node' is not set in .env file"
  exit 1
fi

# Load variables from node_config.json based on the value of $node
JSON_FILE="node_config.json"

# Ensure jq is installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# Extract the node-specific configuration from the JSON file
DATA_PATH=$(jq -r ".[$node].DATA_PATH" $JSON_FILE)
HOME_PATH=$(jq -r ".[$node].HOME_PATH" $JSON_FILE)
PROJECT_NAME=$(jq -r ".[$node].PROJECT_NAME" $JSON_FILE)
BINARY_NAME=$(jq -r ".[$node].BINARY_NAME" $JSON_FILE)
CHAIN_ID=$(jq -r ".[$node].CHAIN_ID" $JSON_FILE)
ALLOCATION=$(jq -r ".[$node].ALLOCATION" $JSON_FILE)
VALIDATOR_KEY=$(jq -r ".[$node].VALIDATOR_KEY" $JSON_FILE)
ORCHESTRATOR_KEY=$(jq -r ".[$node].ORCHESTRATOR_KEY" $JSON_FILE)
MONIKER=$(jq -r ".[$node].MONIKER" $JSON_FILE)
KEYRING=$(jq -r ".[$node].KEYRING" $JSON_FILE)
KEYALGO=$(jq -r ".[$node].KEYALGO" $JSON_FILE)
DENOM=$(jq -r ".[$node].DENOM" $JSON_FILE)
NODE_NAME=$(jq -r ".[$node].NODE_NAME" $JSON_FILE)

# Check if any variable is missing
if [ -z "$DATA_PATH" ] || [ -z "$HOME_PATH" ] || [ -z "$PROJECT_NAME" ] || [ -z "$BINARY_NAME" ] || [ -z "$CHAIN_ID" ] || [ -z "$ALLOCATION" ] || [ -z "$VALIDATOR_KEY" ] || [ -z "$ORCHESTRATOR_KEY" ] || [ -z "$MONIKER" ] || [ -z "$KEYRING" ] || [ -z "$KEYALGO" ] || [ -z "$DENOM" ] || [ -z "$NODE_NAME" ]; then
  echo "Error: One or more required values are missing in the node_config.json file."
  exit 1
fi

# Create a temporary environment file for this iteration
ENV_FILE=$(mktemp)
cat <<EOF > $ENV_FILE
export DATA_PATH=$DATA_PATH
export HOME_PATH=$HOME_PATH
export PROJECT_NAME=$PROJECT_NAME
export BINARY_NAME=$BINARY_NAME
export BIN=$BINARY_NAME
export ARGS="--home $HOME_PATH --keyring-backend test"
export CHAIN_ID=$CHAIN_ID
export ALLOCATION=$ALLOCATION
export VALIDATOR_KEY=$VALIDATOR_KEY
export ORCHESTRATOR_KEY=$ORCHESTRATOR_KEY
export MONIKER=$MONIKER
export KEYRING=$KEYRING
export KEYALGO=$KEYALGO
export DENOM=$DENOM
export NODE_NAME=$NODE_NAME
EOF

# Source the environment variables into the current shell session
source $ENV_FILE

# Debugging: Print environment variables to ensure they're set correctly
echo "Environment Variables:"
echo "DATA_PATH=$DATA_PATH"
echo "HOME_PATH=$HOME_PATH"
echo "PROJECT_NAME=$PROJECT_NAME"
echo "BINARY_NAME=$BINARY_NAME"
echo "CHAIN_ID=$CHAIN_ID"
echo "ALLOCATION=$ALLOCATION"
echo "VALIDATOR_KEY=$VALIDATOR_KEY"
echo "ORCHESTRATOR_KEY=$ORCHESTRATOR_KEY"
echo "MONIKER=$MONIKER"
echo "KEYRING=$KEYRING"
echo "KEYALGO=$KEYALGO"
echo "DENOM=$DENOM"
echo "NODE_NAME=$NODE_NAME"

# Read the template file and replace placeholders with actual values using envsubst
TEMPLATE_FILE="init_node_template.sh"
envsubst '${DATA_PATH} ${HOME_PATH} ${PROJECT_NAME} ${BINARY_NAME} ${CHAIN_ID} ${ALLOCATION} ${VALIDATOR_KEY} ${ORCHESTRATOR_KEY} ${MONIKER} ${KEYRING} ${KEYALGO} ${DENOM} ${NODE_NAME}' < $TEMPLATE_FILE > run_$NODE_NAME.sh

# Make the generated script executable
chmod +x run_$NODE_NAME.sh

# Optionally, you can execute the generated script here:
./run_$NODE_NAME.sh

# Clean up temporary environment file
rm $ENV_FILE
