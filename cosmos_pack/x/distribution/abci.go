package distribution

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// BeginBlocker sets the proposer for determining distribution during endblock BeginBlock设置提议者，以确定endblock期间的分配，并为前一个block分配奖励
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// determine the total power signing the block
	//确定签署区块的总权力 通过LastCommitInfo来统计总的投票权重previousTotalPower，
	//并根据区块中的投票信息计算投票权重的和sumPreviousPrecommitPower  用来计算区块提案者的额外奖励比例
	//var previousTotalPower, previousTotalTatPower, sumPreviousPrecommitPower int64
	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		//previousTotalTatPower += voteInfo.Validator.TatPower
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}

	// TODO this is Tendermint-dependent
	// AllocateTokens()用于进行区块奖励的分发
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposerConsAddr(ctx)
		k.AllocateTokens(ctx, sumPreviousPrecommitPower, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	// record the proposer for when we payout on the next block 记录我们在下一个付款时的proposer
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}
