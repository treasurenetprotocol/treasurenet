#!/bin/bash
input_file="config.toml"
output_file="modified_config.toml"

for node in node{0..3} seednode0; do
    target_dir="/data/$node/.treasurenetd/config/"

      sed -E '/^prometheus[[:space:]]*=[[:space:]]*false$/s/false/true/' $input_file > $output_file
        mv ./$output_file ./$input_file
done
