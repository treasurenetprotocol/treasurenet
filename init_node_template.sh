#!/bin/bash
set -eux

NODE_HOME=$1
MONIKER=$2
CHAIN_ID=$3
KEYRING=$4
KEY1=$5
KEY2=$6
ALLOCATION=$7
BIN=$8
KEYALGO=$9
LOGLEVEL=${10}

command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

$BIN config keyring-backend $KEYRING
$BIN config chain-id $CHAIN_ID --home $NODE_HOME
$BIN eth_keys add --keyring-backend $KEYRING >> $NODE_HOME/validator0-eth-keys

$BIN init $MONIKER --chain-id $CHAIN_ID --keyring-backend $KEYRING --home $NODE_HOME

GENESIS_JSON="${NODE_HOME}/config/genesis.json"
cat $GENESIS_JSON | jq '.app_state["staking"]["params"]["bond_denom"]="aunit"' > ${NODE_HOME}/config/tmp_genesis.json && mv ${NODE_HOME}/config/tmp_genesis.json $GENESIS_JSON
cat $GENESIS_JSON | jq '.app_state["crisis"]["constant_fee"]["denom"]="aunit"' > ${NODE_HOME}/config/tmp_genesis.json && mv ${NODE_HOME}/config/tmp_genesis.json $GENESIS_JSON
cat $GENESIS_JSON | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aunit"' > ${NODE_HOME}/config/tmp_genesis.json && mv ${NODE_HOME}/config/tmp_genesis.json $GENESIS_JSON
cat $GENESIS_JSON | jq '.app_state["mint"]["params"]["mint_denom"]="aunit"' > ${NODE_HOME}/config/tmp_genesis.json && mv ${NODE_HOME}/config/tmp_genesis.json $GENESIS_JSON

if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' ${NODE_HOME}/config/config.toml
else
  sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' ${NODE_HOME}/config/config.toml
fi

jq '.app_state.bank.denom_metadata += [{"name": "Foo Token", "symbol": "FOO", "base": "footoken", "display": "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]},{"name": "Stake Token", "symbol": "STEAK", "base": "aunit", "display": "unit", "description": "A staking test token", "denom_units": [{"denom": "aunit", "exponent": 0}, {"denom": "unit", "exponent": 18}]}]' ${NODE_HOME}/config/genesis.json > ${NODE_HOME}/config/tmp_genesis.json && mv ${NODE_HOME}/config/tmp_genesis.json $GENESIS_JSON

jq '.app_state.bank.denom_metadata += [{"name": "Foo Token2", "symbol": "F20", "base": "footoken2", "display": "mfootoken2", "description": "A second non-staking test token", "denom_units": [{"denom": "footoken2", "exponent": 0}, {"denom": "mfootoken2", "exponent": 6}]}]' ${NODE_HOME}/config/genesis.json > ${NODE_HOME}/config/tmp_genesis.json && mv ${NODE_HOME}/config/tmp_genesis.json ${NODE_HOME}/config/genesis.json

jq '.app_state.bech32ibc.nativeHRP = "treasurenet"' ${NODE_HOME}/config/genesis.json > ${NODE_HOME}/config/tmp_genesis.json && mv ${NODE_HOME}/config/tmp_genesis.json ${NODE_HOME}/config/genesis.json

ETHEREUM_KEY=$(grep address $NODE_HOME/validator0-eth-keys | sed -n "1"p | sed 's/.*://')
echo $ETHEREUM_KEY

$BIN gentx --keyring-backend $KEYRING --home $NODE_HOME --moniker $MONIKER --chain-id=$CHAIN_ID $KEY1 200000000000000000000aunit $ETHEREUM_KEY
$BIN collect-gentxs --home $NODE_HOME

echo "Node initialization complete, all transactions collected."
