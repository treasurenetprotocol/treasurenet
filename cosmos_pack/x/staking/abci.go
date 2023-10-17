package staking

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	k.TrackHistoricalInfo(ctx)
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	return k.BlockValidatorUpdates(ctx)
}
func NewEndBlocker(ctx sdk.Context, k keeper.Keeper, req abci.RequestEndBlock, log sdk.ABCIMessageLogs) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker) //该库提供了一个metrics包，可用于以灵活的方式检测代码、公开应用程序指标和分析运行时性能。

	return k.NewBlockValidatorUpdates(ctx, req, log)
}
