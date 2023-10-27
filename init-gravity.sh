#!/bin/bash
set -eux
# your gaiad binary name
BIN=treasurenetd

#CHAIN_ID="gravity-test"

CHAIN_ID="treasurenet_9000-1"

#NODES=1

#ALLOCATION="100000000000000000000000000aunit"
ALLOCATION="100000000000000000000000000aunit,10000000000stake,10000000000footoken,10000000000footoken2,10000000000ibc/nometadatatoken"
#ALLOCATION="5000000000000000000000aunit,10000000000stake,10000000000footoken,10000000000footoken2,10000000000ibc/nometadatatoken"
#LOGLEVEL="info"

# first we start a genesis.json with validator 1
# validator 1 will also collect the gentx's once gnerated
#STARTING_VALIDATOR=1
#STARTING_VALIDATOR_HOME="--home /validator$STARTING_VALIDATOR"
# todo add git hash to chain name
#$BIN init $STARTING_VALIDATOR_HOME --chain-id=$CHAIN_ID validator1

KEY1="validator"
KEY2="orchestrator"
CHAINID="treasurenet_5005-1"
MONIKER="localtestnet"
KEYRING="test"
#KEYALGO="eth_secp256k1"
KEYALGO="secp256k1"
LOGLEVEL="info"
# to trace evm
TRACE="--trace"
# TRACE=""

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon and client
rm -rf ~/.treasurenet*

make install

$BIN config keyring-backend $KEYRING
$BIN config chain-id $CHAINID

#for i in $(seq 1 $NODES);
#do
GAIA_HOME="--home /root/.treasurenetd"
ARGS="$GAIA_HOME --keyring-backend test"
# Generate a validator key, orchestrator key, and eth key for each validator
$BIN keys add $KEY1 --keyring-backend $KEYRING --algo $KEYALGO 2>> /validator-phrases
$BIN keys add $KEY2 --keyring-backend $KEYRING --algo $KEYALGO 2>> /orchestrator-phrases
$BIN eth_keys add --keyring-backend $KEYRING >> /validator-eth-keys
# if $KEY exists it should be deleted
#treasurenetd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

# Set moniker and chain-id for Treasurenet (Moniker can be anything, chain-id must be an integer)
$BIN init $MONIKER --chain-id $CHAINID 

# Change parameter token denominations to aunit
cat $HOME/.treasurenetd/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aunit"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json
cat $HOME/.treasurenetd/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aunit"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json
cat $HOME/.treasurenetd/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aunit"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json
cat $HOME/.treasurenetd/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aunit"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json

# increase block time (?)
cat $HOME/.treasurenetd/config/genesis.json | jq '.consensus_params["block"]["time_iota_ms"]="1000"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json

# Set gas limit in genesis
cat $HOME/.treasurenetd/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="1000000000000"' > $HOME/.treasurenetd/config/tmp_genesis.json && mv $HOME/.treasurenetd/config/tmp_genesis.json $HOME/.treasurenetd/config/genesis.json

# disable produce empty block
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.treasurenetd/config/config.toml
  else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.treasurenetd/config/config.toml
fi

# if [[ $1 == "pending" ]]; then
#   if [[ "$OSTYPE" == "darwin"* ]]; then
#       sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.treasurenetd/config/config.toml
#   else
#       sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.treasurenetd/config/config.toml
#       sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.treasurenetd/config/config.toml
#   fi
# fi

# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token", "symbol": "FOO", "base": "footoken", display: "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]},{"name": "Stake Token", "symbol": "STEAK", "base": "aunit", display: "unit", "description": "A staking test token", "denom_units": [{"denom": "aunit", "exponent": 0}, {"denom": "unit", "exponent": 18}]}]' /root/.treasurenetd/config/genesis.json > /treasurenet-footoken2-genesis.json
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token2", "symbol": "F20", "base": "footoken2", display: "mfootoken2", "description": "A second non-staking test token", "denom_units": [{"denom": "footoken2", "exponent": 0}, {"denom": "mfootoken2", "exponent": 6}]}]' /treasurenet-footoken2-genesis.json > /treasurenet-bech32ibc-genesis.json

# Set the chain's native bech32 prefix
jq '.app_state.bech32ibc.nativeHRP = "treasurenet"' /treasurenet-bech32ibc-genesis.json > /gov-genesis.json

