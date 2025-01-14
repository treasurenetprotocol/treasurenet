#!/bin/bash
set -eux
$BIN keys add $KEY1 --keyring-backend $KEYRING --algo $KEYALGO 2>> /data/treasurenet/validator0-phrases
$BIN keys add $KEY2 --keyring-backend $KEYRING --algo $KEYALGO 2>> /data/treasurenet/orchestrator0-phrases
$BIN eth_keys add --keyring-backend $KEYRING >> /data/treasurenet/validator0-eth-keys
