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

# 定义节点列表
nodes=("bootnode-1" "bootnode-2")

# 记录节点 ID 到 .env 文件
for node in "${nodes[@]}"; do
  export HOME="/data/$node"
  node_id=$(treasurenetd tendermint show-node-id)
  echo "${node}_address=$node_id" >> /data/actions-runner/_work/treasurenet/treasurenet/.github/scripts/ansible/docker/.env
done

# 恢复 HOME 环境变量
export HOME=/home/ubuntu