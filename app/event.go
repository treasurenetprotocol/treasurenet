package app

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"strings"
	"time" // ★1: Added for individual timeout

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types" // ★2: Added, FilterLogs returns []types.Log
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	store "github.com/treasurenetprotocol/treasurenet/app/contract"
)

type EventLog struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []interface{}
	Err  error `json:"err,omitempty"`
}

func getEvents(ctx sdk.Context, eventSignature []byte, start, end int64) EventLog {
	data1 := make([]interface{}, 0)

	client, err := ethclient.Dial("http://127.0.0.1:8555")
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
	ctx.Logger().Info("eth logs listening window",
		"from_block", start,
		"to_block", end,
		"event_sig", string(eventSignature),
		"topic0", topic,
	)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(start),
		ToBlock:   big.NewInt(end),
		Topics: [][]common.Hash{
			{
				common.HexToHash(topic),
			},
		},
	}

	// ★1: Give FilterLogs an independent context with timeout to avoid being affected by external ctx cancellation
	var logs []types.Log
	{
		reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		lgs, err := client.FilterLogs(reqCtx, query)
		if err != nil {
			return EventLog{
				Code: 400,
				Msg:  "FilterLogs Listening error",
				Data: data1,
				Err:  err,
			}
		}
		logs = lgs
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

	// ★3: Parse event name from event signature (e.g., "BidRecord(address,uint256)" -> "BidRecord")
	sig := string(eventSignature)
	nameEnd := strings.Index(sig, "(")
	if nameEnd <= 0 {
		return EventLog{
			Code: 502,
			Msg:  "invalid event signature",
			Data: data1,
		}
	}
	eventName := sig[:nameEnd]

	ev, ok := contractAbi.Events[eventName]
	if !ok {
		return EventLog{
			Code: 503,
			Msg:  "event not found in ABI: " + eventName,
			Data: data1,
		}
	}

	if len(logs) == 0 {
		return EventLog{
			Code: 600,
			Msg:  "Log is empty",
			Data: data1,
		}
	}

	// ★4: Separate decoding for indexed / non-indexed to avoid length mismatch errors
	dst := make([]interface{}, 0, len(logs))
	for _, vLog := range logs {
		// Optional: Confirm topic0 matches the event to avoid incorrect decoding
		if len(vLog.Topics) == 0 || vLog.Topics[0] != ev.ID {
			// Not this event, skip
			continue
		}

		// Only decode non-indexed parameters (vLog.Data only contains these)
		nonIndexed := ev.Inputs.NonIndexed()
		values, err := nonIndexed.Unpack(vLog.Data)
		if err != nil {
			return EventLog{
				Code: 600,
				Msg:  "Log unpacking error (non-indexed)",
				Data: data1,
				Err:  err,
			}
		}

		// Assemble into map: indexed values from topics[1..], non-indexed from decoded values
		m := map[string]interface{}{}
		ni := 0
		ti := 1 // topics[0] is the event ID, starting from 1 are indexed parameters
		for _, in := range ev.Inputs {
			if in.Indexed {
				if ti >= len(vLog.Topics) {
					return EventLog{
						Code: 601,
						Msg:  "Log topics length insufficient for indexed params",
						Data: data1,
					}
				}
				// Keep Hex value for now; type conversion can be added later if needed
				m[in.Name] = vLog.Topics[ti].Hex()
				ti++
			} else {
				if ni >= len(values) {
					return EventLog{
						Code: 602,
						Msg:  "Log data length insufficient for non-indexed params",
						Data: data1,
					}
				}
				m[in.Name] = values[ni]
				ni++
			}
		}

		dst = append(dst, m)
	}

	if len(dst) == 0 {
		// All logs filtered out by topic verification (rare, mostly due to signature mismatch)
		return EventLog{
			Code: 204,
			Msg:  "no logs matched the event",
			Data: data1,
		}
	}

	return EventLog{
		Code: 200,
		Msg:  "successful",
		Data: dst,
	}
}

func getLogs(ctx sdk.Context, start, end int64) <-chan EventLog {
	results := make(chan EventLog, 1)
	go func() {
		defer close(results)
		eventSignature := []byte("BidRecord(address,uint256)")
		// ★5: Uniformly subtract one from start/end with non-negative protection
		s := start - 1
		e := end - 1
		if s < 0 {
			s = 0
		}
		if e < 0 {
			e = 0
		}
		Even := getEvents(ctx, eventSignature, s, e)
		results <- Even
		ctx.Logger().Info("EventLog received",
			"code", Even.Code,
			"msg", Even.Msg,
			"data_len", len(Even.Data),
			"err", Even.Err,
		)
	}()
	return results
}

func getBidStartLogsNew(ctx sdk.Context, start, end int64) <-chan EventLog {
	resultsNew := make(chan EventLog, 1)
	go func() {
		defer close(resultsNew)
		eventSignature := []byte("BidStart(uint256)")
		// ★5: Synchronized processing
		s := start - 1
		e := end - 1
		if s < 0 {
			s = 0
		}
		if e < 0 {
			e = 0
		}
		EvenNew := getEvents(ctx, eventSignature, s, e)
		resultsNew <- EvenNew
		ctx.Logger().Info("EventLog received",
			"code", EvenNew.Code,
			"msg", EvenNew.Msg,
			"data_len", len(EvenNew.Data),
			"err", EvenNew.Err,
		)
	}()
	return resultsNew
}
