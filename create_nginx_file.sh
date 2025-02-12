# echo $PWD
# sudo mkdir -p /data/ngnix/node1
# sudo mkdir -p /data/ngnix/node2
# sudo mkdir -p /data/ngnix/node3
# sudo mkdir -p /data/ngnix/node4
# sudo mkdir -p /data/ngnix/seednode
# rm -rf /data/ngnix/node1/*
# rm -rf /data/ngnix/node2/*
# rm -rf /data/ngnix/node3/*
# rm -rf /data/ngnix/node4/*
# rm -rf /data/ngnix/seednode/*
# cp -a nginx_backup/* /data/ngnix/node1
# cp -a nginx_backup/* /data/ngnix/node2
# cp -a nginx_backup/* /data/ngnix/node3
# cp -a nginx_backup/* /data/ngnix/node4
# cp -a nginx_backup/* /data/ngnix/seednode
# cd /data/ngnix/node1
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name node1.testnet.treasurenet.io;/' nginx.conf
# cd /data/ngnix/node2
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name node2.testnet.treasurenet.io;/' nginx.conf
# cd /data/ngnix/node3
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name node3.testnet.treasurenet.io;/' nginx.conf
# cd /data/ngnix/node4
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name node4.testnet.treasurenet.io;/' nginx.conf
# cd /data/ngnix/seednode
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name seednode.testnet.treasurenet.io;/' nginx.conf

#!/bin/bash
echo "Print current directory: $PWD" # Print current directory

# mkdir
for dir in node{1..4} seednode; do
    sudo mkdir -p "/data/ngnix/$dir"
    sudo rm -rf "/data/ngnix/$dir"/*
    sudo cp -a nginx_backup/* "/data/ngnix/$dir"
done

# sed -i
nodes=(node1 node2 node3 node4 seednode)
for node in "${nodes[@]}"; do
    target_name="${node/testnet.treasurenet.io}" # 
    new_name="${node/node0/}" # 
    
    # 
    [[ $node == "seednode" ]] && 
        new_server_name="seednode.testnet.treasurenet.io" || 
        new_server_name="${node}.testnet.treasurenet.io"

    sed -i "s|node0.testnet.treasurenet.io|$new_server_name|" \
        "/data/ngnix/$node/nginx.conf"
done
