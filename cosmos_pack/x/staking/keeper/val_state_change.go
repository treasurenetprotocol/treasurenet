package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"

	gogotypes "github.com/gogo/protobuf/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BlockValidatorUpdates calculates the ValidatorUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	// unbond all mature validators from the unbonding queue
	k.UnbondAllMatureValidators(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		addr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress, err := sdk.AccAddressFromBech32(dvPair.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress, err := sdk.AccAddressFromBech32(dvvTriplet.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return validatorUpdates
}

func (k Keeper) NewBlockValidatorUpdates(ctx sdk.Context, req abci.RequestEndBlock, log sdk.ABCIMessageLogs) []abci.ValidatorUpdate {
	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates, err := k.NewApplyAndReturnValidatorSetUpdates(ctx, req, log)
	if err != nil {
		panic(err)
	}

	// unbond all mature validators from the unbonding queue（解除绑定）
	k.UnbondAllMatureValidators(ctx)

	// Remove all mature unbonding delegations from the ubd queue.（）
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		addr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress, err := sdk.AccAddressFromBech32(dvPair.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress, err := sdk.AccAddressFromBech32(dvvTriplet.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return validatorUpdates
}

// ApplyAndReturnValidatorSetUpdates applies and return accumulated updates to the bonded validator set. Also,
// * Updates the active valset as keyed by LastValidatorPowerKey.
// * Updates the total power as keyed by LastTotalPowerKey.
// * Updates validator status' according to updated powers.
// * Updates the fee pool bonded vs not-bonded tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only validators with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the validator set entirely
// are returned to Tendermint.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate, err error) {
	params := k.GetParams(ctx)
	maxValidators := params.MaxValidators
	powerReduction := k.PowerReduction(ctx)
	totalPower := sdk.ZeroInt()
	amtFromBondedToNotBonded, amtFromNotBondedToBonded := sdk.ZeroInt(), sdk.ZeroInt()

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last, err := k.getLastValidatorsByAddr(ctx)
	if err != nil {
		return nil, err
	}

	// Iterate over validators, highest power to lowest.
	iterator := k.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()

	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set
		valAddr := sdk.ValAddress(iterator.Value())
		validator := k.mustGetValidator(ctx, valAddr)
		if validator.Jailed {
			panic("should never retrieve a jailed validator from the power store")
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if validator.PotentialConsensusPower(k.PowerReduction(ctx)) == 0 {
			break
		}

		// apply the appropriate state change if necessary
		switch {
		case validator.IsUnbonded():
			validator, err = k.unbondedToBonded(ctx, validator)
			if err != nil {
				return
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
		case validator.IsUnbonding():
			validator, err = k.unbondingToBonded(ctx, validator)
			if err != nil {
				return
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
		case validator.IsBonded():
			// no state change
		default:
			panic("unexpected validator status")
		}
		// fetch the old power bytes
		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
		if err != nil {
			return nil, err
		}
		oldPowerBytes, found := last[valAddrStr]
		newPower := validator.ConsensusPower(powerReduction)
		newPowerBytes := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: newPower})
		// update the validator set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			updates = append(updates, validator.ABCIValidatorUpdate(powerReduction))

			k.SetLastValidatorPower(ctx, valAddr, newPower)
		}
		delete(last, valAddrStr)
		count++

		totalPower = totalPower.Add(sdk.NewInt(newPower))
	}

	noLongerBonded, err := sortNoLongerBonded(last)
	if err != nil {
		return nil, err
	}

	for _, valAddrBytes := range noLongerBonded {
		validator := k.mustGetValidator(ctx, sdk.ValAddress(valAddrBytes))
		validator, err = k.bondedToUnbonding(ctx, validator)
		if err != nil {
			return
		}
		amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.GetTokens())
		k.DeleteLastValidatorPower(ctx, validator.GetOperator())
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// Update the pools based on the recent updates in the validator set:
	// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
	// to the Bonded pool.
	// - The tokens from the bonded validators that are being kicked out from the validator set
	// need to be transferred to the NotBonded pool.
	switch {
	// Compare and subtract the respective amounts to only perform one transfer.
	// This is done in order to avoid doing multiple updates inside each iterator/loop.
	case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
		k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded))
	case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
		k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded))
	default: // equal amounts of tokens; no update required
	}

	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetLastTotalPower(ctx, totalPower)
	}

	return updates, err
}

