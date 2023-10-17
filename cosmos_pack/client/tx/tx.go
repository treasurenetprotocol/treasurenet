package tx

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/spf13/pflag"

	"github.com/treasurenetprotocol/treasurenet/crypto/ethsecp256k1"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// GenerateOrBroadcastTxCLI will either generate and print and unsigned transaction
// or sign it and broadcast it returning an error upon failure.
func GenerateOrBroadcastTxCLI(clientCtx client.Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg) error {
	txf := NewFactoryCLI(clientCtx, flagSet)
	return GenerateOrBroadcastTxWithFactory(clientCtx, txf, msgs...)
}

// GenerateOrBroadcastTxWithFactory will either generate and print and unsigned transaction
// or sign it and broadcast it returning an error upon failure.
func GenerateOrBroadcastTxWithFactory(clientCtx client.Context, txf Factory, msgs ...sdk.Msg) error {
	// Validate all msgs before generating or broadcasting the tx.
	// We were calling ValidateBasic separately in each CLI handler before.
	// Right now, we're factorizing that call inside this function.
	// ref: https://github.com/cosmos/cosmos-sdk/pull/9236#discussion_r623803504
	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
	}

	if clientCtx.GenerateOnly {
		return GenerateTx(clientCtx, txf, msgs...)
	}

	return BroadcastTx(clientCtx, txf, msgs...)
}

// GenerateTx will generate an unsigned transaction and print it to the writer
// specified by ctx.Output. If simulation was requested, the gas will be
// simulated and also printed to the same writer before the transaction is
// printed.
func GenerateTx(clientCtx client.Context, txf Factory, msgs ...sdk.Msg) error {
	if txf.SimulateAndExecute() {
		if clientCtx.Offline {
			return errors.New("cannot estimate gas in offline mode")
		}

		_, adjusted, err := CalculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return err
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	tx, err := BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return err
	}

	json, err := clientCtx.TxConfig.TxJSONEncoder()(tx.GetTx())
	if err != nil {
		return err
	}

	return clientCtx.PrintString(fmt.Sprintf("%s\n", json))
}

// BroadcastTx attempts to generate, sign and broadcast a transaction with the
// given set of messages. It will also simulate gas requirements if necessary.
// It will return an error upon failure.
func BroadcastTx(clientCtx client.Context, txf Factory, msgs ...sdk.Msg) error {
	txf, err := prepareFactory(clientCtx, txf)
	if err != nil {
		return err
	}

	if txf.SimulateAndExecute() || clientCtx.Simulate {
		_, adjusted, err := CalculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return err
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	if clientCtx.Simulate {
		return nil
	}

	tx, err := BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return err
	}

	if !clientCtx.SkipConfirm {
		out, err := clientCtx.TxConfig.TxJSONEncoder()(tx.GetTx())
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}

	tx.SetFeeGranter(clientCtx.GetFeeGranterAddress())
	err = Sign(txf, clientCtx.GetFromName(), tx, true)
	if err != nil {
		return err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)
}

