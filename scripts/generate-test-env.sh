#!/usr/bin/env bash
set -euo pipefail

output_file="${PWD}/scripts/.env"

if [[ "${1:-}" == "--output" ]]; then
  output_file="${2:?missing path after --output}"
  shift 2
fi

if [[ "${1:-}" == "--force" ]]; then
  force=true
  shift
else
  force=false
fi

if [[ $# -ne 0 ]]; then
  echo "Usage: $0 [--output PATH] [--force]" >&2
  exit 2
fi

if [[ -e "$output_file" && "$force" != true ]]; then
  echo "Refusing to overwrite $output_file; pass --force to replace it." >&2
  exit 1
fi

mkdir -p "$(dirname "$output_file")"
umask 077
program_file="$(mktemp "${TMPDIR:-/tmp}/treasurenet-test-env.XXXXXX.go")"
output_tmp="$(mktemp "${output_file}.tmp.XXXXXX")"
cleanup() {
  rm -f "$program_file" "$output_tmp"
}
trap cleanup EXIT

cat >"$program_file" <<'EOF'
package main

import (
  "crypto/rand"
  "encoding/hex"
  "fmt"
  "os"

  bip39 "github.com/tyler-smith/go-bip39"
)

func fail(err error) {
  fmt.Fprintln(os.Stderr, "could not generate local test credentials:", err)
  os.Exit(1)
}

func mnemonic() string {
  entropy, err := bip39.NewEntropy(256)
  if err != nil {
    fail(err)
  }
  value, err := bip39.NewMnemonic(entropy)
  if err != nil {
    fail(err)
  }
  return value
}

func main() {
  password := make([]byte, 24)
  if _, err := rand.Read(password); err != nil {
    fail(err)
  }

  fmt.Printf("export PASSWORD=%s\n", hex.EncodeToString(password))
  for _, name := range []string{
    "VALIDATOR1_MNEMONIC",
    "VALIDATOR2_MNEMONIC",
    "COMMUNITY_MNEMONIC",
    "SIGNER1_MNEMONIC",
    "SIGNER2_MNEMONIC",
  } {
    fmt.Printf("export %s=%q\n", name, mnemonic())
  }
}
EOF

go run "$program_file" >"$output_tmp"
chmod 600 "$output_tmp"
mv "$output_tmp" "$output_file"
trap - EXIT
rm -f "$program_file"

echo "Generated local-only test credentials at $output_file" >&2
echo "Do not commit or share this file." >&2