// bid v1.0.0

/*****bid 1.0.0*****/
func (k Keeper) NewApplyAndReturnValidatorSetUpdates(ctx sdk.Context, req abci.RequestEndBlock, log sdk.ABCIMessageLogs) (updates []abci.ValidatorUpdate, err error) {
	params := k.GetParams(ctx)
	maxValidators := params.MaxValidators
	powerReduction := k.PowerReduction(ctx)
	// powerReduction2 := k.PowerReduction2(ctx)
	totalPower := sdk.ZeroInt()
	amtFromBondedToNotBonded, amtFromNotBondedToBonded := sdk.ZeroInt(), sdk.ZeroInt()
	last, err := k.getLastValidatorsByAddr(ctx)
	if err != nil {
		return nil, err
	}
	/*
	 * accounts_address  That is, the account address returned from the event log is converted into the corresponding validator_ Addresses, and then replace POS
	 */
	for _, eventlog := range log {
		if eventlog.MsgIndex == 1 {
			iterator := k.ValidatorsPowerStoreIterator(ctx)
			defer iterator.Close()
			// Tatiterator := k.ValidatorsNewPowerStoreIterator(ctx)
			// defer Tatiterator.Close()
			listsupervalidator, listvalidator := k.CombinedSliceList(ctx, iterator, maxValidators, log)
			fmt.Println("listsupervalidatoe:", listsupervalidator)
			fmt.Println("listvalidator", listvalidator)
			newselectlist := SelectList(listsupervalidator, listvalidator, req)
			// newlist := make(map[string]string, len(newselectlist))
			// for _, value := range newselectlist {
			// 	newlist[value] = value
			// }
			k.DeletsNewItreator2(ctx, maxValidators, newselectlist)
			// fmt.Println("删除iterator迭代器:", newselectlist)

			newTatiteratorlist := make(map[string]string)
			NewTatiterator := k.ValidatorsNewPowerStoreIterator(ctx)
			defer NewTatiterator.Close()
			for count := 0; NewTatiterator.Valid() && count < int(maxValidators); NewTatiterator.Next() {
				// everything that is iterated in this loop is becoming or already a
				// part of the bonded validator set
				valAddr := sdk.ValAddress(NewTatiterator.Value())
				validator := k.mustGetValidator(ctx, valAddr)

				if validator.Jailed {
					panic("should never retrieve a jailed validator from the power store")
				}
				// if we get to a zero-power validator (which we don't bond),
				// there are no more possible bonded validators
				if validator.PotentialConsensusPower(k.PowerReduction(ctx)) == 0 {
					break
				}

				// apply the appropriate state change if necessary
				switch {
				case validator.IsUnbonded():
					validator, err = k.unbondedToBonded(ctx, validator)
					if err != nil {
						return
					}
					amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
				case validator.IsUnbonding():
					validator, err = k.unbondingToBonded(ctx, validator)
					if err != nil {
						return
					}
					amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
				case validator.IsBonded():
					// no state change
				default:
					panic("unexpected validator status")
				}
				// fetch the old power bytes
				valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
				if err != nil {
					return nil, err
				}
				oldPowerBytes, found := last[valAddrStr]
				newPower := validator.ConsensusPower(powerReduction)
				newPowerBytes := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: newPower})
				newPower2 := validator.ConsensusTatPower(powerReduction)
				newunitPower := validator.ConsensusNewPower(powerReduction)
				k.SetTatPower(ctx, newPower2, valAddr)
				k.SetNewUnitPower(ctx, newunitPower, valAddr)
				if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
					updates = append(updates, validator.ABCIValidatorUpdate(powerReduction))

					k.SetLastValidatorPower(ctx, valAddr, newPower)
				}
				delete(last, valAddrStr)
				count++
				newTatiteratorlist[valAddrStr] = valAddrStr
				totalPower = totalPower.Add(sdk.NewInt(newPower))
			}

			noLongerBonded, err := sortNoLongerBonded(last)
			if err != nil {
				return nil, err
			}

			for _, valAddrBytes := range noLongerBonded {
				validator := k.mustGetValidator(ctx, sdk.ValAddress(valAddrBytes))
				validator, err = k.bondedToUnbonding(ctx, validator)
				if err != nil {
					return nil, err
				}
				amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.GetTokens())
				k.DeleteLastValidatorPower(ctx, validator.GetOperator())
				updates = append(updates, validator.ABCIValidatorUpdateZero())
			}

			// Update the pools based on the recent updates in the validator set:
			// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
			// to the Bonded pool.
			// - The tokens from the bonded validators that are being kicked out from the validator set
			// need to be transferred to the NotBonded pool.
			switch {
			// Compare and subtract the respective amounts to only perform one transfer.
			// This is done in order to avoid doing multiple updates inside each iterator/loop.
			case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
				k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded))
			case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
				k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded))
			default: // equal amounts of tokens; no update required
			}
		} else if eventlog.MsgIndex == 2 {
			newTatiteratorlist := make(map[string]string)
			NewTatiterator := k.ValidatorsNewPowerStoreIterator(ctx)
			defer NewTatiterator.Close()
			for count := 0; NewTatiterator.Valid() && count < int(maxValidators); NewTatiterator.Next() {
				// everything that is iterated in this loop is becoming or already a
				// part of the bonded validator set
				valAddr := sdk.ValAddress(NewTatiterator.Value())
				// fetch the old power bytes
				valAddrStr, _ := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
				newTatiteratorlist[valAddrStr] = valAddrStr
			}
			iterator := k.ValidatorsPowerStoreIterator(ctx)
			defer iterator.Close()
			// k.DeleteNewIterator(ctx)
			for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
				// everything that is iterated in this loop is becoming or already a
				// part of the bonded validator set
				valAddr := sdk.ValAddress(iterator.Value())
				validator := k.mustGetValidator(ctx, valAddr)
				fmt.Printf("validator eventlog.MsgIndex == 2  :%+v\n", validator)
				if validator.Jailed {
					panic("should never retrieve a jailed validator from the power store")
				}

				// if we get to a zero-power validator (which we don't bond),
				// there are no more possible bonded validators
				if validator.PotentialConsensusPower(k.PowerReduction(ctx)) == 0 {
					break
				}

				valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
				fmt.Println("valAddrStr2:", valAddrStr)
				if err != nil {
					return nil, err
				}
				// oldPowerBytes, found := last[valAddrStr]
				oldPowerBytes := last[valAddrStr]
				_, ok := newTatiteratorlist[valAddrStr]
				newPower := validator.ConsensusPower(powerReduction)
				newPowerBytes := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: newPower})
				// update the validator set if power has changed
				// if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
				// 	updates = append(updates, validator.ABCIValidatorUpdate(powerReduction))

				// 	k.SetLastValidatorPower(ctx, valAddr, newPower)
				// }
				if !bytes.Equal(oldPowerBytes, newPowerBytes) {
					updates = append(updates, validator.ABCIValidatorUpdate(powerReduction))

					k.SetLastValidatorPower(ctx, valAddr, newPower)
				}
				if ok {
					totalPower = totalPower.Add(sdk.NewInt(newPower))
					delete(last, valAddrStr)
				}
				count++
			}
			noLongerBonded, err := sortNoLongerBonded(last)
			if err != nil {
				return nil, err
			}

			for _, valAddrBytes := range noLongerBonded {
				validator := k.mustGetValidator(ctx, sdk.ValAddress(valAddrBytes))
				validator, err = k.bondedToUnbonding(ctx, validator)
				if err != nil {
					return nil, err
				}
				amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.GetTokens())
				k.DeleteLastValidatorPower(ctx, validator.GetOperator())
				updates = append(updates, validator.ABCIValidatorUpdateZero())
			}

			// Update the pools based on the recent updates in the validator set:
			// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
			// to the Bonded pool.
			// - The tokens from the bonded validators that are being kicked out from the validator set
			// need to be transferred to the NotBonded pool.
			switch {
			// Compare and subtract the respective amounts to only perform one transfer.
			// This is done in order to avoid doing multiple updates inside each iterator/loop.
			case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
				k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded))
			case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
				k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded))
			default: // equal amounts of tokens; no update required
			}

			// // set total power on lookup index if there are any updates
		} else {
			// Iterate over validators, highest power to lowest.
			iterator := k.ValidatorsPowerStoreIterator(ctx)
			defer iterator.Close()

			for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
				// everything that is iterated in this loop is becoming or already a
				// part of the bonded validator set
				valAddr := sdk.ValAddress(iterator.Value())
				validator := k.mustGetValidator(ctx, valAddr)
				if validator.Jailed {
					panic("should never retrieve a jailed validator from the power store")
				}
				// if we get to a zero-power validator (which we don't bond),
				// there are no more possible bonded validators
				if validator.PotentialConsensusPower(k.PowerReduction(ctx)) == 0 {
					break
				}

				// apply the appropriate state change if necessary
				switch {
				case validator.IsUnbonded():
					validator, err = k.unbondedToBonded(ctx, validator)
					if err != nil {
						return
					}
					amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
				case validator.IsUnbonding():
					validator, err = k.unbondingToBonded(ctx, validator)
					if err != nil {
						return
					}
					amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
				case validator.IsBonded():
					// no state change
				default:
					panic("unexpected validator status")
				}
				// fetch the old power bytes
				valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
				if err != nil {
					return nil, err
				}
				oldPowerBytes, found := last[valAddrStr]
				newPower := validator.ConsensusPower(powerReduction)
				newPowerBytes := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: newPower})
				// update the validator set if power has changed
				if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
					updates = append(updates, validator.ABCIValidatorUpdate(powerReduction))

					k.SetLastValidatorPower(ctx, valAddr, newPower)
				}
				delete(last, valAddrStr)
				count++

				totalPower = totalPower.Add(sdk.NewInt(newPower))
			}

			noLongerBonded, err := sortNoLongerBonded(last)
			if err != nil {
				return nil, err
			}

			for _, valAddrBytes := range noLongerBonded {
				validator := k.mustGetValidator(ctx, sdk.ValAddress(valAddrBytes))
				validator, err = k.bondedToUnbonding(ctx, validator)
				if err != nil {
					return nil, err
				}
				amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.GetTokens())
				k.DeleteLastValidatorPower(ctx, validator.GetOperator())
				updates = append(updates, validator.ABCIValidatorUpdateZero())
			}

			// Update the pools based on the recent updates in the validator set:
			// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
			// to the Bonded pool.
			// - The tokens from the bonded validators that are being kicked out from the validator set
			// need to be transferred to the NotBonded pool.
			switch {
			// Compare and subtract the respective amounts to only perform one transfer.
			// This is done in order to avoid doing multiple updates inside each iterator/loop.
			case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
				k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded))
			case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
				k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded))
			default: // equal amounts of tokens; no update required
			}
		}
	}
	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetLastTotalPower(ctx, totalPower)
	}
	fmt.Println("updates:", updates)
	return updates, err
}

