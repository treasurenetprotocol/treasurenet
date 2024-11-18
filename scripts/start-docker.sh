#!/bin/bash

echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./treasurenetd validate-genesis --home /treasurenet

echo "starting treasurenet node $ID in background ..."
./treasurenetd start \
--home /treasurenet \
--keyring-backend test

echo "started treasurenet node"
tail -f /dev/null