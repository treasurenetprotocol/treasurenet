#!/bin/sh
set -e

CONFIG=$1
if [ -z $CONFIG ]; then
    echo "No config file supplied"
    exit 1
fi
shift

DATA=$1
if [ -z $DATA ]; then
    echo "No data directory supplied"
    exit 1
fi
shift

pwdfile=$(mktemp /tmp/password.XXXXXX)
tmpfile=$(mktemp /tmp/validator-key.XXXXXX)
genesisfile=$(mktemp /tmp/geth-genesis.XXXXXX)

cat > $pwdfile << EOF
$PASSWORD
EOF

# import validator key
validator_key=$(python -c """
from eth_account import Account
Account.enable_unaudited_hdwallet_features()
print(Account.from_mnemonic('$VALIDATOR1_MNEMONIC').key.hex().replace('0x',''))
""")

validator_address=$(python -c """
from eth_account import Account
Account.enable_unaudited_hdwallet_features()
print(Account.from_mnemonic('$VALIDATOR1_MNEMONIC').address)
""")

cat > $tmpfile << EOF
$validator_key
EOF
geth --datadir $DATA account import $tmpfile --password $pwdfile

# import community key
community_key=$(python -c """
from eth_account import Account
Account.enable_unaudited_hdwallet_features()
print(Account.from_mnemonic('$COMMUNITY_MNEMONIC').key.hex().replace('0x',''))
""")

community_address=$(python -c """
from eth_account import Account
Account.enable_unaudited_hdwallet_features()
print(Account.from_mnemonic('$COMMUNITY_MNEMONIC').address)
""")

cat > $tmpfile << EOF
$community_key
EOF
geth --datadir $DATA account import $tmpfile --password $pwdfile

rm $tmpfile

# Build a local-only genesis with the generated validator as Clique signer and
# both generated accounts funded. Never persist this derived configuration.
jq --arg validator "${validator_address#0x}" --arg community "${community_address#0x}" '
  .extraData = ("0x" + ("0" * 64) + ($validator | ascii_downcase) + ("0" * 130)) |
  .alloc = {($validator): .alloc["57f96e6b86cdefdb3d412547816a82e3e0ebf9d2"], ($community): .alloc["378c50D9264C63F3F92B806d4ee56E9D86FfB3Ec"]}
' "$CONFIG" > "$genesisfile"
geth --datadir $DATA init "$genesisfile"

# start up
geth --networkid 5005 --datadir $DATA --http --http.addr localhost --http.api 'personal,eth,net,web3,txpool,miner' \
-unlock "$validator_address" --password $pwdfile \
--mine --allow-insecure-unlock \
$@

rm $pwdfile
rm $genesisfile
