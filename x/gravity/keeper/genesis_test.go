package keeper

import (
	"fmt"
	"sort"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/treasurenetprotocol/treasurenet/x/gravity/types"
)

// Tests that batches and transactions are preserved during chain restart
func TestBatchAndTxImportExport(t *testing.T) {
	// SETUP ENV + DATA
	// ==================
	input := CreateTestEnv(t)
	defer func() { input.Context.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	ctx := input.Context
	batchSize := 100
	accAddresses := []string{ // Warning: this must match the length of ctrAddresses

		"treasurenet102e577f29shw9ngdp462ml3nj44mdrg4hy2y07",
		"treasurenet1v2xyllxrwd60mfwa6aj6r8fa4xz4235zgd4z75",
		"treasurenet10aays6dtcx7tlwvqrngc06a2rp7jy0cvqfw5vq",
		"treasurenet1sa74q750nrs3729zmd489ae7a3527997au6mtv",
		"treasurenet158ckg6g8l8zy4jckj4hc4tjzx8skeeye4de0k2",
	}
	ethAddresses := []string{
		"0xbcbA9257D4419bcEd1C6E6F06f8f8Beba471B016",
		"0x0F528A4Be8720D2BF71f5A4EbC138e06eCBa289f",
		"0x19Ae0a2f2F47F818017e0EcB3163C04D18b97B65",
		"0x1FBae33071b1CA691B6919fdD138770824Fe676D",
		"0xa09C955532f00b470cF64672232f81A0AbE2A69e",
	}
	ctrAddresses := []string{ // Warning: this must match the length of accAddresses
		"0x7Ab34f792A2C2eE2cD0D0d74aDFe33956Bb68d15",
		"0x628C4FFcC37374Fda5DDd765a19d3Da985554682",
		"0x7F7A4869Abc1bCbFb9801cD187ebaA187D223F0C",
		"0x877d507a8f98E11f28A2DB6a72f73Eec68AF14be",
		"0xa1F1646907f9c44Acb16956f8AaE4231E16ce499",
	}

	// SETUP ACCOUNTS
	// ==================
	senders := make([]*sdk.AccAddress, len(accAddresses))
	for i := range senders {
		sender, err := sdk.AccAddressFromBech32(accAddresses[i])
		require.NoError(t, err)
		senders[i] = &sender
	}
	receivers := make([]*types.EthAddress, len(ethAddresses))
	for i := range receivers {
		receiver, err := types.NewEthAddress(ethAddresses[i])
		require.NoError(t, err)
		receivers[i] = receiver
	}
	contracts := make([]*types.EthAddress, len(ctrAddresses))
	for i := range contracts {
		contract, err := types.NewEthAddress(ctrAddresses[i])
		require.NoError(t, err)
		contracts[i] = contract
	}
	tokens := make([]*types.InternalERC20Token, len(contracts))
	vouchers := make([]*sdk.Coins, len(contracts))
	for i, v := range contracts {
		token, err := types.NewInternalERC20Token(sdk.NewInt(99999999), v.GetAddress().Hex())
		tokens[i] = token
		allVouchers := sdk.NewCoins(token.GravityCoin())
		vouchers[i] = &allVouchers
		require.NoError(t, err)

		// Mint the vouchers
		require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers))
	}

	// give sender i a balance of token i
	for i, v := range senders {
		input.AccountKeeper.NewAccountWithAddress(ctx, *v)
		require.NoError(t, input.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, *v, *vouchers[i]))
	}

	// CREATE TRANSACTIONS
	// ==================
	numTxs := 5000 // should end up with 1000 txs per contract
	txs := make([]*types.InternalOutgoingTransferTx, numTxs)
	fees := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	amounts := []int{51, 52, 53, 54, 55, 56, 57, 58, 59, 60}
	for i := 0; i < numTxs; i++ {
		// Pick fee, amount, sender, receiver, and contract for the ith transaction
		// Sender and contract will always match up (they must since sender i controls the whole balance of the ith token)
		// Receivers should get a balance of many token types since i % len(receivers) is usually different than i % len(contracts)
		fee := fees[i%len(fees)] // fee for this transaction
		amount := amounts[i%len(amounts)]
		sender := senders[i%len(senders)]
		receiver := receivers[i%len(receivers)]
		contract := contracts[i%len(contracts)]
		amountToken, err := types.NewInternalERC20Token(sdk.NewInt(int64(amount)), contract.GetAddress().Hex())
		require.NoError(t, err)
		feeToken, err := types.NewInternalERC20Token(sdk.NewInt(int64(fee)), contract.GetAddress().Hex())
		require.NoError(t, err)

		// add transaction to the pool
		id, err := input.GravityKeeper.AddToOutgoingPool(ctx, *sender, *receiver, amountToken.GravityCoin(), feeToken.GravityCoin())
		require.NoError(t, err)
		ctx.Logger().Info(fmt.Sprintf("Created transaction %v with amount %v and fee %v of contract %v from %v to %v", i, amount, fee, contract, sender, receiver))

		// Record the transaction for later testing
		tx, err := types.NewInternalOutgoingTransferTx(id, sender.String(), receiver.GetAddress().Hex(), amountToken.ToExternal(), feeToken.ToExternal())
		require.NoError(t, err)
		txs[i] = tx
	}

	// when

	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)

	// CREATE BATCHES
	// ==================
	// Want to create batches for half of the transactions for each contract
	// with 100 tx in each batch, 1000 txs per contract, we want 5 batches per contract to batch 500 txs per contract
	batches := make([]*types.InternalOutgoingTxBatch, 5*len(contracts))
	for i, v := range contracts {
		batch, err := input.GravityKeeper.BuildOutgoingTXBatch(ctx, *v, uint(batchSize))
		require.NoError(t, err)
		batches[i] = batch
		ctx.Logger().Info(fmt.Sprintf("Created batch %v for contract %v with %v transactions", i, v.GetAddress(), batchSize))
	}

	checkAllTransactionsExist(t, input.GravityKeeper, ctx, txs)
	exportImport(t, &input)
	checkAllTransactionsExist(t, input.GravityKeeper, ctx, txs)
}

