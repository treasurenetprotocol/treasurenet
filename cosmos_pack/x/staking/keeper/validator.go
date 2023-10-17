package keeper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// get a single validator
func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	// fmt.Println("types.GetValidatorKey(addr)", types.GetValidatorKey(addr))
	value := store.Get(types.GetValidatorKey(addr))
	if value == nil {
		return validator, false
	}

	validator = types.MustUnmarshalValidator(k.cdc, value)
	// TatTokens, _ := k.GetTatTokens(ctx, addr)
	// NewTokens, _ := k.GetNewTokens(ctx, addr)
	TatTokens, _ := k.GetTatTokens2(ctx, addr)
	NewTokens, _ := k.GetNewTokens2(ctx, addr)
	TatPower, _ := k.GetTatPower(ctx, addr)
	NewUnitPower, _ := k.GetNewUnitPower(ctx, addr)
	// validator.TatTokens = sdk.NewInt(TatTokens)
	// validator.NewTokens = sdk.NewInt(NewTokens)
	validator.TatTokens = TatTokens
	validator.NewTokens = NewTokens
	validator.NewUnitPower = NewUnitPower
	validator.TatPower = sdk.NewInt(TatPower)
	return validator, true
}
func (k Keeper) GetTatTokens(ctx sdk.Context, addr sdk.ValAddress) (tattokens int64, found bool) {
	store := ctx.KVStore(k.storeKey)
	// fmt.Println("GetValidator addr:", addr)
	// fmt.Println("types.GetValidatorKey(addr)", types.GetValidatorKey(addr))
	value := store.Get(types.GetTatTokensKey(addr))
	if value == nil {
		return tattokens, false
	}
	newtat, _ := strconv.ParseInt(string(value), 10, 64)
	// strtat, _ := sdk.NewIntFromString(string(value))
	tattokens = newtat
	return tattokens, true
}
func (k Keeper) GetTatTokens2(ctx sdk.Context, addr sdk.ValAddress) (tattokens sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetTatTokensKey(addr))
	if value == nil {
		return sdk.ZeroInt(), false
		// return TatTokens, false
	}
	// newtat, _ := sdk.NewIntFromString(string(value))
	err := tattokens.UnmarshalJSON(value)
	if err != nil {
		return sdk.ZeroInt(), false
	}
	// fmt.Println("newtat:", TatTokens)
	return tattokens, true
	// strtat, _ := sdk.NewIntFromString(string(value))
	// TatTokens = newtat
	// return TatTokens, true
}
func (k Keeper) GetNewTokens(ctx sdk.Context, addr sdk.ValAddress) (newtokens int64, found bool) {
	store := ctx.KVStore(k.storeKey)
	// fmt.Println("GetValidator addr:", addr)
	// fmt.Println("types.GetValidatorKey(addr)", types.GetValidatorKey(addr))
	value := store.Get(types.GetNewTokensKey(addr))
	if value == nil {
		return newtokens, false
	}
	newtat, _ := strconv.ParseInt(string(value), 10, 64)
	// strtat, _ := sdk.NewIntFromString(string(value))
	newtokens = newtat
	return newtokens, true
}
func (k Keeper) GetNewTokens2(ctx sdk.Context, addr sdk.ValAddress) (newtokens sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetNewTokensKey(addr))
	if value == nil {
		return sdk.ZeroInt(), false
		// return NewTokens, false
	}
	// newtat, _ := sdk.NewIntFromString(string(value))
	// strtat, _ := sdk.NewIntFromString(string(value))
	err := newtokens.UnmarshalJSON(value)
	if err != nil {
		return sdk.ZeroInt(), false
	}
	// fmt.Println("newunit:", NewTokens)
	return newtokens, true
	// NewTokens = newtat
	// return NewTokens, true
}
func (k Keeper) GetTatPower(ctx sdk.Context, addr sdk.ValAddress) (tatpower int64, found bool) {
	store := ctx.KVStore(k.storeKey)
	// fmt.Println("GetValidator addr:", addr)
	// fmt.Println("types.GetValidatorKey(addr)", types.GetValidatorKey(addr))
	value := store.Get(types.GetTatPowerKey(addr))
	if value == nil {
		return tatpower, false
	}
	tatpower1, _ := strconv.ParseInt(string(value), 10, 64)
	// strtat, _ := sdk.NewIntFromString(string(value))
	tatpower = tatpower1
	return tatpower, true
}
func (k Keeper) GetTatPower2(ctx sdk.Context, addr sdk.ValAddress) (tatpower sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetTatPowerKey(addr))
	if value == nil {
		return sdk.ZeroInt(), false
		// return TatPower, false
	}
	// tatpower, _ := sdk.NewIntFromString(string(value))
	// TatPower = tatpower
	// return TatPower, true
	err := tatpower.UnmarshalJSON(value)
	if err != nil {
		return sdk.ZeroInt(), false
	}
	fmt.Println("TatPower:", tatpower)
	return tatpower, true
}
func (k Keeper) GetNewUnitPower(ctx sdk.Context, addr sdk.ValAddress) (newunitpower sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetNewUnitPowerKey(addr))
	if value == nil {
		return sdk.ZeroInt(), false
		// return NewUnitPower, false
	}
	// newunitpower, _ := sdk.NewIntFromString(string(value))
	// NewUnitPower = newunitpower
	// return NewUnitPower, true
	newunitpower1, _ := strconv.ParseInt(string(value), 10, 64)
	// strtat, _ := sdk.NewIntFromString(string(value))
	newunitpower = sdk.NewInt(newunitpower1)
	return newunitpower, true
}
func (k Keeper) mustGetValidator(ctx sdk.Context, addr sdk.ValAddress) types.Validator {
	validator, found := k.GetValidator(ctx, addr)
	if !found {
		panic(fmt.Sprintf("validator record not found for address: %X\n", addr))
	}

	return validator
}

