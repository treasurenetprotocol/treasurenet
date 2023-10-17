package types

import (
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMintDenom           = []byte("MintDenom")
	KeyInflationRateChange = []byte("InflationRateChange")
	KeyInflationMax        = []byte("InflationMax")
	KeyInflationMin        = []byte("InflationMin")
	KeyProbability         = []byte("Probability")
	KeyGoalBonded          = []byte("GoalBonded")
	KeyUnitGrant           = []byte("UnitGrant")
	KeyBlocksPerYear       = []byte("BlocksPerYear")
	KeyStartBlock          = []byte("StartBlock")
	KeyEndBlock            = []byte("EndBlock")
	KeyHeightBlock         = []byte("HeightBlock")
	KeyPerReward           = []byte("PerReward")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string, inflationRateChange, inflationMax, inflationMin, probability, goalBonded sdk.Dec, unitGrant uint64, blocksPerYear uint64,
	startblock, endblock, heightblock int64, perreward string) Params {

	return Params{
		MintDenom:           mintDenom,
		InflationRateChange: inflationRateChange,
		InflationMax:        inflationMax,
		InflationMin:        inflationMin,
		Probability:         probability,
		GoalBonded:          goalBonded,
		UnitGrant:           unitGrant,
		BlocksPerYear:       blocksPerYear,
		StartBlock:          startblock,
		EndBlock:            endblock,
		HeightBlock:         heightblock,
		PerReward:           perreward,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:           sdk.DefaultBondDenom,
		InflationRateChange: sdk.NewDecWithPrec(13, 2),
		InflationMax:        sdk.NewDecWithPrec(20, 2),
		InflationMin:        sdk.NewDecWithPrec(7, 2),
		Probability:         sdk.NewDecWithPrec(100, 2),
		GoalBonded:          sdk.NewDecWithPrec(67, 2),
		UnitGrant:           uint64(0),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
		StartBlock:          int64(1),
		EndBlock:            int64(6311520),
		HeightBlock:         int64(12),
		PerReward:           sdk.NewIntWithDecimal(int64(10), 18).String(),
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationRateChange(p.InflationRateChange); err != nil {
		return err
	}
	if err := validateInflationMax(p.InflationMax); err != nil {
		return err
	}
	if err := validateInflationMin(p.InflationMin); err != nil {
		return err
	}
	if err := validateProbability(p.Probability); err != nil {
		return err
	}
	if err := validateGoalBonded(p.GoalBonded); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	if err := validateUnitGrant(p.UnitGrant); err != nil {
		return err
	}
	if err := validateStartBlock(p.StartBlock); err != nil {
		return err
	}
	if err := validateEndBlock(p.EndBlock); err != nil {
		return err
	}
	if err := validateHeightBlock(p.HeightBlock); err != nil {
		return err
	}
	if err := validatePerReward(p.PerReward); err != nil {
		return err
	}
	if p.InflationMax.LT(p.InflationMin) {
		return fmt.Errorf(
			"max inflation (%s) must be greater than or equal to min inflation (%s)",
			p.InflationMax, p.InflationMin,
		)
	}

	return nil

}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyInflationRateChange, &p.InflationRateChange, validateInflationRateChange),
		paramtypes.NewParamSetPair(KeyInflationMax, &p.InflationMax, validateInflationMax),
		paramtypes.NewParamSetPair(KeyInflationMin, &p.InflationMin, validateInflationMin),
		paramtypes.NewParamSetPair(KeyProbability, &p.Probability, validateProbability),
		paramtypes.NewParamSetPair(KeyGoalBonded, &p.GoalBonded, validateGoalBonded),
		paramtypes.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
		paramtypes.NewParamSetPair(KeyUnitGrant, &p.UnitGrant, validateUnitGrant),
		paramtypes.NewParamSetPair(KeyStartBlock, &p.StartBlock, validateStartBlock),
		paramtypes.NewParamSetPair(KeyEndBlock, &p.EndBlock, validateEndBlock),
		paramtypes.NewParamSetPair(KeyHeightBlock, &p.HeightBlock, validateHeightBlock),
		paramtypes.NewParamSetPair(KeyPerReward, &p.PerReward, validatePerReward),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateInflationRateChange(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation rate change cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("inflation rate change too large: %s", v)
	}

	return nil
}

func validateInflationMax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("max inflation cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("max inflation too large: %s", v)
	}

	return nil
}

func validateInflationMin(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min inflation cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("min inflation too large: %s", v)
	}

	return nil
}

func validateProbability(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("probability cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("probability too large: %s", v)
	}

	return nil
}
func validateGoalBonded(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("goal bonded cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("goal bonded too large: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}

func validateUnitGrant(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if v == 0 {
	// 	return fmt.Errorf("blocks unit grant must be positive: %d", v)
	// }

	return nil
}

func validateStartBlock(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if v == 0 {
	// 	return fmt.Errorf("blocks unit grant must be positive: %d", v)
	// }

	return nil
}

func validateEndBlock(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if v == 0 {
	// 	return fmt.Errorf("blocks unit grant must be positive: %d", v)
	// }

	return nil
}

func validateHeightBlock(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if v == 0 {
	// 	return fmt.Errorf("blocks unit grant must be positive: %d", v)
	// }

	return nil
}

func validatePerReward(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if v == 0 {
	// 	return fmt.Errorf("blocks unit grant must be positive: %d", v)
	// }

	return nil
}
