#!/bin/bash  
  
# 定义RPC端点和账户地址  
RPC_ENDPOINT="http://127.0.0.1:26657"  
ADDRESS="treasurenet1u7hutc4r88x7anfyjmqgslmqp57y5j9yd2z0zy"  
  
# 使用treasurenetd命令行工具查询余额  
VALIDATOR_ADDRESS=$(treasurenetd query staking validators --node $RPC_ENDPOINT --output json)  
  
# 输出查询结果    
echo "$VALIDATOR_ADDRESS"