// Validator state transitions

func (k Keeper) bondedToUnbonding(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsBonded() {
		panic(fmt.Sprintf("bad state transition bondedToUnbonding, validator: %v\n", validator))
	}

	return k.beginUnbondingValidator(ctx, validator)
}

func (k Keeper) unbondingToBonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}

	return k.bondValidator(ctx, validator)
}

func (k Keeper) unbondedToBonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonded() {
		panic(fmt.Sprintf("bad state transition unbondedToBonded, validator: %v\n", validator))
	}

	return k.bondValidator(ctx, validator)
}

// UnbondingToUnbonded switches a validator from unbonding state to unbonded state
func (k Keeper) UnbondingToUnbonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if !validator.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}

	return k.completeUnbondingValidator(ctx, validator)
}

// send a validator to jail
func (k Keeper) jailValidator(ctx sdk.Context, validator types.Validator) {
	if validator.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed validator, validator: %v\n", validator))
	}

	validator.Jailed = true
	k.SetValidator(ctx, validator)
	k.DeleteValidatorByPowerIndex(ctx, validator)
	k.DeleteValidatorByTatPowerIndex(ctx, validator)
}

// remove a validator from jail
func (k Keeper) unjailValidator(ctx sdk.Context, validator types.Validator) {
	if !validator.Jailed {
		panic(fmt.Sprintf("cannot unjail already unjailed validator, validator: %v\n", validator))
	}

	validator.Jailed = false
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)
}

