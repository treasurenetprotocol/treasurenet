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

# 定义RPC端点  

RPC_ENDPOINT="http://127.0.0.1:26657"  

  

# 定义账户地址和验证人地址  

DELEGATOR_ADDRESS="treasurenet1u7hutc4r88x7anfyjmqgslmqp57y5j9yd2z0zy"  

VALIDATOR_ADDRESS="treasurenetvaloper1u7hutc4r88x7anfyjmqgslmqp57y5j9yv53wrv"



# validate dependencies are installed
# 定义委托的金额  
AMOUNT="10000000000000000000aunit" # 请根据实际情况修改金额和币种  

DISTRIBUTION_TX=$($BIN tx distribution withdraw-all-rewards  --from $KEY1 --keyring-backend $KEYRING --keyring-dir ~/.treasurenetd --node $RPC_ENDPOINT --chain-id $CHAIN_ID --gas auto --fees 1unit --yes)
         
echo "$DISTRIBUTION_TX"