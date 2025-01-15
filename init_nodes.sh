#!/bin/bash

# Load variables from JSON file
JSON_FILE="variables.json"
TEMPLATE_FILE="template.sh"

# Ensure jq is installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# Read JSON and iterate over each object
for ((i=0; i<$(jq length $JSON_FILE); i++)); do
  # Extract variables for this iteration
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

  # Create a temporary environment file for this iteration
  ENV_FILE=$(mktemp)
  cat <<EOF > $ENV_FILE
export DATA_PATH=$DATA_PATH
export HOME_PATH=$HOME_PATH
export PROJECT_NAME=$PROJECT_NAME
export BINARY_NAME=$BINARY_NAME
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

  # Read the template file and replace placeholders with actual values using envsubst
  envsubst < $TEMPLATE_FILE > run_$NODE_NAME.sh

  # Make the generated script executable
  chmod +x run_$NODE_NAME.sh

  # Optionally, you can execute the generated script here:
  # ./run_$NODE_NAME.sh

  # Clean up temporary environment file
  rm $ENV_FILE
done