# a 60 second voting period to allow us to pass governance proposals in the tests
#jq '.app_state.gov.voting_params.voting_period = "120s"' /gov-genesis.json > /community-pool-genesis.json

# Add some funds to the community pool to test Airdrops, note that the gravity address here is the first 20 bytes
# of the sha256 hash of 'distribution' to create the address of the module
#jq '.app_state.distribution.fee_pool.community_pool = [{"denom": "stake", "amount": "100000000000000000000000000.0"}]' /gov-genesis.json > /community-pool2-genesis.json
#jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.ModuleAccount", "base_account": { "account_number": "0", "address": "gravity1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8r0kyvh","pub_key": null,"sequence": "0"},"name": "distribution","permissions": ["basic"]}]' /community-pool2-genesis.json > /community-pool3-genesis.json
#jq '.app_state.bank.balances += [{"address": "gravity1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8r0kyvh", "coins": [{"amount": "1000000000000000000000000", "denom": "stake"}]}]' /community-pool3-genesis.json > /edited-genesis.json

mv /gov-genesis.json /root/.treasurenetd/config/genesis.json
VALIDATOR_KEY=$($BIN keys show validator -a $ARGS)
ORCHESTRATOR_KEY=$($BIN keys show orchestrator -a $ARGS)
$BIN add-genesis-account $ARGS $VALIDATOR_KEY $ALLOCATION
$BIN add-genesis-account $ARGS $ORCHESTRATOR_KEY $ALLOCATION

ORCHESTRATOR_KEY=$($BIN keys show orchestrator -a $ARGS)
ETHEREUM_KEY=$(grep address /validator-eth-keys | sed -n "1"p | sed 's/.*://')
echo $ETHEREUM_KEY
# the /8 containing 7.7.7.7 is assigned to the DOD and never routable on the public internet
# we're using it in private to prevent gaia from blacklisting it as unroutable
# and allow local pex
#$BIN gentx $ARGS $GAIA_HOME --chain-id=$CHAIN_ID 1000000000000000000000aunit $ETHEREUM_KEY $ORCHESTRATOR_KEY
#$BIN gentx $KEY1 1000000000000000000000aunit $ETHEREUM_KEY $ORCHESTRATOR_KEY --keyring-backend $KEYRING --chain-id $CHAINID 
#$BIN gentx $ARGS --moniker $MONIKER --chain-id=$CHAIN_ID  validator 1000000000000000000000aunit $ETHEREUM_KEY $ORCHESTRATOR_KEY
$BIN gentx $ARGS --moniker $MONIKER --chain-id=$CHAIN_ID validator 258000000000000000000aunit $ETHEREUM_KEY $ORCHESTRATOR_KEY
#$BIN gentx $ETHEREUM_KEY 1000000000000000000000aunit --keyring-backend $KEYRING --chain-id $CHAINID
#$BIN gentx $ORCHESTRATOR_KEY 1000000000000000000000aunit --keyring-backend $KEYRING --chain-id $CHAINID
# Collect genesis tx
$BIN collect-gentxs 

# Run this to ensure everything worked and that the genesis file is setup correctly
#$BIN validate-genesis --keyring-backend $KEYRING

