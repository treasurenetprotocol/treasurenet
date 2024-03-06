package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/treasurenetprotocol/treasurenet/x/gravity/types"
)

// nolint: exhaustruct
func TestPrefixRange(t *testing.T) {
	cases := map[string]struct {
		src      []byte
		expStart []byte
		expEnd   []byte
		expPanic bool
	}{
		"normal":              {src: []byte{1, 3, 4}, expStart: []byte{1, 3, 4}, expEnd: []byte{1, 3, 5}},
		"normal short":        {src: []byte{79}, expStart: []byte{79}, expEnd: []byte{80}},
		"empty case":          {src: []byte{}},
		"roll-over example 1": {src: []byte{17, 28, 255}, expStart: []byte{17, 28, 255}, expEnd: []byte{17, 29, 0}},
		"roll-over example 2": {src: []byte{15, 42, 255, 255},
			expStart: []byte{15, 42, 255, 255}, expEnd: []byte{15, 43, 0, 0}},
		"pathological roll-over": {src: []byte{255, 255, 255, 255}, expStart: []byte{255, 255, 255, 255}},
		"nil prohibited":         {expPanic: true},
	}

	for testName, tc := range cases {
		tc := tc
		t.Run(testName, func(t *testing.T) {
			if tc.expPanic {
				require.Panics(t, func() {
					prefixRange(tc.src)
				})
				return
			}
			start, end := prefixRange(tc.src)
			assert.Equal(t, tc.expStart, start)
			assert.Equal(t, tc.expEnd, end)
		})
	}
}

// Test that valset creation produces the expected normalized power values
// // nolint: exhaustruct
// func TestCurrentValsetNormalization(t *testing.T) {
// 	// Setup the overflow test
// 	maxPower64 := make([]uint64, 64)             // users with max power (approx 2^63)
// 	expPower64 := make([]uint64, 64)             // expected scaled powers
// 	ethAddrs64 := make([]gethcommon.Address, 64) // need 64 eth addresses for this test
// 	for i := 0; i < 64; i++ {
// 		maxPower64[i] = uint64(9223372036854775807)
// 		expPower64[i] = 67108864 // 2^32 split amongst 64 validators
// 		ethAddrs64[i] = gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(i)}, 20))
// 	}

// 	// any lower than this and a validator won't be created
// 	const minStake = 1000000

// 	specs := map[string]struct {
// 		srcPowers []uint64
// 		expPowers []uint64
// 	}{
// 		"one": {
// 			srcPowers: []uint64{minStake},
// 			expPowers: []uint64{4294967296},
// 		},
// 		"two": {
// 			srcPowers: []uint64{minStake * 99, minStake * 1},
// 			expPowers: []uint64{4252017623, 42949672},
// 		},
// 		"four equal": {
// 			srcPowers: []uint64{minStake, minStake, minStake, minStake},
// 			expPowers: []uint64{1073741824, 1073741824, 1073741824, 1073741824},
// 		},
// 		"four equal max power": {
// 			srcPowers: []uint64{4294967296, 4294967296, 4294967296, 4294967296},
// 			expPowers: []uint64{1073741824, 1073741824, 1073741824, 1073741824},
// 		},
// 		"overflow": {
// 			srcPowers: maxPower64,
// 			expPowers: expPower64,
// 		},
// 	}
// 	for msg, spec := range specs {
// 		spec := spec
// 		t.Run(msg, func(t *testing.T) {
// 			input, ctx := SetupTestChain(t, spec.srcPowers, true)
// 			r, err := input.GravityKeeper.GetCurrentValset(ctx)
// 			require.NoError(t, err)
// 			rMembers, err := types.BridgeValidators(r.Members).ToInternal()
// 			require.NoError(t, err)
// 			assert.Equal(t, spec.expPowers, rMembers.GetPowers())
// 		})
// 	}
// }

