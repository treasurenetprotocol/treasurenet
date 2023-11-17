package app

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v3/modules/core/02-client/client"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	// unnamed import of statik for swagger UI support

	_ "github.com/treasurenetprotocol/treasurenet/client/docs/statik"

	"github.com/treasurenetprotocol/treasurenet/app/ante"
	srvflags "github.com/treasurenetprotocol/treasurenet/server/flags"
	treasurenet "github.com/treasurenetprotocol/treasurenet/types"
	"github.com/treasurenetprotocol/treasurenet/x/evm"
	evmrest "github.com/treasurenetprotocol/treasurenet/x/evm/client/rest"
	evmkeeper "github.com/treasurenetprotocol/treasurenet/x/evm/keeper"
	evmtypes "github.com/treasurenetprotocol/treasurenet/x/evm/types"
	"github.com/treasurenetprotocol/treasurenet/x/feemarket"
	feemarketkeeper "github.com/treasurenetprotocol/treasurenet/x/feemarket/keeper"
	feemarkettypes "github.com/treasurenetprotocol/treasurenet/x/feemarket/types"

	// Osmosis-Labs Bech32-IBC
	"github.com/osmosis-labs/bech32-ibc/x/bech32ibc"
	bech32ibckeeper "github.com/osmosis-labs/bech32-ibc/x/bech32ibc/keeper"
	bech32ibctypes "github.com/osmosis-labs/bech32-ibc/x/bech32ibc/types"

	gravity "github.com/treasurenetprotocol/treasurenet/x/gravity"
	"github.com/treasurenetprotocol/treasurenet/x/gravity/keeper"
	gravitytypes "github.com/treasurenetprotocol/treasurenet/x/gravity/types"

	// Force-load the tracer engines to trigger registration due to Go-Ethereum v1.10.15 changes
	_ "github.com/ethereum/go-ethereum/eth/tracers/js"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
)

const appName = "treasurenetd"

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.ProposalHandler, upgradeclient.CancelProposalHandler,
			ibcclientclient.UpdateClientProposalHandler, ibcclientclient.UpgradeProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		gravity.AppModuleBasic{},
		bech32ibc.AppModuleBasic{},
		// Treasurenet modules
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner}, // used for secure addition and subtraction of balance using module account
		gravitytypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}
)

var _ simapp.App = (*TreasurenetApp)(nil)

var _ servertypes.Application = (*TreasurenetApp)(nil)

// TreasurenetApp implements an extended ABCI application. It is an application
// that may process transactions through Ethereum's EVM running atop of
// Tendermint consensus.
type TreasurenetApp struct {
	*baseapp.BaseApp

	// encoding
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper authkeeper.AccountKeeper
	// BankKeeper       bankkeeper.Keeper
	BankKeeper       bankkeeper.BaseKeeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper
	GravityKeeper    keeper.Keeper
	Bech32IbcKeeper  bech32ibckeeper.Keeper
	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	// Treasurenet keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// the configurator
	configurator module.Configurator
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultNodeHome = filepath.Join(userHomeDir, ".treasurenetd")
}

