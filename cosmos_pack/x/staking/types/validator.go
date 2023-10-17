package types

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
	"gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// TODO: Why can't we just have one string description which can be JSON by convention
	MaxMonikerLength         = 70
	MaxIdentityLength        = 3000
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

var (
	BondStatusUnspecified = BondStatus_name[int32(Unspecified)]
	BondStatusUnbonded    = BondStatus_name[int32(Unbonded)]
	BondStatusUnbonding   = BondStatus_name[int32(Unbonding)]
	BondStatusBonded      = BondStatus_name[int32(Bonded)]
)

var _ ValidatorI = Validator{}

// NewValidator constructs a new Validator
//nolint:interfacer
func NewValidator(operator sdk.ValAddress, pubKey cryptotypes.PubKey, description Description) (Validator, error) {
	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return Validator{}, err
	}

	return Validator{
		OperatorAddress: operator.String(),
		ConsensusPubkey: pkAny,
		Jailed:          false,
		Status:          Unbonded,
		Tokens:          sdk.ZeroInt(),
		DelegatorShares: sdk.ZeroDec(),
		Description:     description,
		UnbondingHeight: int64(0),
		UnbondingTime:   time.Unix(0, 0).UTC(),
		Commission:      NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		// MinSelfDelegation: sdk.OneInt(),
		MinSelfDelegation: sdk.NewIntWithDecimal(int64(158), 18),
		// MinSelfDelegation: sdk.TwoInt(),
		TatTokens:    sdk.ZeroInt(),
		NewTokens:    sdk.ZeroInt(),
		TatPower:     sdk.ZeroInt(),
		NewUnitPower: sdk.ZeroInt(),
	}, nil
}

// String implements the Stringer interface for a Validator object.
func (v Validator) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

// Validators is a collection of Validator
type Validators []Validator

func (v Validators) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// ToSDKValidators -  convenience function convert []Validator to []sdk.ValidatorI
func (v Validators) ToSDKValidators() (validators []ValidatorI) {
	for _, val := range v {
		validators = append(validators, val)
	}

	return validators
}

// Sort Validators sorts validator array in ascending operator address order
func (v Validators) Sort() {
	sort.Sort(v)
}

// Implements sort interface
func (v Validators) Len() int {
	return len(v)
}

// Implements sort interface
func (v Validators) Less(i, j int) bool {
	return bytes.Compare(v[i].GetOperator().Bytes(), v[j].GetOperator().Bytes()) == -1
}

// Implements sort interface
func (v Validators) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// ValidatorsByVotingPower implements sort.Interface for []Validator based on
// the VotingPower and Address fields.
// The validators are sorted first by their voting power (descending). Secondary index - Address (ascending).
// Copied from tendermint/types/validator_set.go
type ValidatorsByVotingPower []Validator

func (valz ValidatorsByVotingPower) Len() int { return len(valz) }

func (valz ValidatorsByVotingPower) Less(i, j int, r sdk.Int) bool {
	if valz[i].ConsensusPower(r) == valz[j].ConsensusPower(r) {
		addrI, errI := valz[i].GetConsAddr()
		addrJ, errJ := valz[j].GetConsAddr()
		// If either returns error, then return false
		if errI != nil || errJ != nil {
			return false
		}
		return bytes.Compare(addrI, addrJ) == -1
	}
	return valz[i].ConsensusPower(r) > valz[j].ConsensusPower(r)
}

func (valz ValidatorsByVotingPower) Swap(i, j int) {
	valz[i], valz[j] = valz[j], valz[i]
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v Validators) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range v {
		if err := v[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}

// return the redelegation
func MustMarshalValidator(cdc codec.BinaryCodec, validator *Validator) []byte {
	return cdc.MustMarshal(validator)
}

// unmarshal a redelegation from a store value
func MustUnmarshalValidator(cdc codec.BinaryCodec, value []byte) Validator {
	validator, err := UnmarshalValidator(cdc, value)
	if err != nil {
		panic(err)
	}

	return validator
}

// unmarshal a redelegation from a store value
func UnmarshalValidator(cdc codec.BinaryCodec, value []byte) (v Validator, err error) {
	err = cdc.Unmarshal(value, &v)
	return v, err
}

// IsBonded checks if the validator status equals Bonded
func (v Validator) IsBonded() bool {
	return v.GetStatus() == Bonded
}

// IsUnbonded checks if the validator status equals Unbonded
func (v Validator) IsUnbonded() bool {
	return v.GetStatus() == Unbonded
}

// IsUnbonding checks if the validator status equals Unbonding
func (v Validator) IsUnbonding() bool {
	return v.GetStatus() == Unbonding
}

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

func NewDescription(moniker, identity, website, securityContact, details string) Description {
	return Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: securityContact,
		Details:         details,
	}
}

