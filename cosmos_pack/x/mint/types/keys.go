package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// MinterKey is the key to use for the keeper store.
var MinterKey = []byte{0x00}
var TatAllTokensKey = []byte{0x01}
var TatAllTokensEndKey = []byte{0x02}
var TatAllTokensYear = []byte{0x03}
var TatAllTokensAllKey = []byte{0x04}
var TatNumberKey = []byte{0x05}

const (
	// module name
	ModuleName = "mint"

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParameters       = "parameters"
	QueryInflation        = "inflation"
	QueryAnnualProvisions = "annual_provisions"
)

// GetValidatorKey creates the key for the validator with address
// VALUE: staking/Validator
func GetTatAllTokensKey(tokenyears sdk.ValAddress) []byte {
	return append(TatAllTokensKey, address.MustLengthPrefix(tokenyears)...)
}