func WritePubkeyTxResponse(
	clientCtx client.Context, w http.ResponseWriter, br rest.BaseReq, dr sdk.AccAddress, msgs ...sdk.Msg,
) {
	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	gasSetting, err := flags.ParseGasSetting(br.Gas)
	if rest.CheckBadRequestError(w, err) {
		return
	}

	txf := Factory{fees: br.Fees, gasPrices: br.GasPrices}.
		WithAccountNumber(br.AccountNumber).
		WithSequence(br.Sequence).
		WithGas(gasSetting.Gas).
		WithGasAdjustment(gasAdj).
		WithMemo(br.Memo).
		WithChainID(br.ChainID).
		WithSimulateAndExecute(br.Simulate).
		WithTxConfig(clientCtx.TxConfig).
		WithTimeoutHeight(br.TimeoutHeight)

	if br.Simulate || gasSetting.Simulate {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkerrors.ErrorInvalidGasAdjustment.Error())
			return
		}

		_, adjusted, err := CalculateGas(clientCtx, txf, msgs...)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		txf = txf.WithGas(adjusted)

		if br.Simulate {
			rest.WriteSimulationResponse(w, clientCtx.LegacyAmino, txf.Gas())
			return
		}
	}
	tx1, err := BuildUnsignedTx(txf, msgs...)
	if rest.CheckBadRequestError(w, err) {
		return
	}

	tmpKey := make([]byte, ethsecp256k1.PubKeySize)
	//tmpKey := make([]byte, secp256k1.PubKeySize)
	hexPK1 := "02"
	pubbyte := []byte(br.From)
	//accbyte := []byte(msgs.AccountAddress)
	pubhex := string(pubbyte[2:])
	//accaddr := string(accbyte[2:])
	hexPK1 += pubhex
	//AccAddr := strings.ToUpper(accaddr)
	bz, _ := hex.DecodeString(hexPK1)
	copy(tmpKey[:], bz)
	PubKey_Hex := hex.EncodeToString(tmpKey)
	fmt.Printf("PubKey_Hex=%+v\n", PubKey_Hex)
	pubB := &ethsecp256k1.PubKey{Key: tmpKey}
	//pubB := &secp256k1.PubKey{Key: tmpKey}
	fmt.Printf("pubB=%+v\n", pubB)
	fmt.Printf("pubB_address:=%+v\n", pubB.Address())
	NewAddress, _ := sdk.AccAddressFromHex(pubB.Address().String())
	address := NewAddress.String()
	//address, _ := sdk.AccAddressFromBech32(pubB.Address())
	fmt.Printf("address is :%+v\n", address)
	fmt.Printf("delegaor_address is :%+v\ntype:%T\n", dr.String(), dr.String())
	if dr.String() != address {
		hexPK2 := "03"
		tmpKey2 := make([]byte, ethsecp256k1.PubKeySize)
		pubbyte2 := []byte(br.From)
		pubhex2 := string(pubbyte2[2:])
		hexPK2 += pubhex2
		bz, _ := hex.DecodeString(hexPK2)
		copy(tmpKey2[:], bz)
		PubKey_Hex2 := hex.EncodeToString(tmpKey2)
		fmt.Printf("PubKey_Hex=%+v\n", PubKey_Hex2)
		pubB = &ethsecp256k1.PubKey{Key: tmpKey2}
		//pubB = &secp256k1.PubKey{Key: tmpKey2}
		// NewAddress, _ := sdk.AccAddressFromHex(pubB.Address().String())
		// address := NewAddress.String()
		// //address, _ := sdk.AccAddressFromBech32(pubB.Address())
		// fmt.Printf("address is :%+v\n", address)
	}
	ptr, _ := clientCtx.Codec.MarshalInterfaceJSON(pubB)
	fmt.Println("ptr=\n", ptr)
	//fmt.Println("ptr2=\n", pubB.GetKey())
	var pk cryptotypes.PubKey
	_ = clientCtx.Codec.UnmarshalInterfaceJSON(ptr, &pk)
	fmt.Println("pk=\n", pk)
	signMode := txf.txConfig.SignModeHandler().DefaultMode()
	fmt.Printf("signMode:=%+v\n", signMode)
	// signerData := authsigning.SignerData{
	// 	ChainID:       txf.chainID,
	// 	AccountNumber: txf.accountNumber,
	// 	Sequence:      txf.sequence,
	// }
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	fmt.Printf("tx1.GetTx()1=%+v\n", tx1.GetTx())
	sig := signing.SignatureV2{
		PubKey:   pk,
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}

	if err := tx1.SetSignatures(sig); err != nil {
		return
	}
	txJSONBytes1, err := clientCtx.TxConfig.TxJSONEncoder()(tx1.GetTx())
	if err != nil {
		return
	}
	fmt.Printf("txJSONBytesNew1=%+v\ntype:%T\n", txJSONBytes1, txJSONBytes1)
	txJSON1 := string(txJSONBytes1)
	fmt.Println("txJSONBytesNew2=\n", txJSON1)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(txJSONBytes1)
}

