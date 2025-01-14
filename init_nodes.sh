#!/bin/bash
set -eux


CONFIG_FILE="node_config.json"
NODES=$(jq -c '.nodes[]' $CONFIG_FILE)

echo "Cleaning up previous files..."
rm -rf /data/node/* 

for NODE in $NODES; do
  NODE_HOME=$(echo $NODE | jq -r '.NODE_HOME')
  MONIKER=$(echo $NODE | jq -r '.MONIKER')  
  CHAIN_ID=$(echo $NODE | jq -r '.CHAIN_ID')
  KEYRING=$(echo $NODE | jq -r '.KEYRING')
  KEY1=$(echo $NODE | jq -r '.KEY1')
  KEY2=$(echo $NODE | jq -r '.KEY2')
  ALLOCATION=$(echo $NODE | jq -r '.ALLOCATION')
  BIN=$(echo $NODE | jq -r '.BIN')
  KEYALGO=$(echo $NODE | jq -r '.KEYALGO')  
  LOGLEVEL=$(echo $NODE | jq -r '.LOGLEVEL') 

  # 执行模板脚本，传递参数
  ./init_node_template.sh $NODE_HOME $MONIKER $CHAIN_ID $KEYRING $KEY1 $KEY2 $ALLOCATION $BIN $KEYALGO $LOGLEVEL
done
