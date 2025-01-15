
accounts00=$(jq -c -r '.app_state.auth.accounts[0]' /data/node0/.treasurenet/config/genesis.json)

accounts01=$(jq -c -r '.app_state.auth.accounts[1]' /data/node0/.treasurenet/config/genesis.json)

accounts10=$(jq -c -r '.app_state.auth.accounts[0]' /data/node1/.treasurenet/config/genesis.json)

accounts11=$(jq -c -r '.app_state.auth.accounts[1]' /data/node1/.treasurenet/config/genesis.json)

accounts20=$(jq -c -r '.app_state.auth.accounts[0]' /data/node2/.treasurenet/config/genesis.json)

accounts21=$(jq -c -r '.app_state.auth.accounts[1]' /data/node2/.treasurenet/config/genesis.json)

accounts30=$(jq -c -r '.app_state.auth.accounts[0]' /data/node3/.treasurenet/config/genesis.json)

accounts31=$(jq -c -r '.app_state.auth.accounts[1]' /data/node3/.treasurenet/config/genesis.json)




balances00=$(jq -c -r '.app_state.bank.balances[0].address' /data/node0/.treasurenet/config/genesis.json)

balances01=$(jq -c -r '.app_state.bank.balances[1].address' /data/node0/.treasurenet/config/genesis.json)
#echo $balances01
balances10=$(jq -c -r '.app_state.bank.balances[0].address' /data/node1/.treasurenet/config/genesis.json)

balances11=$(jq -c -r '.app_state.bank.balances[1].address' /data/node1/.treasurenet/config/genesis.json)

balances20=$(jq -c -r '.app_state.bank.balances[0].address' /data/node2/.treasurenet/config/genesis.json)

balances21=$(jq -c -r '.app_state.bank.balances[1].address' /data/node2/.treasurenet/config/genesis.json)

balances30=$(jq -c -r '.app_state.bank.balances[0].address' /data/node3/.treasurenet/config/genesis.json)

balances31=$(jq -c -r '.app_state.bank.balances[1].address' /data/node3/.treasurenet/config/genesis.json)

    
    
gen_txs00=$(jq -c -r '.app_state.genutil.gen_txs[0]' /data/node0/.treasurenet/config/genesis.json)
#echo "$gen_txs00"
gen_txs10=$(jq -c -r '.app_state.genutil.gen_txs[0]' /data/node1/.treasurenet/config/genesis.json)
#echo "$gen_txs10"
gen_txs20=$(jq -c -r '.app_state.genutil.gen_txs[0]' /data/node2/.treasurenet/config/genesis.json)

gen_txs30=$(jq -c -r '.app_state.genutil.gen_txs[0]' /data/node3/.treasurenet/config/genesis.json)


jq --argjson accounts00 "$accounts00" --argjson accounts01 "$accounts01" \
   --argjson accounts10 "$accounts10" --argjson accounts11 "$accounts11" \
   --argjson accounts20 "$accounts20" --argjson accounts21 "$accounts21" \
   --argjson accounts30 "$accounts30" --argjson accounts31 "$accounts31" \
   --arg balances00 "$balances00" --arg balances01 "$balances01" \
   --arg balances10 "$balances10" --arg balances11 "$balances11" \
   --arg balances20 "$balances20" --arg balances21 "$balances21" \
   --arg balances30 "$balances30" --arg balances31 "$balances31" \
   --argjson gen_txs00 "$gen_txs00" --argjson gen_txs10 "$gen_txs10" \
   --argjson gen_txs20 "$gen_txs20" --argjson gen_txs30 "$gen_txs30" \
   '
   (
     .app_state.auth.accounts |= map(
       if . == "accounts00" then $accounts00 
       elif . == "accounts01" then $accounts01
       elif . == "accounts10" then $accounts10
       elif . == "accounts11" then $accounts11
       elif . == "accounts20" then $accounts20
       elif . == "accounts21" then $accounts21
       elif . == "accounts30" then $accounts30
       elif . == "accounts31" then $accounts31
       else . end
     )
   ) |
   (
     .app_state.bank.balances |= map(
       if .address == "balances00" then .address = $balances00
       elif .address == "balances01" then .address = $balances01
       elif .address == "balances10" then .address = $balances10
       elif .address == "balances11" then .address = $balances11
       elif .address == "balances20" then .address = $balances20
       elif .address == "balances21" then .address = $balances21
       elif .address == "balances30" then .address = $balances30
       elif .address == "balances31" then .address = $balances31
       else . end
     )
   ) |
   (
     .app_state.genutil.gen_txs |= map(
       if . == "gen_txs00" then $gen_txs00 
       elif . == "gen_txs10" then $gen_txs10
       elif . == "gen_txs20" then $gen_txs20
       elif . == "gen_txs30" then $gen_txs30
       else . end
     )
   )
' template.json > genesis.json
cat genesis.json
mv genesis.json /data/node/


