#!/bin/bash


cd /data/node0/.treasurenetd/config/


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
mv ./$input_file /data/node1/.treasurenetd/config/$input_file
mv ./$input_file /data/node2/.treasurenetd/config/$input_file
mv ./$input_file /data/node3/.treasurenetd/config/$input_file
mv ./$input_file /data/seednode0/.treasurenetd/config/$input_file