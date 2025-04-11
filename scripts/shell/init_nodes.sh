#!/bin/bash

# This script initializes multiple blockchain nodes from a JSON configuration file
# It reads node parameters from JSON, generates environment variables,
# creates executable scripts from a template, and optionally runs them

# Path to JSON configuration file containing node parameters
JSON_FILE="../../node_config.json"

# Path to template script that will be customized for each node
TEMPLATE_FILE="init_node_template.sh"

# Check if jq (JSON processor) is installed
# Exit with error message if jq is not available
command -v jq > /dev/null 2>&1 || { 
    echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; 
    exit 1; 
}

# Loop through each node configuration in the JSON array
for ((i=0; i<$(jq length $JSON_FILE); i++)); do
  # Extract all node parameters from JSON configuration
  DATA_PATH=$(jq -r ".[$i].DATA_PATH" $JSON_FILE)
  HOME_PATH=$(jq -r ".[$i].HOME_PATH" $JSON_FILE)
  PROJECT_NAME=$(jq -r ".[$i].PROJECT_NAME" $JSON_FILE)
  BINARY_NAME=$(jq -r ".[$i].BINARY_NAME" $JSON_FILE)
  CHAIN_ID=$(jq -r ".[$i].CHAIN_ID" $JSON_FILE)
  ALLOCATION=$(jq -r ".[$i].ALLOCATION" $JSON_FILE)
  VALIDATOR_KEY=$(jq -r ".[$i].VALIDATOR_KEY" $JSON_FILE)
  ORCHESTRATOR_KEY=$(jq -r ".[$i].ORCHESTRATOR_KEY" $JSON_FILE)
  MONIKER=$(jq -r ".[$i].MONIKER" $JSON_FILE)
  KEYRING=$(jq -r ".[$i].KEYRING" $JSON_FILE)
  KEYALGO=$(jq -r ".[$i].KEYALGO" $JSON_FILE)
  DENOM=$(jq -r ".[$i].DENOM" $JSON_FILE)
  NODE_NAME=$(jq -r ".[$i].NODE_NAME" $JSON_FILE)

  # Create temporary environment file to store extracted variables
  ENV_FILE=$(mktemp)
  
  # Write all environment variables to temporary file
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

  # Load environment variables into current shell session
  source $ENV_FILE

  # Debug output - display all loaded variables
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

  # Generate node initialization script by substituting variables in template
  # envsubst replaces ${VAR} placeholders with actual values
  envsubst '${DATA_PATH} ${HOME_PATH} ${PROJECT_NAME} ${BINARY_NAME} ${CHAIN_ID} ${ALLOCATION} ${VALIDATOR_KEY} ${ORCHESTRATOR_KEY} ${MONIKER} ${KEYRING} ${KEYALGO} ${DENOM} ${NODE_NAME}' < $TEMPLATE_FILE > run_$NODE_NAME.sh

  # Make generated script executable
  chmod +x run_$NODE_NAME.sh

  # Optional: Execute the generated script immediately
  bash ./run_$NODE_NAME.sh

  # Clean up temporary environment file
  rm $ENV_FILE
done