// perform all the store operations for when a validator status becomes bonded 当验证器状态变为已绑定时，执行所有存储操作
func (k Keeper) bondValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)
	// k.DeleteValidatorByTatPowerIndex(ctx, validator)

	validator = validator.UpdateStatus(types.Bonded)

	// save the now bonded validator record to the two referenced stores
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	// k.SetNewValidatorByPowerIndex(ctx, validator)
	// delete from queue if present
	k.DeleteValidatorQueue(ctx, validator)

	// trigger hook
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return validator, err
	}
	k.AfterValidatorBonded(ctx, consAddr, validator.GetOperator())

	return validator, err
}

// perform all the store operations for when a validator begins unbonding 当验证器开始解除绑定时，执行所有存储操作
func (k Keeper) beginUnbondingValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	params := k.GetParams(ctx)

	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)
	// k.DeleteValidatorByTatPowerIndex(ctx, validator)
	// sanity check
	if validator.Status != types.Bonded {
		panic(fmt.Sprintf("should not already be unbonded or unbonding, validator: %v\n", validator))
	}

	validator = validator.UpdateStatus(types.Unbonding)

	// set the unbonding completion time and completion height appropriately
	validator.UnbondingTime = ctx.BlockHeader().Time.Add(params.UnbondingTime)
	validator.UnbondingHeight = ctx.BlockHeader().Height

	// save the now unbonded validator record and power index
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	// k.SetNewValidatorByPowerIndex(ctx, validator)
	// Adds to unbonding validator queue
	k.InsertUnbondingValidatorQueue(ctx, validator)

	// trigger hook
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return validator, err
	}
	k.AfterValidatorBeginUnbonding(ctx, consAddr, validator.GetOperator())

	return validator, nil
}

