{ pkgs
, config
, ethermint ? (import ../. { inherit pkgs; })
}: rec {
  start-ethermint = pkgs.writeShellScriptBin "start-ethermint" ''
    # rely on environment to provide ethermintd
    export PATH=${pkgs.test-env}/bin:$PATH
    dotenv="$PWD/scripts/.env"
    test -f "$dotenv" || { echo "Missing $dotenv; run ./scripts/generate-test-env.sh first." >&2; exit 1; }
    ${../scripts/start-ethermint.sh} ${config.ethermint-config} "$dotenv" $@
  '';
  start-geth = pkgs.writeShellScriptBin "start-geth" ''
    export PATH=${pkgs.test-env}/bin:${pkgs.go-ethereum}/bin:$PATH
    dotenv="$PWD/scripts/.env"
    test -f "$dotenv" || { echo "Missing $dotenv; run ./scripts/generate-test-env.sh first." >&2; exit 1; }
    source "$dotenv"
    ${../scripts/start-geth.sh} ${config.geth-genesis} $@
  '';
  start-scripts = pkgs.symlinkJoin {
    name = "start-scripts";
    paths = [ start-ethermint start-geth ];
  };
}
