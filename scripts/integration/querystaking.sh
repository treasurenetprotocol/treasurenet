#!/bin/bash  
  
# define rpc endpoints and account addresses  
RPC_ENDPOINT="http://127.0.0.1:26657"  
ADDRESS="treasurenet1u7hutc4r88x7anfyjmqgslmqp57y5j9yd2z0zy"  
  
# use trendenenetd command-line tool to query balance  
VALIDATOR_ADDRESS=$(treasurenetd query staking validators --node $RPC_ENDPOINT --output json)  
  
# output query results    
echo "$VALIDATOR_ADDRESS"