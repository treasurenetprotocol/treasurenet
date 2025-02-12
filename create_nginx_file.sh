# echo $PWD
# sudo mkdir -p /data/nginx1/node1
# sudo mkdir -p /data/nginx1/node2
# sudo mkdir -p /data/nginx1/node3
# sudo mkdir -p /data/nginx1/node4
# sudo mkdir -p /data/nginx1/seednode
# rm -rf /data/nginx1/node1/*
# rm -rf /data/nginx1/node2/*
# rm -rf /data/nginx1/node3/*
# rm -rf /data/nginx1/node4/*
# rm -rf /data/nginx1/seednode/*
# cp -a nginx_backup/* /data/nginx1/node1
# cp -a nginx_backup/* /data/nginx1/node2
# cp -a nginx_backup/* /data/nginx1/node3
# cp -a nginx_backup/* /data/nginx1/node4
# cp -a nginx_backup/* /data/nginx1/seednode
# cd /data/nginx1/node1
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name node1.testnet.treasurenet.io;/' nginx.conf
# cd /data/nginx1/node2
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name node2.testnet.treasurenet.io;/' nginx.conf
# cd /data/nginx1/node3
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name node3.testnet.treasurenet.io;/' nginx.conf
# cd /data/nginx1/node4
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name node4.testnet.treasurenet.io;/' nginx.conf
# cd /data/nginx1/seednode
# sed -i 's/    server_name node0.testnet.treasurenet.io;/    server_name seednode.testnet.treasurenet.io;/' nginx.conf

#!/bin/bash
echo "Print current directory: $PWD" # Print current directory

# mkdir
for dir in node{0..3} seednode0; do
    sudo mkdir -p "/data/nginx1/$dir"
    sudo rm -rf "/data/nginx1/$dir"/*
    sudo cp -a nginx_backup/* "/data/nginx1/$dir"
done

# sed -i
nodes=(node0 node1 node2 node3 seednode0)
for node in "${nodes[@]}"; do
    target_name="${node/testnet.treasurenet.io}" # 
    new_name="${node/node0/}" # 
    
    # 
    [[ $node == "seednode0" ]] && 
        new_server_name="seednode0.testnet.treasurenet.io" || 
        new_server_name="${node}.testnet.treasurenet.io"

    sudo  sed -i "s|node0.testnet.treasurenet.io|$new_server_name|" \
        "/data/nginx1/$node/nginx.conf"
done
