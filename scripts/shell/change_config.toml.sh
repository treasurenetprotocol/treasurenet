#!/bin/bash

##############################################################################
# Treasurenet Monitoring Configuration Script
#
# Purpose: Enables Prometheus monitoring across all Treasurenet nodes by:
#          - Modifying config.toml files
#          - Changing prometheus=false to prometheus=true
#          - Processing both validator nodes (node0-node3) and seednode0
#
# Usage: Run once to enable metrics collection for monitoring systems
# Warning: Modifies node configuration files directly
##############################################################################

# Define input and output filenames
input_file="config.toml"
output_file="modified_config.toml"

# Process all nodes including seed node
for node in node{0..3} seednode0; do
    # Target configuration directory for each node
    target_dir="/data/$node/.treasurenetd/config/"
    
    # Change to node's config directory
    cd $target_dir
    
    # Enable Prometheus monitoring by replacing 'false' with 'true'
    # -E flag enables extended regular expressions
    sed -E '/^prometheus[[:space:]]*=[[:space:]]*false$/s/false/true/' $input_file > $output_file
    
    # Replace original config file with modified version
    mv ./$output_file ./$input_file
done