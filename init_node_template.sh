#!/bin/bash
set -eux

# Clear any previous data (make sure to do this carefully)
sudo rm -rf $DATA_PATH/*
sudo rm -rf $DATA_PATH/.*
sudo rm -rf $HOME_PATH/*  # 清理默认目录

# Ensure necessary directories exist
sudo mkdir -p $DATA_PATH/$PROJECT_NAME
sudo mkdir -p $DATA_PATH/.$PROJECT_NAME
sudo chown -R $USER:$USER $DATA_PATH


# your gaiad binary name
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

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

$BIN config keyring-backend $KEYRING
$BIN config chain-id $CHAIN_ID

# Define home directory for current user
GAIA_HOME="--home $HOME_PATH"
ARGS="--home $HOME_PATH --keyring-backend test"

# Add keys for validator and orchestrator
$BIN keys add $VALIDATOR_KEY --keyring-backend $KEYRING --algo $KEYALGO 2>> $DATA_PATH/$PROJECT_NAME/$KEY1-phrases
$BIN keys add $ORCHESTRATOR_KEY --keyring-backend $KEYRING --algo $KEYALGO 2>> $DATA_PATH/$PROJECT_NAME/$KEY2-phrases
$BIN eth_keys add --keyring-backend $KEYRING >> $DATA_PATH/$PROJECT_NAME/$KEY1-eth-keys

# Initialize the node
$BIN init $MONIKER --chain-id $CHAIN_ID

# Modify genesis.json using jq
cat $HOME_PATH/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="$DENOM"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json
cat $HOME_PATH/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="$DENOM"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json
cat $HOME_PATH/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="$DENOM"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json
cat $HOME_PATH/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="$DENOM"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json

# Modify consensus params in genesis.json
cat $HOME_PATH/config/genesis.json | jq '.consensus_params["block"]["time_iota_ms"]="1000"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json
cat $HOME_PATH/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="1000000000000"' > $HOME_PATH/config/tmp_genesis.json && mv $HOME_PATH/config/tmp_genesis.json $HOME_PATH/config/genesis.json

# Update the config.toml to disable empty block creation
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME_PATH/config/config.toml
else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME_PATH/config/config.toml
fi

# Modify genesis.json for new tokens
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token", "symbol": "FOO", "base": "footoken", "display": "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]}, {"name": "Stake Token", "symbol": "UNIT", "base": "aunit", "display": "unit", "description": "A staking test token", "denom_units": [{"denom": "aunit", "exponent": 0}, {"denom": "unit", "exponent": 18}]}]' $HOME_PATH/config/genesis.json > /tmp/treasurenet-footoken2-genesis.json
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token2", "symbol": "F20", "base": "footoken2", "display": "mfootoken2", "description": "A second non-staking test token", "denom_units": [{"denom": "footoken2", "exponent": 0}, {"denom": "mfootoken2", "exponent": 6}]}]' /tmp/treasurenet-footoken2-genesis.json > /tmp/treasurenet-bech32ibc-genesis.json

jq '.app_state.bech32ibc.nativeHRP = "treasurenet"' /tmp/treasurenet-bech32ibc-genesis.json > /tmp/gov-genesis.json

# Move final genesis.json back to correct location
mv /tmp/gov-genesis.json $HOME_PATH/config/genesis.json

# Add accounts to genesis
VALIDATOR_KEY1=$($BIN keys show $KEY1 -a $ARGS)
ORCHESTRATOR_KEY1=$($BIN keys show $KEY2 -a $ARGS)
$BIN add-genesis-account $ARGS $VALIDATOR_KEY1 $ALLOCATION
$BIN add-genesis-account $ARGS $ORCHESTRATOR_KEY1 $ALLOCATION

# Generate transaction
ETHEREUM_KEY=$(grep address $DATA_PATH/$PROJECT_NAME/$KEY1-eth-keys | sed -n "1"p | sed 's/.*://')
echo $ETHEREUM_KEY

$BIN gentx $ARGS --moniker $MONIKER --chain-id=$CHAIN_ID $KEY1 258000000000000000000aunit $ETHEREUM_KEY $ORCHESTRATOR_KEY1

# Collect transactions
$BIN collect-gentxs

sudo mv $HOME_PATH/* $DATA_PATH/.$PROJECT_NAME/
