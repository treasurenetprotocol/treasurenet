package app

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	store "github.com/treasurenetprotocol/treasurenet/app/contract"
)

type EventLog struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data [][]interface{} `json:"data"`
}

type EventLogNew struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data [][]interface{} `json:"data"`
}

func getLogs(ctx context.Context, start, end int64) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := fetchLogs(ctx, start, end, "BidRecord"); err != nil {
				log.Println("Error fetching logs:", err)
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func getBidStartLogsNew(ctx context.Context, start, end int64) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := fetchLogs(ctx, start, end, "BidStart"); err != nil {
				log.Println("Error fetching BidStart logs:", err)
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func fetchLogs(ctx context.Context, start, end int64, eventName string) error {
	client, err := ethclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		return fmt.Errorf("ethclient Listening error: %w", err)
	}
	defer client.Close()

	eventSignature := []byte(eventName)
	hash := crypto.Keccak256Hash(eventSignature)
	topic := hash.Hex()
	fmt.Println("Event title new", eventName, topic)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(start),
		ToBlock:   big.NewInt(end),
		Topics: [][]common.Hash{
			{
				hash,
			},
		},
	}
	fmt.Println(eventName, "Listening start：", start)
	fmt.Println(eventName, "Listening end：", end)

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("%s FilterLogs Listening error: %w", eventName, err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(store.ContractABI)))
	if err != nil {
		return fmt.Errorf("%s contractAbi Listening error: %w", eventName, err)
	}
	fmt.Printf("%s contractAbi:%+v\n", eventName, contractAbi)

	if len(logs) == 0 {
		return fmt.Errorf("%s Log is empty", eventName)
	}

	dst := make([][]interface{}, len(logs))
	for index, vLog := range logs {
		LogAbi, err := contractAbi.Unpack(eventName, vLog.Data)
		if err != nil {
			return fmt.Errorf("unpacking %s log data error: %w", eventName, err)
		}
		dst[index] = LogAbi
	}

	eventLog := EventLog{
		Code: 200,
		Msg:  "successful",
		Data: dst,
	}
	fmt.Printf("%s logabi:%+v\n", eventName, dst)

	return nil
}