// String implements the Stringer interface for a Description object.
func (d Description) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}

	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}

	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}

	if d2.SecurityContact == DoNotModifyDesc {
		d2.SecurityContact = d.SecurityContact
	}

	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return NewDescription(
		d2.Moniker,
		d2.Identity,
		d2.Website,
		d2.SecurityContact,
		d2.Details,
	).EnsureLength()
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, error) {
	if len(d.Moniker) > MaxMonikerLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid moniker length; got: %d, max: %d", len(d.Moniker), MaxMonikerLength)
	}

	if len(d.Identity) > MaxIdentityLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid identity length; got: %d, max: %d", len(d.Identity), MaxIdentityLength)
	}

	if len(d.Website) > MaxWebsiteLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid website length; got: %d, max: %d", len(d.Website), MaxWebsiteLength)
	}

	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid security contact length; got: %d, max: %d", len(d.SecurityContact), MaxSecurityContactLength)
	}

	if len(d.Details) > MaxDetailsLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid details length; got: %d, max: %d", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorUpdate(r sdk.Int) abci.ValidatorUpdate {
	tmProtoPk, err := v.TmConsPublicKey()
	if err != nil {
		panic(err)
	}

	return abci.ValidatorUpdate{
		PubKey:   tmProtoPk,
		Power:    v.ConsensusPower(r),
		TatPower: v.ConsensusTatPower(r),
	}
}
func (v Validator) ABCIBidValidatorUpdate(r sdk.Int) abci.ValidatorUpdate {
	tmProtoPk, err := v.TmConsPublicKey()
	if err != nil {
		panic(err)
	}

	return abci.ValidatorUpdate{
		PubKey:   tmProtoPk,
		Power:    v.ConsensusBidPower(r),
		TatPower: v.ConsensusTatPower(r),
	}
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Validator) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	tmProtoPk, err := v.TmConsPublicKey()
	if err != nil {
		panic(err)
	}

	return abci.ValidatorUpdate{
		PubKey:   tmProtoPk,
		Power:    0,
		TatPower: 0,
	}
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Validator) ABCIValidatorNewUpdate(r sdk.Int) abci.ValidatorUpdate {
	tmProtoPk, err := v.TmConsPublicKey()
	if err != nil {
		panic(err)
	}

	return abci.ValidatorUpdate{
		PubKey:   tmProtoPk,
		Power:    v.ConsensusNewsPower(r),
		TatPower: v.ConsensusTatPower(r),
	}
}

// SetInitialCommission attempts to set a validator's initial commission. An
// error is returned if the commission is invalid.
func (v Validator) SetInitialCommission(commission Commission) (Validator, error) {
	if err := commission.Validate(); err != nil {
		return v, err
	}

	v.Commission = commission

	return v, nil
}

// In some situations, the exchange rate becomes invalid, e.g. if
// Validator loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (v Validator) InvalidExRate() bool {
	return v.Tokens.IsZero() && v.DelegatorShares.IsPositive()
}

// calculate the token worth of provided shares
func (v Validator) TokensFromShares(shares sdk.Dec) sdk.Dec {
	// return (shares.MulInt(v.Tokens)).Quo(v.DelegatorShares)
	return (shares.Mul(v.DelegatorShares)).Quo(v.DelegatorShares)
}

// calculate the token worth of provided shares, truncated
func (v Validator) TokensFromSharesTruncated(shares sdk.Dec) sdk.Dec {
	// return (shares.MulInt(v.Tokens)).QuoTruncate(v.DelegatorShares)
	return (shares.Mul(v.DelegatorShares)).QuoTruncate(v.DelegatorShares)
}

// TokensFromSharesRoundUp returns the token worth of provided shares, rounded
// up.
func (v Validator) TokensFromSharesRoundUp(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).QuoRoundUp(v.DelegatorShares)
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the validator has no tokens.
func (v Validator) SharesFromTokens(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}
	fmt.Printf("SharesFromTokens的值:%+v\n", v.GetTokens().String())
	fmt.Printf("v.GetDelegatorShares()的值:%+v\n", v.GetDelegatorShares().String())
	// return v.GetDelegatorShares().MulInt(amt).QuoInt(v.GetTokens()), nil
	return v.GetDelegatorShares().MulInt(amt).Quo(v.GetDelegatorShares()), nil
}

