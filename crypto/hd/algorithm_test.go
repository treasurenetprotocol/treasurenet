package hd

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/stretchr/testify/require"

	cryptocodec "github.com/treasurenetprotocol/treasurenet/crypto/codec"
	ethermint "github.com/treasurenetprotocol/treasurenet/types"
)

func init() {
	amino := codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(amino)
}

func newTestMnemonic(t testing.TB) string {
	t.Helper()
	entropy, err := bip39.NewEntropy(256)
	require.NoError(t, err)
	mnemonic, err := bip39.NewMnemonic(entropy)
	require.NoError(t, err)
	return mnemonic
}

func TestKeyring(t *testing.T) {
	dir := t.TempDir()
	mockIn := strings.NewReader("")

	kr, err := keyring.New("foo", keyring.BackendTest, dir, mockIn, EthSecp256k1Option())
	require.NoError(t, err)

	info, err := kr.Key("foo")
	require.Error(t, err)
	require.Nil(t, info)

	mockIn.Reset("foo\n")
	info, mnemonic, err := kr.NewMnemonic("foo", keyring.English, ethermint.BIP44HDPath, keyring.DefaultBIP39Passphrase, EthSecp256k1)
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)
	require.Equal(t, "foo", info.GetName())
	require.Equal(t, "local", info.GetType().String())
	require.Equal(t, EthSecp256k1Type, info.GetAlgo())

	hdPath := ethermint.BIP44HDPath
	bz, err := EthSecp256k1.Derive()(mnemonic, keyring.DefaultBIP39Passphrase, hdPath)
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	wrongBz, err := EthSecp256k1.Derive()(mnemonic, keyring.DefaultBIP39Passphrase, "44'/60'/0'/0/")
	require.Error(t, err)
	require.Empty(t, wrongBz)

	privkey := EthSecp256k1.Generate()(bz)
	addr := common.BytesToAddress(privkey.PubKey().Address().Bytes())
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	require.NoError(t, err)
	path := hdwallet.MustParseDerivationPath(hdPath)
	account, err := wallet.Derive(path, false)
	require.NoError(t, err)
	require.Equal(t, addr.String(), account.Address.String())
}

func TestDerivation(t *testing.T) {
	mnemonic := newTestMnemonic(t)
	bz, err := EthSecp256k1.Derive()(mnemonic, keyring.DefaultBIP39Passphrase, ethermint.BIP44HDPath)
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	badBz, err := EthSecp256k1.Derive()(mnemonic, keyring.DefaultBIP39Passphrase, "m/44'/60'/0'/1/0")
	require.NoError(t, err)
	require.NotEmpty(t, badBz)
	require.NotEqual(t, bz, badBz)

	privkey := EthSecp256k1.Generate()(bz)
	badPrivKey := EthSecp256k1.Generate()(badBz)
	require.False(t, privkey.Equals(badPrivKey))

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	require.NoError(t, err)
	path := hdwallet.MustParseDerivationPath(ethermint.BIP44HDPath)
	account, err := wallet.Derive(path, false)
	require.NoError(t, err)
	badPath := hdwallet.MustParseDerivationPath("m/44'/60'/0'/1/0")
	badAccount, err := wallet.Derive(badPath, false)
	require.NoError(t, err)

	require.NotEmpty(t, account.Address.String())
	require.NotEmpty(t, badAccount.Address.String())
	require.NotEqual(t, account.Address.String(), badAccount.Address.String())
	require.Equal(t, common.BytesToAddress(privkey.PubKey().Address().Bytes()).String(), account.Address.String())
	require.Equal(t, common.BytesToAddress(badPrivKey.PubKey().Address().Bytes()).String(), badAccount.Address.String())
	require.NotEqual(t, common.BytesToAddress(privkey.PubKey().Address()).String(), badAccount.Address.String())
	require.NotEqual(t, common.BytesToAddress(badPrivKey.PubKey().Address()).String(), account.Address.Hex())
}
