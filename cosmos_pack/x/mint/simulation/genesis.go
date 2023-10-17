package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// Simulation parameter constants
const (
	Inflation           = "inflation"
	InflationRateChange = "inflation_rate_change"
	InflationMax        = "inflation_max"
	InflationMin        = "inflation_min"
	Probability         = "pro_bability"
	GoalBonded          = "goal_bonded"
)

// GenInflation randomized Inflation
func GenInflation(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationRateChange randomized InflationRateChange
func GenInflationRateChange(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationMax randomized InflationMax
func GenInflationMax(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(20, 2)
}

// GenInflationMin randomized InflationMin
func GenInflationMin(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(7, 2)
}

// GenProbability randomized GenProbability
func GenProbability(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 2)
}

// GenGoalBonded randomized GoalBonded
func GenGoalBonded(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(67, 2)
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation sdk.Dec
	var tatprobability sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	// params
	var inflationRateChange sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationRateChange, &inflationRateChange, simState.Rand,
		func(r *rand.Rand) { inflationRateChange = GenInflationRateChange(r) },
	)

	var inflationMax sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationMax, &inflationMax, simState.Rand,
		func(r *rand.Rand) { inflationMax = GenInflationMax(r) },
	)

	var inflationMin sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationMin, &inflationMin, simState.Rand,
		func(r *rand.Rand) { inflationMin = GenInflationMin(r) },
	)

	var probabiLity sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, Probability, &probabiLity, simState.Rand,
		func(r *rand.Rand) { probabiLity = GenProbability(r) },
	)

	var goalBonded sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, GoalBonded, &goalBonded, simState.Rand,
		func(r *rand.Rand) { goalBonded = GenGoalBonded(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)
	unitGrant := uint64(0)
	startBlock := int64(1)
	endBlock := int64(6511520)
	heightBlock := int64(12)
	perreward := sdk.NewIntWithDecimal(int64(10), 18).String()
	// perreward := int64(5000000000000000000)
	params := types.NewParams(mintDenom, inflationRateChange, inflationMax, inflationMin, probabiLity, goalBonded, unitGrant, blocksPerYear, startBlock, endBlock, heightBlock, perreward)

	mintGenesis := types.NewGenesisState(types.InitialMinter(inflation, tatprobability), params)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