// 由于slashing模块的接入当validator被惩罚时，tokens进行了衰减，导致后期delegator_shares的值和token不一致，最后staker进行委托后导致值出现小数的问题影响使用
func (v Validator) SharesFromTokensNew(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}
	fmt.Printf("SharesFromTokens的值:%+v\n", v.GetTokens().String())
	fmt.Printf("v.GetDelegatorShares()的值:%+v\n", v.GetDelegatorShares().String())
	return v.GetDelegatorShares().MulInt(amt).Quo(v.GetDelegatorShares()), nil
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the validator has no tokens.
func (v Validator) SharesFromTatTokens(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}

	return v.GetDelegatorShares().MulInt(amt).QuoInt(v.GetTatTokens()), nil
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given
// a bond amount. It returns an error if the validator has no tokens.
func (v Validator) SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}

	// return v.GetDelegatorShares().MulInt(amt).QuoTruncate(v.GetTokens().ToDec()), nil
	return v.GetDelegatorShares().MulInt(amt).QuoTruncate(v.GetDelegatorShares()), nil
}

// get the bonded tokens which the validator holds
func (v Validator) BondedTokens() sdk.Int {
	if v.IsBonded() {
		return v.Tokens
	}

	return sdk.ZeroInt()
}

// ConsensusPower gets the consensus-engine power. Aa reduction of 10^6 from
// validator tokens is applied
func (v Validator) ConsensusPower(r sdk.Int) int64 {
	if v.IsBonded() {
		return v.PotentialConsensusPower(r)
	}

	return 0
}
func (v Validator) ConsensusBidPower(r sdk.Int) int64 {
	return v.PotentialConsensusPower(r)
}

func (v Validator) ConsensusTatPower(r sdk.Int) int64 {
	if v.IsBonded() {
		return v.PotentialConsensusTatPower(r)
	}

	return 0
}

func (v Validator) ConsensusNewPower(r sdk.Int) int64 {
	if v.IsBonded() {
		return v.PotentialConsensusNewPower(r)
	}

	return 0
}

func (v Validator) ConsensusNewsPower(r sdk.Int) int64 {
	if v.IsBonded() {
		return v.PotentialConsensusNewsPower(r)
	}

	return 0
}

// PotentialConsensusPower returns the potential consensus-engine power.
func (v Validator) PotentialConsensusPower(r sdk.Int) int64 {
	return sdk.TokensToConsensusPower(v.Tokens, r)
}

func (v Validator) PotentialConsensusTatPower(r sdk.Int) int64 {
	return sdk.TokensToConsensusPower(v.TatTokens, r)
}

func (v Validator) PotentialConsensusNewPower(r sdk.Int) int64 {
	return sdk.TokensToConsensusPower(v.NewTokens, r)
}
func (v Validator) PotentialConsensusNewsPower(r sdk.Int) int64 {
	TatTokens := sdk.TokensToConsensusPower(v.Tokens, r)
	NewTokens := sdk.TokensToConsensusPower(v.NewTokens, r)
	Tokens := NewTokens + TatTokens
	return Tokens
}

// UpdateStatus updates the location of the shares within a validator
// to reflect the new status
func (v Validator) UpdateStatus(newStatus BondStatus) Validator {
	v.Status = newStatus
	return v
}

// AddTokensFromDel adds tokens to a validator
func (v Validator) AddTokensFromDel(amount sdk.Int) (Validator, sdk.Dec) {
	// calculate the shares to issue
	var issuedShares sdk.Dec
	if v.DelegatorShares.IsZero() {
		// the first delegation to a validator sets the exchange rate to one
		issuedShares = amount.ToDec()
	} else {
		shares, err := v.SharesFromTokens(amount)
		fmt.Printf("shares的值:%+v\n", shares.String())
		fmt.Printf("shares中Token的值:%+v\n", amount.String())
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}
	// fmt.Printf("v的值:%+v\n", v)
	v.Tokens = v.Tokens.Add(amount)
	v.DelegatorShares = v.DelegatorShares.Add(issuedShares)
	// 将tattokens移参数的形式存储起来，本来直接存在validator中就可以，但是setvalidator存的是validator的序列化键值对，但是因为编码问题validator的序列化大小已经确定，更改起来太麻烦，只能用参数的形式来存储
	return v, issuedShares
}
func (v Validator) AddTatTokens(amount sdk.Int) sdk.Int {
	// fmt.Printf("v的值:%+v\n", v)
	TatTokens := v.TatTokens.Add(amount)
	return TatTokens
}

func (v Validator) AddTokens(amount sdk.Int) sdk.Int {
	// fmt.Printf("v的值:%+v\n", v)
	TotalTokens := v.Tokens.Add(amount)
	return TotalTokens
}

// AddTokensFromDel adds tat tokens to a validator
func (v Validator) AddTatTokensFromDel(amount sdk.Int) (Validator, sdk.Dec) {
	// calculate the shares to issue
	var issuedShares sdk.Dec
	if v.DelegatorShares.IsZero() {
		// the first delegation to a validator sets the exchange rate to one
		issuedShares = amount.ToDec()
	} else {
		shares, err := v.SharesFromTatTokens(amount)
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}

	v.TatTokens = v.Tokens.Add(amount)
	v.DelegatorShares = v.DelegatorShares.Add(issuedShares)

	return v, issuedShares
}

