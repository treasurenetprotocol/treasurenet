package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AllocateTokens handles distribution of the collected fees
// bondedVotes is a list of (validator address, validator voted on last block flag) for all
// validators in the bonded set.
func (k Keeper) AllocateTokens(
	ctx sdk.Context, sumPreviousPrecommitPower, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, bondedVotes []abci.VoteInfo,
) {
	votetatnum := int64(0)
	logger := k.Logger(ctx)

	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	// fmt.Println("feeCollectorName:", k.feeCollectorName)
	// fmt.Println("feeCollector:", feeCollector)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)
	fmt.Println("第一次feesCollected:", feesCollected)
	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	// temporary workaround to keep CanWithdrawInvariant happy
	// general discussions here: https://github.com/cosmos/cosmos-sdk/issues/2906#issuecomment-441867634
	feePool := k.GetFeePool(ctx)
	if totalPreviousPower == 0 {
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected...)
		k.SetFeePool(ctx, feePool)
		return
	}

	// calculate fraction votes  sumPreviousPrecommitPower   totalPreviousPower->Count the total voting weight through lastcommitinfo
	previousFractionVotes := sdk.NewDec(sumPreviousPrecommitPower).Quo(sdk.NewDec(totalPreviousPower))
	fmt.Println("proposer比例:", previousFractionVotes)
	// calculate previous proposer reward
	baseProposerReward := k.GetBaseProposerReward(ctx) // Proposer rate block proponents receive a fixed proportion of the current block reward as the basic reward. The default value is% 1 (set to% 30 according to project requirements)
	fmt.Println("baseProposerReward:", baseProposerReward)
	bonusProposerReward := k.GetBonusProposerReward(ctx) // Current payout proposer reward rate     Additional incentives for block proponents --> When all active verifiers vote and all votes are packaged into blocks, the maximum proportion of additional rewards that block proponents can receive is% 4 by default
	// proposermutiplierThis value is used to calculate the additional reward. The basic reward proportion of the block proponent is fixed, but the additional reward proportion is floating. This method is to calculate the sum of the proportion of the two rewards
	// proposerMultiplier = baseProposerReward+bonusProposerReward*（sumPrecommitPower/totalPower）sumPrecommitPower/totalPower  这个值就是previousFractionVotes
	proposerMultiplier := baseProposerReward.Add(bonusProposerReward.MulTruncate(previousFractionVotes))
	fmt.Println("proposerMultiplier:", proposerMultiplier)
	// Modify the proposer's reward proportion. This block reward * proposer's reward proportion (30%)
	// proposerMultiplier
	// NewproposerReward := feesCollected.MulDecTruncate(0.300000000000000000)
	proposerReward := feesCollected.MulDecTruncate(proposerMultiplier)
	fmt.Println("proposerReward:", proposerReward)
	// pay previous proposer
	remaining := feesCollected
	proposerValidator := k.stakingKeeper.ValidatorByConsAddr(ctx, previousProposer)
	fmt.Println("proposerValidator:", proposerValidator)
	if proposerValidator != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposerReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, proposerReward.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, proposerValidator.GetOperator().String()),
			),
		)

		k.AllocateTokensToValidator(ctx, proposerValidator, proposerReward)

		remaining = remaining.Sub(proposerReward)
	} else {
		// previous proposer can be unknown if say, the unbonding period is 1 block, so
		// e.g. a validator undelegates at block X, it's removed entirely by
		// block X+1's endblock, then X+2 we need to refer to the previous
		// proposer for X+1, but we've forgotten about them.
		logger.Error(fmt.Sprintf(
			"WARNING: Attempt to allocate proposer rewards to unknown proposer %s. "+
				"This should happen only if the proposer unbonded completely within a single block, "+
				"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
				"We recommend you investigate immediately.",
			previousProposer.String()))
	}

	// calculate fraction allocated to validators
	communityTax := k.GetCommunityTax(ctx) // Community tax
	tatReward := k.GetTatReward(ctx)       // Tat reward rate
	fmt.Println("tatReward:", tatReward)
	// allocate tokens proportionally to voting power
	// TODO consider parallelizing later, ref https://github.com/cosmos/cosmos-sdk/pull/3099#discussion_r246276376
	var voteMultiplier sdk.Dec
	var tattotalpower int64
	var newunitallpower int64
	var alltokenpower int64
	powerReduction := k.stakingKeeper.GetPowerReduction(ctx)
	// Tat reward allocation
	for _, vote := range bondedVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
		TatPower := validator.GetTatPower()
		newtatpower := TatPower.Int64()
		NewUnitPower := validator.GetNewUnitPower()
		NewUnitPower2 := NewUnitPower.Int64()
		TokenPower := validator.GetConsensusPower(powerReduction)
		tattotalpower += newtatpower
		newunitallpower += NewUnitPower2
		alltokenpower += TokenPower
		if newtatpower != int64(0) {
			votetatnum++
		}
	}
	fmt.Println("tattotalpower：", tattotalpower)
	fmt.Println("alltokenpower：", alltokenpower)
	fmt.Println("votetatnum++：", votetatnum)
	var votetatnumFraction sdk.Dec
	if tattotalpower != int64(0) {
		for _, vote := range bondedVotes {
			validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
			// Get the total amount of tatpower
			// params := k.stakingKeeper.GetParams(ctx)
			TatPower := validator.GetTatPower()
			newtatpower := TatPower.Int64()
			// TODO consider microslashing for missing votes.
			// ref https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
			tatpowerFraction := sdk.NewDec(newtatpower).QuoTruncate(sdk.NewDec(tattotalpower))
			votetatnumFraction = sdk.NewDec(votetatnum).QuoTruncate(sdk.NewDec(int64(len(bondedVotes))))
			fmt.Printf("tatpowerFraction:%+v\n", tatpowerFraction)
			fmt.Println("第二步bid-tat feesCollected:", feesCollected)
			fmt.Println("测试supervalidator在active validator中的比例:", votetatnumFraction)
			tatreward := feesCollected.MulDecTruncate(tatReward).MulDecTruncate(tatpowerFraction).MulDecTruncate(votetatnumFraction)
			fmt.Printf("reward:%+v\n", tatreward)
			// k.AllocateTokensToValidator(ctx, validator, tatreward)
			k.AllocateTokensToValidatorTat(ctx, validator, tatreward)
			remaining = remaining.Sub(tatreward)
		}
		voteMultiplier = sdk.OneDec().Sub(proposerMultiplier).Sub(communityTax).Sub(tatReward.MulTruncate(votetatnumFraction))
		fmt.Println("voteMultiplier:", voteMultiplier)
	} else {
		voteMultiplier = sdk.OneDec().Sub(proposerMultiplier).Sub(communityTax)
		fmt.Println("NOT BID TAT voteMultiplier:", voteMultiplier)
	}
	fmt.Println("newunitallpower：", newunitallpower)
	// Get the previous totalPreviousPower
	// totalallpower := k.stakingKeeper.GetTotalAllPower(ctx)
	// k.stakingKeeper.SetTotalAllPower(ctx, totalPreviousPower)
	// if totalPreviousPower == alltokenpower {
	for _, vote := range bondedVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
		fmt.Println("voteMultiplier:", voteMultiplier)
		// TODO consider microslashing for missing votes.
		// ref https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
		fmt.Println("vote.Validator.Power:", vote.Validator.Power)
		powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(totalPreviousPower))
		fmt.Printf("powerFraction:%+v\n", powerFraction)
		fmt.Println("不质押TAT feesCollected:", feesCollected)
		reward := feesCollected.MulDecTruncate(voteMultiplier).MulDecTruncate(powerFraction)
		fmt.Printf("reward:%+v\n", reward)
		k.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
	}
	// } else {
	// 	for _, vote := range bondedVotes {
	// 		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
	// 		fmt.Println("voteMultiplier:", voteMultiplier)
	// 		// TODO consider microslashing for missing votes.
	// 		// ref https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
	// 		//To calculate the reward of stacking, you need to discard the pledge of unit to isolate stacking and bid from each other
	// 		newunit := validator.GetNewUnitPower().Int64()
	// 		newpower := vote.Validator.Power - newunit
	// 		newtotalPreviousPower := totalPreviousPower - newunitallpower
	// 		//powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(totalPreviousPower))
	// 		powerFraction := sdk.NewDec(newpower).QuoTruncate(sdk.NewDec(newtotalPreviousPower))
	// 		reward := feesCollected.MulDecTruncate(voteMultiplier).MulDecTruncate(powerFraction)
	// 		k.AllocateTokensToValidator(ctx, validator, reward)
	// 		remaining = remaining.Sub(reward)
	// 	}
	// }
	// allocate community funding
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.SetFeePool(ctx, feePool)
}

