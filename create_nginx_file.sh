#!/bin/bash
echo "Print current directory: $PWD" # Print current directory

# mkdir
for dir in node{0..3} seednode0; do
    sudo mkdir -p "/data/nginx1/$dir"
    sudo rm -rf "/data/nginx1/$dir"/*
    sudo cp -a nginx_backup/* "/data/nginx1/$dir"
done

nodes=(node0 node1 node2 node3 seednode0)
for node in "${nodes[@]}"; do
    target_name="${node/testnet.treasurenet.io}" # 
    new_name="${node/node0/}" #
    
    # 
    [[ $node == "seednode0" ]] && 
        new_server_name="seednode0.testnet.treasurenet.io" || 
        new_server_name="${node}.testnet.treasurenet.io"

    # 
    sudo sed -i "s|monitoring.node0.testnet.treasurenet.io|monitoring.$new_server_name|" \
        "/data/nginx1/$node/nginx.conf"
    sudo sed -i "s|cosmosapi.node0.testnet.treasurenet.io|cosmosapi.$new_server_name|" \
        "/data/nginx1/$node/nginx.conf"
    sudo sed -i "s|tm-mtrcs.node0.testnet.treasurenet.io|tm-mtrcs.$new_server_name|" \
        "/data/nginx1/$node/nginx.conf"
done