// NewTreasurenetApp returns a reference to a new initialized Treasurenet application.
func NewTreasurenetApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *TreasurenetApp {
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	// NOTE we use custom transaction decoder that supports the sdk.Tx interface instead of sdk.StdTx
	bApp := baseapp.NewBaseApp(
		appName,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		// SDK keys
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, capabilitytypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey,
		// ibc keys
		ibchost.StoreKey, ibctransfertypes.StoreKey,
		gravitytypes.StoreKey,
		bech32ibctypes.StoreKey,
		// treasurenet keys
		evmtypes.StoreKey, feemarkettypes.StoreKey,
	)

	// Add the EVM transient store key
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &TreasurenetApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// init params keeper and subspaces
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()

	// use custom Treasurenet account for contracts
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), treasurenet.ProtoAccount, maccPerms,
	)

	// app.AccountKeeper = authkeeper.NewAccountKeeper(
	// 	appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	// )

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName), &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp)

	// // register the staking hooks
	// // NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	// app.StakingKeeper = *stakingKeeper.SetHooks(
	// 	stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	// )

	app.AuthzKeeper = authzkeeper.NewKeeper(keys[authzkeeper.StoreKey], appCodec, app.BaseApp.MsgServiceRouter())

	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), stakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	// Create Treasurenet keepers
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, app.GetSubspace(feemarkettypes.ModuleName), keys[feemarkettypes.StoreKey], tkeys[feemarkettypes.TransientKey],
	)

	app.Bech32IbcKeeper = *bech32ibckeeper.NewKeeper(
		app.IBCKeeper.ChannelKeeper,
		appCodec,
		keys[bech32ibctypes.StoreKey],
		app.TransferKeeper,
	)

	app.GravityKeeper = keeper.NewKeeper(
		keys[gravitytypes.StoreKey],
		app.GetSubspace(gravitytypes.ModuleName),
		appCodec,
		&app.BankKeeper,
		&stakingKeeper,
		&app.SlashingKeeper,
		&app.DistrKeeper,
		&app.AccountKeeper,
		&app.TransferKeeper,
		&app.Bech32IbcKeeper,
	)

	// app.GravityKeeper = keeper.NewKeeper(
	// 	appCodec,
	// 	keys[gravitytypes.StoreKey],
	// 	app.GetSubspace(gravitytypes.ModuleName),
	// 	app.StakingKeeper,
	// 	app.BankKeeper,
	// 	app.SlashingKeeper,
	// )

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks(), app.GravityKeeper.Hooks()),
	)
	// app.StakingKeeper = *stakingKeeper.SetHooks(
	// 	stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	// )
	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec, keys[evmtypes.StoreKey], tkeys[evmtypes.TransientKey], app.GetSubspace(evmtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.FeeMarketKeeper,
		tracer,
	)

	// // Create IBC Keeper
	// app.IBCKeeper = ibckeeper.NewKeeper(
	// 	appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	// )

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(gravitytypes.RouterKey, keeper.NewGravityProposalHandler(app.GravityKeeper)).
		AddRoute(bech32ibctypes.RouterKey, bech32ibc.NewBech32IBCProposalHandler(app.Bech32IbcKeeper))

	govKeeper := govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter,
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		// SDK app modules
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		// auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),

		// ibc modules
		ibc.NewAppModule(app.IBCKeeper),
		transferModule,
		// Treasurenet app modules
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		gravity.NewAppModule(
			app.GravityKeeper,
			app.BankKeeper,
		),
		bech32ibc.NewAppModule(
			appCodec,
			app.Bech32IbcKeeper,
		),
		feemarket.NewAppModule(app.FeeMarketKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: upgrade module must go first to handle software upgrades.
	// NOTE: staking module is required if HistoricalEntries param > 0.
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibchost.ModuleName,
		// no-op modules
		ibctransfertypes.ModuleName,
		bech32ibctypes.ModuleName,
		gravitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
	)

	// NOTE: fee market module must go last in order to retrieve the block gas used.
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		gravitytypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		// no-op modules
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		bech32ibctypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		// SDK modules
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		ibchost.ModuleName,
		gravitytypes.ModuleName,
		// evm module denomination is used by the feemarket module, in AnteHandle
		evmtypes.ModuleName,
		bech32ibctypes.ModuleName,
		// NOTE: feemarket need to be initialized before genutil module:
		// gentx transactions use MinGasPriceDecorator.AnteHandle
		feemarkettypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// add test gRPC service for testing gRPC queries in isolation
	// testdata.RegisterTestServiceServer(app.GRPCQueryRouter(), testdata.TestServiceImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations

	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		// Use custom RandomGenesisAccounts so that auth module could create random EthAccounts in genesis state when genesis.json not specified
		auth.NewAppModule(appCodec, app.AccountKeeper, RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		transferModule,
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	// use Treasurenet's custom AnteHandler

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	anteHandler, err := ante.NewAnteHandler(ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeegrantKeeper:  app.FeeGrantKeeper,
		IBCKeeper:       app.IBCKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  maxGasWanted,
	})
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}

