#!/bin/bash
set -eux
###############################################################################
# Treasurenet Node Initialization Script
#
# This script automates the setup of a Treasurenet blockchain node with:
# - Complete genesis file configuration
# - Validator and orchestrator key generation
# - Custom token and staking parameters
# - Ethereum-compatible key generation
#
# The script:
# 1. Cleans existing data directories
# 2. Generates cryptographic keys
# 3. Configures the genesis file
# 4. Sets up initial accounts and balances
# 5. Prepares the node for network participation
#
# Requirements:
# - jq (for JSON processing)
# - Treasurenet binary (gaiad or other specified binary)
# - Properly configured node_config.json
###############################################################################
# Clear any previous blockchain data (caution: this is destructive)
sudo rm -rf $DATA_PATH/*
sudo rm -rf $DATA_PATH/.*
sudo rm -rf $HOME_PATH/*  

# Create required directories with proper permissions
sudo mkdir -p $DATA_PATH/$PROJECT_NAME
sudo mkdir -p $DATA_PATH/.$PROJECT_NAME
sudo chown -R $USER:$USER $DATA_PATH

# Set binary name variable
BIN=$BINARY_NAME

# Configuration variables from environment
CHAIN_ID="$CHAIN_ID"
ALLOCATION="$ALLOCATION"
KEY1="$VALIDATOR_KEY"
KEY2="$ORCHESTRATOR_KEY"
MONIKER="$MONIKER"
KEYRING="$KEYRING"
KEYALGO="$KEYALGO"
LOGLEVEL="info"
TRACE="--trace"

# Verify jq is installed for JSON processing
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# Configure keyring backend and chain ID
$BIN config keyring-backend $KEYRING
$BIN config chain-id $CHAIN_ID

# Define home directory flags
GAIA_HOME="--home $HOME_PATH"
ARGS="--home $HOME_PATH --keyring-backend test"

# Generate validator and orchestrator keys
$BIN keys add $VALIDATOR_KEY --keyring-backend $KEYRING --algo $KEYALGO 2>> $DATA_PATH/.$PROJECT_NAME/$KEY1-phrases
$BIN keys add $ORCHESTRATOR_KEY --keyring-backend $KEYRING --algo $KEYALGO 2>> $DATA_PATH/.$PROJECT_NAME/$KEY2-phrases
$BIN eth_keys add --keyring-backend $KEYRING >> $DATA_PATH/.$PROJECT_NAME/$KEY1-eth-keys

# Initialize blockchain node
$BIN init $MONIKER --chain-id $CHAIN_ID

# Update genesis.json with custom denominations
cat $HOME_PATH/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="$DENOM"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json
cat $HOME_PATH/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="$DENOM"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json
cat $HOME_PATH/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="$DENOM"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json
cat $HOME_PATH/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="$DENOM"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json

# Modify consensus parameters
cat $HOME_PATH/config/genesis.json | jq '.consensus_params["block"]["time_iota_ms"]="1000"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json
cat $HOME_PATH/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="1000000000000"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json

# Disable empty blocks in config.toml (OS-specific sed syntax)
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME_PATH/config/config.toml
else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME_PATH/config/config.toml
fi

# Add custom token metadata to genesis
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token", "symbol": "FOO", "base": "footoken", "display": "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]}, {"name": "Stake Token", "symbol": "UNIT", "base": "aunit", "display": "unit", "description": "A staking test token", "denom_units": [{"denom": "aunit", "exponent": 0}, {"denom": "unit", "exponent": 18}]}]' $HOME_PATH/config/genesis.json > /tmp/treasurenet-footoken2-genesis.json
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token2", "symbol": "F20", "base": "footoken2", "display": "mfootoken2", "description": "A second non-staking test token", "denom_units": [{"denom": "footoken2", "exponent": 0}, {"denom": "mfootoken2", "exponent": 6}]}]' /tmp/treasurenet-footoken2-genesis.json > /tmp/treasurenet-bech32ibc-genesis.json

# Set native HRP for bech32 addresses
jq '.app_state.bech32ibc.nativeHRP = "treasurenet"' /tmp/treasurenet-bech32ibc-genesis.json > /tmp/gov-genesis.json

# Finalize genesis.json
mv /tmp/gov-genesis.json $HOME_PATH/config/genesis.json

# Add initial accounts to genesis
VALIDATOR_KEY1=$($BIN keys show $KEY1 -a $ARGS)
ORCHESTRATOR_KEY1=$($BIN keys show $KEY2 -a $ARGS)
$BIN add-genesis-account $ARGS $VALIDATOR_KEY1 $ALLOCATION
$BIN add-genesis-account $ARGS $ORCHESTRATOR_KEY1 $ALLOCATION

# Update account.json with validator addresses
FILE="/data/account.json"
sudo jq --arg key1 "$KEY1" --arg key2 "$KEY2" --arg validator_key "$VALIDATOR_KEY1" --arg orchestrator_key "$ORCHESTRATOR_KEY1" \
    '. + {($key1): $validator_key, ($key2): $orchestrator_key}' "$FILE" > tmp.json && sudo mv tmp.json "$FILE"

# Extract Ethereum address from generated keys
ETHEREUM_KEY=$(grep address $DATA_PATH/.$PROJECT_NAME/$KEY1-eth-keys | sed -n "1"p | sed 's/.*://')
echo $ETHEREUM_KEY

# Generate genesis transaction
$BIN gentx $ARGS --moniker $MONIKER --chain-id=$CHAIN_ID $KEY1 258000000000000000000aunit $ETHEREUM_KEY $ORCHESTRATOR_KEY1

# Collect all genesis transactions
$BIN collect-gentxs

# Move final configuration to data directory
sudo mv $HOME_PATH/* $DATA_PATH/.$PROJECT_NAME/