func WriteGeneratedTxResponse(
	clientCtx client.Context, w http.ResponseWriter, br rest.BaseReq, msgs ...sdk.Msg,
) {
	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	gasSetting, err := flags.ParseGasSetting(br.Gas)
	if rest.CheckBadRequestError(w, err) {
		return
	}

	txf := Factory{fees: br.Fees, gasPrices: br.GasPrices}.
		WithAccountNumber(br.AccountNumber).
		WithSequence(br.Sequence).
		WithGas(gasSetting.Gas).
		WithGasAdjustment(gasAdj).
		WithMemo(br.Memo).
		WithChainID(br.ChainID).
		WithSimulateAndExecute(br.Simulate).
		WithTxConfig(clientCtx.TxConfig).
		WithTimeoutHeight(br.TimeoutHeight)

	if br.Simulate || gasSetting.Simulate {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkerrors.ErrorInvalidGasAdjustment.Error())
			return
		}
		_, adjusted, err := CalculateGas(clientCtx, txf, msgs...)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		txf = txf.WithGas(adjusted)
		if br.Simulate {
			rest.WriteSimulationResponse(w, clientCtx.LegacyAmino, txf.Gas())
			return
		}
	}
	tx, err := BuildUnsignedTx(txf, msgs...)
	if rest.CheckBadRequestError(w, err) {
		return
	}

	stdTx, err := ConvertTxToStdTx(clientCtx.LegacyAmino, tx.GetTx())
	if rest.CheckInternalServerError(w, err) {
		return
	}

	output, err := clientCtx.LegacyAmino.MarshalJSON(stdTx)
	if rest.CheckInternalServerError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(output)
}

func WriteGeneratedTxResponse2(
	clientCtx client.Context, w http.ResponseWriter, br rest.BaseReq, msgs ...sdk.Msg,
) {
	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	gasSetting, err := flags.ParseGasSetting(br.Gas)
	if rest.CheckBadRequestError(w, err) {
		return
	}

	txf := Factory{fees: br.Fees, gasPrices: br.GasPrices}.
		WithAccountNumber(br.AccountNumber).
		WithSequence(br.Sequence).
		WithGas(gasSetting.Gas).
		WithGasAdjustment(gasAdj).
		WithMemo(br.Memo).
		WithChainID(br.ChainID).
		WithSimulateAndExecute(br.Simulate).
		WithTxConfig(clientCtx.TxConfig).
		WithTimeoutHeight(br.TimeoutHeight)

	if br.Simulate || gasSetting.Simulate {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkerrors.ErrorInvalidGasAdjustment.Error())
			return
		}

		_, adjusted, err := CalculateGas(clientCtx, txf, msgs...)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		txf = txf.WithGas(adjusted)

		if br.Simulate {
			rest.WriteSimulationResponse(w, clientCtx.LegacyAmino, txf.Gas())
			return
		}
	}
	tx, err := BuildUnsignedTx(txf, msgs...)
	if rest.CheckBadRequestError(w, err) {
		return
	}

	stdTx, err := ConvertTxToStdTx(clientCtx.LegacyAmino, tx.GetTx())
	if rest.CheckInternalServerError(w, err) {
		return
	}

	output, err := clientCtx.LegacyAmino.MarshalJSON(stdTx)
	if rest.CheckInternalServerError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(output)
}

// BuildUnsignedTx builds a transaction to be signed given a set of messages. The
// transaction is initially created via the provided factory's generator. Once
// created, the fee, memo, and messages are set.
func BuildUnsignedTx(txf Factory, msgs ...sdk.Msg) (client.TxBuilder, error) {
	if txf.chainID == "" {
		return nil, fmt.Errorf("chain ID required but not specified")
	}

	fees := txf.fees

	if !txf.gasPrices.IsZero() {
		if !fees.IsZero() {
			return nil, errors.New("cannot provide both fees and gas prices")
		}

		glDec := sdk.NewDec(int64(txf.gas))

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		fees = make(sdk.Coins, len(txf.gasPrices))

		for i, gp := range txf.gasPrices {
			fee := gp.Amount.Mul(glDec)
			fees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}
	}

	tx := txf.txConfig.NewTxBuilder()

	if err := tx.SetMsgs(msgs...); err != nil {
		return nil, err
	}

	tx.SetMemo(txf.memo)
	tx.SetFeeAmount(fees)
	tx.SetGasLimit(txf.gas)
	tx.SetTimeoutHeight(txf.TimeoutHeight())

	return tx, nil
}

// BuildSimTx creates an unsigned tx with an empty single signature and returns
// the encoded transaction or an error if the unsigned transaction cannot be
// built.
func BuildSimTx(txf Factory, msgs ...sdk.Msg) ([]byte, error) {
	txb, err := BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}

	// Create an empty signature literal as the ante handler will populate with a
	// sentinel pubkey.
	sig := signing.SignatureV2{
		PubKey: &secp256k1.PubKey{},
		Data: &signing.SingleSignatureData{
			SignMode: txf.signMode,
		},
		Sequence: txf.Sequence(),
	}
	if err := txb.SetSignatures(sig); err != nil {
		return nil, err
	}

	return txf.txConfig.TxEncoder()(txb.GetTx())
}

// CalculateGas simulates the execution of a transaction and returns the
// simulation response obtained by the query and the adjusted gas amount.
func CalculateGas(
	clientCtx gogogrpc.ClientConn, txf Factory, msgs ...sdk.Msg,
) (*tx.SimulateResponse, uint64, error) {
	txBytes, err := BuildSimTx(txf, msgs...)
	if err != nil {
		return nil, 0, err
	}
	txSvcClient := tx.NewServiceClient(clientCtx)
	simRes, err := txSvcClient.Simulate(context.Background(), &tx.SimulateRequest{
		TxBytes: txBytes,
	})
	//fmt.Printf("gas used:= %+v\n", simRes.GasInfo.GasUsed)
	if err != nil {
		return nil, 0, err
	}

	return simRes, uint64(txf.GasAdjustment() * float64(simRes.GasInfo.GasUsed)), nil
}

// prepareFactory ensures the account defined by ctx.GetFromAddress() exists and
// if the account number and/or the account sequence number are zero (not set),
// they will be queried for and set on the provided Factory. A new Factory with
// the updated fields will be returned.
func prepareFactory(clientCtx client.Context, txf Factory) (Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.accountRetriever.EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.accountNumber, txf.sequence
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.accountRetriever.GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}

// SignWithPrivKey signs a given tx with the given private key, and returns the
// corresponding SignatureV2 if the signing is successful.
func SignWithPrivKey(
	signMode signing.SignMode, signerData authsigning.SignerData,
	txBuilder client.TxBuilder, priv cryptotypes.PrivKey, txConfig client.TxConfig,
	accSeq uint64,
) (signing.SignatureV2, error) {
	var sigV2 signing.SignatureV2

	// Generate the bytes to be signed.
	signBytes, err := txConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return sigV2, err
	}

	// Sign those bytes
	signature, err := priv.Sign(signBytes)
	if err != nil {
		return sigV2, err
	}

	// Construct the SignatureV2 struct
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: signature,
	}

	sigV2 = signing.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &sigData,
		Sequence: accSeq,
	}

	return sigV2, nil
}

func checkMultipleSigners(mode signing.SignMode, tx authsigning.Tx) error {
	if mode == signing.SignMode_SIGN_MODE_DIRECT &&
		len(tx.GetSigners()) > 1 {
		return sdkerrors.Wrap(sdkerrors.ErrNotSupported, "Signing in DIRECT mode is only supported for transactions with one signer only")
	}
	return nil
}