// Name returns the name of the App
func (app *TreasurenetApp) Name() string { return app.BaseApp.Name() }

// int64 return []byte
func Int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// BeginBlocker updates every begin block
func (app *TreasurenetApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	var delta sdk.Dec
	paramsnew := app.MintKeeper.GetParams(ctx)
	// stakingparamsnew := app.StakingKeeper.GetParams(ctx)
	// stakingparamsnew.UnbondingTime = time.Hour * 24 * 7 * 3
	// app.StakingKeeper.SetParams(ctx, stakingparamsnew)
	// if paramsnew.EndBlock < int64(6311521) {
	// 	// paramsnew.EndBlock = int64(6311520)
	// 	paramsnew.EndBlock = int64(20) + req.Header.Height
	// 	app.MintKeeper.SetParams(ctx, paramsnew)
	// }
	if paramsnew.EndBlock < req.Header.Height-int64(2) || paramsnew.EndBlock == int64(6311520) {
		// paramsnew.EndBlock = int64(6311520)
		paramsnew.EndBlock = int64(120) + req.Header.Height
		app.MintKeeper.SetParams(ctx, paramsnew)
	}
	reqnew := req.Header.Height

	if paramsnew.EndBlock == reqnew {
		nowTime := time.Now().Add(2 * time.Second)
		ctx1, _ := context.WithDeadline(context.Background(), nowTime)
		go getMintat(ctx1)
	}
	fmt.Println(" paramsnew.EndBlock: ", paramsnew.EndBlock)
	if paramsnew.EndBlock+int64(1) == reqnew && !tat.IsNil() {
		year, _ := app.MintKeeper.GettMinterYear(ctx)
		TatAll, _ := app.MintKeeper.GettMinterTatAll(ctx)
		newtatstart, _ := tat.Sub(TatAll).MarshalJSON()
		OldTatAll, _ := tat.MarshalJSON()
		app.MintKeeper.SetMinterTat(ctx, Int64ToBytes(year.Int64()), newtatstart)
		app.MintKeeper.SetMinterTatAll(ctx, OldTatAll)
	}
	if paramsnew.EndBlock+int64(2) == reqnew {
		newyear, _ := app.MintKeeper.GettMinterYear(ctx)
		fmt.Println("UNIT reduction algorithm ---  Current year : ", newyear.Int64())
		// newyearend, _ := newyear.Marshal()
		newTatEnd, _ := app.MintKeeper.GetMinterTatNew(ctx, Int64ToBytes(newyear.Int64()))
		fmt.Printf("UNIT reduction algorithm ---  Last year's TAT production :%+v\n, time :%+v\n", newTatEnd.ToDec().String(), newyear.Int64())
		newyearstart := newyear.Sub(sdk.OneInt())
		newTat, _ := app.MintKeeper.GetMinterTatNew(ctx, Int64ToBytes(newyearstart.Int64()))
		fmt.Printf("UNIT reduction algorithm ---   The year before last  TAT production :%+v\n, time: %+v\n", newTat.ToDec().String(), newyearstart.Int64())
		if newyear.Int64() > int64(1) {
			if newTatEnd.IsZero() || newTat.IsZero() {
				delta = sdk.NewDecWithPrec(50, 2)
				fmt.Printf("UNIT reduction algorithm ---  One year's TAT production was zero delta :%+v\n, time: %+v\n", delta, newyear.Int64())
				NewPerReward, _ := sdk.NewDecFromStr(paramsnew.PerReward)
				paramsnew.PerReward = delta.Mul(NewPerReward).String()
			} else {
				delta = sdk.NewDecFromInt(newTatEnd).Quo(sdk.NewDecFromInt(newTat)).Quo(sdk.NewDecWithPrec(110, 2))
				fmt.Printf("UNIT reduction algorithm ---  delta :%+v, time: %+v\n", delta, newyear.Int64())
				// if delta.RoundInt64() > int64(0) && delta.RoundInt64() < int64(1) {
				if sdk.MaxDec(delta, sdk.NewDecWithPrec(0, 2)) == delta && sdk.MinDec(delta, sdk.NewDecWithPrec(100, 2)) == delta {
					newdelta := sdk.NewDecWithPrec(90, 2).Sub(sdk.NewDecWithPrec(75, 2)).Mul(delta)
					newdeltatat := sdk.NewDecWithPrec(75, 2).Add(newdelta)
					NewPerReward, _ := sdk.NewDecFromStr(paramsnew.PerReward)
					paramsnew.PerReward = newdeltatat.Mul(NewPerReward).String()
				}
				if sdk.MaxDec(delta, sdk.NewDecWithPrec(100, 2)) == delta || sdk.NewDecWithPrec(100, 2) == delta {
					// math.Floor(Float64())
					newdelta1 := delta.Sub(sdk.OneDec())
					newdelta2, err := newdelta1.Float64()
					if err != nil {
						panic(err)
					}
					newdeltatats := sdk.MinDec(sdk.NewDecWithPrec(1, 2).Mul(sdk.MustNewDecFromStr(strconv.FormatFloat(math.Floor(newdelta2), 'f', -1, 64))).Add(sdk.NewDecWithPrec(90, 2)), sdk.NewDecWithPrec(99, 2))
					NewPerReward, _ := sdk.NewDecFromStr(paramsnew.PerReward)
					paramsnew.PerReward = newdeltatats.Mul(NewPerReward).String()
				}
			}
		}
		paramsnew.EndBlock += int64(120)
		// paramsnew.EndBlock += int64(6311520)
		app.MintKeeper.SetParams(ctx, paramsnew)
		fmt.Println("UNIT reduction algorithm ---  Current  Block rewards ： ", paramsnew.PerReward)
		year, _ := app.MintKeeper.GettMinterYear(ctx)
		yearnew, _ := year.Add(sdk.OneInt()).MarshalJSON()
		app.MintKeeper.SetMinterYear(ctx, yearnew)
	}
	return app.mm.BeginBlock(ctx, req)
}
func (app *TreasurenetApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	var msgLog sdk.ABCIMessageLog
	var msgLogs sdk.ABCIMessageLogs
	params := app.MintKeeper.GetParams(ctx)
	events := sdk.Events{sdk.NewEvent("transfer", sdk.NewAttribute("sender", "foo"))}
	// if params.HeightBlock == int64(12) {
	// 	params.HeightBlock = int64(2)
	// 	app.MintKeeper.SetParams(ctx, params)
	// }
	StartBlock := params.StartBlock
	EndBlock := params.StartBlock + params.HeightBlock
	HeightBlock := params.HeightBlock
	// EndBlock := params.StartBlock + params.HeightBlock
	// fmt.Println("loger:", app.Logger().With())
	newreq := req.Height
	if EndBlock < newreq {
		params.StartBlock = newreq
		app.MintKeeper.SetParams(ctx, params)
	}
	// if EndBlock != newreq && HeightBlock == int64(2) {
	// 	msgLog = sdk.NewABCIMessageLog(uint32(2), "  We haven't started bidding yet EndBlock != newreq height=2 ", events)
	// 	msgLogs = sdk.ABCIMessageLogs{msgLog}
	// }
	if EndBlock == newreq && HeightBlock == int64(2) {
		nowTime := time.Now().Add(2 * time.Second)
		ctx1, _ := context.WithDeadline(context.Background(), nowTime)
		// go getLogsNew(ctx1, StartBlock, EndBlock)
		go getBidStartLogsNew(ctx1, StartBlock, EndBlock)
		params.StartBlock = EndBlock
		app.MintKeeper.SetParams(ctx, params)
		msgLog = sdk.NewABCIMessageLog(uint32(3), "  We haven't started bidding yet EndBlock == newreq height=2 ", events)
		msgLogs = sdk.ABCIMessageLogs{msgLog}
	}
	if EndBlock != newreq && HeightBlock == int64(2) {
		msgLog = sdk.NewABCIMessageLog(uint32(3), "  We haven't started bidding yet EndBlock != newreq height=2 ", events)
		msgLogs = sdk.ABCIMessageLogs{msgLog}
	}
	// if EndBlock != newreq && HeightBlock == int64(60) && StartBlock+int64(1) != newreq {
	// 	msgLog = sdk.NewABCIMessageLog(uint32(2), " We haven't started a new round of bidding yet EndBlock != newreq height=60", events)
	// 	msgLogs = sdk.ABCIMessageLogs{msgLog}
	// 	// return app.mm.EndBlock(ctx, req)
	// }
	if EndBlock == newreq && HeightBlock == int64(60) {
		nowTime := time.Now().Add(2 * time.Second)
		ctx1, _ := context.WithDeadline(context.Background(), nowTime)
		go getLogs(ctx1, StartBlock, EndBlock)
		// go getLogsNew(ctx1, StartBlock, EndBlock)
		params.StartBlock = EndBlock
		app.MintKeeper.SetParams(ctx, params)
		msgLog = sdk.NewABCIMessageLog(uint32(2), " We haven't started a new round of bidding yet ", events)
		msgLogs = sdk.ABCIMessageLogs{msgLog}
		// return app.mm.EndBlock(ctx, req)
	}
	if StartBlock+int64(1) == newreq && HeightBlock != int64(60) {
		if Even.Code == 200 && len(Even.Data) > 0 {
			res1 := Even.Data[0]
			// fmt.Println(reflect.TypeOf(Even.Data[0]))
			// fmt.Println("res1:", res1[0])
			res2, _ := json.Marshal(res1[0])
			// fmt.Println("res2:", string(res2))
			heigehtnew := string(res2)
			heigehtnew1, _ := sdk.NewIntFromString(heigehtnew)
			NewHeight := heigehtnew1.Int64()
			// fmt.Println("NewHeight:", NewHeight)
			params.HeightBlock = int64(60)
			// params.StartBlock = NewHeight3 + int64(60)
			params.StartBlock = NewHeight
			app.MintKeeper.SetParams(ctx, params)
		} else {
			params.HeightBlock = int64(2)
			app.MintKeeper.SetParams(ctx, params)
		}
		msgLog = sdk.NewABCIMessageLog(uint32(3), " We haven't started bidding yetStartBlock+int64(1) == newreq height!=60 ", events)
		msgLogs = sdk.ABCIMessageLogs{msgLog}
	}
	if StartBlock+int64(1) == newreq && HeightBlock == int64(60) {
		fmt.Println("Log acquisition succeeded")
		res, _ := json.Marshal(EvenNew.Data)
		msgLog = sdk.NewABCIMessageLog(uint32(1), string(res), events)
		msgLogs = sdk.ABCIMessageLogs{msgLog}
		// return app.mm.NewEndBlock(ctx, req, msgLogs)
	}
	if StartBlock+int64(1) != newreq && HeightBlock == int64(60) {
		msgLog = sdk.NewABCIMessageLog(uint32(2), " We haven't started a new round of bidding yet ", events)
		msgLogs = sdk.ABCIMessageLogs{msgLog}
	}
	return app.mm.NewEndBlock(ctx, req, msgLogs)
}

// InitChainer updates at chain initialization
func (app *TreasurenetApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads state at a particular height
func (app *TreasurenetApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *TreasurenetApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *TreasurenetApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns TreasurenetApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *TreasurenetApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns TreasurenetApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *TreasurenetApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns TreasurenetApp's InterfaceRegistry
func (app *TreasurenetApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TreasurenetApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TreasurenetApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *TreasurenetApp) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TreasurenetApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *TreasurenetApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *TreasurenetApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)

	evmrest.RegisterTxRoutes(clientCtx, apiSvr.Router)

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

func (app *TreasurenetApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *TreasurenetApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(
	appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// SDK subspaces
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(gravitytypes.ModuleName)
	// treasurenet subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	return paramsKeeper
}
