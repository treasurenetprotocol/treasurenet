#!/bin/bash
input_file="config.toml"
output_file="modified_config.toml"

for node in genesis-validator-{1..6} bootnode-{1..2} rpc-{1..2}; do
    target_dir="/data/$node/.treasurenetd/config/"
    cd "$target_dir" || exit  # 确保目录存在

    # 同时修改 prometheus 和 laddr
    sed -E -e '/^prometheus[[:space:]]*=[[:space:]]*false$/s/false/true/' \
           -e 's|^(laddr[[:space:]]*=[[:space:]]*"tcp://)127.0.0.1:26657"|\10.0.0.0:26657"|' \
           "$input_file" > "$output_file"

    mv -- "$output_file" "$input_file"
done