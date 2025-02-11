sudo rm -rf /home/ubuntu/.treasurenetd

treasurenetd init seednode --chain-id "treasurenet_5005-1"

sudo mkdir -p /data/seednode/.treasurenetd

sudo rm -rf /data/seednode/.treasurenetd

sudo mv /home/ubuntu/.treasurenetd  /data/seednode/.treasurenetd

cd /data/node1/.treasurenetd/config/

cp -a genesis.json /data/seednode/.treasurenetd/config/genesis.json
