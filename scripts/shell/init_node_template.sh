#!/bin/bash
set -eux

#############################################################################
# Treasurenet Node Initialization Script
#
# Purpose: Fully configures a new Treasurenet blockchain node including:
# - Key generation (validator/orchestrator/ethereum)
# - Genesis file customization
# - Consensus parameters setup
# - Token metadata configuration
# - Account initialization
#
# Features:
# - Secure key management with keyring
# - Custom token definitions
# - Automated genesis configuration
# - Multi-chain support
#############################################################################

# Clean previous node data (WARNING: Destructive operation)
sudo rm -rf $DATA_PATH/*
sudo rm -rf $DATA_PATH/.*
sudo rm -rf $HOME_PATH/*  # Clean default directories

# Ensure required directories exist with proper permissions
sudo mkdir -p $DATA_PATH/.$PROJECT_NAME
sudo chown -R $USER:$USER $DATA_PATH

# Binary and configuration variables
BIN=$BINARY_NAME
CHAIN_ID="$CHAIN_ID"
ALLOCATION="$ALLOCATION"
KEY1="$VALIDATOR_KEY"
KEY2="$ORCHESTRATOR_KEY"
MONIKER="$MONIKER"
KEYRING="$KEYRING"
KEYALGO="$KEYALGO"
LOGLEVEL="info"
TRACE="--trace"
KEYRING_SECRET="$KEYRING_SECRET"

# Validate required dependencies
command -v jq > /dev/null 2>&1 || { 
    echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/";
    exit 1; 
}

# Configure keyring backend and chain ID
$BIN config keyring-backend $KEYRING
$BIN config chain-id $CHAIN_ID

# Define home directory flags
GAIA_HOME="--home $HOME_PATH"
ARGS="--home $HOME_PATH --keyring-backend file"

# Generate cryptographic keys
# Using printf for secure password piping
printf "$KEYRING_SECRET\n$KEYRING_SECRET\n" | $BIN keys add $VALIDATOR_KEY \
    --keyring-backend $KEYRING --algo $KEYALGO 2>> $DATA_PATH/.$PROJECT_NAME/$KEY1-phrases

printf "$KEYRING_SECRET\n" | $BIN keys add $ORCHESTRATOR_KEY \
    --keyring-backend $KEYRING --algo $KEYALGO 2>> $DATA_PATH/.$PROJECT_NAME/$KEY2-phrases

printf "$KEYRING_SECRET\n" | $BIN eth_keys add --keyring-backend $KEYRING \
    >> $DATA_PATH/.$PROJECT_NAME/$KEY1-eth-keys

# Initialize blockchain node
$BIN init $MONIKER --chain-id $CHAIN_ID

# Configure genesis parameters
update_genesis() {
    local file="$HOME_PATH/config/genesis.json"
    local tmp="/tmp/genesis_tmp.json"
    
    jq "$1" "$file" > "$tmp" && mv "$tmp" "$file"
}

# Update staking and token parameters
update_genesis '.app_state["staking"]["params"]["bond_denom"]="$DENOM"'
update_genesis '.app_state["crisis"]["constant_fee"]["denom"]="$DENOM"'
update_genesis '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="$DENOM"'
update_genesis '.app_state["mint"]["params"]["mint_denom"]="$DENOM"'

# Set consensus parameters
update_genesis '.consensus_params["block"]["time_iota_ms"]="1000"'
update_genesis '.consensus_params["block"]["max_gas"]="1000000000000"'

# Disable empty blocks in config.toml
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME_PATH/config/config.toml
else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME_PATH/config/config.toml
fi

# Configure token metadata
jq '.app_state.bank.denom_metadata = [
    {
        "name": "Stake Token",
        "symbol": "Unit",
        "base": "aunit",
        "display": "unit",
        "description": "A staking test token",
        "denom_units": [
            {"denom": "aunit", "exponent": 0},
            {"denom": "unit", "exponent": 18}
        ]
    }
]' $HOME_PATH/config/genesis.json > $HOME_PATH/config/modified_genesis.json && mv $HOME_PATH/config/modified_genesis.json $HOME_PATH/config/genesis.json

jq '.app_state.bech32ibc.nativeHRP = "treasurenet"' $HOME_PATH/config/genesis.json > /tmp/treasurenet-bech32ibc-genesis.json && mv /tmp/treasurenet-bech32ibc-genesis.json $HOME_PATH/config/genesis.json
# Initialize accounts
VALIDATOR_KEY1=$(printf "$KEYRING_SECRET\n" | $BIN keys show $KEY1 -a $ARGS)
ORCHESTRATOR_KEY1=$(printf "$KEYRING_SECRET\n" | $BIN keys show $KEY2 -a $ARGS)

$BIN add-genesis-account $ARGS $VALIDATOR_KEY1 $ALLOCATION
$BIN add-genesis-account $ARGS $ORCHESTRATOR_KEY1 $ALLOCATION

# Update account registry
FILE="/data/account.json"
sudo jq --arg key1 "$KEY1" --arg key2 "$KEY2" \
    --arg validator_key "$VALIDATOR_KEY1" \
    --arg orchestrator_key "$ORCHESTRATOR_KEY1" \
    '. + {($key1): $validator_key, ($key2): $orchestrator_key}' \
    "$FILE" > tmp.json && sudo mv tmp.json "$FILE"

# Generate genesis transaction
ETHEREUM_KEY=$(grep address $DATA_PATH/.$PROJECT_NAME/$KEY1-eth-keys | head -1 | cut -d':' -f2 | tr -d ' ')
printf "$KEYRING_SECRET\n" | $BIN gentx $ARGS \
    --moniker $MONIKER \
    --chain-id=$CHAIN_ID \
    $KEY1 158000000000000000000aunit \
    $ETHEREUM_KEY $ORCHESTRATOR_KEY1

# Finalize genesis
$BIN collect-gentxs

# Move configuration to final location
sudo mv $HOME_PATH/* $DATA_PATH/.$PROJECT_NAME/