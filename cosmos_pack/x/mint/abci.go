package mint

import (
	// "context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	// fmt.Println("mint模块开始监听event")
	// paramsnew := k.GetParams(ctx)
	// // StartBlock := paramsnew.StartBlock
	// EndBlock := paramsnew.StartBlock + paramsnew.HeightBlock
	// req := ctx.BlockHeight()
	// if EndBlock == req {
	// 	fmt.Println("mint模块开始监听req:", req)
	// 	nowTime := time.Now().Add(2 * time.Second)
	// 	ctx1, _ := context.WithDeadline(context.Background(), nowTime)
	// 	// go getLogsNew(ctx1, StartBlock, EndBlock)
	// 	go getMintat(ctx1)
	// }

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)
	// recalculate inflation rate
	totalStakingSupply := k.StakingTokenSupply(ctx)
	fmt.Printf("测试mint发币模块的params:%+v\n", params)
	fmt.Printf("测试mint发币模块的minter:%+v\n", minter)
	//fmt.Printf("totalStakingSupply = %+v\n", totalStakingSupply)
	//gasUsed := ctx.BlockGasMeter().GasConsumed()
	//mainGasUsed := ctx.GasMeter().GasConsumed()
	bondedRatio := k.BondedRatio(ctx)

	/*
			  Recalculate block rewards

		The block reward for the first two years was fixed at 5 units, i.e. 5 * 10 * 18

		From the 3rd year of the third year, the actual annual growth rate of Tat in 2nd year B = total amount of 2nd Tat / total amount of 1st tat, and the deviation between the actual growth rate and the target growth rate

		delta = B/（1+10%）； When delta > = 1, do not adjust block reward; When delta < 1, block reward in 3rd year = delta * N0 unit per block;

		That is to say, starting from the third year, divide the total amount of Tat in the first two years / the total amount of Tat in the first year by (1 + 10%). If this value > = 1, the reward of each block is still 5unit. If this value is less than 1, the reward of the block in the third year is 5unit *
	*/
	//Cumulative Tat production by last month
	//Unit reward
	AccumulateTat := params.PerReward
	fmt.Printf("AccumulateTat:%+v\n", AccumulateTat)
	// NewAcc := params.PerReward
	NewAcc, _ := sdk.NewIntFromString(AccumulateTat)
	NewBlockReward, _ := sdk.NewDecFromStr(AccumulateTat)
	fmt.Printf("NewAcc:%+v\n", NewAcc)
	fmt.Printf("NewBlockReward:%+v\n", NewBlockReward.RoundInt())
	// fmt.Printf("AccumulateTat:%+v\n", NewAcc.String())
	// NewAcc := sdk.NewInt(AccumulateTat)
	// fmt.Println("NewAcc:", NewAcc)
	// NewAccumulateTat := int(float32(AccumulateTat) * params.Probability)
	// fmt.Println(reflect.TypeOf(params.Probability))
	fmt.Printf("Probability = %+v\n", params.Probability)
	// Probability := sdk.NewDecWithPrec(1, 2)
	// NewAccumulateTat := float32(AccumulateTat) * Probability
	// fmt.Println("NewAccumulateTat:", NewAccumulateTat)
	fmt.Println("bondedRatio:", bondedRatio)                                          //Bondedratio current asset mortgage ratio in the chain
	minter.TatProbability = minter.NextProbabilityRate(params)                        //Tatprobability indicates the proportion of Tat
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)                  //The inflation field represents the annual inflation rate of the current block
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply) //AnnualProvisions It refers to the number of newly cast assets on the chain each year under the current annual inflation rate calculated according to the applicable annual inflation rate of the current block and the total amount of assets on the chain
	// NewAccumulateTat := minter.NewNextAnnualProvisions(params, NewAcc)
	// minter.NewAnnualProvisions = NewAccumulateTat
	// minter.AnnualProvisions = NewAccumulateTat //  treasuenetd query mint annual-provisions
	minter.NewAnnualProvisions = NewBlockReward
	minter.AnnualProvisions = NewBlockReward //  treasuenetd query mint annual-provisions
	// newann := NewAccumulateTat.TruncateInt64() - UnitGrant
	// fmt.Println("newann:", newann)
	// minter.NewAnnualProvisions = sdk.NewDec(AccumulateTat)
	// fmt.Printf("NewAccumulateTat = %+v\n", NewAccumulateTat)
	// NewTwoNumber := 365 * 2 * AfterWeek
	// fmt.Println("NewTwoNumber.unit64:", uint64(NewTwoNumber))
	// mint coins, update supply
	// mintedCoin := minter.BlockProvision(params)
	// mintedCoins := sdk.NewCoins(mintedCoin)
	// New chain assets
	// NewmintedCoin := sdk.NewCoin(params.MintDenom, NewAccumulateTat.TruncateInt())
	NewmintedCoin := sdk.NewCoin(params.MintDenom, NewBlockReward.RoundInt())
	//The reward of each block previously calculated according to the number of blocks produced each year is now a fixed value of 5unit
	//NewmintedCoin := minter.NewBlockProvision(params, uint64(params.BlocksPerYear))
	fmt.Println("NewmintedCoin:", NewmintedCoin)
	//The cumulative output of Tat in one year unitgrant monitors the contract to obtain Tat output when the height of the block reaches the height of one year. When it reaches the height of the second year, it monitors again. At this time, the Tat output is the total output of two years. The total output of two years - the output of the first year (unitgrant) is used to calculate the inflation rate of the third year
	//UnitGrant += NewmintedCoin.Amount.Int64()
	//params.UnitGrant = uint64(UnitGrant)
	//Update cumulative value
	k.SetParams(ctx, params)
	NewmintedCoins := sdk.NewCoins(NewmintedCoin)
	//minter.UnitGrant = sdk.NewDec(UnitGrant)
	// fmt.Println("mintedCoin", mintedCoin)
	// fmt.Println("mintedCoins", mintedCoins)
	// fmt.Println("NewmintedCoin:", NewmintedCoin)
	// fmt.Println("minter.UnitGrant:", minter.UnitGrant)
	// err := k.MintCoins(ctx, mintedCoins)
	// if err != nil {
	// 	panic(err)
	// }
	//fmt.Printf("minter(new) = %+v\n", minter)
	k.SetMinter(ctx, minter)
	//minternew := k.GetMinter(ctx)
	//fmt.Printf("minternew:%+v\n", minternew)
	err := k.MintCoins(ctx, NewmintedCoins)
	if err != nil {
		panic(err)
	}
	// send the minted coins to the fee collector account
	// err = k.AddCollectedFees(ctx, mintedCoins)
	// if err != nil {
	// 	panic(err)
	// }
	err = k.AddCollectedFees(ctx, NewmintedCoins)
	if err != nil {
		panic(err)
	}

	// if mintedCoin.Amount.IsInt64() {
	// 	defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	// }

	if NewmintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(NewmintedCoin.Amount.Int64()), "minted_tokens")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(types.AttributeKeyNewAnnualProvisions, minter.NewAnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, NewmintedCoin.Amount.String()),
		),
	)
}
