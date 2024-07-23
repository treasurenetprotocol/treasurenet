package app

import (
	"context"
	"fmt"

	// "log"
	"strings"
	// "time"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	store "github.com/treasurenetprotocol/treasurenet/app/contract"
)

type EventLog struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []interface{}
	Err  error `json:"err,omitempty"` // 添加错误信息字段
}

func getEvents(ctx context.Context, eventSignature []byte, start, end int64) EventLog {
	// var Even EventLog
	data1 := make([]interface{}, 0)
	client, err := ethclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		return EventLog{
			Code: 300,
			Msg:  "ethclient Listening error",
			Data: data1,
			Err:  err,
		}
	}
	defer client.Close()

	hash := crypto.Keccak256Hash(eventSignature)
	topic := hash.Hex()
	fmt.Println("Listening start：", start)
	fmt.Println("Listening end：", end)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(start),
		ToBlock:   big.NewInt(end),
		Topics: [][]common.Hash{
			{
				common.HexToHash(topic),
			},
		},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return EventLog{
			Code: 400,
			Msg:  "FilterLogs Listening error",
			Data: data1,
			Err:  err,
		}
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(store.ContractABI)))
	if err != nil {
		return EventLog{
			Code: 500,
			Msg:  "contractAbi Listening error",
			Data: data1,
			Err:  err,
		}
	}

	dst := make([]interface{}, len(logs))
	for index, vLog := range logs {
		var eventData map[string]interface{}
		err := contractAbi.UnpackIntoInterface(&eventData, "BidRecord", vLog.Data)
		if err != nil {
			return EventLog{
				Code: 600,
				Msg:  "Log unpacking error",
				Data: data1,
				Err:  err,
			}
		}
		dst[index] = eventData
	}

	if len(logs) == 0 {
		return EventLog{
			Code: 600,
			Msg:  "Log is empty",
			Data: data1,
		}
	}

	return EventLog{
		Code: 200,
		Msg:  "successful",
		Data: dst,
	}
}

func getLogs(ctx context.Context, start, end int64) <-chan EventLog {
	results := make(chan EventLog, 1)
	go func() {
		defer close(results)
		eventSignature := []byte("BidRecord(address,uint256)")
		Even := getEvents(ctx, eventSignature, start, end)
		results <- Even
		fmt.Printf("EventLog: %+v\n", Even)
	}()
	return results
}

func getBidStartLogsNew(ctx context.Context, start, end int64) <-chan EventLog {
	resultsNew := make(chan EventLog, 1)
	go func() {
		defer close(resultsNew)
		eventSignature := []byte("BidStart(uint256)")
		EvenNew := getEvents(ctx, eventSignature, start, end)
		resultsNew <- EvenNew
		fmt.Printf("EventLog: %+v\n", EvenNew)
	}()
	return resultsNew
}
