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

DELEGATOR_ADDRESS="treasurenet1u2nl4da6rpkag4wymwwy8c673nmqgpfq2e7fmx"  

VALIDATOR_ADDRESS="treasurenetvaloper1u2nl4da6rpkag4wymwwy8c673nmqgpfqt8dg6w"



# validate dependencies are installed
# 定义委托的金额  
AMOUNT="10000000000000000000aunit" # 请根据实际情况修改金额和币种  
   
  
# 使用evmosd命令行工具执行staking委托操作  
DELEGATE_TX=$($BIN tx staking delegate $VALIDATOR_ADDRESS $AMOUNT \  
    --from $KEY1 \  
    --keyring-backend $KEYRING \  
    --keyring-dir ~/.treasureenetd \  
    --node $RPC_ENDPOINT \  
    --chain-id $CHAIN_ID \  
    --gas auto \  
    --gas-prices 10aunit \  
    --yes \  
    --trust-node=true)  
  
# 输出委托操作的结果  
echo "委托操作结果："  
echo "$DELEGATE_TX"