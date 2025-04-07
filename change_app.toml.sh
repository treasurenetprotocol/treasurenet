#!/bin/bash


cd /data/genesis-validator-1/.treasurenetd/config/


input_file="app.toml"
output_file="modified_app.toml"


sed -E '

/^\[telemetry\]/,/^\[/ {
    s/^(enabled) = false/\1 = true/
    s/^(enable-hostname) = false/\1 = true/
    s/^(enable-hostname-label) = false/\1 = true/
    s/^(enable-service-label) = false/\1 = true/
    s/^(prometheus-retention-time) = 0/\1 = 100/
}


/^\[api\]/,/^\[/ {
    s/^(enable) = false/\1 = true/
    s/^(swagger) = false/\1 = true/
    s/^(enabled-unsafe-cors) = false/\1 = true/
}' "$input_file" > "$output_file"


mv ./$output_file ./$input_file
cp ./$input_file /data/genesis-validator-2/.treasurenetd/config/$input_file
cp ./$input_file /data/genesis-validator-3/.treasurenetd/config/$input_file
cp ./$input_file /data/genesis-validator-4/.treasurenetd/config/$input_file
cp ./$input_file /data/genesis-validator-5/.treasurenetd/config/$input_file
cp ./$input_file /data/genesis-validator-6/.treasurenetd/config/$input_file
cp ./$input_file /data/rpc-1/.treasurenetd/config/$input_file
cp ./$input_file /data/rpc-2/.treasurenetd/config/$input_file
cp ./$input_file /data/bootnode-1/.treasurenetd/config/$input_file
cp ./$input_file /data/bootnode-2/.treasurenetd/config/$input_file