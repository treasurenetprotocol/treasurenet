import pytest

from .network import setup_treasurenet, setup_treasurenet_rocksdb, setup_geth


@pytest.fixture(scope="session")
def treasurenet(tmp_path_factory):
    path = tmp_path_factory.mktemp("treasurenet")
    yield from setup_treasurenet(path, 26650)


@pytest.fixture(scope="session")
def treasurenet_rocksdb(tmp_path_factory):
    path = tmp_path_factory.mktemp("treasurenet-rocksdb")
    yield from setup_treasurenet_rocksdb(path, 20650)


@pytest.fixture(scope="session")
def geth(tmp_path_factory):
    path = tmp_path_factory.mktemp("geth")
    yield from setup_geth(path, 8545)


@pytest.fixture(scope="session", params=["treasurenet", "treasurenet-ws"])
def treasurenet_rpc_ws(request, treasurenet):
    """
    run on both treasurenet and treasurenet websocket
    """
    provider = request.param
    if provider == "treasurenet":
        yield treasurenet
    elif provider == "treasurenet-ws":
        treasurenet_ws = treasurenet.copy()
        treasurenet_ws.use_websocket()
        yield treasurenet_ws
    else:
        raise NotImplementedError


@pytest.fixture(scope="module", params=["treasurenet", "treasurenet-ws", "treasurenet-rocksdb", "geth"])
def cluster(request, treasurenet, treasurenet_rocksdb, geth):
    """
    run on treasurenet, treasurenet websocket,
    treasurenet built with rocksdb (memIAVL + versionDB)
    and geth
    """
    provider = request.param
    if provider == "treasurenet":
        yield treasurenet
    elif provider == "treasurenet-ws":
        treasurenet_ws = treasurenet.copy()
        treasurenet_ws.use_websocket()
        yield treasurenet_ws
    elif provider == "geth":
        yield geth
    elif provider == "treasurenet-rocksdb":
        yield treasurenet_rocksdb
    else:
        raise NotImplementedError


@pytest.fixture(scope="module", params=["treasurenet", "treasurenet-rocksdb"])
def treasurenet_cluster(request, treasurenet, treasurenet_rocksdb):
    """
    run on treasurenet default build &
    treasurenet with rocksdb build and memIAVL + versionDB
    """
    provider = request.param
    if provider == "treasurenet":
        yield treasurenet
    elif provider == "treasurenet-rocksdb":
        yield treasurenet_rocksdb
    else:
        raise NotImplementedError