// RemoveTokens removes tokens from a validator
func (v Validator) RemoveTokens(tokens sdk.Int) Validator {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}

	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}

	v.Tokens = v.Tokens.Sub(tokens)

	return v
}

// RemoveDelShares removes delegator shares from a validator.
// NOTE: because token fractions are left in the valiadator,
//       the exchange rate of future shares of this validator can increase.
func (v Validator) RemoveDelShares(delShares sdk.Dec) (Validator, sdk.Int) {
	remainingShares := v.DelegatorShares.Sub(delShares)

	var issuedTokens sdk.Int
	if remainingShares.IsZero() {
		// last delegation share gets any trimmings
		issuedTokens = v.Tokens
		v.Tokens = sdk.ZeroInt()
	} else {
		// leave excess tokens in the validator
		// however fully use all the delegator shares
		issuedTokens = v.TokensFromShares(delShares).TruncateInt()
		v.Tokens = v.Tokens.Sub(issuedTokens)

		if v.Tokens.IsNegative() {
			panic("attempting to remove more tokens than available in validator")
		}
	}

	v.DelegatorShares = remainingShares

	return v, issuedTokens
}

// MinEqual defines a more minimum set of equality conditions when comparing two
// validators.
func (v *Validator) MinEqual(other *Validator) bool {
	return v.OperatorAddress == other.OperatorAddress &&
		v.Status == other.Status &&
		v.Tokens.Equal(other.Tokens) &&
		v.DelegatorShares.Equal(other.DelegatorShares) &&
		v.Description.Equal(other.Description) &&
		v.Commission.Equal(other.Commission) &&
		v.Jailed == other.Jailed &&
		v.MinSelfDelegation.Equal(other.MinSelfDelegation) &&
		v.ConsensusPubkey.Equal(other.ConsensusPubkey) &&
		v.TatTokens.Equal(other.Tokens)

}

// Equal checks if the receiver equals the parameter
func (v *Validator) Equal(v2 *Validator) bool {
	return v.MinEqual(v2) &&
		v.UnbondingHeight == v2.UnbondingHeight &&
		v.UnbondingTime.Equal(v2.UnbondingTime)
}

func (v Validator) IsJailed() bool        { return v.Jailed }
func (v Validator) GetMoniker() string    { return v.Description.Moniker }
func (v Validator) GetStatus() BondStatus { return v.Status }
func (v Validator) GetOperator() sdk.ValAddress {
	if v.OperatorAddress == "" {
		return nil
	}
	addr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// ConsPubKey returns the validator PubKey as a cryptotypes.PubKey.
func (v Validator) ConsPubKey() (cryptotypes.PubKey, error) {
	pk, ok := v.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expecting cryptotypes.PubKey, got %T", pk)
	}

	return pk, nil

}

// TmConsPublicKey casts Validator.ConsensusPubkey to tmprotocrypto.PubKey.
func (v Validator) TmConsPublicKey() (tmprotocrypto.PublicKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return tmprotocrypto.PublicKey{}, err
	}

	tmPk, err := cryptocodec.ToTmProtoPublicKey(pk)
	if err != nil {
		return tmprotocrypto.PublicKey{}, err
	}

	return tmPk, nil
}

// GetConsAddr extracts Consensus key address
func (v Validator) GetConsAddr() (sdk.ConsAddress, error) {
	pk, ok := v.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expecting cryptotypes.PubKey, got %T", pk)
	}

	return sdk.ConsAddress(pk.Address()), nil
}

func (v Validator) GetTokens() sdk.Int       { return v.Tokens }
func (v Validator) GetTatTokens() sdk.Int    { return v.TatTokens }
func (v Validator) GetNewTokens() sdk.Int    { return v.NewTokens }
func (v Validator) GetTatPower() sdk.Int     { return v.TatPower }
func (v Validator) GetNewUnitPower() sdk.Int { return v.NewUnitPower }
func (v Validator) GetBondedTokens() sdk.Int { return v.BondedTokens() }
func (v Validator) GetConsensusPower(r sdk.Int) int64 {
	return v.ConsensusPower(r)
}
func (v Validator) GetCommission() sdk.Dec        { return v.Commission.Rate }
func (v Validator) GetMinSelfDelegation() sdk.Int { return v.MinSelfDelegation }
func (v Validator) GetDelegatorShares() sdk.Dec   { return v.DelegatorShares }

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v Validator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.ConsensusPubkey, &pk)
}
