#!/bin/bash

#############################################################################
# NGINX Configuration Deployment Script
#
# Purpose: Deploys and customizes NGINX configurations for multiple nodes
#          in a Treasurenet testnet environment.
#
# Actions:
# 1. Creates directories for each node
# 2. Cleans existing configurations
# 3. Copies base NGINX configurations
# 4. Customizes server names in each node's configuration
#
# Targets: node0-node3 and seednode0
#############################################################################

# Debug: Print current working directory
echo "Print current directory: $PWD"

# Create and prepare directories for each node
for dir in node{0..3} seednode0; do
    # Create directory with parent directories if needed (-p)
    sudo mkdir -p "/data/nginx1/$dir"
    
    # Remove any existing content in directory
    sudo rm -rf "/data/nginx1/$dir"/*
    
    # Copy base NGINX configurations to node directory
    # -a preserves all file attributes and permissions
    sudo cp -a nginx_backup/* "/data/nginx1/$dir"
done

# Define all nodes as an array
nodes=(node0 node1 node2 node3 seednode0)

# Customize configurations for each node
for node in "${nodes[@]}"; do
    # Generate server name based on node type
    if [[ $node == "seednode0" ]]; then
        new_server_name="seednode0.testnet.treasurenet.io"
    else
        new_server_name="${node}.testnet.treasurenet.io"
    fi

    # Replace monitoring server name in nginx.conf
    sudo sed -i "s|monitoring.node0.testnet.treasurenet.io|monitoring.$new_server_name|" \
        "/data/nginx1/$node/nginx.conf"
    
    # Replace cosmosapi server name in nginx.conf
    sudo sed -i "s|cosmosapi.node0.testnet.treasurenet.io|cosmosapi.$new_server_name|" \
        "/data/nginx1/$node/nginx.conf"
    
    # Replace metrics server name in nginx.conf
    sudo sed -i "s|tm-mtrcs.node0.testnet.treasurenet.io|tm-mtrcs.$new_server_name|" \
        "/data/nginx1/$node/nginx.conf"
done