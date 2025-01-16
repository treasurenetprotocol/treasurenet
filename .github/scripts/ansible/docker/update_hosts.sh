#!/bin/bash



#Declare an associative array
declare -A private_ip_mapping

#  Read node and private IP mappings from ip_mapping.json
while read -r node public_ip private_ip; do
  private_ip_mapping["$node"]="$private_ip"
done < <(jq -r 'to_entries[] | "\(.key) \(.value.public_ip) \(.value.private_ip)"' ip_mapping.json)

# Declare an array to collect all host entries
hosts_entries=()

#  Iterate over each node and generate host entries
for node in "${!private_ip_mapping[@]}"; do
  private_ip="${private_ip_mapping[$node]}"
  hostname="treasurenet-${node}"
  hosts_entries+=("$private_ip $hostname")
done

# Convert all host entries to a string separated by newlines
hosts_content=$(printf "%s\n" "${hosts_entries[@]}")

echo "(Generated hosts content):"
echo "$hosts_content"

#  /etc/hosts Backup /etc/hosts
sudo cp /etc/hosts /etc/hosts.bak

#  /etc/hosts
# Iterate over all host entries and add them to /etc/hosts
for entry in "${hosts_entries[@]}"; do
  # Check if the entry already exists
  if ! grep -q "^$entry\$" /etc/hosts; then
    echo "$entry" | sudo tee -a /etc/hosts > /dev/null
    echo "(Added): $entry"
  else
    echo " (Entry already exists): $entry"
  fi
done

echo "hosts  (Hosts file updated successfully)."
