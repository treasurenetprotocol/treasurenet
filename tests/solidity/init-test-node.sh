#!/bin/bash
set -euo pipefail

CHAINID="ethermint_5005-1"
MONIKER="localtestnet"

VAL_KEY="localkey"
USER1_KEY="user1"
USER2_KEY="user2"
USER3_KEY="user3"
USER4_KEY="user4"

# remove existing daemon and client
rm -rf ~/.ethermint*

# Generate fresh, local-only keys after cleaning the test keyring. Do not
# replace these with committed mnemonics: genesis accounts are derived from
# this local keyring below.
for key in "$VAL_KEY" "$USER1_KEY" "$USER2_KEY" "$USER3_KEY" "$USER4_KEY"; do
  ethermintd keys add "$key" --keyring-backend test --algo "eth_secp256k1" >/dev/null
done

ethermintd init $MONIKER --chain-id $CHAINID

# Set gas limit in genesis
cat $HOME/.ethermintd/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json

# Reduce the block time to 1s
sed -i -e '/^timeout_commit =/ s/= .*/= "850ms"/' $HOME/.ethermintd/config/config.toml

# Allocate genesis accounts (cosmos formatted addresses)
ethermintd add-genesis-account "$(ethermintd keys show $VAL_KEY   -a --keyring-backend test)" 1000000000000000000000aphoton,1000000000000000000stake --keyring-backend test
ethermintd add-genesis-account "$(ethermintd keys show $USER1_KEY -a --keyring-backend test)" 1000000000000000000000aphoton,1000000000000000000stake --keyring-backend test
ethermintd add-genesis-account "$(ethermintd keys show $USER2_KEY -a --keyring-backend test)" 1000000000000000000000aphoton,1000000000000000000stake --keyring-backend test
ethermintd add-genesis-account "$(ethermintd keys show $USER3_KEY -a --keyring-backend test)" 1000000000000000000000aphoton,1000000000000000000stake --keyring-backend test
ethermintd add-genesis-account "$(ethermintd keys show $USER4_KEY -a --keyring-backend test)" 1000000000000000000000aphoton,1000000000000000000stake --keyring-backend test

# Sign genesis transaction
ethermintd gentx $VAL_KEY 1000000000000000000stake --amount=1000000000000000000000aphoton --chain-id $CHAINID --keyring-backend test

# Collect genesis tx
ethermintd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
ethermintd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
ethermintd start --pruning=nothing --rpc.unsafe --keyring-backend test --log_level info --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable
