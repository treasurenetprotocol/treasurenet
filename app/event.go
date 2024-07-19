package app

import (
	"context"
	"log"
	"strings"
	"time"

	// "encoding/json"
	"fmt"
	"math/big"

	// store "./contract" // for demo

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	// "github.com/treasurenet/x/evm/types"
	store "github.com/treasurenetprotocol/treasurenet/app/contract"
)

type EventLog struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data [][]interface{} `json:"data"` // `json:"-"`No serialization
}

// type EventLog struct {
// 	Code int           `json:"code"`
// 	Msg  string        `json:"msg"`
// 	Data []interface{} `json:"data"` // `json:"-"`No serialization
// }
type EventLogNew struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data [][]interface{} `json:"data"` // `json:"-"`No serialization
}

var Even EventLog

var EvenNew EventLogNew

func getLogs(ctx context.Context, start, end int64) {
	data1 := make([][]interface{}, 0)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Even = new(EventLog)
			client, err := ethclient.Dial("ws://127.0.0.1:8546")
			if err != nil {
				fmt.Println("No listening:", err)
				EvenNew = EventLogNew{
					Code: 300,
					Msg:  "ethclient Listening error",
					Data: data1,
				}
			} else {
				// eventSignature := []byte("ItemSet(bytes32,bytes32)")
				// eventSignature := []byte("bidList(address,uint256,uint256)")
				// dongqi(address indexed account, uint256 indexed amount, uint256 variety);
				eventSignature := []byte("BidRecord(address,uint256)")
				hash := crypto.Keccak256Hash(eventSignature)
				topic := hash.Hex()
				fmt.Println("Test get log title new", topic)
				query := ethereum.FilterQuery{
					FromBlock: big.NewInt(start),
					ToBlock:   big.NewInt(end),
					Topics: [][]common.Hash{
						{
							hash,
						},
					},
				}
				fmt.Println("Listening start：", start)
				fmt.Println("Listening end：", end)
				logs1, err := client.FilterLogs(context.Background(), query)
				if err != nil {
					EvenNew = EventLogNew{
						Code: 400,
						Msg:  "FilterLogs Listening error",
						Data: data1,
					}
				}
				contractAbi, err := abi.JSON(strings.NewReader(string(store.ContractABI)))
				if err != nil {
					EvenNew = EventLogNew{
						Code: 500,
						Msg:  "contractAbi Listening error",
						Data: data1,
					}
				}
				fmt.Printf("contractAbi:%+v\n", contractAbi)
				if len(logs1) == 0 {
					EvenNew = EventLogNew{
						Code: 600,
						Msg:  "Log is empty",
						Data: data1,
					}
				}
				dst := make([][]interface{}, len(logs1))
				for index, vLog := range logs1 {
					LogAbi, err := contractAbi.Unpack("BidRecord", vLog.Data)
					if err != nil {
						log.Fatal(err)
					}
					dst[index] = LogAbi
				}
				EvenNew = EventLogNew{
					Code: 200,
					Msg:  "successful",
					Data: dst,
				}
				fmt.Printf("logabi:%+v\n", dst)
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func getBidStartLogsNew(ctx context.Context, start, end int64) {
	BidStartdata := make([][]interface{}, 0)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			client, err := ethclient.Dial("ws://127.0.0.1:8546")
			if err != nil {
				fmt.Println("getBidStartLogs No listening:", err)
				Even = EventLog{
					Code: 300,
					Msg:  "getBidStartLogs ethclient Listening error",
					Data: BidStartdata,
				}
			} else {
				eventSignature := []byte("BidStart(uint256)")
				hash := crypto.Keccak256Hash(eventSignature)
				topic := hash.Hex()
				fmt.Println("getBidStartLogs title new", topic)
				query := ethereum.FilterQuery{
					FromBlock: big.NewInt(start),
					ToBlock:   big.NewInt(end),
					Topics: [][]common.Hash{
						{
							hash,
						},
					},
				}
				fmt.Println("getBidStartLogs Listening start：", start)
				fmt.Println("getBidStartLogs Listening end：", end)
				logs1, err := client.FilterLogs(context.Background(), query)
				// client.Close()
				if err != nil {
					Even = EventLog{
						Code: 400,
						Msg:  "getBidStartLogs FilterLogs Listening error",
						Data: BidStartdata,
					}
				}
				contractAbi, err := abi.JSON(strings.NewReader(string(store.ContractABI)))
				if err != nil {
					Even = EventLog{
						Code: 500,
						Msg:  "getBidStartLogs contractAbi Listening error",
						Data: BidStartdata,
					}
				}
				fmt.Printf("getBidStartLogs contractAbi:%+v\n", contractAbi)
				if len(logs1) == 0 {
					Even = EventLog{
						Code: 600,
						Msg:  "getBidStartLogs is empty",
						Data: BidStartdata,
					}
				}
				dst := make([][]interface{}, len(logs1))
				for index, vLog := range logs1 {
					LogAbi, err := contractAbi.Unpack("BidStart", vLog.Data)
					if err != nil {
						log.Fatal(err)
					}
					dst[index] = LogAbi
				}
				Even = EventLog{
					Code: 200,
					Msg:  "successful",
					Data: dst,
				}
				fmt.Printf("logabi:%+v\n", dst)
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}