// get a single validator by consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	opAddr := store.Get(types.GetValidatorByConsAddrKey(consAddr))
	if opAddr == nil {
		return validator, false
	}
	// fmt.Println("GetValidatorByConsAddr opAddr:", opAddr)
	return k.GetValidator(ctx, opAddr)
}

func (k Keeper) mustGetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) types.Validator {
	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		panic(fmt.Errorf("validator with consensus-Address %s not found", consAddr))
	}

	return validator
}

// set the main record holding validator details
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	// v := &validator
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalValidator(k.cdc, &validator)
	store.Set(types.GetValidatorKey(validator.GetOperator()), bz)
}

// tattoken
func (k Keeper) SetTat(ctx sdk.Context, tatToken int64, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	// tat := sdk.NewInt(tatToken)
	// bz, _ := tat.MarshalJSON()
	bz, _ := json.Marshal(tatToken)
	// bz, _ := json.Marshal(&validator)
	// unbz := types.MustUnmarshalValidator(k.cdc, bz)
	// fmt.Printf("unbz:%+v\n", unbz)
	// UnmarshalValidator
	// fmt.Printf("bz:%+v\n", bz)
	// fmt.Printf("validator.GetOperator():%+v\n", validator.GetOperator())
	// fmt.Printf("types.GetValidatorKey(validator.GetOperator()):%+v\n", types.GetValidatorKey(validator.GetOperator()))
	store.Set(types.GetTatTokensKey(addr), bz)
}

// tattoken
func (k Keeper) SetTat2(ctx sdk.Context, tatToken []byte, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetTatTokensKey(addr), tatToken)
}
func (k Keeper) SetNewToken(ctx sdk.Context, newToken int64, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := json.Marshal(newToken)
	store.Set(types.GetNewTokensKey(addr), bz)
}
func (k Keeper) SetNewToken2(ctx sdk.Context, newToken []byte, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetNewTokensKey(addr), newToken)
}
func (k Keeper) SetNewUnitPower(ctx sdk.Context, newunitPower int64, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := json.Marshal(newunitPower)
	store.Set(types.GetNewUnitPowerKey(addr), bz)
}

// tatpower
func (k Keeper) SetTatPower(ctx sdk.Context, tatPower int64, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := json.Marshal(tatPower)
	store.Set(types.GetTatPowerKey(addr), bz)
}
func (k Keeper) SetTatPower2(ctx sdk.Context, tatPower []byte, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetTatPowerKey(addr), tatPower)
}

// validator index
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, validator types.Validator) error {
	consPk, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorByConsAddrKey(consPk), validator.GetOperator())
	return nil
}

// validator index
func (k Keeper) SetValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	// jailed validators are not kept in the power index
	if validator.Jailed {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorsByPowerIndexKey(validator, k.PowerReduction(ctx)), validator.GetOperator())
}

// new validator index
func (k Keeper) SetNewValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	// jailed validators are not kept in the power index
	if validator.Jailed {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorsByNewPowerIndexKey(validator, k.PowerReduction(ctx)), validator.GetOperator())
}

// validator index
func (k Keeper) DeleteValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorsByPowerIndexKey(validator, k.PowerReduction(ctx)))
}