# if [[ $1 == "pending" ]]; then
#   echo "pending mode is on, please wait for the first block committed."
# fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed) --json-rpc.ws-address localhost:8546
#./treasurenetd query staking validator ethvaloper1s8znfe26xwq8lm5dhzxh4rjr7n90wrl2hn4dwj -o json --chain-id treasurenet_4143179869527-1 --home /data/mytreasurenet/node8//treasurenetd
#./treasurenetd tx slashing unjail --from node8 --chain-id treasurenet_4143179869527-1 --home /data/mytreasurenet/node8/treasurenetd --keyring-backend test --gas-prices 20aunit
#treasurenetd start --pruning=nothing --trace --log_level info --minimum-gas-prices=0.0001aunit --json-rpc.api eth,txpool,personal,net,debug,web3,miner --p2p.persistent_peers=a2f4cd782091b4429efc0748109d326053c99a59@54.241.61.93:26656,56490415fafa170557e410069faf05e3077500fd@54.193.100.97:26656,a82f9a917b35c77c572f6b9139635cec3c5f0317@18.144.45.142:26656
# --json-rpc.address 0.0.0.0:8555  为了和以太坊启动的端口不冲突
#$BIN start --pruning=nothing --log_level $LOGLEVEL --json-rpc.api eth,txpool,personal,net,debug,web3,miner --trace --json-rpc.address 0.0.0.0:8555 
#gravity start  --log_level $LOGLEVEL  --trace --home /validator1
# nohup ./treasurenetd10 start --home /data/.treasurenetd --evm.tracer=json  --keyring-backend file --json-rpc.api eth,txpool,net,web3,miner --pruning=nothing --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants --chain-id treasurenet_5002-1 < /root/psw.txt >> /data/node0.log 2>&1 &
# nohup ./treasurenetd7 start --home /data/.treasurenetd --evm.tracer=json  --keyring-backend file --json-rpc.api eth,txpool,net,web3,miner --pruning=nothing --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants --chain-id treasurenet_5002-1 < /root/psw.txt >> /data/node1.log 2>&1 &
# nohup ./treasurenetd7 start --home /data/.treasurenetd --evm.tracer=json  --keyring-backend file --json-rpc.api eth,txpool,net,debug,web3,miner --pruning=nothing --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants --chain-id treasurenet_5002-1 < /root/psw.txt >> /data/node2.log 2>&1 &
# nohup ./treasurenetd7 start --home /data/.treasurenetd --evm.tracer=json  --keyring-backend file --json-rpc.api eth,txpool,net,web3,miner --pruning=nothing --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants --chain-id treasurenet_5002-1 < /root/psw.txt >> /data/node3.log 2>&1 &
# nohup ./treasurenetd7 start --home /data/.treasurenetd --evm.tracer=json  --keyring-backend file --json-rpc.api eth,txpool,net,web3,miner --pruning=nothing --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants --chain-id treasurenet_5002-1 < /root/psw.txt >> /data/node4.log 2>&1 &
# nohup ./treasurenetd7 start --home /data/.treasurenetd --evm.tracer=json  --keyring-backend file --json-rpc.api eth,txpool,net,web3,miner --pruning=nothing --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants --chain-id treasurenet_5002-1 < /root/psw.txt >> /data/node5.log 2>&1 &
# nohup ./treasurenetd9 start --home /data/.treasurenetd --keyring-backend file --json-rpc.api eth,txpool,net,debug,web3,miner --pruning=nothing --evm.tracer=json --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants < /data/psw.txt >> /data/node0.log 2>&1 &
# nohup ./treasurenetd21 start --home /data/.treasurenetd --keyring-backend file --json-rpc.api eth,txpool,net,debug,web3,miner --pruning=nothing --evm.tracer=json --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants < /data/psw.txt >> /data/node1.log 2>&1 &
# nohup ./treasurenetd22 start --home /data/.treasurenetd --keyring-backend file --json-rpc.api eth,txpool,net,debug,web3,miner --pruning=nothing --evm.tracer=json --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants < /data/psw.txt >> /data/node2.log 2>&1 &
# nohup ./treasurenetd22 start --home /data/.treasurenetd --keyring-backend file --json-rpc.api eth,txpool,net,debug,web3,miner --pruning=nothing --evm.tracer=json --trace --json-rpc.address 0.0.0.0:8555 --rpc.laddr tcp://0.0.0.0:26657 --x-crisis-skip-assert-invariants < /data/psw.txt >> /data/node3.log 2>&1 &
# cosmovisor run tx gov submit-proposal software-upgrade test1 --title test1 --description Upgrade --upgrade-height 250 --upgrade-info '{"title":"test1","description":"Upgrade to version 1.0.0","plan":{"name":"test1","time":"2023-06-30T12:00:00Z","height":250,"info":"Please follow the upgrade instructions on our website.","epoch":{"chain-id":"treasurenet_8000-1","from-genesis":{"chain_id":"treasurenet_8000-1"}}}}' --from validator --yes --fees 1unit


# 1993d4d8bfe3c99187a2bd8be737e5ba4a6381c4@node0.testnet.treasurenet.io:26656,99b14853bae1aa3f0f111f1a5d5d78a109cb94c1@node1.testnet.treasurenet.io:26656,cdf12ee68327006cd952c1850ef9ddc8dc0a87a9@node2.testnet.treasurenet.io:26656,f75b821fe80cbf835355e1e4e2ebdf0b4001e18f@node3.testnet.treasurenet.io:26656,9852c4976d89d11ad2412c433995503140c47eaf@node4.testnet.treasurenet.io:26656