// perform all the store operations for when a validator status becomes unbonded
func (k Keeper) completeUnbondingValidator(ctx sdk.Context, validator types.Validator) types.Validator {
	validator = validator.UpdateStatus(types.Unbonded)
	k.SetValidator(ctx, validator)

	return validator
}

// map of operator bech32-addresses to serialized power
// We use bech32 strings here, because we can't have slices as keys: map[[]byte][]byte
type validatorsByAddr map[string][]byte

// get the last validator set
func (k Keeper) getLastValidatorsByAddr(ctx sdk.Context) (validatorsByAddr, error) {
	last := make(validatorsByAddr)

	iterator := k.LastValidatorsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// extract the validator address from the key (prefix is 1-byte, addrLen is 1-byte)
		valAddr := types.AddressFromLastValidatorPowerKey(iterator.Key())
		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
		if err != nil {
			return nil, err
		}

		powerBytes := iterator.Value()
		last[valAddrStr] = make([]byte, len(powerBytes))
		copy(last[valAddrStr], powerBytes)
	}

	return last, nil
}

// get the last validator tat set
func (k Keeper) getLastValidatorsTatByAddr(ctx sdk.Context) (validatorsByAddr, error) {
	last := make(validatorsByAddr)

	tatiterator := k.LastValidatorsTatIterator(ctx)
	defer tatiterator.Close()

	for ; tatiterator.Valid(); tatiterator.Next() {
		// extract the validator address from the key (prefix is 1-byte, addrLen is 1-byte)
		valAddr := types.AddressFromLastValidatorPowerKey(tatiterator.Key())
		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
		if err != nil {
			return nil, err
		}

		powerBytes := tatiterator.Value()
		last[valAddrStr] = make([]byte, len(powerBytes))
		copy(last[valAddrStr], powerBytes)
	}

	return last, nil
}