// validator index tat
func (k Keeper) DeleteValidatorByTatPowerIndex(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorsByNewPowerIndexKey(validator, k.PowerReduction(ctx)))
}

// validator index tat
func (k Keeper) DeleteValidatorByTatPowerIndex2(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorsByNewPowerTatIndexKey(validator, k.PowerReduction(ctx)))
}

// validator index
// func (k Keeper) SetNewValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
// 	store := ctx.KVStore(k.storeKey)
// 	store.Set(types.GetValidatorsByPowerIndexKey(validator, k.PowerReduction(ctx)), validator.GetOperator())
// }

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) AddValidatorTokensAndShares(ctx sdk.Context, validator types.Validator,
	tokensToAdd sdk.Int) (valOut types.Validator, addedShares sdk.Dec) {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	k.DeleteValidatorByTatPowerIndex(ctx, validator)
	validator, addedShares = validator.AddTokensFromDel(tokensToAdd)
	// TatTokens, _ := k.GetTatTokens(ctx, validator.GetOperator())
	// NewTokens, _ := k.GetNewTokens(ctx, validator.GetOperator())
	TatTokens, _ := k.GetTatTokens2(ctx, validator.GetOperator())
	NewTokens, _ := k.GetNewTokens2(ctx, validator.GetOperator())
	TatPower, _ := k.GetTatPower(ctx, validator.GetOperator())
	NewUnitPower, _ := k.GetNewUnitPower(ctx, validator.GetOperator())
	// TatTokens := params.TatTokens
	// validator.TatTokens = sdk.NewInt(TatTokens)
	// validator.NewTokens = sdk.NewInt(NewTokens)
	validator.TatTokens = TatTokens
	validator.NewTokens = NewTokens
	validator.TatPower = sdk.NewInt(TatPower)
	validator.NewUnitPower = NewUnitPower
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)
	return validator, addedShares
}

// Update the tokens of an existing validator, update the validators tat power index key
func (k Keeper) AddValidatorTatTokensAndShares(ctx sdk.Context, validator types.Validator,
	tokensToAdd sdk.Int) (valOut types.Validator, addedShares sdk.Dec) {

	validator, addedShares = validator.AddTatTokensFromDel(tokensToAdd)
	k.SetValidator(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	return validator, addedShares
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) RemoveValidatorTokensAndShares(ctx sdk.Context, validator types.Validator,
	sharesToRemove sdk.Dec) (valOut types.Validator, removedTokens sdk.Int) {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	k.DeleteValidatorByTatPowerIndex(ctx, validator)
	validator, removedTokens = validator.RemoveDelShares(sharesToRemove)
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)
	return validator, removedTokens
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) RemoveValidatorTokens(ctx sdk.Context,
	validator types.Validator, tokensToRemove sdk.Int) types.Validator {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	k.DeleteValidatorByTatPowerIndex(ctx, validator)
	validator = validator.RemoveTokens(tokensToRemove)
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)
	return validator
}

// UpdateValidatorCommission attempts to update a validator's commission rate.
// An error is returned if the new commission rate is invalid.
func (k Keeper) UpdateValidatorCommission(ctx sdk.Context,
	validator types.Validator, newRate sdk.Dec) (types.Commission, error) {
	commission := validator.Commission
	blockTime := ctx.BlockHeader().Time

	if err := commission.ValidateNewRate(newRate, blockTime); err != nil {
		return commission, err
	}

	commission.Rate = newRate
	commission.UpdateTime = blockTime

	return commission, nil
}

// remove the validator record and associated indexes
// except for the bonded validator index which is only handled in ApplyAndReturnTendermintUpdates
// TODO, this function panics, and it's not good.
func (k Keeper) RemoveValidator(ctx sdk.Context, address sdk.ValAddress) {
	// first retrieve the old validator record
	validator, found := k.GetValidator(ctx, address)
	if !found {
		return
	}

	if !validator.IsUnbonded() {
		panic("cannot call RemoveValidator on bonded or unbonding validators")
	}

	if validator.Tokens.IsPositive() {
		panic("attempting to remove a validator which still contains tokens")
	}

	valConsAddr, err := validator.GetConsAddr()
	if err != nil {
		panic(err)
	}

	// delete the old validator record
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorKey(address))
	store.Delete(types.GetValidatorByConsAddrKey(valConsAddr))
	store.Delete(types.GetValidatorsByPowerIndexKey(validator, k.PowerReduction(ctx)))

	// call hooks
	k.AfterValidatorRemoved(ctx, valConsAddr, validator.GetOperator())
}

// get groups of validators

