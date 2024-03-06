package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	treasurenet "github.com/treasurenetprotocol/treasurenet/types"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSetBech32Prefixes(t *testing.T) {
	config := sdk.GetConfig()

	require.Equal(t, Bech32PrefixAccAddr, config.GetBech32AccountAddrPrefix())
	require.Equal(t, Bech32PrefixAccPub, config.GetBech32AccountPubPrefix())
	require.Equal(t, Bech32PrefixValAddr, config.GetBech32ValidatorAddrPrefix())
	require.Equal(t, Bech32PrefixValPub, config.GetBech32ValidatorPubPrefix())
	require.Equal(t, Bech32PrefixConsAddr, config.GetBech32ConsensusAddrPrefix())
	require.Equal(t, Bech32PrefixConsPub, config.GetBech32ConsensusPubPrefix())
}

func TestSetCoinType(t *testing.T) {
	config := sdk.GetConfig()
	require.Equal(t, sdk.CoinType, int(config.GetCoinType()))
	require.Equal(t, sdk.FullFundraiserPath, config.GetFullBIP44Path())
}

func TestHDPath(t *testing.T) {
	params := *hd.NewFundraiserParams(0, treasurenet.Bip44CoinType, 0)
	hdPath := params.String()
	require.Equal(t, "m/44'/60'/0'/0/0", hdPath)
	require.Equal(t, hdPath, treasurenet.BIP44HDPath)
}
