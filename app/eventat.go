package app

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	mintat "github.com/treasurenetprotocol/treasurenet/app/contractmint"
)

var newtat sdk.Dec

var tat sdk.Int

func getMintat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Even = new(EventLog)
			// client, err := ethclient.Dial("ws://127.0.0.1:8546")
			conn, err := ethclient.Dial("http://127.0.0.1:8555") // ws://127.0.0.1:8546
			if err != nil {
				fmt.Println("No listening:", err)
			}
			defer conn.Close()
			// contractAddress
			// tokenAddress := common.HexToAddress("0x83754343cDc9dDC8A5FcDD283d0aeaF689Af6b8d")
			// tokenAddress := common.HexToAddress("0xeC9B8297aa88603004c1aB91b10B0220C704BcC3")
			tokenAddress := common.HexToAddress("0x465C5ed965692F850f0a3Df1aA29955953a53714")
			instance, err := mintat.NewContract(tokenAddress, conn)
			if err != nil {
				fmt.Println("NewMint is error:", err)
			}

			// call token name
			// name, err := instance.Maintat(&bind.CallOpts{})
			name, err := instance.TotalSupply(&bind.CallOpts{})
			if err != nil {
				fmt.Println("instance.Maintat error:", err)
			} else {
				fmt.Println("name is:", name)
				tat = sdk.NewIntFromBigInt(name)
				newtat = tat.ToDec()
				fmt.Println("newtat is:", newtat)
				NewmintedCoin := sdk.NewCoin("TAT", tat)
				fmt.Println("tatvall NewmintedCoin is:", NewmintedCoin)
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}