// get the last validator set
func (k Keeper) getLastValidatorsNewByAddr(ctx sdk.Context) (validatorsByAddr, error) {
	last := make(validatorsByAddr)

	iterator := k.LastValidatorsNewIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// extract the validator address from the key (prefix is 1-byte, addrLen is 1-byte)
		valAddr := types.AddressFromLastValidatorPowerKey(iterator.Key())
		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
		if err != nil {
			return nil, err
		}

		powerBytes := iterator.Value()
		last[valAddrStr] = make([]byte, len(powerBytes))
		copy(last[valAddrStr], powerBytes)
	}

	return last, nil
}

// given a map of remaining validators to previous bonded power
// returns the list of validators to be unbonded, sorted by operator address
func sortNoLongerBonded(last validatorsByAddr) ([][]byte, error) {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0

	for valAddrStr := range last {
		valAddrBytes, err := sdk.ValAddressFromBech32(valAddrStr)
		if err != nil {
			return nil, err
		}
		noLongerBonded[index] = valAddrBytes
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})

	return noLongerBonded, nil
}

// Handle the list of tatvalidator and validator
func CombinedSlice(iterator sdk.Iterator, maxValidators uint32, validatorsByAddr string) ([]string, []string) {
	var ListSuperValidator []string
	var ListValidator []string
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		valAddr := sdk.ValAddress(iterator.Value())
		// fmt.Println("valAddr:", valAddr)
		validatorstring := valAddr.String()
		if validatorstring != validatorsByAddr {
			ListValidator = append(ListValidator, validatorstring)
		} else {
			ListSuperValidator = append(ListSuperValidator, validatorstring)
		}
	}
	return ListSuperValidator, ListValidator
}
func (k Keeper) CombinedSliceList(ctx sdk.Context, iterator sdk.Iterator, maxValidators uint32, log sdk.ABCIMessageLogs) ([]string, []string) {
	var Data [][]interface{}
	ListSuperValidator := []string{}
	ListValidator := []string{}
	var tat sdk.Int
	var newunit sdk.Int
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		valAddr := sdk.ValAddress(iterator.Value())
		// fmt.Println("valAddr:", valAddr)
		validatorstring := valAddr.String()
		for _, eventlog := range log {
			if eventlog.MsgIndex == 1 {
				asslog := []byte(eventlog.Log)
				err := json.Unmarshal(asslog, &Data)
				if err != nil {
					fmt.Println("error:", err)
				}
				if len(Data) == 0 {
					Zero := sdk.ZeroInt()
					NewZero, _ := Zero.MarshalJSON()
					k.SetTat2(ctx, NewZero, valAddr)
					k.SetNewToken2(ctx, NewZero, valAddr)
				} else {
					for _, vlog := range Data {
						fmt.Printf("Conversion of account address to verifier address :%+v\n", vlog[0].(string))
						a := []byte(vlog[0].(string))
						c := string(a[2:])
						s := strings.ToUpper(c)
						NewValidatoradd, _ := sdk.ValAddressFromHex(s)
						NewValidatoraddstring := NewValidatoradd.String()
						fmt.Printf("CombinedSliceList validatorstring :%+v\n", validatorstring)
						fmt.Printf("CombinedSliceList NewValidatoraddstring :%+v\n", NewValidatoraddstring)
						if validatorstring == NewValidatoraddstring {
							//	ListSuperValidator = append(ListValidator, validatorstring)
							ListSuperValidator = append(ListSuperValidator, validatorstring)
							fmt.Printf("CombinedSliceList ListSuperValidator :%+v\n", ListSuperValidator)
						}
						fmt.Println(reflect.TypeOf(vlog[1]))
						stringtat := strconv.FormatFloat(vlog[1].(float64), 'f', -1, 64)
						fmt.Println("stringtat:", stringtat)
						tat, _ = sdk.NewIntFromString(stringtat)
						stringunit := strconv.FormatFloat(vlog[1].(float64), 'f', -1, 64)
						fmt.Println("stringunit:", stringunit)
						newunit, _ = sdk.NewIntFromString(stringunit)
						newtat, _ := tat.MarshalJSON()
						newunitbyte, _ := newunit.MarshalJSON()
						fmt.Println("newtat:", newtat)
						fmt.Println("newunitbyte:", newunitbyte)
						k.SetTat2(ctx, newtat, NewValidatoradd)
						k.SetNewToken2(ctx, newunitbyte, NewValidatoradd)
					}
				}

			}
		}
		ListValidator = append(ListValidator, validatorstring)
	}
	return ListSuperValidator, ListValidator
}
func SelectList(listsupervalidator []string, listvalidator []string, req abci.RequestEndBlock) []string {
	var Activelent int
	var newlistsupervalidator []string
	var newlistvalidator []string
	var newvalidator []string
	listval := SubtrDemo(listsupervalidator, listvalidator)
	switch {
	case len(listvalidator) >= 200:
		Activelent = 100
	case len(listvalidator) >= 8 && len(listvalidator) < 200:
		Activelent = len(listvalidator) / 2
	default:
		Activelent = len(listvalidator)
	}
	if len(listsupervalidator) >= 100 && len(listval) >= 100 {
		newlistsupervalidator = listsupervalidator[:100]
		newlistvalidator = listval[:100]
		newvalidator = append(newlistsupervalidator, newlistvalidator...)
	} else if len(listsupervalidator) < 100 && len(listval) < 100 {
		newvalidator = append(listsupervalidator, listval...)
	} else if len(listsupervalidator) < 100 && len(listval) > 100 {
		num := 2*Activelent - len(listsupervalidator)
		newlistvalidator = listval[:num]
		newvalidator = append(listsupervalidator, newlistvalidator...)
	} else if len(listsupervalidator) > 100 && len(listval) < 100 {
		num := 2*Activelent - len(listval)
		newlistsupervalidator = listsupervalidator[:num]
		newvalidator = append(newlistsupervalidator, listval...)
	}
	MicsSliceList := MicsSlice(newvalidator, Activelent, req)
	return MicsSliceList
}
func SubtrDemo(listsupervalidator []string, listvalidator []string) []string {
	var removal []string
	temp := map[string]struct{}{} // map[string]struct{}{}创建了一个key类型为String值类型为空struct的map，Equal -> make(map[string]struct{})

	for _, val := range listsupervalidator {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{} // 空struct 不占内存空间
		}
	}
	for _, val := range listvalidator {
		if _, ok := temp[val]; !ok {
			removal = append(removal, val)
		}
	}

	return removal
}

