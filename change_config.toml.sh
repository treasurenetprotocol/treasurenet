#!/bin/bash
input_file="config.toml"
output_file="modified_config.toml"

# 修改循环范围为 genesis-validator-1到6 和 bootnode-1到2
for node in genesis-validator-{1..6} bootnode-{1..2}; do
    target_dir="/data/$node/.treasurenetd/config/"
    # 进入目标目录并执行替换操作
    cd "$target_dir" || exit  # 添加错误处理，目录不存在则退出
    sed -E '/^prometheus[[:space:]]*=[[:space:]]*false$/s/false/true/' "$input_file" > "$output_file"
    mv -- "$output_file" "$input_file"  # 使用 -- 避免文件名解析问题
done