// get the set of all validators with no limits, used during genesis dump
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}

	return validators
}

// return a given amount of all the validators
func (k Keeper) GetValidators(ctx sdk.Context, maxRetrieve uint32) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	validators = make([]types.Validator, maxRetrieve)

	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators[i] = validator
		i++
	}

	return validators[:i] // trim if the array length < maxRetrieve
}

// get the current group of bonded validators sorted by power-rank
func (k Keeper) GetBondedValidatorsByPower(ctx sdk.Context) []types.Validator {
	maxValidators := k.MaxValidators(ctx)
	validators := make([]types.Validator, maxValidators)
	fmt.Println("测试GetBondedValidatorsByPower")
	iterator := k.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
		address := iterator.Value()
		validator := k.mustGetValidator(ctx, address)

		if validator.IsBonded() {
			validators[i] = validator
			i++
		}
	}

	return validators[:i] // trim
}

// returns an iterator for the current validator power store
func (k Keeper) ValidatorsPowerStoreIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.ValidatorsByPowerIndexKey)
}

// returns an iterator for the current validator newpower store
func (k Keeper) ValidatorsNewPowerStoreIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.ValidatorsByNewPowerIndexKey)
}

// Last Validator Index

// Load the last validator power.
// Returns zero if the operator was not a validator last block.
func (k Keeper) GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) (power int64) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetLastValidatorPowerKey(operator))
	if bz == nil {
		return 0
	}

	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshal(bz, &intV)

	return intV.GetValue()
}

// Set the last validator power.
func (k Keeper) SetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress, power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: power})
	store.Set(types.GetLastValidatorPowerKey(operator), bz)
}

// Set the last validator tatpower.
func (k Keeper) SetLastValidatorTatPower(ctx sdk.Context, operator sdk.ValAddress, tatpower int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: tatpower})
	store.Set(types.GetLastValidatorTatPowerKey(operator), bz)
}

// Set the last validator unitpower.
func (k Keeper) SetLastValidatorUnitPower(ctx sdk.Context, operator sdk.ValAddress, newpower int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: newpower})
	store.Set(types.GetLastValidatorUnitPowerKey(operator), bz)
}

// Delete the last validator power.
func (k Keeper) DeleteLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastValidatorPowerKey(operator))
}

// Delete the last validator tatpower.
func (k Keeper) DeleteLastValidatorTatPower(ctx sdk.Context, operator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastValidatorTatPowerKey(operator))
}

// returns an iterator for the consensus validators in the last block
func (k Keeper) LastValidatorsIterator(ctx sdk.Context) (iterator sdk.Iterator) {
	store := ctx.KVStore(k.storeKey)
	iterator = sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)

	return iterator
}

// returns an iterator for the consensus validators tat in the last block
func (k Keeper) LastValidatorsTatIterator(ctx sdk.Context) (iterator sdk.Iterator) {
	store := ctx.KVStore(k.storeKey)
	iterator = sdk.KVStorePrefixIterator(store, types.LastValidatorTatPowerKey)

	return iterator
}
func (k Keeper) LastValidatorsNewIterator(ctx sdk.Context) (iterator sdk.Iterator) {
	store := ctx.KVStore(k.storeKey)
	iterator = sdk.KVStorePrefixIterator(store, types.LastValidatorUnitPowerKey)

	return iterator
}

// Iterate over last validator powers. 迭代最后的验证程序权限
func (k Keeper) IterateLastValidatorPowers(ctx sdk.Context, handler func(operator sdk.ValAddress, power int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(types.AddressFromLastValidatorPowerKey(iter.Key()))
		intV := &gogotypes.Int64Value{}

		k.cdc.MustUnmarshal(iter.Value(), intV)

		if handler(addr, intV.GetValue()) {
			break
		}
	}
}

// get the group of the bonded validators
func (k Keeper) GetLastValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)

	// add the actual validator power sorted store
	maxValidators := k.MaxValidators(ctx)
	validators = make([]types.Validator, maxValidators)

	iterator := sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		// sanity check
		if i >= int(maxValidators) {
			panic("more validators than maxValidators found")
		}

		address := types.AddressFromLastValidatorPowerKey(iterator.Key())
		validator := k.mustGetValidator(ctx, address)

		validators[i] = validator
		i++
	}

	return validators[:i] // trim
}

// GetUnbondingValidators returns a slice of mature validator addresses that
// complete their unbonding at a given time and height.
func (k Keeper) GetUnbondingValidators(ctx sdk.Context, endTime time.Time, endHeight int64) []string {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetValidatorQueueKey(endTime, endHeight))
	if bz == nil {
		return []string{}
	}

	addrs := types.ValAddresses{}
	k.cdc.MustUnmarshal(bz, &addrs)

	return addrs.Addresses
}