// AllocateTokensToValidator allocate tokens to a particular validator, splitting according to commission
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commission
	commission := tokens.MulDec(val.GetCommission())
	fmt.Printf("commission:%+v\n", commission)
	shared := tokens.Sub(commission)
	// update current commission
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
	currentCommission := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), currentCommission)

	// update current rewards
	currentRewards := k.GetValidatorCurrentRewards(ctx, val.GetOperator())
	fmt.Printf("currentRewards:%+v\n", currentRewards)
	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	fmt.Printf("shared:%+v\n", shared)
	fmt.Printf("currentRewards.Rewards:%+v\n", currentRewards.Rewards)
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), currentRewards)

	// update outstanding rewards
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
	fmt.Printf("outstanding:%+v\n", outstanding)
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	fmt.Printf("outstanding.Rewards:%+v\n", outstanding.Rewards)
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)
}

// AllocateTokensToValidatorTat allocate tokens to a particular validator, splitting according to Tatreward
func (k Keeper) AllocateTokensToValidatorTat(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commission
	// update current tatreward
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTatreward,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
	currentTatreward := k.GetValidatorAccumulatedTatreward(ctx, val.GetOperator())
	currentTatreward.Tatreward = currentTatreward.Tatreward.Add(tokens...)
	k.SetValidatorAccumulatedTatreward(ctx, val.GetOperator(), currentTatreward)
}
