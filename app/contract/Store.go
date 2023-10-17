// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// BidBiderList is an auto generated low-level Go binding around an user-defined struct.
type BidBiderList struct {
	Bider  common.Address
	Amount *big.Int
	Time   *big.Int
}

// ContractABI is the input ABI used to generate the binding from.
const ContractABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BidRecord\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"BidStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"TATBiders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"bidTAT\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bidderList\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"bider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"internalType\":\"structBid.BiderList[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tatAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isTATBider\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mybidAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roundTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// ContractBin is the compiled bytecode used for deploying new contracts.
var ContractBin = "0x608060405234801561001057600080fd5b50611a3a806100206000396000f3fe6080604052600436106100955760003560e01c80638da5cb5b116100595780638da5cb5b1461018a578063bd2ea1f3146101b5578063c4d66de8146101f2578063ee5f8c8f1461021b578063f2fde38b146102465761009c565b80632b231df3146100a15780632c63b824146100de5780634b954a1d14610109578063704416b414610146578063715018a6146101735761009c565b3661009c57005b600080fd5b3480156100ad57600080fd5b506100c860048036038101906100c391906111ba565b61026f565b6040516100d59190611202565b60405180910390f35b3480156100ea57600080fd5b506100f3610723565b604051610100919061122c565b60405180910390f35b34801561011557600080fd5b50610130600480360381019061012b91906111ba565b6107db565b60405161013d9190611288565b60405180910390f35b34801561015257600080fd5b5061015b61081a565b60405161016a939291906113b2565b60405180910390f35b34801561017f57600080fd5b50610188610b41565b005b34801561019657600080fd5b5061019f610b55565b6040516101ac9190611288565b60405180910390f35b3480156101c157600080fd5b506101dc60048036038101906101d7919061141c565b610b7f565b6040516101e99190611202565b60405180910390f35b3480156101fe57600080fd5b506102196004803603810190610214919061141c565b610bf6565b005b34801561022757600080fd5b50610230610dbc565b60405161023d919061122c565b60405180910390f35b34801561025257600080fd5b5061026d6004803603810190610268919061141c565b610e5d565b005b6000670de0b6b3a76400008210156102bc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102b3906114a6565b60405180910390fd5b61012c606c546102cc91906114f5565b4211156104d65761012c80606c54426102e5919061154b565b6102ef91906115ae565b6102f991906115df565b606c5461030691906114f5565b606c81905550600060698190555060005b6065805490508110156104c657606660006065838154811061033c5761033b611639565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81549060ff021916905560676000606583815481106103c9576103c8611639565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055606860006065838154811061044a57610449611639565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000905580806104be90611668565b915050610317565b50606560006104d5919061110a565b5b606b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663adc9772e33846040518363ffffffff1660e01b81526004016105339291906116b1565b600060405180830381600087803b15801561054d57600080fd5b505af1158015610561573d6000803e3d6000fd5b5050505061056e33610b7f565b61062e576001606660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506065339080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b81606760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461067d91906114f5565b9250508190555042606860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555081606960008282546106da91906114f5565b925050819055507f76de802f12cc4b4f411dd2288c3a054af917c4c026a01f3130ae5d3c9d0c59ba33836040516107129291906116b1565b60405180910390a160019050919050565b600080606760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054141561077557600090506107d8565b61012c606c5461078591906114f5565b42111561079557600090506107d8565b606760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490505b90565b606581815481106107eb57600080fd5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60606000806000606c54905061012c606c5461083691906114f5565b4211156108db5760008067ffffffffffffffff811115610859576108586116da565b5b60405190808252806020026020018201604052801561089257816020015b61087f61112b565b8152602001906001900390816108775790505b50905061012c80606c54426108a7919061154b565b6108b191906115ae565b6108bb91906115df565b606c546108c891906114f5565b9150806000839450945094505050610b3c565b600060658054905067ffffffffffffffff8111156108fc576108fb6116da565b5b60405190808252806020026020018201604052801561093557816020015b61092261112b565b81526020019060019003908161091a5790505b50905060005b606580549050811015610b2d576065818154811061095c5761095b611639565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682828151811061099a57610999611639565b5b60200260200101516000019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505060676000606583815481106109f0576109ef611639565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054828281518110610a6957610a68611639565b5b602002602001015160200181815250506068600060658381548110610a9157610a90611639565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054828281518110610b0a57610b09611639565b5b602002602001015160400181815250508080610b2590611668565b91505061093b565b50806069548394509450945050505b909192565b610b49610ee1565b610b536000610f5f565b565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600061012c606c54610b9191906114f5565b421115610ba15760009050610bf1565b606660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690505b919050565b60008060019054906101000a900460ff16159050808015610c275750600160008054906101000a900460ff1660ff16105b80610c545750610c3630611025565b158015610c535750600160008054906101000a900460ff1660ff16145b5b610c93576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c8a9061177b565b60405180910390fd5b60016000806101000a81548160ff021916908360ff1602179055508015610cd0576001600060016101000a81548160ff0219169083151502179055505b610cd8611048565b81606b60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555042606c8190555060006069819055507f51fe6b55de8012ece012a5079ce14024411f6fa0692b05a177b94591c6ba3df043604051610d57919061122c565b60405180910390a18015610db85760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024986001604051610daf91906117ed565b60405180910390a15b5050565b600080606c5411610e02576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610df990611854565b60405180910390fd5b6000606c54905061012c606c54610e1991906114f5565b421115610e565761012c80606c5442610e32919061154b565b610e3c91906115ae565b610e4691906115df565b606c54610e5391906114f5565b90505b8091505090565b610e65610ee1565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610ed5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ecc906118e6565b60405180910390fd5b610ede81610f5f565b50565b610ee96110a1565b73ffffffffffffffffffffffffffffffffffffffff16610f07610b55565b73ffffffffffffffffffffffffffffffffffffffff1614610f5d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f5490611952565b60405180910390fd5b565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600060019054906101000a900460ff16611097576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161108e906119e4565b60405180910390fd5b61109f6110a9565b565b600033905090565b600060019054906101000a900460ff166110f8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110ef906119e4565b60405180910390fd5b6111086111036110a1565b610f5f565b565b50805460008255906000526020600020908101906111289190611162565b50565b6040518060600160405280600073ffffffffffffffffffffffffffffffffffffffff16815260200160008152602001600081525090565b5b8082111561117b576000816000905550600101611163565b5090565b600080fd5b6000819050919050565b61119781611184565b81146111a257600080fd5b50565b6000813590506111b48161118e565b92915050565b6000602082840312156111d0576111cf61117f565b5b60006111de848285016111a5565b91505092915050565b60008115159050919050565b6111fc816111e7565b82525050565b600060208201905061121760008301846111f3565b92915050565b61122681611184565b82525050565b6000602082019050611241600083018461121d565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061127282611247565b9050919050565b61128281611267565b82525050565b600060208201905061129d6000830184611279565b92915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6112d881611267565b82525050565b6112e781611184565b82525050565b60608201600082015161130360008501826112cf565b50602082015161131660208501826112de565b50604082015161132960408501826112de565b50505050565b600061133b83836112ed565b60608301905092915050565b6000602082019050919050565b600061135f826112a3565b61136981856112ae565b9350611374836112bf565b8060005b838110156113a557815161138c888261132f565b975061139783611347565b925050600181019050611378565b5085935050505092915050565b600060608201905081810360008301526113cc8186611354565b90506113db602083018561121d565b6113e8604083018461121d565b949350505050565b6113f981611267565b811461140457600080fd5b50565b600081359050611416816113f0565b92915050565b6000602082840312156114325761143161117f565b5b600061144084828501611407565b91505092915050565b600082825260208201905092915050565b7f544154206e6f7420656e6f75676820746f206269640000000000000000000000600082015250565b6000611490601583611449565b915061149b8261145a565b602082019050919050565b600060208201905081810360008301526114bf81611483565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061150082611184565b915061150b83611184565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156115405761153f6114c6565b5b828201905092915050565b600061155682611184565b915061156183611184565b925082821015611574576115736114c6565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006115b982611184565b91506115c483611184565b9250826115d4576115d361157f565b5b828204905092915050565b60006115ea82611184565b91506115f583611184565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561162e5761162d6114c6565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600061167382611184565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156116a6576116a56114c6565b5b600182019050919050565b60006040820190506116c66000830185611279565b6116d3602083018461121d565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b6000611765602e83611449565b915061177082611709565b604082019050919050565b6000602082019050818103600083015261179481611758565b9050919050565b6000819050919050565b600060ff82169050919050565b6000819050919050565b60006117d76117d26117cd8461179b565b6117b2565b6117a5565b9050919050565b6117e7816117bc565b82525050565b600060208201905061180260008301846117de565b92915050565b7f626f6e7573207374616b6520686173206e6f7420737461727465642079657400600082015250565b600061183e601f83611449565b915061184982611808565b602082019050919050565b6000602082019050818103600083015261186d81611831565b9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b60006118d0602683611449565b91506118db82611874565b604082019050919050565b600060208201905081810360008301526118ff816118c3565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061193c602083611449565b915061194782611906565b602082019050919050565b6000602082019050818103600083015261196b8161192f565b9050919050565b7f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960008201527f6e697469616c697a696e67000000000000000000000000000000000000000000602082015250565b60006119ce602b83611449565b91506119d982611972565b604082019050919050565b600060208201905081810360008301526119fd816119c1565b905091905056fea26469706673582212201185f7d44f81abb89ed03eb60a2a6919949527e86a584806674e46ff7ed764bc64736f6c63430008090033"

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// TATBiders is a free data retrieval call binding the contract method 0x4b954a1d.
//
// Solidity: function TATBiders(uint256 ) view returns(address)
func (_Contract *ContractCaller) TATBiders(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "TATBiders", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TATBiders is a free data retrieval call binding the contract method 0x4b954a1d.
//
// Solidity: function TATBiders(uint256 ) view returns(address)
func (_Contract *ContractSession) TATBiders(arg0 *big.Int) (common.Address, error) {
	return _Contract.Contract.TATBiders(&_Contract.CallOpts, arg0)
}

// TATBiders is a free data retrieval call binding the contract method 0x4b954a1d.
//
// Solidity: function TATBiders(uint256 ) view returns(address)
func (_Contract *ContractCallerSession) TATBiders(arg0 *big.Int) (common.Address, error) {
	return _Contract.Contract.TATBiders(&_Contract.CallOpts, arg0)
}

// BidderList is a free data retrieval call binding the contract method 0x704416b4.
//
// Solidity: function bidderList() view returns((address,uint256,uint256)[], uint256, uint256)
func (_Contract *ContractCaller) BidderList(opts *bind.CallOpts) ([]BidBiderList, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "bidderList")

	if err != nil {
		return *new([]BidBiderList), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]BidBiderList)).(*[]BidBiderList)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// BidderList is a free data retrieval call binding the contract method 0x704416b4.
