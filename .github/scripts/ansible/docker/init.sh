#!/bin/bash
set -eux

# your gaiad binary name
BIN=./treasurenetd

CHAIN_ID="treasurenet_5005-1"
ALLOCATION="10000000000000000000000aunit,10000000000stake,10000000000footoken,10000000000footoken2,10000000000ibc/nometadatatoken"
KEY1="validator0"
KEY2="orchestrator0"
CHAINID="treasurenet_5005-1"
MONIKER="node0"
KEYRING="file"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
TRACE="--trace"

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

$BIN config keyring-backend $KEYRING
$BIN config chain-id $CHAINID

GAIA_HOME="--home /home/ec2-user/.treasurenetd"
ARGS="$GAIA_HOME --keyring-backend file"

# Generate a validator1 key, orchestrator1 key, and eth key for each validator1
$BIN keys add $KEY1 --keyring-backend $KEYRING --algo $KEYALGO 2>> /data/treasurenet/validator0-phrases
$BIN keys add $KEY2 --keyring-backend $KEYRING --algo $KEYALGO 2>> /data/treasurenet/orchestrator0-phrases
$BIN eth_keys add --keyring-backend $KEYRING >> /data/treasurenet/validator0-eth-keys

$BIN init $MONIKER --chain-id $CHAINID --keyring-backend $KEYRING --home /home/ec2-user/.treasurenetd

# Modify genesis.json file to update staking, crisis, and gov parameters
cat $HOME/.treasurenetd/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aunit"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json
cat $HOME/.treasurenetd/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aunit"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json
cat $HOME/.treasurenetd/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aunit"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json
cat $HOME/.treasurenetd/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aunit"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json

# Modify consensus params in genesis.json
cat $HOME/.treasurenetd/config/genesis.json | jq '.consensus_params["block"]["time_iota_ms"]="1000"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json
cat $HOME/.treasurenetd/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json

# Disable produce empty block if running on macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.treasurenetd/config/config.toml
else
  sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.treasurenetd/config/config.toml
fi

# Add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token", "symbol": "FOO", "base": "footoken", display: "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]},{"name": "Stake Token", "symbol": "STEAK", "base": "aunit", display: "unit", "description": "A staking test token", "denom_units": [{"denom": "aunit", "exponent": 0}, {"denom": "unit", "exponent": 18}]}]' /home/ec2-user/.treasurenetd/config/genesis.json > /home/ec2-user/treasurenet-footoken2-genesis.json
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token2", "symbol": "F20", "base": "footoken2", display: "mfootoken2", "description": "A second non-staking test token", "denom_units": [{"denom": "footoken2", "exponent": 0}, {"denom": "mfootoken2", "exponent": 6}]}]' /home/ec2-user/treasurenet-footoken2-genesis.json > /home/ec2-user/treasurenet-bech32ibc-genesis.json

# Set the chain's native bech32 prefix
jq '.app_state.bech32ibc.nativeHRP = "treasurenet"' /home/ec2-user/treasurenet-bech32ibc-genesis.json > /home/ec2-user/gov-genesis.json

mv /home/ec2-user/gov-genesis.json /home/ec2-user/.treasurenetd/config/genesis.json

# Get the validator and orchestrator keys and add genesis accounts
validator1_KEY=$($BIN keys show validator0 -a $ARGS)
orchestrator1_KEY=$($BIN keys show orchestrator0 -a $ARGS)
$BIN add-genesis-account $ARGS $validator1_KEY $ALLOCATION
$BIN add-genesis-account $ARGS $orchestrator1_KEY $ALLOCATION

# Generate tx and collect gentxs
ETHEREUM_KEY=$(grep address /data/treasurenet/validator0-eth-keys | sed -n "1"p | sed 's/.*://')
echo $ETHEREUM_KEY

$BIN gentx $ARGS --moniker $MONIKER --chain-id=$CHAIN_ID validator0 200000000000000000000aunit $ETHEREUM_KEY $orchestrator1_KEY
$BIN collect-gentxs
