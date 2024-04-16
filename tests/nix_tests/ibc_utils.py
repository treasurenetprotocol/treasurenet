import json
import subprocess
from pathlib import Path
from typing import Any, Dict, List, NamedTuple

from pystarport import ports

from .network import (
    CosmosChain,
    Hermes,
    build_patched_treasurenetd,
    create_snapshots_dir,
    setup_custom_treasurenet,
)
from .utils import (
    ADDRS,
    eth_to_bech32,
    memiavl_config,
    setup_stride,
    update_treasurenet_bin,
    update_treasurenetd_and_setup_stride,
    wait_for_port,
)

# aunit IBC representation on another chain connected via channel-0.
TREASURENET_IBC_DENOM = "ibc/8EAC8061F4499F03D2D1419A3E73D346289AE9DB89CAB1486B72539572B1915E"
# uosmo IBC representation on the Treasurenet chain.
OSMO_IBC_DENOM = "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518"

RATIO = 10**10
# IBC_CHAINS_META metadata of cosmos chains to setup these for IBC tests
IBC_CHAINS_META = {
    "treasurenet": {
        "chain_name": "treasurenet_5005-1",
        "bin": "treasurenetd",
        "denom": "aunit",
    },
    "treasurenet-rocksdb": {
        "chain_name": "treasurenet_5005-1",
        "bin": "treasurenetd-rocksdb",
        "denom": "aunit",
    },
    "chainmain": {
        "chain_name": "chainmain-1",
        "bin": "chain-maind",
        "denom": "basecro",
    },
    "stride": {
        "chain_name": "stride-1",
        "bin": "strided",
        "denom": "ustrd",
    },
    "osmosis": {
        "chain_name": "osmosis-1",
        "bin": "osmosisd",
        "denom": "uosmo",
    },
    "gaia": {
        "chain_name": "cosmoshub-1",
        "bin": "gaiad",
        "denom": "uatom",
    },
}
EVM_CHAINS = ["treasurenet_5005", "chainmain-1"]


class IBCNetwork(NamedTuple):
    chains: Dict[str, Any]
    hermes: Hermes


def get_treasurenet_generator(
    tmp_path: Path,
    file: str,
    is_rocksdb: bool = False,
    stride_included: bool = False,
    custom_scenario: str | None = None,
):
    """
    setup treasurenet with custom config
    depending on the build
    """
    post_init_func = None
    if is_rocksdb:
        file = memiavl_config(tmp_path, file)
        gen = setup_custom_treasurenet(
            tmp_path,
            26710,
            Path(__file__).parent / file,
            chain_binary="treasurenetd-rocksdb",
            post_init=create_snapshots_dir,
        )
    else:
        file = f"configs/{file}.jsonnet"
        if custom_scenario:
            # build the binary modified for a custom scenario
            modified_bin = build_patched_treasurenetd(custom_scenario)
            post_init_func = update_treasurenet_bin(modified_bin)
            if stride_included:
                post_init_func = update_treasurenetd_and_setup_stride(modified_bin)
            gen = setup_custom_treasurenet(
                tmp_path,
                26700,
                Path(__file__).parent / file,
                post_init=post_init_func,
                chain_binary=modified_bin,
            )
        else:
            if stride_included:
                post_init_func = setup_stride()
            gen = setup_custom_treasurenet(
                tmp_path,
                28700,
                Path(__file__).parent / file,
                post_init=post_init_func,
            )

    return gen


def prepare_network(
    tmp_path: Path,
    file: str,
    chain_names: List[str],
    custom_scenario=None,
):
    chains_to_connect = []

    # set up the chains
    for chain in chain_names:
        meta = IBC_CHAINS_META[chain]
        chain_name = meta["chain_name"]
        chains_to_connect.append(chain_name)

        # treasurenet is the first chain
        # set it up and the relayer
        if "treasurenet" in chain_name:
            # setup treasurenet with the custom config
            # depending on the build
            gen = get_treasurenet_generator(
                tmp_path,
                file,
                "-rocksdb" in chain,
                "stride" in chain_names,
                custom_scenario,
            )
            treasurenet = next(gen)
            # wait for grpc ready
            wait_for_port(ports.grpc_port(treasurenet.base_port(0)))  # treasurenet grpc
            # setup relayer
            hermes = Hermes(tmp_path / "relayer.toml")
            chains = {"treasurenet": treasurenet}
            continue

        chain_instance = CosmosChain(tmp_path / chain_name, meta["bin"])
        # wait for grpc ready in other_chains
        wait_for_port(ports.grpc_port(chain_instance.base_port()))

        chains[chain] = chain_instance
        # pystarport (used to start the setup), by default uses ethereum
        # hd-path to create the relayers keys on hermes.
        # If this is not needed (e.g. in Cosmos chains like Stride, Osmosis, etc.)
        # then overwrite the relayer key
        if chain_name not in EVM_CHAINS:
            subprocess.run(
                [
                    "hermes",
                    "--config",
                    hermes.configpath,
                    "keys",
                    "add",
                    "--chain",
                    chain_name,
                    "--mnemonic-file",
                    tmp_path / "relayer.env",
                    "--overwrite",
                ],
                check=True,
            )

    # Nested loop to connect all chains with each other
    for i, chain_a in enumerate(chains_to_connect):
        for chain_b in chains_to_connect[i + 1 :]:
            subprocess.check_call(
                [
                    "hermes",
                    "--config",
                    hermes.configpath,
                    "create",
                    "channel",
                    "--a-port",
                    "transfer",
                    "--b-port",
                    "transfer",
                    "--a-chain",
                    chain_a,
                    "--b-chain",
                    chain_b,
                    "--new-client-connection",
                    "--yes",
                ]
            )

    treasurenet.supervisorctl("start", "relayer-demo")
    wait_for_port(hermes.port)
    yield IBCNetwork(chains, hermes)


def assert_ready(ibc):
    # wait for hermes
    output = subprocess.getoutput(
        f"curl -s -X GET 'http://127.0.0.1:{ibc.hermes.port}/state' | jq"
    )
    assert json.loads(output)["status"] == "success"


def hermes_transfer(ibc, other_chain_name="chainmain-1", other_chain_denom="basecro"):
    assert_ready(ibc)
    # chainmain-1 -> treasurenet_5005-1
    my_ibc0 = other_chain_name
    my_ibc1 = "treasurenet_5005-1"
    my_channel = "channel-0"
    dst_addr = eth_to_bech32(ADDRS["signer2"])
    src_amount = 10
    src_denom = other_chain_denom
    # dstchainid srcchainid srcportid srchannelid
    cmd = (
        f"hermes --config {ibc.hermes.configpath} tx ft-transfer "
        f"--dst-chain {my_ibc1} --src-chain {my_ibc0} --src-port transfer "
        f"--src-channel {my_channel} --amount {src_amount} "
        f"--timeout-height-offset 1000 --number-msgs 1 "
        f"--denom {src_denom} --receiver {dst_addr} --key-name relayer"
    )
    subprocess.run(cmd, check=True, shell=True)
    return src_amount


def get_balance(chain, addr, denom):
    balance = chain.cosmos_cli().balance(addr, denom)
    print("balance", balance, addr, denom)
    return balance


def get_balances(chain, addr):
    print("Addr: ", addr)
    balance = chain.cosmos_cli().balances(addr)
    print("balance", balance, addr)
    return balance
