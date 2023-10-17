package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         sdk.StoreKey
	paramSpace       paramtypes.Subspace
	stakingKeeper    types.StakingKeeper
	bankKeeper       types.BankKeeper
	feeCollectorName string
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	sk types.StakingKeeper, ak types.AccountKeeper, bk types.BankKeeper,
	feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace,
		stakingKeeper:    sk,
		bankKeeper:       bk,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// get the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// get the tat
func (k Keeper) GetMinterTat(ctx sdk.Context) (newtokens sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.TatAllTokensKey)
	if value == nil {
		return sdk.ZeroInt(), false
		// return NewTokens, false
	}
	// newtat, _ := sdk.NewIntFromString(string(value))
	// strtat, _ := sdk.NewIntFromString(string(value))
	err := newtokens.UnmarshalJSON(value)
	if err != nil {
		return sdk.ZeroInt(), false
	}
	// fmt.Println("newunit:", NewTokens)
	return newtokens, true
}

func (k Keeper) GetMinterTatNew(ctx sdk.Context, year []byte) (newtokens sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetTatAllTokensKey(year))
	if value == nil {
		return sdk.ZeroInt(), false
		// return NewTokens, false
	}
	// newtat, _ := sdk.NewIntFromString(string(value))
	// strtat, _ := sdk.NewIntFromString(string(value))
	err := newtokens.UnmarshalJSON(value)
	if err != nil {
		return sdk.ZeroInt(), false
	}
	// fmt.Println("newunit:", NewTokens)
	return newtokens, true
}

// get the tatend
func (k Keeper) GettMinterTatEnd(ctx sdk.Context) (newtokens sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.TatAllTokensEndKey)
	if value == nil {
		return sdk.ZeroInt(), false
		// return NewTokens, false
	}
	// newtat, _ := sdk.NewIntFromString(string(value))
	// strtat, _ := sdk.NewIntFromString(string(value))
	err := newtokens.UnmarshalJSON(value)
	if err != nil {
		return sdk.ZeroInt(), false
	}
	// fmt.Println("newunit:", NewTokens)
	return newtokens, true
}

// get the tatendAll
func (k Keeper) GettMinterTatAll(ctx sdk.Context) (newtokens sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.TatAllTokensAllKey)
	if value == nil {
		return sdk.ZeroInt(), false
		// return NewTokens, false
	}
	// newtat, _ := sdk.NewIntFromString(string(value))
	// strtat, _ := sdk.NewIntFromString(string(value))
	err := newtokens.UnmarshalJSON(value)
	if err != nil {
		return sdk.ZeroInt(), false
	}
	// fmt.Println("newunit:", NewTokens)
	return newtokens, true
}

// get the tatendAll
func (k Keeper) GettMinterTatNum(ctx sdk.Context) (newtokens sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.TatNumberKey)
	if value == nil {
		return sdk.ZeroInt(), false
		// return NewTokens, false
	}
	// newtat, _ := sdk.NewIntFromString(string(value))
	// strtat, _ := sdk.NewIntFromString(string(value))
	err := newtokens.UnmarshalJSON(value)
	if err != nil {
		return sdk.ZeroInt(), false
	}
	// fmt.Println("newunit:", NewTokens)
	return newtokens, true
}

// get the tatyear
func (k Keeper) GettMinterYear(ctx sdk.Context) (year sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.TatAllTokensYear)
	if value == nil {
		return sdk.OneInt(), false
		// return NewTokens, false
	}
	// newtat, _ := sdk.NewIntFromString(string(value))
	// strtat, _ := sdk.NewIntFromString(string(value))
	err := year.UnmarshalJSON(value)
	if err != nil {
		return sdk.OneInt(), false
	}
	// fmt.Println("newunit:", NewTokens)
	return year, true
}

// set the minter(minter实现了ProtoMarshaler interface里的所有方法，所以MustMarshal(o ProtoMarshaler) 用到了多态性 ProtoMarshaler=&minter)
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	// fmt.Printf("keeper_minter:%+v", minter)
	// fmt.Printf("k.storeKey:%+v\n", k.storeKey)
	store := ctx.KVStore(k.storeKey)
	// fmt.Printf("store:%+v\n", store)
	b := k.cdc.MustMarshal(&minter)
	// fmt.Printf("b:%+v\n", string(b))
	store.Set(types.MinterKey, b)
}

func (k Keeper) SetMinterTat(ctx sdk.Context, year, tatTokens []byte) {
	// fmt.Printf("keeper_minter:%+v", minter)
	// fmt.Printf("k.storeKey:%+v\n", k.storeKey)
	store := ctx.KVStore(k.storeKey)
	fmt.Printf("types.GetTatAllTokensKey(year):%+v\n", types.GetTatAllTokensKey(year))
	store.Set(types.GetTatAllTokensKey(year), tatTokens)
}

func (k Keeper) SetMinterTatEnd(ctx sdk.Context, tatTokens []byte) {
	// fmt.Printf("keeper_minter:%+v", minter)
	// fmt.Printf("k.storeKey:%+v\n", k.storeKey)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TatAllTokensEndKey, tatTokens)
}

func (k Keeper) SetMinterTatAll(ctx sdk.Context, tatTokens []byte) {
	// fmt.Printf("keeper_minter:%+v", minter)
	// fmt.Printf("k.storeKey:%+v\n", k.storeKey)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TatAllTokensAllKey, tatTokens)
}

func (k Keeper) SetMinterTatNum(ctx sdk.Context, tatTokens []byte) {
	// fmt.Printf("keeper_minter:%+v", minter)
	// fmt.Printf("k.storeKey:%+v\n", k.storeKey)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TatNumberKey, tatTokens)
}

func (k Keeper) SetMinterYear(ctx sdk.Context, year []byte) {
	// fmt.Printf("keeper_minter:%+v", minter)
	// fmt.Printf("k.storeKey:%+v\n", k.storeKey)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TatAllTokensYear, year)
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx sdk.Context) sdk.Int {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker. 要在 BeginBlocker 中使用的 MintCoin
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's  AddCollectedFees 实现对底层供应保持者的别名调用
// AddCollectedFees to be used in BeginBlocker. AddCollectedFees 用于 BeginBlocker。
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}
