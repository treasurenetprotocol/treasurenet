
# accounts00=$(jq -c -r '.app_state.auth.accounts[0]' /data/node0/.treasurenetd/config/genesis.json)

# accounts01=$(jq -c -r '.app_state.auth.accounts[1]' /data/node0/.treasurenetd/config/genesis.json)

# accounts10=$(jq -c -r '.app_state.auth.accounts[0]' /data/node1/.treasurenetd/config/genesis.json)

# accounts11=$(jq -c -r '.app_state.auth.accounts[1]' /data/node1/.treasurenetd/config/genesis.json)

# accounts20=$(jq -c -r '.app_state.auth.accounts[0]' /data/node2/.treasurenetd/config/genesis.json)

# accounts21=$(jq -c -r '.app_state.auth.accounts[1]' /data/node2/.treasurenetd/config/genesis.json)

# accounts30=$(jq -c -r '.app_state.auth.accounts[0]' /data/node3/.treasurenetd/config/genesis.json)

# accounts31=$(jq -c -r '.app_state.auth.accounts[1]' /data/node3/.treasurenetd/config/genesis.json)




# balances00=$(jq -c -r '.app_state.bank.balances[0].address' /data/node0/.treasurenetd/config/genesis.json)

# balances01=$(jq -c -r '.app_state.bank.balances[1].address' /data/node0/.treasurenetd/config/genesis.json)
# #echo $balances01
# balances10=$(jq -c -r '.app_state.bank.balances[0].address' /data/node1/.treasurenetd/config/genesis.json)

# balances11=$(jq -c -r '.app_state.bank.balances[1].address' /data/node1/.treasurenetd/config/genesis.json)

# balances20=$(jq -c -r '.app_state.bank.balances[0].address' /data/node2/.treasurenetd/config/genesis.json)

# balances21=$(jq -c -r '.app_state.bank.balances[1].address' /data/node2/.treasurenetd/config/genesis.json)

# balances30=$(jq -c -r '.app_state.bank.balances[0].address' /data/node3/.treasurenetd/config/genesis.json)

# balances31=$(jq -c -r '.app_state.bank.balances[1].address' /data/node3/.treasurenetd/config/genesis.json)

    
    
# gen_txs00=$(jq '.' /data/node0/.treasurenetd/config/gentx/gen*.json)
# #echo "$gen_txs00"
# gen_txs10=$(jq '.' /data/node1/.treasurenetd/config/gentx/gen*.json)
# #echo "$gen_txs10"
# gen_txs20=$(jq '.' /data/node2/.treasurenetd/config/gentx/gen*.json)

# gen_txs30=$(jq '.' /data/node3/.treasurenetd/config/gentx/gen*.json)


# jq --argjson accounts00 "$accounts00" --argjson accounts01 "$accounts01" \
#    --argjson accounts10 "$accounts10" --argjson accounts11 "$accounts11" \
#    --argjson accounts20 "$accounts20" --argjson accounts21 "$accounts21" \
#    --argjson accounts30 "$accounts30" --argjson accounts31 "$accounts31" \
#    --arg balances00 "$balances00" --arg balances01 "$balances01" \
#    --arg balances10 "$balances10" --arg balances11 "$balances11" \
#    --arg balances20 "$balances20" --arg balances21 "$balances21" \
#    --arg balances30 "$balances30" --arg balances31 "$balances31" \
#    --argjson gen_txs00 "$gen_txs00" --argjson gen_txs10 "$gen_txs10" \
#    --argjson gen_txs20 "$gen_txs20" --argjson gen_txs30 "$gen_txs30" \
#    '
#    (
#      .app_state.auth.accounts |= map(
#        if . == "accounts00" then $accounts00 
#        elif . == "accounts01" then $accounts01
#        elif . == "accounts10" then $accounts10
#        elif . == "accounts11" then $accounts11
#        elif . == "accounts20" then $accounts20
#        elif . == "accounts21" then $accounts21
#        elif . == "accounts30" then $accounts30
#        elif . == "accounts31" then $accounts31
#        else . end
#      )
#    ) |
#    (
#      .app_state.bank.balances |= map(
#        if .address == "balances00" then .address = $balances00
#        elif .address == "balances01" then .address = $balances01
#        elif .address == "balances10" then .address = $balances10
#        elif .address == "balances11" then .address = $balances11
#        elif .address == "balances20" then .address = $balances20
#        elif .address == "balances21" then .address = $balances21
#        elif .address == "balances30" then .address = $balances30
#        elif .address == "balances31" then .address = $balances31
#        else . end
#      )
#    ) |
#    (
#      .app_state.genutil.gen_txs |= map(
#        if . == "gen_txs00" then $gen_txs00 
#        elif . == "gen_txs10" then $gen_txs10
#        elif . == "gen_txs20" then $gen_txs20
#        elif . == "gen_txs30" then $gen_txs30
#        else . end
#      )
#    )
# ' template.json > genesis.json



sudo rm -rf /data/node/gen_txs/*
sudo mkdir -p /data/node/gen_txs
sudo chown -R $USER:$USER /data/node

cp -a /data/node1/.treasurenetd/config/gentx/* /data/node/gen_txs/
cp -a /data/node2/.treasurenetd/config/gentx/* /data/node/gen_txs/
cp -a /data/node3/.treasurenetd/config/gentx/* /data/node/gen_txs/
cp -a /data/node4/.treasurenetd/config/gentx/* /data/node/gen_txs/

export HOME=/data/node1

# 读取 JSON 文件并提取除 validator1 和 orchestrator1 外的地址
json_file="/data/test.json"

# 使用 jq 逐一提取每个地址并执行命令
for key in $(jq -r 'keys_unsorted[]' "$json_file"); do
  if [[ "$key" != "validator1" && "$key" != "orchestrator1" ]]; then
    ACCOUNT=$(jq -r ".${key}" "$json_file")
    echo "Adding genesis account for $key with address $ACCOUNT"
    
    # 执行命令
    treasurenetd add-genesis-account --trace --keyring-backend test $ACCOUNT 10000000000000000000000aunit,10000000000stake,10000000000footoken,10000000000footoken2,10000000000ibc/nometadatatoken
  fi
done
cd /data/node1/.treasurenetd/config/gentx
treasurenetd collect-gentxs

cd /data/node1/.treasurenetd/config/
mv config.toml config.toml1
cp /data/node2/.treasurenetd/config/config.toml ./

cp -a genesis.json /data/node1/.treasurenetd/config/genesis.json
cp -a genesis.json /data/node2/.treasurenetd/config/genesis.json
cp -a genesis.json /data/node3/.treasurenetd/config/genesis.json
cp -a genesis.json /data/node4/.treasurenetd/config/genesis.json

export HOME=/home/ubuntu

cp -a /data/node/gen_txs/* /data/node1/.treasurenetd/config/gentx/
cp -a /data/node/gen_txs/* /data/node2/.treasurenetd/config/gentx/
cp -a /data/node/gen_txs/* /data/node3/.treasurenetd/config/gentx/
cp -a /data/node/gen_txs/* /data/node4/.treasurenetd/config/gentx/

