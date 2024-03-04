package client

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func TestInitConfigNonNotExistError(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "nonPerms")
	fmt.Println("tempDir is", tempDir)
	fmt.Println("subDir is", subDir)
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub directory: %v", err)
	}
	cmd := &cobra.Command{}
	cmd.PersistentFlags().String(flags.FlagHome, "", "")
	if err := cmd.PersistentFlags().Set(flags.FlagHome, subDir); err != nil {
		t.Fatalf("Could not set home flag [%T] %v", err, err)
	}
	cmd.PersistentFlags().String(flags.FlagChainID, "treasurenet_5005-1", "Specify Chain ID for sending Tx")
	err := InitConfig(cmd)
	fmt.Println("os.IsPermission is", os.IsPermission(err))
	if err := InitConfig(cmd); os.IsPermission(err) {
		t.Fatalf("Failed to catch permissions error, got: [%T] %v", err, err)
	}
}