// Sign signs a given tx with a named key. The bytes signed over are canconical.
// The resulting signature will be added to the transaction builder overwriting the previous
// ones if overwrite=true (otherwise, the signature will be appended).
// Signing a transaction with mutltiple signers in the DIRECT mode is not supprted and will
// return an error.
// An error is returned upon failure.
func Sign(txf Factory, name string, txBuilder client.TxBuilder, overwriteSig bool) error {
	if txf.keybase == nil {
		return errors.New("keybase must be set prior to signing a transaction")
	}

	signMode := txf.signMode
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		// use the SignModeHandler's default mode if unspecified
		signMode = txf.txConfig.SignModeHandler().DefaultMode()
	}
	if err := checkMultipleSigners(signMode, txBuilder.GetTx()); err != nil {
		return err
	}

	key, err := txf.keybase.Key(name)
	if err != nil {
		return err
	}
	pubKey := key.GetPubKey()
	signerData := authsigning.SignerData{
		ChainID:       txf.chainID,
		AccountNumber: txf.accountNumber,
		Sequence:      txf.sequence,
	}

	// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
	// TxBuilder under the hood, and SignerInfos is needed to generated the
	// sign bytes. This is the reason for setting SetSignatures here, with a
	// nil signature.
	//
	// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
	// also doesn't affect its generated sign bytes, so for code's simplicity
	// sake, we put it here.
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}
	var prevSignatures []signing.SignatureV2
	if !overwriteSig {
		prevSignatures, err = txBuilder.GetTx().GetSignaturesV2()
		if err != nil {
			return err
		}
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		return err
	}

	// Generate the bytes to be signed.
	bytesToSign, err := txf.txConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return err
	}

	// Sign those bytes
	sigBytes, _, err := txf.keybase.Sign(name, bytesToSign)
	if err != nil {
		return err
	}

	// Construct the SignatureV2 struct
	sigData = signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}

	if overwriteSig {
		return txBuilder.SetSignatures(sig)
	}
	prevSignatures = append(prevSignatures, sig)
	return txBuilder.SetSignatures(prevSignatures...)
}

// GasEstimateResponse defines a response definition for tx gas estimation.
type GasEstimateResponse struct {
	GasEstimate uint64 `json:"gas_estimate" yaml:"gas_estimate"`
}

func (gr GasEstimateResponse) String() string {
	return fmt.Sprintf("gas estimate: %d", gr.GasEstimate)
}

//Signing a Transaction

// func sendTx() error {
//     // --snip--

//     privs := []cryptotypes.PrivKey{priv1, priv2}
//     accNums:= []uint64{..., ...} // The accounts' account numbers
//     accSeqs:= []uint64{..., ...} // The accounts' sequence numbers

//     // First round: we gather all the signer infos. We use the "set empty
//     // signature" hack to do that.
//     var sigsV2 []signing.SignatureV2
//     for i, priv := range privs {
//         sigV2 := signing.SignatureV2{
//             PubKey: priv.PubKey(),
//             Data: &signing.SingleSignatureData{
//                 SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
//                 Signature: nil,
//             },
//             Sequence: accSeqs[i],
//         }

//         sigsV2 = append(sigsV2, sigV2)
//     }
//     err := txBuilder.SetSignatures(sigsV2...)
//     if err != nil {
//         return err
//     }

//     // Second round: all signer infos are set, so each signer can sign.
//     sigsV2 = []signing.SignatureV2{}
//     for i, priv := range privs {
//         signerData := xauthsigning.SignerData{
//             ChainID:       chainID,
//             AccountNumber: accNums[i],
//             Sequence:      accSeqs[i],
//         }
//         sigV2, err := tx.SignWithPrivKey(
//             encCfg.TxConfig.SignModeHandler().DefaultMode(), signerData,
//             txBuilder, priv, encCfg.TxConfig, accSeqs[i])
//         if err != nil {
//             return nil, err
//         }

//         sigsV2 = append(sigsV2, sigV2)
//     }
//     err = txBuilder.SetSignatures(sigsV2...)
//     if err != nil {
//         return err
//     }
// }

type bech32Output struct {
	Formats []string `json:"formats"`
}

func newBech32Output(args string) bech32Output {
	bech32Prefixes := []string{"treasurenet", "treasurenetpub"}
	addr := strings.TrimSpace(args)
	bz, _ := hex.DecodeString(addr)
	out := bech32Output{Formats: make([]string, len(bech32Prefixes))}

	for i, prefix := range bech32Prefixes {
		bech32Addr, err := bech32.ConvertAndEncode(prefix, bz)
		if err != nil {
			panic(err)
		}

		out.Formats[i] = bech32Addr
	}

	return out
}
