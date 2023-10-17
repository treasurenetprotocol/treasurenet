package app

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/treasurenetprotocol/treasurenet/encoding"
	treasurenet "github.com/treasurenetprotocol/treasurenet/types"
	evmtypes "github.com/treasurenetprotocol/treasurenet/x/evm/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/treasurenetprotocol/treasurenet/crypto/ethsecp256k1"
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// TreasurenetApp testing.
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1, // no limit
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

// Setup initializes a new TreasurenetApp. A Nop logger is set in TreasurenetApp.
func Setup(isCheckTx bool, patchGenesis func(*TreasurenetApp, simapp.GenesisState) simapp.GenesisState) *TreasurenetApp {
	db := dbm.NewMemDB()
	app := NewTreasurenetApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		if patchGenesis != nil {
			genesisState = patchGenesis(app, genesisState)
		}

		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				ChainId:         "treasurenet_8000-1",
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

// RandomGenesisAccounts is used by the auth module to create random genesis accounts in simulation when a genesis.json is not specified.
// In contrast, the default auth module's RandomGenesisAccounts implementation creates only base accounts and vestings accounts.
func RandomGenesisAccounts(simState *module.SimulationState) authtypes.GenesisAccounts {
	emptyCodeHash := crypto.Keccak256(nil)
	genesisAccs := make(authtypes.GenesisAccounts, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		bacc := authtypes.NewBaseAccountWithAddress(acc.Address)

		ethacc := &treasurenet.EthAccount{
			BaseAccount: bacc,
			CodeHash:    common.BytesToHash(emptyCodeHash).String(),
		}
		// Only consider making a vesting account once the initial bonded validator
		// set is exhausted due to needing to track DelegatedVesting.
		if !(int64(i) > simState.NumBonded && simState.Rand.Intn(100) < 50) {
			genesisAccs[i] = ethacc
			continue
		}
		initialVesting := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, simState.Rand.Int63n(simState.InitialStake)))
		var endTime int64

		startTime := simState.GenTimestamp.Unix()

		// Allow for some vesting accounts to vest very quickly while others very slowly.
		if simState.Rand.Intn(100) < 50 {
			endTime = int64(simulation.RandIntBetween(simState.Rand, int(startTime)+1, int(startTime+(60*60*24*30))))
		} else {
			endTime = int64(simulation.RandIntBetween(simState.Rand, int(startTime)+1, int(startTime+(60*60*12))))
		}

		bva := vestingtypes.NewBaseVestingAccount(bacc, initialVesting, endTime)

		if simState.Rand.Intn(100) < 50 {
			genesisAccs[i] = vestingtypes.NewContinuousVestingAccountRaw(bva, startTime)
		} else {
			genesisAccs[i] = vestingtypes.NewDelayedVestingAccountRaw(bva)
		}
		// genesisAccs[i] = ethacc
	}

	return genesisAccs
}

// RandomAccounts creates random accounts with an ethsecp256k1 private key
// TODO: replace secp256k1.GenPrivKeyFromSecret() with similar function in go-ethereum
func RandomAccounts(r *rand.Rand, n int) []simtypes.Account {
	accs := make([]simtypes.Account, n)

	for i := 0; i < n; i++ {
		// don't need that much entropy for simulation
		privkeySeed := make([]byte, 15)
		_, _ = r.Read(privkeySeed)

		prv := secp256k1.GenPrivKeyFromSecret(privkeySeed)
		ethPrv := &ethsecp256k1.PrivKey{}
		_ = ethPrv.UnmarshalAmino(prv.Bytes()) // UnmarshalAmino simply copies the bytes and assigns them to ethPrv.Key
		accs[i].PrivKey = ethPrv
		accs[i].PubKey = accs[i].PrivKey.PubKey()
		accs[i].Address = sdk.AccAddress(accs[i].PubKey.Address())

		accs[i].ConsKey = ed25519.GenPrivKeyFromSecret(privkeySeed)
	}

	return accs
}

// StateFn returns the initial application state using a genesis or the simulation parameters.
// It is a wrapper of simapp.AppStateFn to replace evm param EvmDenom with staking param BondDenom.
func StateFn(cdc codec.JSONCodec, simManager *module.SimulationManager) simtypes.AppStateFn {
	return func(r *rand.Rand, accs []simtypes.Account, config simtypes.Config,
	) (appState json.RawMessage, simAccs []simtypes.Account, chainID string, genesisTimestamp time.Time) {
		appStateFn := simapp.AppStateFn(cdc, simManager)
		appState, simAccs, chainID, genesisTimestamp = appStateFn(r, accs, config)

		rawState := make(map[string]json.RawMessage)
		err := json.Unmarshal(appState, &rawState)
		if err != nil {
			panic(err)
		}

		stakingStateBz, ok := rawState[stakingtypes.ModuleName]
		if !ok {
			panic("staking genesis state is missing")
		}

		stakingState := new(stakingtypes.GenesisState)
		cdc.MustUnmarshalJSON(stakingStateBz, stakingState)

		// we should get the BondDenom and make it the evmdenom.
		// thus simulation accounts could have positive amount of gas token.
		bondDenom := stakingState.Params.BondDenom

		evmStateBz, ok := rawState[evmtypes.ModuleName]
		if !ok {
			panic("staking genesis state is missing")
		}

		evmState := new(evmtypes.GenesisState)
		cdc.MustUnmarshalJSON(evmStateBz, evmState)

		// we should replace the EvmDenom with BondDenom
		evmState.Params.EvmDenom = bondDenom

		// change appState back
		rawState[evmtypes.ModuleName] = cdc.MustMarshalJSON(evmState)

		// replace appstate
		appState, err = json.Marshal(rawState)
		if err != nil {
			panic(err)
		}
		return appState, simAccs, chainID, genesisTimestamp
	}
}
