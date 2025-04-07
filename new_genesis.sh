
# 定义节点列表
nodes=("genesis-validator-1" "genesis-validator-2" "genesis-validator-3" "genesis-validator-4" "genesis-validator-5" "genesis-validator-6" "rpc-1" "rpc-2")

# 循环执行 cp 命令
for node in "${nodes[@]}"; do
  cp -auv "/data/${node}/.treasurenetd/config/gentx/"*.json \
    /data/genesis-validator-1/.treasurenetd/config/gentx/
done

# 设置 HOME 环境变量
export HOME=/data/genesis-validator-1

# JSON 文件路径
json_file="/data/account.json"

# 遍历 JSON 文件中的键并添加 genesis 账户
for key in $(jq -r 'keys_unsorted[]' "$json_file"); do
  if [[ "$key" != "validator1" && "$key" != "orchestrator1" ]]; then
    ACCOUNT=$(jq -r ".${key}" "$json_file")
    echo "Adding genesis account for $key with address $ACCOUNT"
    
    treasurenetd add-genesis-account --trace --keyring-backend test $ACCOUNT 10000000000000000000000aunit,10000000000stake,10000000000footoken,10000000000footoken2,10000000000ibc/nometadatatoken
  fi
done

# 进入 gentx 目录并收集 gentx 文件
cd /data/genesis-validator-1/.treasurenetd/config/gentx
treasurenetd collect-gentxs

# 备份并替换 config.toml 文件
cd /data/genesis-validator-1/.treasurenetd/config/
mv config.toml config.toml1
cp /data/node1/.treasurenetd/config/config.toml ./

# 记录节点 ID 到 .env 文件
for node in "${nodes[@]}"; do
  export HOME="/data/$node"
  node_id=$(treasurenetd tendermint show-node-id)
  echo "${node}_address=$node_id" >> /data/actions-runner/_work/treasurenet/treasurenet/.github/scripts/ansible/docker/.env
done

# 输出 .env 文件内容
cat /data/actions-runner/_work/treasurenet/treasurenet/.github/scripts/ansible/docker/.env
echo "Node IDs appended to .env file."

# 复制 genesis.json 文件到其他节点
for node in "${nodes[@]}"; do
  cd /data/genesis-validator-1/.treasurenetd/config/
  cp -a genesis.json "/data/${node}/.treasurenetd/config/genesis.json"
done

# 恢复 HOME 环境变量
export HOME=/home/ubuntu

# 复制 gentx 文件到其他节点
for node in "${nodes[@]}"; do
 cd /data/genesis-validator-1/.treasurenetd/config/
  cp -auv "/data/genesis-validator-1/.treasurenetd/config/gentx/"*.json \
    /data/${node}/.treasurenetd/config/gentx/
done
