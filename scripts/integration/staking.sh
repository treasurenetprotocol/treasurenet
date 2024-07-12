#!/bin/bash
set -eux
# your gaiad binary name
BIN=treasurenetd


CHAIN_ID="treasurenet_5005-1"

KEY1="validator"
KEY2="orchestrator"
CHAINID="treasurenet_5005-1"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# to trace evm
TRACE="--trace"
# TRACE=""

#define rpc endpoints  

RPC_ENDPOINT="http://127.0.0.1:26657"  

  

# define account address and verifier address
DELEGATOR_ADDRESS="treasurenet1u7hutc4r88x7anfyjmqgslmqp57y5j9yd2z0zy"  

VALIDATOR_ADDRESS="treasurenetvaloper1u7hutc4r88x7anfyjmqgslmqp57y5j9yv53wrv"



# validate dependencies are installed
# define amount of commission  
AMOUNT="10000000000000000000aunit" # please modify amount and currency according to actual situation  

# execute staking delegation operation using treasurenetd command-line tool 
# DELEGATE_TX=$($BIN tx staking delegate $VALIDATOR_ADDRESS $AMOUNT  
#     --from $KEY1 \  
#     --keyring-backend $KEYRING \  
#     --keyring-dir ~/.treasureenetd \  
#     --node $RPC_ENDPOINT \  
#     --chain-id $CHAIN_ID \  
#     --gas auto \  
#     --gas-prices 10aunit \  
#     --yes)  
DELEGATE_TX=$($BIN tx staking delegate $VALIDATOR_ADDRESS $AMOUNT  --from $KEY1 --keyring-backend $KEYRING --keyring-dir ~/.treasurenetd --node $RPC_ENDPOINT --chain-id $CHAIN_ID --gas auto --fees 1unit --yes)
        
# output result of delegated operation 
echo "commissioned operation result："  
echo "$DELEGATE_TX"