// Requires that all transactions in txs exist in keeper
func checkAllTransactionsExist(t *testing.T, keeper Keeper, ctx sdk.Context, txs []*types.InternalOutgoingTransferTx) {
	unbatched := keeper.GetUnbatchedTransactions(ctx)
	batches := keeper.GetOutgoingTxBatches(ctx)
	// Collect all txs into an array
	var gotTxs []*types.InternalOutgoingTransferTx
	gotTxs = append(gotTxs, unbatched...)
	for _, batch := range batches {
		gotTxs = append(gotTxs, batch.Transactions...)
	}
	require.Equal(t, len(txs), len(gotTxs))
	// Sort both arrays for simple searching
	sort.Slice(gotTxs, func(i, j int) bool {
		return gotTxs[i].Id < gotTxs[j].Id
	})
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Id < txs[j].Id
	})
	// Actually check that the txs all exist, iterate on txs in case some got lost in the import/export step
	for i, exp := range txs {
		require.Equal(t, exp.Id, gotTxs[i].Id)
		require.Equal(t, exp.Erc20Fee, gotTxs[i].Erc20Fee)
		require.Equal(t, exp.Erc20Token, gotTxs[i].Erc20Token)
		require.Equal(t, exp.DestAddress.GetAddress(), gotTxs[i].DestAddress.GetAddress())
		require.Equal(t, exp.Sender.String(), gotTxs[i].Sender.String())
	}
}

// Exports and then imports all bridge state, overwrites the `input` test environment to simulate chain restart
func exportImport(t *testing.T, input *TestInput) {
	genesisState := ExportGenesis(input.Context, input.GravityKeeper)
	newEnv := CreateTestEnv(t)
	input = &newEnv
	unbatched := input.GravityKeeper.GetUnbatchedTransactions(input.Context)
	require.Empty(t, unbatched)
	batches := input.GravityKeeper.GetOutgoingTxBatches(input.Context)
	require.Empty(t, batches)
	InitGenesis(input.Context, input.GravityKeeper, genesisState)
}
