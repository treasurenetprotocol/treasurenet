// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mint

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MintABI is the input ABI used to generate the binding from.
const MintABI = "[{\"inputs\":[],\"name\":\"getCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maintat\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// MintBin is the compiled bytecode used for deploying new contracts.
var MintBin = "0x6080604052660e35fa931a000060005534801561001b57600080fd5b5060e08061002a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80639005d8a7146037578063a87d942c146051575b600080fd5b603d606b565b604051604891906091565b60405180910390f35b60576071565b604051606291906091565b60405180910390f35b60005481565b60008054905090565b6000819050919050565b608b81607a565b82525050565b600060208201905060a460008301846084565b9291505056fea264697066735822122084e035097cd827aaebe70783de90bf6505f03b59f4611fa99d831f6378ca320164736f6c63430008090033"

// DeployMint deploys a new Ethereum contract, binding an instance of Mint to it.
func DeployMint(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Mint, error) {
	parsed, err := abi.JSON(strings.NewReader(MintABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MintBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Mint{MintCaller: MintCaller{contract: contract}, MintTransactor: MintTransactor{contract: contract}, MintFilterer: MintFilterer{contract: contract}}, nil
}

// Mint is an auto generated Go binding around an Ethereum contract.
type Mint struct {
	MintCaller     // Read-only binding to the contract
	MintTransactor // Write-only binding to the contract
	MintFilterer   // Log filterer for contract events
}

// MintCaller is an auto generated read-only Go binding around an Ethereum contract.
type MintCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MintTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MintFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MintSession struct {
	Contract     *Mint             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MintCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MintCallerSession struct {
	Contract *MintCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MintTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MintTransactorSession struct {
	Contract     *MintTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MintRaw is an auto generated low-level Go binding around an Ethereum contract.
type MintRaw struct {
	Contract *Mint // Generic contract binding to access the raw methods on
}

// MintCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MintCallerRaw struct {
	Contract *MintCaller // Generic read-only contract binding to access the raw methods on
}

// MintTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MintTransactorRaw struct {
	Contract *MintTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMint creates a new instance of Mint, bound to a specific deployed contract.
func NewMint(address common.Address, backend bind.ContractBackend) (*Mint, error) {
	contract, err := bindMint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mint{MintCaller: MintCaller{contract: contract}, MintTransactor: MintTransactor{contract: contract}, MintFilterer: MintFilterer{contract: contract}}, nil
}

// NewMintCaller creates a new read-only instance of Mint, bound to a specific deployed contract.
func NewMintCaller(address common.Address, caller bind.ContractCaller) (*MintCaller, error) {
	contract, err := bindMint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MintCaller{contract: contract}, nil
}

// NewMintTransactor creates a new write-only instance of Mint, bound to a specific deployed contract.
func NewMintTransactor(address common.Address, transactor bind.ContractTransactor) (*MintTransactor, error) {
	contract, err := bindMint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MintTransactor{contract: contract}, nil
}

// NewMintFilterer creates a new log filterer instance of Mint, bound to a specific deployed contract.
func NewMintFilterer(address common.Address, filterer bind.ContractFilterer) (*MintFilterer, error) {
	contract, err := bindMint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MintFilterer{contract: contract}, nil
}

// bindMint binds a generic wrapper to an already deployed contract.
func bindMint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MintABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mint *MintRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mint.Contract.MintCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mint *MintRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mint.Contract.MintTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mint *MintRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mint.Contract.MintTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mint *MintCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mint *MintTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mint *MintTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mint.Contract.contract.Transact(opts, method, params...)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Mint *MintCaller) GetCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mint.contract.Call(opts, &out, "getCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Mint *MintSession) GetCount() (*big.Int, error) {
	return _Mint.Contract.GetCount(&_Mint.CallOpts)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Mint *MintCallerSession) GetCount() (*big.Int, error) {
	return _Mint.Contract.GetCount(&_Mint.CallOpts)
}

// Maintat is a free data retrieval call binding the contract method 0x9005d8a7.
//
// Solidity: function maintat() view returns(uint256)
func (_Mint *MintCaller) Maintat(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mint.contract.Call(opts, &out, "maintat")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Maintat is a free data retrieval call binding the contract method 0x9005d8a7.
//
// Solidity: function maintat() view returns(uint256)
func (_Mint *MintSession) Maintat() (*big.Int, error) {
	return _Mint.Contract.Maintat(&_Mint.CallOpts)
}

// Maintat is a free data retrieval call binding the contract method 0x9005d8a7.
//
// Solidity: function maintat() view returns(uint256)
func (_Mint *MintCallerSession) Maintat() (*big.Int, error) {
	return _Mint.Contract.Maintat(&_Mint.CallOpts)
}
