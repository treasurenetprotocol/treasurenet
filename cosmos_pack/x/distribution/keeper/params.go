package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// GetParams returns the total set of distribution parameters.
func (k Keeper) GetParams(clientCtx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(clientCtx, &params)
	return params
}

// SetParams sets the distribution parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetCommunityTax returns the current distribution community tax.
func (k Keeper) GetCommunityTax(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyCommunityTax, &percent)
	return percent
}

// GetBaseProposerReward returns the current distribution base proposer rate. 返回当前分发基本提案人费率。
func (k Keeper) GetBaseProposerReward(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyBaseProposerReward, &percent)
	return percent
}

// GetBonusProposerReward returns the current distribution bonus proposer reward 回当前分发奖金提议者奖励率
// rate.
func (k Keeper) GetBonusProposerReward(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyBonusProposerReward, &percent)
	return percent
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx sdk.Context) (enabled bool) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWithdrawAddrEnabled, &enabled)
	return enabled
}

// GetTatReward returns the current distribution bonus proposer reward 返回当前提议者TAT奖励率
// rate.
func (k Keeper) GetTatReward(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyTatReward, &percent)
	return percent
}
