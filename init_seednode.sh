rm -rf /home/ubuntu/.treasurenetd

treasurenetd init seednode --chain-id "treasurenet_5005-1"

mv /home/ubuntu/.treasurenetd  /data/seednode/.treasurenetd/config/genesis.json

cd /data/node1/.treasurenetd/config/

cp -a genesis.json /data/node2/.treasurenetd/config/genesis.json