//
// Solidity: function bidderList() view returns((address,uint256,uint256)[], uint256, uint256)
func (_Contract *ContractSession) BidderList() ([]BidBiderList, *big.Int, *big.Int, error) {
	return _Contract.Contract.BidderList(&_Contract.CallOpts)
}

// BidderList is a free data retrieval call binding the contract method 0x704416b4.
//
// Solidity: function bidderList() view returns((address,uint256,uint256)[], uint256, uint256)
func (_Contract *ContractCallerSession) BidderList() ([]BidBiderList, *big.Int, *big.Int, error) {
	return _Contract.Contract.BidderList(&_Contract.CallOpts)
}

// IsTATBider is a free data retrieval call binding the contract method 0xbd2ea1f3.
//
// Solidity: function isTATBider(address account) view returns(bool)
func (_Contract *ContractCaller) IsTATBider(opts *bind.CallOpts, account common.Address) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "isTATBider", account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTATBider is a free data retrieval call binding the contract method 0xbd2ea1f3.
//
// Solidity: function isTATBider(address account) view returns(bool)
func (_Contract *ContractSession) IsTATBider(account common.Address) (bool, error) {
	return _Contract.Contract.IsTATBider(&_Contract.CallOpts, account)
}

// IsTATBider is a free data retrieval call binding the contract method 0xbd2ea1f3.
//
// Solidity: function isTATBider(address account) view returns(bool)
func (_Contract *ContractCallerSession) IsTATBider(account common.Address) (bool, error) {
	return _Contract.Contract.IsTATBider(&_Contract.CallOpts, account)
}

