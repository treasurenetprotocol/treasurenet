package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// initialize rewards for a new validator
func (k Keeper) initializeValidator(ctx sdk.Context, val stakingtypes.ValidatorI) {
	// set initial historical rewards (period 0) with reference count of 1   设置初始历史奖励（期间 0），引用计数为 1
	k.SetValidatorHistoricalRewards(ctx, val.GetOperator(), 0, types.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))

	// set current rewards (starting at period 1) 设置当前奖励（从第 1 期开始）
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), types.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))

	// set accumulated commission 设置累计的佣金收益
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), types.InitialValidatorAccumulatedCommission())

	// set outstanding rewards  设置待分配的累积收益
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), types.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})
}

// increment validator period, returning the period just ended 增加验证器周期，返回刚刚结束的周期
// 时期递增时当前时期收益信息会被转存到历史时期收益信息当中
// 这个方法首先根据累计的收益和当前抵押中的链上资产总量计算当前时期单位链上资产的累计收益  current = rewards/token
func (k Keeper) IncrementValidatorPeriod(ctx sdk.Context, val stakingtypes.ValidatorI) uint64 {
	// fetch current rewards
	rewards := k.GetValidatorCurrentRewards(ctx, val.GetOperator())

	// calculate current ratio  计算当前比例
	var current sdk.DecCoins
	if val.GetTokens().IsZero() {

		// can't calculate ratio for zero-token validators
		// ergo we instead add to the community pool
		feePool := k.GetFeePool(ctx)
		outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
		feePool.CommunityPool = feePool.CommunityPool.Add(rewards.Rewards...)
		outstanding.Rewards = outstanding.GetRewards().Sub(rewards.Rewards)
		k.SetFeePool(ctx, feePool)
		k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)

		current = sdk.DecCoins{}
	} else {
		// note: necessary to truncate so we don't allow withdrawing more rewards than owed 注意：截取操作是有必要的，防止取回的收益高于实际拥有的
		current = rewards.Rewards.QuoDecTruncate(val.GetTokens().ToDec())
		// current = rewards.Rewards.QuoDecTruncate(val.GetDelegatorShares())
	}

	// fetch historical rewards for last period
	historical := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), rewards.Period-1).CumulativeRewardRatio

	// decrement reference count
	k.decrementReferenceCount(ctx, val.GetOperator(), rewards.Period-1)

	// set new historical rewards with reference count of 1 设置新的历史奖励，引用计数为 1
	k.SetValidatorHistoricalRewards(ctx, val.GetOperator(), rewards.Period, types.NewValidatorHistoricalRewards(historical.Add(current...), 1))

	// set current rewards, incrementing period by 1 设置当前奖励，周期加 1
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), types.NewValidatorCurrentRewards(sdk.DecCoins{}, rewards.Period+1))

	return rewards.Period
}

// increment the reference count for a historical rewards value
func (k Keeper) incrementReferenceCount(ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k Keeper) decrementReferenceCount(ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	} else {
		k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
	}
}

func (k Keeper) updateValidatorSlashFraction(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	if fraction.GT(sdk.OneDec()) || fraction.IsNegative() {
		panic(fmt.Sprintf("fraction must be >=0 and <=1, current fraction: %v", fraction))
	}

	val := k.stakingKeeper.Validator(ctx, valAddr)

	// increment current period
	newPeriod := k.IncrementValidatorPeriod(ctx, val)

	// increment reference count on period we need to track
	k.incrementReferenceCount(ctx, valAddr, newPeriod)

	slashEvent := types.NewValidatorSlashEvent(newPeriod, fraction)
	height := uint64(ctx.BlockHeight())

	k.SetValidatorSlashEvent(ctx, valAddr, height, newPeriod, slashEvent)
}
