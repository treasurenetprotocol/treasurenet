#!/bin/bash

# Load variables from JSON file
JSON_FILE="node_config.json"
TEMPLATE_FILE="init_node_template.sh"

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
export KEYRING_SECRET=$KEYRING_SECRET
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
  echo "KEYRING_SECRET=$KEYRING_SECRET"

  # Read the template file and replace placeholders with actual values using envsubst
  envsubst '${DATA_PATH} ${HOME_PATH} ${PROJECT_NAME} ${BINARY_NAME} ${CHAIN_ID} ${ALLOCATION} ${VALIDATOR_KEY} ${ORCHESTRATOR_KEY} ${MONIKER} ${KEYRING} ${KEYALGO} ${DENOM} ${NODE_NAME} ${KEYRING_SECRET}' < $TEMPLATE_FILE > run_$NODE_NAME.sh

  # Make the generated script executable
  chmod +x run_$NODE_NAME.sh

  # Optionally, you can execute the generated script here:
  bash ./run_$NODE_NAME.sh

  # Clean up temporary environment file
  rm $ENV_FILE
done