// SetUnbondingValidatorsQueue sets a given slice of validator addresses into
// the unbonding validator queue by a given height and time.
func (k Keeper) SetUnbondingValidatorsQueue(ctx sdk.Context, endTime time.Time, endHeight int64, addrs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.ValAddresses{Addresses: addrs})
	store.Set(types.GetValidatorQueueKey(endTime, endHeight), bz)
}

// InsertUnbondingValidatorQueue inserts a given unbonding validator address into
// the unbonding validator queue for a given height and time.
func (k Keeper) InsertUnbondingValidatorQueue(ctx sdk.Context, val types.Validator) {
	addrs := k.GetUnbondingValidators(ctx, val.UnbondingTime, val.UnbondingHeight)
	addrs = append(addrs, val.OperatorAddress)
	k.SetUnbondingValidatorsQueue(ctx, val.UnbondingTime, val.UnbondingHeight, addrs)
}

// DeleteValidatorQueueTimeSlice deletes all entries in the queue indexed by a
// given height and time.
func (k Keeper) DeleteValidatorQueueTimeSlice(ctx sdk.Context, endTime time.Time, endHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorQueueKey(endTime, endHeight))
}

// DeleteValidatorQueue removes a validator by address from the unbonding queue
// indexed by a given height and time.
func (k Keeper) DeleteValidatorQueue(ctx sdk.Context, val types.Validator) {
	addrs := k.GetUnbondingValidators(ctx, val.UnbondingTime, val.UnbondingHeight)
	newAddrs := []string{}

	for _, addr := range addrs {
		if addr != val.OperatorAddress {
			newAddrs = append(newAddrs, addr)
		}
	}

	if len(newAddrs) == 0 {
		k.DeleteValidatorQueueTimeSlice(ctx, val.UnbondingTime, val.UnbondingHeight)
	} else {
		k.SetUnbondingValidatorsQueue(ctx, val.UnbondingTime, val.UnbondingHeight, newAddrs)
	}
}

// ValidatorQueueIterator returns an interator ranging over validators that are
// unbonding whose unbonding completion occurs at the given height and time.
func (k Keeper) ValidatorQueueIterator(ctx sdk.Context, endTime time.Time, endHeight int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.ValidatorQueueKey, sdk.InclusiveEndBytes(types.GetValidatorQueueKey(endTime, endHeight)))
}

// UnbondAllMatureValidators unbonds all the mature unbonding validators that
// have finished their unbonding period.
func (k Keeper) UnbondAllMatureValidators(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	blockTime := ctx.BlockTime()
	blockHeight := ctx.BlockHeight()

	// unbondingValIterator will contains all validator addresses indexed under
	// the ValidatorQueueKey prefix. Note, the entire index key is composed as
	// ValidatorQueueKey | timeBzLen (8-byte big endian) | timeBz | heightBz (8-byte big endian),
	// so it may be possible that certain validator addresses that are iterated
	// over are not ready to unbond, so an explicit check is required.
	unbondingValIterator := k.ValidatorQueueIterator(ctx, blockTime, blockHeight)
	defer unbondingValIterator.Close()

	for ; unbondingValIterator.Valid(); unbondingValIterator.Next() {
		key := unbondingValIterator.Key()
		keyTime, keyHeight, err := types.ParseValidatorQueueKey(key)
		if err != nil {
			panic(fmt.Errorf("failed to parse unbonding key: %w", err))
		}

		// All addresses for the given key have the same unbonding height and time.
		// We only unbond if the height and time are less than the current height
		// and time.
		if keyHeight <= blockHeight && (keyTime.Before(blockTime) || keyTime.Equal(blockTime)) {
			addrs := types.ValAddresses{}
			k.cdc.MustUnmarshal(unbondingValIterator.Value(), &addrs)

			for _, valAddr := range addrs.Addresses {
				addr, err := sdk.ValAddressFromBech32(valAddr)
				if err != nil {
					panic(err)
				}
				val, found := k.GetValidator(ctx, addr)
				if !found {
					panic("validator in the unbonding queue was not found")
				}

				if !val.IsUnbonding() {
					panic("unexpected validator in unbonding queue; status was not unbonding")
				}

				val = k.UnbondingToUnbonded(ctx, val)
				if val.GetDelegatorShares().IsZero() {
					k.RemoveValidator(ctx, val.GetOperator())
				}
			}

			store.Delete(key)
		}
	}
}
