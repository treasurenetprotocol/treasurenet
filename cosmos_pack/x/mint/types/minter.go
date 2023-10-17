package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(inflation, annualProvisions, tatprobability, newannualProvisions, unitGrant sdk.Dec) Minter {
	return Minter{
		Inflation:           inflation,
		AnnualProvisions:    annualProvisions,
		TatProbability:      tatprobability,
		NewAnnualProvisions: newannualProvisions,
		UnitGrant:           unitGrant,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation, tatprobability sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
		tatprobability,
		sdk.NewDec(0),
		sdk.NewDec(0),
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 13%.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(13, 2),
		sdk.NewDecWithPrec(1, 2),
	)
}

// validate minter
func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			minter.Inflation.String())
	}
	return nil
}

// NextInflationRate returns the new inflation rate for the next hour.
func (m Minter) NextInflationRate(params Params, bondedRatio sdk.Dec) sdk.Dec {
	// The target annual inflation rate is recalculated for each previsions cycle. The
	// inflation is also subject to a rate change (positive or negative) depending on
	// the distance from the desired ratio (67%). The maximum rate change possible is
	// defined to be 13% per year, however the annual inflation is capped as between
	// 7% and 20%.

	// (1 - bondedRatio/GoalBonded) * InflationRateChange
	inflationRateChangePerYear := sdk.OneDec().
		Sub(bondedRatio.Quo(params.GoalBonded)).
		Mul(params.InflationRateChange)
	inflationRateChange := inflationRateChangePerYear.Quo(sdk.NewDec(int64(params.BlocksPerYear)))

	// adjust the new annual inflation for this next cycle
	inflation := m.Inflation.Add(inflationRateChange) // note inflationRateChange may be negative
	if inflation.GT(params.InflationMax) {
		inflation = params.InflationMax
	}
	if inflation.LT(params.InflationMin) {
		inflation = params.InflationMin
	}

	return inflation
}

// 获取TAT占比
func (m Minter) NextProbabilityRate(params Params) sdk.Dec {
	probability := params.Probability
	return probability
}

// NextAnnualProvisions returns the annual provisions based on current total
// supply and inflation rate.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply sdk.Int) sdk.Dec {
	return m.Inflation.MulInt(totalSupply)
}

// 计算unit的新的产量
func (m Minter) NewNextAnnualProvisions(_ Params, totalSupply sdk.Int) sdk.Dec {
	return m.TatProbability.MulInt(totalSupply)
}

// BlockProvision returns the provisions for a block based on the annual 根据年度返回区块的规定
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerYear)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}

// BlockProvision returns the provisions for a block based on the annual
func (m Minter) NewBlockProvision(params Params, newnumber uint64) sdk.Coin {
	provisionAmt := m.NewAnnualProvisions.QuoInt(sdk.NewInt(int64(newnumber)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
