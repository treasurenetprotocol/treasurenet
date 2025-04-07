sudo rm -rf /home/ubuntu/.treasurenetd

treasurenetd init bootnode-1 --chain-id "treasurenet_5005-1"

sudo mkdir -p /data/bootnode-1/.treasurenetd

sudo rm -rf /data/bootnode-1/.treasurenetd

sudo mv /home/ubuntu/.treasurenetd  /data/bootnode-1/.treasurenetd

cd /data/genesis-validator-1/.treasurenetd/config/

cp -a genesis.json /data/bootnode-1/.treasurenetd/config/genesis.json

sudo rm -rf /home/ubuntu/.treasurenetd

treasurenetd init bootnode-2 --chain-id "treasurenet_5005-1"

sudo mkdir -p /data/bootnode-2/.treasurenetd

sudo rm -rf /data/bootnode-2/.treasurenetd

sudo mv /home/ubuntu/.treasurenetd  /data/bootnode-2/.treasurenetd

cd /data/genesis-validator-1/.treasurenetd/config/

cp -a genesis.json /data/bootnode-2/.treasurenetd/config/genesis.json