// random number
func MicsSlice(origin []string, count int, req abci.RequestEndBlock) []string {
	tmpOrigin := make([]string, len(origin))
	copy(tmpOrigin, origin)
	// 一定要seed
	newheight := req.Height
	// rand.Seed(time.Now().Unix())
	rand.Seed(newheight)
	rand.Shuffle(len(tmpOrigin), func(i int, j int) {
		tmpOrigin[i], tmpOrigin[j] = tmpOrigin[j], tmpOrigin[i]
	})
	fmt.Println(tmpOrigin)
	result := make([]string, 0, count)
	for index, value := range tmpOrigin {
		if index == count {
			break
		}
		result = append(result, value)
	}
	return result
}

func (k Keeper) DeleteNewIterator(ctx sdk.Context) {
	TatIterator := k.ValidatorsNewPowerStoreIterator(ctx)
	// iterator := k.ValidatorsNewPowerStoreIterator(ctx)
	defer TatIterator.Close()
	for ; TatIterator.Valid(); TatIterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set
		// fmt.Printf("iterator.Value:%v\n", iterator.Value())
		valAddr := sdk.ValAddress(TatIterator.Value())
		// fmt.Println("valAddr:", valAddr)
		validator := k.mustGetValidator(ctx, valAddr)
		k.SetTatPower(ctx, int64(0), valAddr)
		k.DeleteValidatorByTatPowerIndex(ctx, validator)
	}
}
func (k Keeper) DeletsNewItreator2(ctx sdk.Context, maxValidators uint32, newselectlist []string) {
	newlist := make(map[string]string, len(newselectlist))
	newTatiteratorlist := make(map[string]string)
	// iterator := k.ValidatorsPowerStoreIterator(ctx)
	// defer iterator.Close()
	for _, value := range newselectlist {
		newlist[value] = value
	}
	NewTatiterator := k.ValidatorsNewPowerStoreIterator(ctx)
	defer NewTatiterator.Close()
	for count := 0; NewTatiterator.Valid() && count < int(maxValidators); NewTatiterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set
		valAddr := sdk.ValAddress(NewTatiterator.Value())
		validator := k.mustGetValidator(ctx, valAddr)
		// fetch the old power bytes
		valAddrStr, _ := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
		newTatiteratorlist[valAddrStr] = valAddrStr
		if len(newselectlist) > 0 {
			_, found := newlist[valAddrStr]
			if !found {
				// valAddrBytes, _ := sdk.ValAddressFromBech32(valAddrStr)
				// validator := k.mustGetValidator(ctx, sdk.ValAddress(valAddrBytes))
				k.DeleteValidatorByTatPowerIndex(ctx, validator)
			}
		}
	}
	for _, value := range newTatiteratorlist {
		delete(newlist, value)
	}
	if len(newlist) > 0 {
		for _, val := range newlist {
			valAddrBytes, _ := sdk.ValAddressFromBech32(val)
			validator := k.mustGetValidator(ctx, valAddrBytes)
			k.SetNewValidatorByPowerIndex(ctx, validator)
		}
	}
}

