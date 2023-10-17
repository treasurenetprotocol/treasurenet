package mint

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getMintat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("退出协程getMintatceshi！！！")
			return
		default:
			fmt.Println("监控中 getMintatceshi ！！！")
			// Even = new(EventLog)
			// client, err := ethclient.Dial("ws://127.0.0.1:8546")
			conn, err := ethclient.Dial("http://127.0.0.1:8555") // ws://127.0.0.1:8546
			if err != nil {
				fmt.Println("No listening:", err)
			}
			defer conn.Close()
			// contractAddress
			tokenAddress := common.HexToAddress("0x83754343cDc9dDC8A5FcDD283d0aeaF689Af6b8d")
			instance, err := NewMint(tokenAddress, conn)
			if err != nil {
				fmt.Println("NewMint is error:", err)
			}

			// call token name
			name, err := instance.Maintat(&bind.CallOpts{})
			if err != nil {
				fmt.Println("instance.Maintat error:", err)
			}
			fmt.Println("name is:", name)
			tat := sdk.NewIntFromBigInt(name)
			newtat := tat.ToDec()
			fmt.Println("newtat is:", newtat)
			NewmintedCoin := sdk.NewCoin("TAT", tat)
			fmt.Println("监听tat总量NewmintedCoin is:", NewmintedCoin)
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}
