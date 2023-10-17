package main_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/treasurenetprotocol/treasurenet/app"
	ethermintd "github.com/treasurenetprotocol/treasurenet/cmd/treasurenetd"
)

func TestInitCmd(t *testing.T) {
	rootCmd, _ := ethermintd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",          // Test the init cmd
		"etherminttest", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
		fmt.Sprintf("--%s=%s", flags.FlagChainID, "ethermint_5005-1"),
	})

	err := svrcmd.Execute(rootCmd, app.DefaultNodeHome)
	require.NoError(t, err)
}
