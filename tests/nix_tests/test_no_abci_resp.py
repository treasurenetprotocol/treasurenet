from pathlib import Path

import pytest

from .network import create_snapshots_dir, setup_custom_treasurenet
from .utils import memiavl_config, wait_for_block


@pytest.fixture(scope="module")
def custom_treasurenet(tmp_path_factory):
    path = tmp_path_factory.mktemp("no-abci-resp")
    yield from setup_custom_treasurenet(
        path,
        26260,
        Path(__file__).parent / "configs/discard-abci-resp.jsonnet",
    )


@pytest.fixture(scope="module")
def custom_treasurenet_rocksdb(tmp_path_factory):
    path = tmp_path_factory.mktemp("no-abci-resp-rocksdb")
    yield from setup_custom_treasurenet(
        path,
        26810,
        memiavl_config(path, "discard-abci-resp"),
        post_init=create_snapshots_dir,
        chain_binary="treasurenetd-rocksdb",
    )


@pytest.fixture(scope="module", params=["treasurenet", "treasurenet-rocksdb"])
def treasurenet_cluster(request, custom_treasurenet, custom_treasurenet_rocksdb):
    """
    run on treasurenet and
    treasurenet built with rocksdb (memIAVL + versionDB)
    """
    provider = request.param
    if provider == "treasurenet":
        yield custom_treasurenet
    elif provider == "treasurenet-rocksdb":
        yield custom_treasurenet_rocksdb
    else:
        raise NotImplementedError


def test_gas_eth_tx(treasurenet_cluster):
    """
    When node does not persist ABCI responses
    eth_gasPrice should return an error instead of crashing
    """
    wait_for_block(treasurenet_cluster.cosmos_cli(), 3)
    try:
        treasurenet_cluster.w3.eth.gas_price
        raise Exception("This query should have failed")
    except Exception as error:
        assert "node is not persisting abci responses" in error.args[0]["message"]