// nolint: exhaustruct
func TestAttestationIterator(t *testing.T) {
	input := CreateTestEnv(t)
	defer func() { input.Context.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	ctx := input.Context
	// add some attestations to the store

	att1 := &types.Attestation{
		Observed: true,
		Votes:    []string{},
	}
	dep1 := &types.MsgSendToCosmosClaim{
		EventNonce:     1,
		TokenContract:  TokenContractAddrs[0],
		Amount:         sdk.NewInt(100),
		EthereumSender: EthAddrs[0].String(),
		CosmosReceiver: AccAddrs[0].String(),
		Orchestrator:   AccAddrs[0].String(),
	}
	att2 := &types.Attestation{
		Observed: true,
		Votes:    []string{},
	}
	dep2 := &types.MsgSendToCosmosClaim{
		EventNonce:     2,
		TokenContract:  TokenContractAddrs[0],
		Amount:         sdk.NewInt(100),
		EthereumSender: EthAddrs[0].String(),
		CosmosReceiver: AccAddrs[0].String(),
		Orchestrator:   AccAddrs[0].String(),
	}
	hash1, err := dep1.ClaimHash()
	require.NoError(t, err)
	hash2, err := dep2.ClaimHash()
	require.NoError(t, err)

	input.GravityKeeper.SetAttestation(ctx, dep1.EventNonce, hash1, att1)
	input.GravityKeeper.SetAttestation(ctx, dep2.EventNonce, hash2, att2)

	atts := []types.Attestation{}
	input.GravityKeeper.IterateAttestations(ctx, false, func(_ []byte, att types.Attestation) bool {
		atts = append(atts, att)
		return false
	})

	require.Len(t, atts, 2)
}

// nolint: exhaustruct
func TestDelegateKeys(t *testing.T) {
	input := CreateTestEnv(t)
	defer func() { input.Context.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	ctx := input.Context
	k := input.GravityKeeper
	var (
		ethAddrs = []string{"0x0F528A4Be8720D2BF71f5A4EbC138e06eCBa289f",
			"0x19Ae0a2f2F47F818017e0EcB3163C04D18b97B65", "0x1FBae33071b1CA691B6919fdD138770824Fe676D",
			"0xa09C955532f00b470cF64672232f81A0AbE2A69e"}

		valAddrs = []string{"treasurenetvaloper1v2xyllxrwd60mfwa6aj6r8fa4xz4235zfnxrlu",
			"treasurenetvaloper10aays6dtcx7tlwvqrngc06a2rp7jy0cvpha4dg", "treasurenetvaloper1sa74q750nrs3729zmd489ae7a3527997uzf62y",
			"treasurenetvaloper158ckg6g8l8zy4jckj4hc4tjzx8skeeye5n2whz"}

		orchAddrs = []string{"treasurenet1jz68cpsjqkstc2pdxlpkdfcw35uwp0uyxrgvac", "treasurenet1rrc6tv02rx9jqfncmyteghj35xv3d2yskgle0s",
			"treasurenet1rfanxd2agtt3a9xcpsgge5ad03hsfj4z0qd5zl", "treasurenet1gez03q60xkff4qa5w7ceckmnqemyxgftqx5tp0"}
	)

	for i := range ethAddrs {
		// set some addresses
		val, err1 := sdk.ValAddressFromBech32(valAddrs[i])
		orch, err2 := sdk.AccAddressFromBech32(orchAddrs[i])
		require.NoError(t, err1)
		require.NoError(t, err2)
		// set the orchestrator address
		k.SetOrchestratorValidator(ctx, val, orch)
		// set the ethereum address
		ethAddr, err := types.NewEthAddress(ethAddrs[i])
		require.NoError(t, err)
		k.SetEthAddressForValidator(ctx, val, *ethAddr)
	}

	addresses := k.GetDelegateKeys(ctx)
	for i := range addresses {
		res := addresses[i]
		assert.Equal(t, valAddrs[i], res.Validator)
		assert.Equal(t, orchAddrs[i], res.Orchestrator)
		assert.Equal(t, ethAddrs[i], res.EthAddress)
	}

}

// nolint: exhaustruct
func TestLastSlashedValsetNonce(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	defer func() { input.Context.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	k := input.GravityKeeper

	vs, err := k.GetCurrentValset(ctx)
	require.NoError(t, err)

	i := 1
	for ; i < 10; i++ {
		vs.Height = uint64(i)
		vs.Nonce = uint64(i)
		k.StoreValset(ctx, vs)
		k.SetLatestValsetNonce(ctx, vs.Nonce)
	}

	latestValsetNonce := k.GetLatestValsetNonce(ctx)
	assert.Equal(t, latestValsetNonce, uint64(i-1))

	latestValset := k.GetLatestValset(ctx)
	assert.Equal(t, uint64(i-1), latestValset.Nonce)

	// lastSlashedValsetNonce should be zero initially.
	lastSlashedValsetNonce := k.GetLastSlashedValsetNonce(ctx)
	assert.Equal(t, lastSlashedValsetNonce, uint64(0))
	unslashedValsets := k.GetUnSlashedValsets(ctx, uint64(12))
	assert.Equal(t, len(unslashedValsets), 9)

	// check if last Slashed Valset nonce is set properly or not
	k.SetLastSlashedValsetNonce(ctx, uint64(3))
	lastSlashedValsetNonce = k.GetLastSlashedValsetNonce(ctx)
	assert.Equal(t, lastSlashedValsetNonce, uint64(3))

	lastSlashedValset := k.GetValset(ctx, lastSlashedValsetNonce)

	// when valset height + signedValsetsWindow > current block height, len(unslashedValsets) should be zero
	unslashedValsets = k.GetUnSlashedValsets(ctx, uint64(ctx.BlockHeight()))
	assert.Equal(t, len(unslashedValsets), 0)

	// when lastSlashedValset height + signedValsetsWindow == BlockHeight, len(unslashedValsets) should be zero
	heightDiff := uint64(ctx.BlockHeight()) - lastSlashedValset.Height
	unslashedValsets = k.GetUnSlashedValsets(ctx, heightDiff)
	assert.Equal(t, len(unslashedValsets), 0)

	// when signedValsetsWindow is between lastSlashedValset height and latest valset's height
	unslashedValsets = k.GetUnSlashedValsets(ctx, heightDiff-2)
	assert.Equal(t, len(unslashedValsets), 2)

	// when signedValsetsWindow > latest valset's height
	unslashedValsets = k.GetUnSlashedValsets(ctx, heightDiff-6)
	assert.Equal(t, len(unslashedValsets), 6)
	fmt.Println("unslashedValsetsRange", unslashedValsets)
}