func (k Keeper) AddNewIterator(ctx sdk.Context, validator types.Validator) {
	k.SetNewValidatorByPowerIndex(ctx, validator)
}

// func (k Keeper) AddNewIterator2(ctx sdk.Context, iterator sdk.Iterator, maxValidators uint32, log sdk.ABCIMessageLogs) {
// 	listsupervalidator, listvalidator := k.CombinedSliceList(ctx, iterator, maxValidators, log)
// 	fmt.Println("listsupervalidatoe:", listsupervalidator)
// 	fmt.Println("listvalidator", listvalidator)
// 	newselectlist := SelectList(listsupervalidator, listvalidator)
// 	fmt.Println("生成新的迭代器:", 123)
// 	k.DeleteNewIterator(ctx)
// 	if len(newselectlist) > 0 {
// 		fmt.Println("newselectlist:", newselectlist)
// 		for _, value := range newselectlist {
// 			Tatvaladdr, _ := sdk.ValAddressFromBech32(value)
// 			fmt.Println("生成新的迭代器22222222:", Tatvaladdr)
// 			Tatvalidator := k.mustGetValidator(ctx, Tatvaladdr)
// 			k.AddNewIterator(ctx, Tatvalidator)
// 		}
// 	}
// }
// func (k Keeper) AddNewIterator3(ctx sdk.Context, iterator sdk.Iterator, maxValidators uint32, log sdk.ABCIMessageLogs) {
// 	listsupervalidator, listvalidator := k.CombinedSliceList(ctx, iterator, maxValidators, log)
// 	fmt.Println("listsupervalidatoe:", listsupervalidator)
// 	fmt.Println("listvalidator", listvalidator)
// 	newselectlist := SelectList(listsupervalidator, listvalidator)
// 	fmt.Println("生成新的迭代器:", 123)
// 	k.DeleteNewIterator(ctx)
// 	if len(newselectlist) > 0 {
// 		fmt.Println("newselectlist:", newselectlist)
// 		for _, value := range newselectlist {
// 			Tatvaladdr, _ := sdk.ValAddressFromBech32(value)
// 			fmt.Println("生成新的迭代器22222222:", Tatvaladdr)
// 			Tatvalidator := k.mustGetValidator(ctx, Tatvaladdr)
// 			k.AddNewIterator(ctx, Tatvalidator)
// 		}
// 	}
// }
