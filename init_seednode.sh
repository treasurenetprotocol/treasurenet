sudo rm -rf /home/ubuntu/.treasurenetd

treasurenetd init seednode0 --chain-id "treasurenet_5005-1"

sudo mkdir -p /data/seednode0/.treasurenetd

sudo rm -rf /data/seednode0/.treasurenetd

sudo mv /home/ubuntu/.treasurenetd  /data/seednode0/.treasurenetd

cd /data/node0/.treasurenetd/config/

cp -a genesis.json /data/seednode0/.treasurenetd/config/genesis.json