// MybidAmount is a free data retrieval call binding the contract method 0x2c63b824.
//
// Solidity: function mybidAmount() view returns(uint256)
func (_Contract *ContractCaller) MybidAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "mybidAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MybidAmount is a free data retrieval call binding the contract method 0x2c63b824.
//
// Solidity: function mybidAmount() view returns(uint256)
func (_Contract *ContractSession) MybidAmount() (*big.Int, error) {
	return _Contract.Contract.MybidAmount(&_Contract.CallOpts)
}

// MybidAmount is a free data retrieval call binding the contract method 0x2c63b824.
//
// Solidity: function mybidAmount() view returns(uint256)
func (_Contract *ContractCallerSession) MybidAmount() (*big.Int, error) {
	return _Contract.Contract.MybidAmount(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCallerSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// RoundTime is a free data retrieval call binding the contract method 0xee5f8c8f.
//
// Solidity: function roundTime() view returns(uint256)
func (_Contract *ContractCaller) RoundTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "roundTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RoundTime is a free data retrieval call binding the contract method 0xee5f8c8f.
//
// Solidity: function roundTime() view returns(uint256)
func (_Contract *ContractSession) RoundTime() (*big.Int, error) {
	return _Contract.Contract.RoundTime(&_Contract.CallOpts)
}

// RoundTime is a free data retrieval call binding the contract method 0xee5f8c8f.
//
// Solidity: function roundTime() view returns(uint256)
func (_Contract *ContractCallerSession) RoundTime() (*big.Int, error) {
	return _Contract.Contract.RoundTime(&_Contract.CallOpts)
}

// BidTAT is a paid mutator transaction binding the contract method 0x2b231df3.
//
// Solidity: function bidTAT(uint256 amount) returns(bool)
func (_Contract *ContractTransactor) BidTAT(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "bidTAT", amount)
}

// BidTAT is a paid mutator transaction binding the contract method 0x2b231df3.
//
// Solidity: function bidTAT(uint256 amount) returns(bool)
func (_Contract *ContractSession) BidTAT(amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.BidTAT(&_Contract.TransactOpts, amount)
}

// BidTAT is a paid mutator transaction binding the contract method 0x2b231df3.
//
// Solidity: function bidTAT(uint256 amount) returns(bool)
func (_Contract *ContractTransactorSession) BidTAT(amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.BidTAT(&_Contract.TransactOpts, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _tatAddress) returns()
func (_Contract *ContractTransactor) Initialize(opts *bind.TransactOpts, _tatAddress common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "initialize", _tatAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _tatAddress) returns()
func (_Contract *ContractSession) Initialize(_tatAddress common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Initialize(&_Contract.TransactOpts, _tatAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _tatAddress) returns()
func (_Contract *ContractTransactorSession) Initialize(_tatAddress common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Initialize(&_Contract.TransactOpts, _tatAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractSession) Receive() (*types.Transaction, error) {
	return _Contract.Contract.Receive(&_Contract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractTransactorSession) Receive() (*types.Transaction, error) {
	return _Contract.Contract.Receive(&_Contract.TransactOpts)
}

// ContractBidRecordIterator is returned from FilterBidRecord and is used to iterate over the raw logs and unpacked data for BidRecord events raised by the Contract contract.
type ContractBidRecordIterator struct {
	Event *ContractBidRecord // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractBidRecordIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBidRecord)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractBidRecord)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractBidRecordIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBidRecordIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBidRecord represents a BidRecord event raised by the Contract contract.
type ContractBidRecord struct {
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBidRecord is a free log retrieval operation binding the contract event 0x76de802f12cc4b4f411dd2288c3a054af917c4c026a01f3130ae5d3c9d0c59ba.
//
// Solidity: event BidRecord(address account, uint256 amount)
func (_Contract *ContractFilterer) FilterBidRecord(opts *bind.FilterOpts) (*ContractBidRecordIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "BidRecord")
	if err != nil {
		return nil, err
	}
	return &ContractBidRecordIterator{contract: _Contract.contract, event: "BidRecord", logs: logs, sub: sub}, nil
}

// WatchBidRecord is a free log subscription operation binding the contract event 0x76de802f12cc4b4f411dd2288c3a054af917c4c026a01f3130ae5d3c9d0c59ba.
//
// Solidity: event BidRecord(address account, uint256 amount)
func (_Contract *ContractFilterer) WatchBidRecord(opts *bind.WatchOpts, sink chan<- *ContractBidRecord) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "BidRecord")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBidRecord)
				if err := _Contract.contract.UnpackLog(event, "BidRecord", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBidRecord is a log parse operation binding the contract event 0x76de802f12cc4b4f411dd2288c3a054af917c4c026a01f3130ae5d3c9d0c59ba.
//
// Solidity: event BidRecord(address account, uint256 amount)
func (_Contract *ContractFilterer) ParseBidRecord(log types.Log) (*ContractBidRecord, error) {
	event := new(ContractBidRecord)
	if err := _Contract.contract.UnpackLog(event, "BidRecord", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBidStartIterator is returned from FilterBidStart and is used to iterate over the raw logs and unpacked data for BidStart events raised by the Contract contract.
type ContractBidStartIterator struct {
	Event *ContractBidStart // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractBidStartIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBidStart)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractBidStart)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractBidStartIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBidStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBidStart represents a BidStart event raised by the Contract contract.
type ContractBidStart struct {
	Height *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBidStart is a free log retrieval operation binding the contract event 0x51fe6b55de8012ece012a5079ce14024411f6fa0692b05a177b94591c6ba3df0.
//
// Solidity: event BidStart(uint256 height)
func (_Contract *ContractFilterer) FilterBidStart(opts *bind.FilterOpts) (*ContractBidStartIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "BidStart")
	if err != nil {
		return nil, err
	}
	return &ContractBidStartIterator{contract: _Contract.contract, event: "BidStart", logs: logs, sub: sub}, nil
}

// WatchBidStart is a free log subscription operation binding the contract event 0x51fe6b55de8012ece012a5079ce14024411f6fa0692b05a177b94591c6ba3df0.
//
// Solidity: event BidStart(uint256 height)
func (_Contract *ContractFilterer) WatchBidStart(opts *bind.WatchOpts, sink chan<- *ContractBidStart) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "BidStart")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBidStart)
				if err := _Contract.contract.UnpackLog(event, "BidStart", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBidStart is a log parse operation binding the contract event 0x51fe6b55de8012ece012a5079ce14024411f6fa0692b05a177b94591c6ba3df0.
//
// Solidity: event BidStart(uint256 height)
func (_Contract *ContractFilterer) ParseBidStart(log types.Log) (*ContractBidStart, error) {
	event := new(ContractBidStart)
	if err := _Contract.contract.UnpackLog(event, "BidStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Contract contract.
type ContractInitializedIterator struct {
	Event *ContractInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractInitialized represents a Initialized event raised by the Contract contract.
type ContractInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Contract *ContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractInitializedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractInitializedIterator{contract: _Contract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Contract *ContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractInitialized) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractInitialized)
				if err := _Contract.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Contract *ContractFilterer) ParseInitialized(log types.Log) (*ContractInitialized, error) {
	event := new(ContractInitialized)
	if err := _Contract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Contract contract.
type ContractOwnershipTransferredIterator struct {
	Event *ContractOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOwnershipTransferred represents a OwnershipTransferred event raised by the Contract contract.
type ContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractOwnershipTransferredIterator{contract: _Contract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOwnershipTransferred)
				if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) ParseOwnershipTransferred(log types.Log) (*ContractOwnershipTransferred, error) {
	event := new(ContractOwnershipTransferred)
	if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
