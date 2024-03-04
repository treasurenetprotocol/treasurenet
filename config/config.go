package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// bech32PrefixAccAddr defines the bech32 prefix of an account's address
	Bech32PrefixAccAddr = "treasurenet"
	// bech32PrefixAccPub defines the bech32 prefix of an account's public key
	Bech32PrefixAccPub = "treasurenetpub"
	// bech32PrefixValAddr defines the bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = "treasurenetvaloper"
	// bech32PrefixValPub defines the bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = "treasurenetvaloperpub"
	// bech32PrefixConsAddr defines the bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = "treasurenetvalcons"
	// bech32PrefixConsPub defines the bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = "treasurenetvalconspub"
)

func SetBech32Prefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}
