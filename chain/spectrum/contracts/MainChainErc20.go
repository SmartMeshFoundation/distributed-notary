// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LockedSpectrumABI is the input ABI used to generate the binding from.
const LockedSpectrumABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"lockout\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"cancleLockOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"secret_hash\",\"type\":\"bytes32\"},{\"name\":\"expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"prepareLockoutHTLC\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lockin_htlc\",\"outputs\":[{\"name\":\"SecretHash\",\"type\":\"bytes32\"},{\"name\":\"Expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"queryLockin\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"cancelLockin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"lockin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"queryLockout\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lockout_htlc\",\"outputs\":[{\"name\":\"SecretHash\",\"type\":\"bytes32\"},{\"name\":\"Expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"secret_hash\",\"type\":\"bytes32\"},{\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"prepareLockin\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"PrepareLockin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"LockoutSecret\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"PrepareLockout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"Lockin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"CancelLockin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"CancelLockout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"}]"

// LockedSpectrumBin is the compiled bytecode used for deploying new contracts.
const LockedSpectrumBin = `0x60018054600160a060020a031916905560c0604052601d60808190527f4c6f636b6564537065637472756d20666f722061746d6f73706865726500000060a090815261004e91600291906100b8565b506040805180820190915260048082527f76302e31000000000000000000000000000000000000000000000000000000006020909201918252610093916003916100b8565b503480156100a057600080fd5b5060008054600160a060020a03191633179055610153565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100f957805160ff1916838001178555610126565b82800160010185558215610126579182015b8281111561012657825182559160200191906001019061010b565b50610132929150610136565b5090565b61015091905b80821115610132576000815560010161013c565b90565b610be1806101626000396000f3006080604052600436106100e55763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663043d918081146100f757806306fdde031461011d57806310a276eb146101a75780631a10a238146101c85780631e0ef9a4146101f257806354fd4d501461023157806357e1ee591461024657806376188aa51461026757806379ba5097146102885780637fd408d21461029d5780638caa80f7146102c15780638da5cb5b146102e257806395d89b4114610313578063a6f9dae114610328578063b852876114610349578063e0ae1a811461036a575b3480156100f157600080fd5b50600080fd5b34801561010357600080fd5b5061011b600160a060020a0360043516602435610378565b005b34801561012957600080fd5b506101326104f3565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561016c578181015183820152602001610154565b50505050905090810190601f1680156101995780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101b357600080fd5b5061011b600160a060020a036004351661057e565b3480156101d457600080fd5b5061011b600160a060020a0360043516602435604435606435610617565b3480156101fe57600080fd5b50610213600160a060020a03600435166106e4565b60408051938452602084019290925282820152519081900360600190f35b34801561023d57600080fd5b50610132610705565b34801561025257600080fd5b50610213600160a060020a0360043516610760565b34801561027357600080fd5b5061011b600160a060020a0360043516610788565b34801561029457600080fd5b5061011b610858565b3480156102a957600080fd5b5061011b600160a060020a03600435166024356108ef565b3480156102cd57600080fd5b50610213600160a060020a0360043516610a37565b3480156102ee57600080fd5b506102f7610a5f565b60408051600160a060020a039092168252519081900360200190f35b34801561031f57600080fd5b50610132610a6e565b34801561033457600080fd5b5061011b600160a060020a0360043516610aa5565b34801561035557600080fd5b50610213600160a060020a0360043516610b06565b61011b600435602435610b27565b600160a060020a03821660009081526005602052604081206002810154909181116103a257600080fd5b600182015443106103b257600080fd5b604080516020808201869052825180830382018152918301928390528151600293918291908401908083835b602083106103fd5780518252601f1990920191602091820191016103de565b51815160209384036101000a600019018019909216911617905260405191909301945091925050808303816000865af115801561043e573d6000803e3d6000fd5b5050506040513d602081101561045357600080fd5b505182541461046157600080fd5b60006002830181905580835560018301819055604051600160a060020a0386169183156108fc02918491818181858888f193505050501580156104a8573d6000803e3d6000fd5b5060408051600160a060020a03861681526020810185905281517fa0cfd4562aeab0234916ed60532417d84246c70a7f817dfc44e9c3d3423a84d3929181900390910190a150505050565b6002805460408051602060018416156101000260001901909316849004601f810184900484028201840190925281815292918301828280156105765780601f1061054b57610100808354040283529160200191610576565b820191906000526020600020905b81548152906001019060200180831161055957829003601f168201915b505050505081565b600160a060020a0381166000908152600560205260408120600281015490918082116105a957600080fd5b600183015443116105b957600080fd5b508154600060028401819055808455600184015560408051600160a060020a03861681526020810183905281517f625a628f697109c4cbee05890a8e5accf4b75c40503b6d8480b11715d148c2db929181900390910190a150505050565b60008054600160a060020a0316331461062f57600080fd5b600160a060020a038516151561064457600080fd5b50600160a060020a038416600090815260056020526040812090821161066957600080fd5b60028101541561067857600080fd5b61012c4301831161068857600080fd5b600281018290558381556001810183905560408051600160a060020a03871681526020810184905281517fbbae3304c67c8fbb052efa093374fc235534c3d862512a40007e7e35062a0475929181900390910190a15050505050565b60046020526000908152604090208054600182015460029092015490919083565b6003805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156105765780601f1061054b57610100808354040283529160200191610576565b600160a060020a03166000908152600460205260409020805460018201546002909201549092565b600160a060020a0381166000908152600460205260408120600281015490918082116107b357600080fd5b600183015443116107c357600080fd5b50815460006002840181905580845560018401819055604051600160a060020a0386169184156108fc02918591818181858888f1935050505015801561080d573d6000803e3d6000fd5b5060408051600160a060020a03861681526020810183905281517f026b98a8ac743c75f99f54b50949aa5e66574f9b73738858c62935046e4aa6c9929181900390910190a150505050565b600154600160a060020a0316331461086f57600080fd5b60005460015460408051600160a060020a03938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a1600180546000805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600160a060020a03821660009081526004602052604081206002810154909190811061091a57600080fd5b604080516020808201869052825180830382018152918301928390528151600293918291908401908083835b602083106109655780518252601f199092019160209182019101610946565b51815160209384036101000a600019018019909216911617905260405191909301945091925050808303816000865af11580156109a6573d6000803e3d6000fd5b5050506040513d60208110156109bb57600080fd5b50518254146109c957600080fd5b600182015443106109d957600080fd5b508054600060028301819055808355600183015560408051600160a060020a03861681526020810183905281517f0c89a242247566f6482a4febbbda97a1676fb18de194f38bf8f53d2d7a792c15929181900390910190a150505050565b600160a060020a03166000908152600560205260409020805460018201546002909201549092565b600054600160a060020a031681565b60408051808201909152600481527f4c534d5400000000000000000000000000000000000000000000000000000000602082015281565b600054600160a060020a03163314610abc57600080fd5b600054600160a060020a0382811691161415610ad757600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b60056020526000908152604090208054600182015460029092015490919083565b3360009081526004602052604081206002015415610b4457600080fd5b60003411610b5157600080fd5b5033600081815260046020908152604091829020858155600181018590553460028201819055835194855291840191909152815190927f1cc3ff93fb861f5fb2869fc15945f233d14ea7a4afa5721ad3c9804be90f3c6a92908290030190a15050505600a165627a7a7230582060cf399bb18209f28beca6abeea11ab5509126f8d9130976899fe43cbb21dc400029`

// DeployLockedSpectrum deploys a new Ethereum contract, binding an instance of LockedSpectrum to it.
func DeployLockedSpectrum(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LockedSpectrum, error) {
	parsed, err := abi.JSON(strings.NewReader(LockedSpectrumABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LockedSpectrumBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LockedSpectrum{LockedSpectrumCaller: LockedSpectrumCaller{contract: contract}, LockedSpectrumTransactor: LockedSpectrumTransactor{contract: contract}, LockedSpectrumFilterer: LockedSpectrumFilterer{contract: contract}}, nil
}

// LockedSpectrum is an auto generated Go binding around an Ethereum contract.
type LockedSpectrum struct {
	LockedSpectrumCaller     // Read-only binding to the contract
	LockedSpectrumTransactor // Write-only binding to the contract
	LockedSpectrumFilterer   // Log filterer for contract events
}

// LockedSpectrumCaller is an auto generated read-only Go binding around an Ethereum contract.
type LockedSpectrumCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockedSpectrumTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LockedSpectrumTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockedSpectrumFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LockedSpectrumFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockedSpectrumSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LockedSpectrumSession struct {
	Contract     *LockedSpectrum   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LockedSpectrumCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LockedSpectrumCallerSession struct {
	Contract *LockedSpectrumCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// LockedSpectrumTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LockedSpectrumTransactorSession struct {
	Contract     *LockedSpectrumTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// LockedSpectrumRaw is an auto generated low-level Go binding around an Ethereum contract.
type LockedSpectrumRaw struct {
	Contract *LockedSpectrum // Generic contract binding to access the raw methods on
}

// LockedSpectrumCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LockedSpectrumCallerRaw struct {
	Contract *LockedSpectrumCaller // Generic read-only contract binding to access the raw methods on
}

// LockedSpectrumTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LockedSpectrumTransactorRaw struct {
	Contract *LockedSpectrumTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLockedSpectrum creates a new instance of LockedSpectrum, bound to a specific deployed contract.
func NewLockedSpectrum(address common.Address, backend bind.ContractBackend) (*LockedSpectrum, error) {
	contract, err := bindLockedSpectrum(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LockedSpectrum{LockedSpectrumCaller: LockedSpectrumCaller{contract: contract}, LockedSpectrumTransactor: LockedSpectrumTransactor{contract: contract}, LockedSpectrumFilterer: LockedSpectrumFilterer{contract: contract}}, nil
}

// NewLockedSpectrumCaller creates a new read-only instance of LockedSpectrum, bound to a specific deployed contract.
func NewLockedSpectrumCaller(address common.Address, caller bind.ContractCaller) (*LockedSpectrumCaller, error) {
	contract, err := bindLockedSpectrum(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumCaller{contract: contract}, nil
}

// NewLockedSpectrumTransactor creates a new write-only instance of LockedSpectrum, bound to a specific deployed contract.
func NewLockedSpectrumTransactor(address common.Address, transactor bind.ContractTransactor) (*LockedSpectrumTransactor, error) {
	contract, err := bindLockedSpectrum(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumTransactor{contract: contract}, nil
}

// NewLockedSpectrumFilterer creates a new log filterer instance of LockedSpectrum, bound to a specific deployed contract.
func NewLockedSpectrumFilterer(address common.Address, filterer bind.ContractFilterer) (*LockedSpectrumFilterer, error) {
	contract, err := bindLockedSpectrum(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumFilterer{contract: contract}, nil
}

// bindLockedSpectrum binds a generic wrapper to an already deployed contract.
func bindLockedSpectrum(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LockedSpectrumABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockedSpectrum *LockedSpectrumRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LockedSpectrum.Contract.LockedSpectrumCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockedSpectrum *LockedSpectrumRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.LockedSpectrumTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockedSpectrum *LockedSpectrumRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.LockedSpectrumTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockedSpectrum *LockedSpectrumCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LockedSpectrum.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockedSpectrum *LockedSpectrumTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockedSpectrum *LockedSpectrumTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.contract.Transact(opts, method, params...)
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc(address ) constant returns(bytes32 SecretHash, uint256 Expiration, uint256 value)
func (_LockedSpectrum *LockedSpectrumCaller) LockinHtlc(opts *bind.CallOpts, arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	ret := new(struct {
		SecretHash [32]byte
		Expiration *big.Int
		Value      *big.Int
	})
	out := ret
	err := _LockedSpectrum.contract.Call(opts, out, "lockin_htlc", arg0)
	return *ret, err
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc(address ) constant returns(bytes32 SecretHash, uint256 Expiration, uint256 value)
func (_LockedSpectrum *LockedSpectrumSession) LockinHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _LockedSpectrum.Contract.LockinHtlc(&_LockedSpectrum.CallOpts, arg0)
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc(address ) constant returns(bytes32 SecretHash, uint256 Expiration, uint256 value)
func (_LockedSpectrum *LockedSpectrumCallerSession) LockinHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _LockedSpectrum.Contract.LockinHtlc(&_LockedSpectrum.CallOpts, arg0)
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc(address ) constant returns(bytes32 SecretHash, uint256 Expiration, uint256 value)
func (_LockedSpectrum *LockedSpectrumCaller) LockoutHtlc(opts *bind.CallOpts, arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	ret := new(struct {
		SecretHash [32]byte
		Expiration *big.Int
		Value      *big.Int
	})
	out := ret
	err := _LockedSpectrum.contract.Call(opts, out, "lockout_htlc", arg0)
	return *ret, err
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc(address ) constant returns(bytes32 SecretHash, uint256 Expiration, uint256 value)
func (_LockedSpectrum *LockedSpectrumSession) LockoutHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _LockedSpectrum.Contract.LockoutHtlc(&_LockedSpectrum.CallOpts, arg0)
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc(address ) constant returns(bytes32 SecretHash, uint256 Expiration, uint256 value)
func (_LockedSpectrum *LockedSpectrumCallerSession) LockoutHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _LockedSpectrum.Contract.LockoutHtlc(&_LockedSpectrum.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_LockedSpectrum *LockedSpectrumCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LockedSpectrum.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_LockedSpectrum *LockedSpectrumSession) Name() (string, error) {
	return _LockedSpectrum.Contract.Name(&_LockedSpectrum.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_LockedSpectrum *LockedSpectrumCallerSession) Name() (string, error) {
	return _LockedSpectrum.Contract.Name(&_LockedSpectrum.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_LockedSpectrum *LockedSpectrumCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LockedSpectrum.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_LockedSpectrum *LockedSpectrumSession) Owner() (common.Address, error) {
	return _LockedSpectrum.Contract.Owner(&_LockedSpectrum.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_LockedSpectrum *LockedSpectrumCallerSession) Owner() (common.Address, error) {
	return _LockedSpectrum.Contract.Owner(&_LockedSpectrum.CallOpts)
}

// QueryLockin is a free data retrieval call binding the contract method 0x57e1ee59.
//
// Solidity: function queryLockin(address account) constant returns(bytes32, uint256, uint256)
func (_LockedSpectrum *LockedSpectrumCaller) QueryLockin(opts *bind.CallOpts, account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _LockedSpectrum.contract.Call(opts, out, "queryLockin", account)
	return *ret0, *ret1, *ret2, err
}

// QueryLockin is a free data retrieval call binding the contract method 0x57e1ee59.
//
// Solidity: function queryLockin(address account) constant returns(bytes32, uint256, uint256)
func (_LockedSpectrum *LockedSpectrumSession) QueryLockin(account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	return _LockedSpectrum.Contract.QueryLockin(&_LockedSpectrum.CallOpts, account)
}

// QueryLockin is a free data retrieval call binding the contract method 0x57e1ee59.
//
// Solidity: function queryLockin(address account) constant returns(bytes32, uint256, uint256)
func (_LockedSpectrum *LockedSpectrumCallerSession) QueryLockin(account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	return _LockedSpectrum.Contract.QueryLockin(&_LockedSpectrum.CallOpts, account)
}

// QueryLockout is a free data retrieval call binding the contract method 0x8caa80f7.
//
// Solidity: function queryLockout(address account) constant returns(bytes32, uint256, uint256)
func (_LockedSpectrum *LockedSpectrumCaller) QueryLockout(opts *bind.CallOpts, account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _LockedSpectrum.contract.Call(opts, out, "queryLockout", account)
	return *ret0, *ret1, *ret2, err
}

// QueryLockout is a free data retrieval call binding the contract method 0x8caa80f7.
//
// Solidity: function queryLockout(address account) constant returns(bytes32, uint256, uint256)
func (_LockedSpectrum *LockedSpectrumSession) QueryLockout(account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	return _LockedSpectrum.Contract.QueryLockout(&_LockedSpectrum.CallOpts, account)
}

// QueryLockout is a free data retrieval call binding the contract method 0x8caa80f7.
//
// Solidity: function queryLockout(address account) constant returns(bytes32, uint256, uint256)
func (_LockedSpectrum *LockedSpectrumCallerSession) QueryLockout(account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	return _LockedSpectrum.Contract.QueryLockout(&_LockedSpectrum.CallOpts, account)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_LockedSpectrum *LockedSpectrumCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LockedSpectrum.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_LockedSpectrum *LockedSpectrumSession) Symbol() (string, error) {
	return _LockedSpectrum.Contract.Symbol(&_LockedSpectrum.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_LockedSpectrum *LockedSpectrumCallerSession) Symbol() (string, error) {
	return _LockedSpectrum.Contract.Symbol(&_LockedSpectrum.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_LockedSpectrum *LockedSpectrumCaller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LockedSpectrum.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_LockedSpectrum *LockedSpectrumSession) Version() (string, error) {
	return _LockedSpectrum.Contract.Version(&_LockedSpectrum.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_LockedSpectrum *LockedSpectrumCallerSession) Version() (string, error) {
	return _LockedSpectrum.Contract.Version(&_LockedSpectrum.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_LockedSpectrum *LockedSpectrumTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockedSpectrum.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_LockedSpectrum *LockedSpectrumSession) AcceptOwnership() (*types.Transaction, error) {
	return _LockedSpectrum.Contract.AcceptOwnership(&_LockedSpectrum.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_LockedSpectrum *LockedSpectrumTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LockedSpectrum.Contract.AcceptOwnership(&_LockedSpectrum.TransactOpts)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(address account) returns()
func (_LockedSpectrum *LockedSpectrumTransactor) CancelLockin(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.contract.Transact(opts, "cancelLockin", account)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(address account) returns()
func (_LockedSpectrum *LockedSpectrumSession) CancelLockin(account common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.CancelLockin(&_LockedSpectrum.TransactOpts, account)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(address account) returns()
func (_LockedSpectrum *LockedSpectrumTransactorSession) CancelLockin(account common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.CancelLockin(&_LockedSpectrum.TransactOpts, account)
}

// CancleLockOut is a paid mutator transaction binding the contract method 0x10a276eb.
//
// Solidity: function cancleLockOut(address account) returns()
func (_LockedSpectrum *LockedSpectrumTransactor) CancleLockOut(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.contract.Transact(opts, "cancleLockOut", account)
}

// CancleLockOut is a paid mutator transaction binding the contract method 0x10a276eb.
//
// Solidity: function cancleLockOut(address account) returns()
func (_LockedSpectrum *LockedSpectrumSession) CancleLockOut(account common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.CancleLockOut(&_LockedSpectrum.TransactOpts, account)
}

// CancleLockOut is a paid mutator transaction binding the contract method 0x10a276eb.
//
// Solidity: function cancleLockOut(address account) returns()
func (_LockedSpectrum *LockedSpectrumTransactorSession) CancleLockOut(account common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.CancleLockOut(&_LockedSpectrum.TransactOpts, account)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address _newOwner) returns()
func (_LockedSpectrum *LockedSpectrumTransactor) ChangeOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.contract.Transact(opts, "changeOwner", _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address _newOwner) returns()
func (_LockedSpectrum *LockedSpectrumSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.ChangeOwner(&_LockedSpectrum.TransactOpts, _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address _newOwner) returns()
func (_LockedSpectrum *LockedSpectrumTransactorSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.ChangeOwner(&_LockedSpectrum.TransactOpts, _newOwner)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(address account, bytes32 secret) returns()
func (_LockedSpectrum *LockedSpectrumTransactor) Lockin(opts *bind.TransactOpts, account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _LockedSpectrum.contract.Transact(opts, "lockin", account, secret)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(address account, bytes32 secret) returns()
func (_LockedSpectrum *LockedSpectrumSession) Lockin(account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.Lockin(&_LockedSpectrum.TransactOpts, account, secret)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(address account, bytes32 secret) returns()
func (_LockedSpectrum *LockedSpectrumTransactorSession) Lockin(account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.Lockin(&_LockedSpectrum.TransactOpts, account, secret)
}

// Lockout is a paid mutator transaction binding the contract method 0x043d9180.
//
// Solidity: function lockout(address account, bytes32 secret) returns()
func (_LockedSpectrum *LockedSpectrumTransactor) Lockout(opts *bind.TransactOpts, account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _LockedSpectrum.contract.Transact(opts, "lockout", account, secret)
}

// Lockout is a paid mutator transaction binding the contract method 0x043d9180.
//
// Solidity: function lockout(address account, bytes32 secret) returns()
func (_LockedSpectrum *LockedSpectrumSession) Lockout(account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.Lockout(&_LockedSpectrum.TransactOpts, account, secret)
}

// Lockout is a paid mutator transaction binding the contract method 0x043d9180.
//
// Solidity: function lockout(address account, bytes32 secret) returns()
func (_LockedSpectrum *LockedSpectrumTransactorSession) Lockout(account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.Lockout(&_LockedSpectrum.TransactOpts, account, secret)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0xe0ae1a81.
//
// Solidity: function prepareLockin(bytes32 secret_hash, uint256 expiration) returns()
func (_LockedSpectrum *LockedSpectrumTransactor) PrepareLockin(opts *bind.TransactOpts, secret_hash [32]byte, expiration *big.Int) (*types.Transaction, error) {
	return _LockedSpectrum.contract.Transact(opts, "prepareLockin", secret_hash, expiration)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0xe0ae1a81.
//
// Solidity: function prepareLockin(bytes32 secret_hash, uint256 expiration) returns()
func (_LockedSpectrum *LockedSpectrumSession) PrepareLockin(secret_hash [32]byte, expiration *big.Int) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.PrepareLockin(&_LockedSpectrum.TransactOpts, secret_hash, expiration)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0xe0ae1a81.
//
// Solidity: function prepareLockin(bytes32 secret_hash, uint256 expiration) returns()
func (_LockedSpectrum *LockedSpectrumTransactorSession) PrepareLockin(secret_hash [32]byte, expiration *big.Int) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.PrepareLockin(&_LockedSpectrum.TransactOpts, secret_hash, expiration)
}

// PrepareLockoutHTLC is a paid mutator transaction binding the contract method 0x1a10a238.
//
// Solidity: function prepareLockoutHTLC(address account, bytes32 secret_hash, uint256 expiration, uint256 value) returns()
func (_LockedSpectrum *LockedSpectrumTransactor) PrepareLockoutHTLC(opts *bind.TransactOpts, account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _LockedSpectrum.contract.Transact(opts, "prepareLockoutHTLC", account, secret_hash, expiration, value)
}

// PrepareLockoutHTLC is a paid mutator transaction binding the contract method 0x1a10a238.
//
// Solidity: function prepareLockoutHTLC(address account, bytes32 secret_hash, uint256 expiration, uint256 value) returns()
func (_LockedSpectrum *LockedSpectrumSession) PrepareLockoutHTLC(account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.PrepareLockoutHTLC(&_LockedSpectrum.TransactOpts, account, secret_hash, expiration, value)
}

// PrepareLockoutHTLC is a paid mutator transaction binding the contract method 0x1a10a238.
//
// Solidity: function prepareLockoutHTLC(address account, bytes32 secret_hash, uint256 expiration, uint256 value) returns()
func (_LockedSpectrum *LockedSpectrumTransactorSession) PrepareLockoutHTLC(account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _LockedSpectrum.Contract.PrepareLockoutHTLC(&_LockedSpectrum.TransactOpts, account, secret_hash, expiration, value)
}

// LockedSpectrumCancelLockinIterator is returned from FilterCancelLockin and is used to iterate over the raw logs and unpacked data for CancelLockin events raised by the LockedSpectrum contract.
type LockedSpectrumCancelLockinIterator struct {
	Event *LockedSpectrumCancelLockin // Event containing the contract specifics and raw log

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
func (it *LockedSpectrumCancelLockinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockedSpectrumCancelLockin)
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
		it.Event = new(LockedSpectrumCancelLockin)
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
func (it *LockedSpectrumCancelLockinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockedSpectrumCancelLockinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockedSpectrumCancelLockin represents a CancelLockin event raised by the LockedSpectrum contract.
type LockedSpectrumCancelLockin struct {
	Account    common.Address
	SecretHash [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancelLockin is a free log retrieval operation binding the contract event 0x026b98a8ac743c75f99f54b50949aa5e66574f9b73738858c62935046e4aa6c9.
//
// Solidity: event CancelLockin(address account, bytes32 secretHash)
func (_LockedSpectrum *LockedSpectrumFilterer) FilterCancelLockin(opts *bind.FilterOpts) (*LockedSpectrumCancelLockinIterator, error) {

	logs, sub, err := _LockedSpectrum.contract.FilterLogs(opts, "CancelLockin")
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumCancelLockinIterator{contract: _LockedSpectrum.contract, event: "CancelLockin", logs: logs, sub: sub}, nil
}

// WatchCancelLockin is a free log subscription operation binding the contract event 0x026b98a8ac743c75f99f54b50949aa5e66574f9b73738858c62935046e4aa6c9.
//
// Solidity: event CancelLockin(address account, bytes32 secretHash)
func (_LockedSpectrum *LockedSpectrumFilterer) WatchCancelLockin(opts *bind.WatchOpts, sink chan<- *LockedSpectrumCancelLockin) (event.Subscription, error) {

	logs, sub, err := _LockedSpectrum.contract.WatchLogs(opts, "CancelLockin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockedSpectrumCancelLockin)
				if err := _LockedSpectrum.contract.UnpackLog(event, "CancelLockin", log); err != nil {
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

// LockedSpectrumCancelLockoutIterator is returned from FilterCancelLockout and is used to iterate over the raw logs and unpacked data for CancelLockout events raised by the LockedSpectrum contract.
type LockedSpectrumCancelLockoutIterator struct {
	Event *LockedSpectrumCancelLockout // Event containing the contract specifics and raw log

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
func (it *LockedSpectrumCancelLockoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockedSpectrumCancelLockout)
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
		it.Event = new(LockedSpectrumCancelLockout)
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
func (it *LockedSpectrumCancelLockoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockedSpectrumCancelLockoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockedSpectrumCancelLockout represents a CancelLockout event raised by the LockedSpectrum contract.
type LockedSpectrumCancelLockout struct {
	Account    common.Address
	SecretHash [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancelLockout is a free log retrieval operation binding the contract event 0x625a628f697109c4cbee05890a8e5accf4b75c40503b6d8480b11715d148c2db.
//
// Solidity: event CancelLockout(address account, bytes32 secretHash)
func (_LockedSpectrum *LockedSpectrumFilterer) FilterCancelLockout(opts *bind.FilterOpts) (*LockedSpectrumCancelLockoutIterator, error) {

	logs, sub, err := _LockedSpectrum.contract.FilterLogs(opts, "CancelLockout")
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumCancelLockoutIterator{contract: _LockedSpectrum.contract, event: "CancelLockout", logs: logs, sub: sub}, nil
}

// WatchCancelLockout is a free log subscription operation binding the contract event 0x625a628f697109c4cbee05890a8e5accf4b75c40503b6d8480b11715d148c2db.
//
// Solidity: event CancelLockout(address account, bytes32 secretHash)
func (_LockedSpectrum *LockedSpectrumFilterer) WatchCancelLockout(opts *bind.WatchOpts, sink chan<- *LockedSpectrumCancelLockout) (event.Subscription, error) {

	logs, sub, err := _LockedSpectrum.contract.WatchLogs(opts, "CancelLockout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockedSpectrumCancelLockout)
				if err := _LockedSpectrum.contract.UnpackLog(event, "CancelLockout", log); err != nil {
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

// LockedSpectrumLockinIterator is returned from FilterLockin and is used to iterate over the raw logs and unpacked data for Lockin events raised by the LockedSpectrum contract.
type LockedSpectrumLockinIterator struct {
	Event *LockedSpectrumLockin // Event containing the contract specifics and raw log

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
func (it *LockedSpectrumLockinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockedSpectrumLockin)
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
		it.Event = new(LockedSpectrumLockin)
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
func (it *LockedSpectrumLockinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockedSpectrumLockinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockedSpectrumLockin represents a Lockin event raised by the LockedSpectrum contract.
type LockedSpectrumLockin struct {
	Account    common.Address
	SecretHash [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLockin is a free log retrieval operation binding the contract event 0x0c89a242247566f6482a4febbbda97a1676fb18de194f38bf8f53d2d7a792c15.
//
// Solidity: event Lockin(address account, bytes32 secretHash)
func (_LockedSpectrum *LockedSpectrumFilterer) FilterLockin(opts *bind.FilterOpts) (*LockedSpectrumLockinIterator, error) {

	logs, sub, err := _LockedSpectrum.contract.FilterLogs(opts, "Lockin")
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumLockinIterator{contract: _LockedSpectrum.contract, event: "Lockin", logs: logs, sub: sub}, nil
}

// WatchLockin is a free log subscription operation binding the contract event 0x0c89a242247566f6482a4febbbda97a1676fb18de194f38bf8f53d2d7a792c15.
//
// Solidity: event Lockin(address account, bytes32 secretHash)
func (_LockedSpectrum *LockedSpectrumFilterer) WatchLockin(opts *bind.WatchOpts, sink chan<- *LockedSpectrumLockin) (event.Subscription, error) {

	logs, sub, err := _LockedSpectrum.contract.WatchLogs(opts, "Lockin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockedSpectrumLockin)
				if err := _LockedSpectrum.contract.UnpackLog(event, "Lockin", log); err != nil {
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

// LockedSpectrumLockoutSecretIterator is returned from FilterLockoutSecret and is used to iterate over the raw logs and unpacked data for LockoutSecret events raised by the LockedSpectrum contract.
type LockedSpectrumLockoutSecretIterator struct {
	Event *LockedSpectrumLockoutSecret // Event containing the contract specifics and raw log

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
func (it *LockedSpectrumLockoutSecretIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockedSpectrumLockoutSecret)
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
		it.Event = new(LockedSpectrumLockoutSecret)
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
func (it *LockedSpectrumLockoutSecretIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockedSpectrumLockoutSecretIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockedSpectrumLockoutSecret represents a LockoutSecret event raised by the LockedSpectrum contract.
type LockedSpectrumLockoutSecret struct {
	Account common.Address
	Secret  [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLockoutSecret is a free log retrieval operation binding the contract event 0xa0cfd4562aeab0234916ed60532417d84246c70a7f817dfc44e9c3d3423a84d3.
//
// Solidity: event LockoutSecret(address account, bytes32 secret)
func (_LockedSpectrum *LockedSpectrumFilterer) FilterLockoutSecret(opts *bind.FilterOpts) (*LockedSpectrumLockoutSecretIterator, error) {

	logs, sub, err := _LockedSpectrum.contract.FilterLogs(opts, "LockoutSecret")
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumLockoutSecretIterator{contract: _LockedSpectrum.contract, event: "LockoutSecret", logs: logs, sub: sub}, nil
}

// WatchLockoutSecret is a free log subscription operation binding the contract event 0xa0cfd4562aeab0234916ed60532417d84246c70a7f817dfc44e9c3d3423a84d3.
//
// Solidity: event LockoutSecret(address account, bytes32 secret)
func (_LockedSpectrum *LockedSpectrumFilterer) WatchLockoutSecret(opts *bind.WatchOpts, sink chan<- *LockedSpectrumLockoutSecret) (event.Subscription, error) {

	logs, sub, err := _LockedSpectrum.contract.WatchLogs(opts, "LockoutSecret")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockedSpectrumLockoutSecret)
				if err := _LockedSpectrum.contract.UnpackLog(event, "LockoutSecret", log); err != nil {
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

// LockedSpectrumOwnerUpdateIterator is returned from FilterOwnerUpdate and is used to iterate over the raw logs and unpacked data for OwnerUpdate events raised by the LockedSpectrum contract.
type LockedSpectrumOwnerUpdateIterator struct {
	Event *LockedSpectrumOwnerUpdate // Event containing the contract specifics and raw log

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
func (it *LockedSpectrumOwnerUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockedSpectrumOwnerUpdate)
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
		it.Event = new(LockedSpectrumOwnerUpdate)
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
func (it *LockedSpectrumOwnerUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockedSpectrumOwnerUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockedSpectrumOwnerUpdate represents a OwnerUpdate event raised by the LockedSpectrum contract.
type LockedSpectrumOwnerUpdate struct {
	PrevOwner common.Address
	NewOwner  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnerUpdate is a free log retrieval operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: event OwnerUpdate(address _prevOwner, address _newOwner)
func (_LockedSpectrum *LockedSpectrumFilterer) FilterOwnerUpdate(opts *bind.FilterOpts) (*LockedSpectrumOwnerUpdateIterator, error) {

	logs, sub, err := _LockedSpectrum.contract.FilterLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumOwnerUpdateIterator{contract: _LockedSpectrum.contract, event: "OwnerUpdate", logs: logs, sub: sub}, nil
}

// WatchOwnerUpdate is a free log subscription operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: event OwnerUpdate(address _prevOwner, address _newOwner)
func (_LockedSpectrum *LockedSpectrumFilterer) WatchOwnerUpdate(opts *bind.WatchOpts, sink chan<- *LockedSpectrumOwnerUpdate) (event.Subscription, error) {

	logs, sub, err := _LockedSpectrum.contract.WatchLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockedSpectrumOwnerUpdate)
				if err := _LockedSpectrum.contract.UnpackLog(event, "OwnerUpdate", log); err != nil {
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

// LockedSpectrumPrepareLockinIterator is returned from FilterPrepareLockin and is used to iterate over the raw logs and unpacked data for PrepareLockin events raised by the LockedSpectrum contract.
type LockedSpectrumPrepareLockinIterator struct {
	Event *LockedSpectrumPrepareLockin // Event containing the contract specifics and raw log

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
func (it *LockedSpectrumPrepareLockinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockedSpectrumPrepareLockin)
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
		it.Event = new(LockedSpectrumPrepareLockin)
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
func (it *LockedSpectrumPrepareLockinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockedSpectrumPrepareLockinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockedSpectrumPrepareLockin represents a PrepareLockin event raised by the LockedSpectrum contract.
type LockedSpectrumPrepareLockin struct {
	Account common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPrepareLockin is a free log retrieval operation binding the contract event 0x1cc3ff93fb861f5fb2869fc15945f233d14ea7a4afa5721ad3c9804be90f3c6a.
//
// Solidity: event PrepareLockin(address account, uint256 value)
func (_LockedSpectrum *LockedSpectrumFilterer) FilterPrepareLockin(opts *bind.FilterOpts) (*LockedSpectrumPrepareLockinIterator, error) {

	logs, sub, err := _LockedSpectrum.contract.FilterLogs(opts, "PrepareLockin")
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumPrepareLockinIterator{contract: _LockedSpectrum.contract, event: "PrepareLockin", logs: logs, sub: sub}, nil
}

// WatchPrepareLockin is a free log subscription operation binding the contract event 0x1cc3ff93fb861f5fb2869fc15945f233d14ea7a4afa5721ad3c9804be90f3c6a.
//
// Solidity: event PrepareLockin(address account, uint256 value)
func (_LockedSpectrum *LockedSpectrumFilterer) WatchPrepareLockin(opts *bind.WatchOpts, sink chan<- *LockedSpectrumPrepareLockin) (event.Subscription, error) {

	logs, sub, err := _LockedSpectrum.contract.WatchLogs(opts, "PrepareLockin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockedSpectrumPrepareLockin)
				if err := _LockedSpectrum.contract.UnpackLog(event, "PrepareLockin", log); err != nil {
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

// LockedSpectrumPrepareLockoutIterator is returned from FilterPrepareLockout and is used to iterate over the raw logs and unpacked data for PrepareLockout events raised by the LockedSpectrum contract.
type LockedSpectrumPrepareLockoutIterator struct {
	Event *LockedSpectrumPrepareLockout // Event containing the contract specifics and raw log

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
func (it *LockedSpectrumPrepareLockoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockedSpectrumPrepareLockout)
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
		it.Event = new(LockedSpectrumPrepareLockout)
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
func (it *LockedSpectrumPrepareLockoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockedSpectrumPrepareLockoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockedSpectrumPrepareLockout represents a PrepareLockout event raised by the LockedSpectrum contract.
type LockedSpectrumPrepareLockout struct {
	Account common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPrepareLockout is a free log retrieval operation binding the contract event 0xbbae3304c67c8fbb052efa093374fc235534c3d862512a40007e7e35062a0475.
//
// Solidity: event PrepareLockout(address account, uint256 _value)
func (_LockedSpectrum *LockedSpectrumFilterer) FilterPrepareLockout(opts *bind.FilterOpts) (*LockedSpectrumPrepareLockoutIterator, error) {

	logs, sub, err := _LockedSpectrum.contract.FilterLogs(opts, "PrepareLockout")
	if err != nil {
		return nil, err
	}
	return &LockedSpectrumPrepareLockoutIterator{contract: _LockedSpectrum.contract, event: "PrepareLockout", logs: logs, sub: sub}, nil
}

// WatchPrepareLockout is a free log subscription operation binding the contract event 0xbbae3304c67c8fbb052efa093374fc235534c3d862512a40007e7e35062a0475.
//
// Solidity: event PrepareLockout(address account, uint256 _value)
func (_LockedSpectrum *LockedSpectrumFilterer) WatchPrepareLockout(opts *bind.WatchOpts, sink chan<- *LockedSpectrumPrepareLockout) (event.Subscription, error) {

	logs, sub, err := _LockedSpectrum.contract.WatchLogs(opts, "PrepareLockout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockedSpectrumPrepareLockout)
				if err := _LockedSpectrum.contract.UnpackLog(event, "PrepareLockout", log); err != nil {
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

// OwnedABI is the input ABI used to generate the binding from.
const OwnedABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"}]"

// OwnedBin is the compiled bytecode used for deploying new contracts.
const OwnedBin = `0x608060405260018054600160a060020a031916905534801561002057600080fd5b5060008054600160a060020a031916331790556101f7806100426000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166379ba5097811461005b5780638da5cb5b14610072578063a6f9dae1146100a3575b600080fd5b34801561006757600080fd5b506100706100c4565b005b34801561007e57600080fd5b5061008761015b565b60408051600160a060020a039092168252519081900360200190f35b3480156100af57600080fd5b50610070600160a060020a036004351661016a565b600154600160a060020a031633146100db57600080fd5b60005460015460408051600160a060020a03938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a1600180546000805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600054600160a060020a031681565b600054600160a060020a0316331461018157600080fd5b600054600160a060020a038281169116141561019c57600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a72305820e0929f13473862b23b70b7c53be5bdfe3c689e570bf1565d7cd30816352724c60029`

// DeployOwned deploys a new Ethereum contract, binding an instance of Owned to it.
func DeployOwned(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Owned, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// Owned is an auto generated Go binding around an Ethereum contract.
type Owned struct {
	OwnedCaller     // Read-only binding to the contract
	OwnedTransactor // Write-only binding to the contract
	OwnedFilterer   // Log filterer for contract events
}

// OwnedCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedSession struct {
	Contract     *Owned            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedCallerSession struct {
	Contract *OwnedCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OwnedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedTransactorSession struct {
	Contract     *OwnedTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedRaw struct {
	Contract *Owned // Generic contract binding to access the raw methods on
}

// OwnedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedCallerRaw struct {
	Contract *OwnedCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedTransactorRaw struct {
	Contract *OwnedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwned creates a new instance of Owned, bound to a specific deployed contract.
func NewOwned(address common.Address, backend bind.ContractBackend) (*Owned, error) {
	contract, err := bindOwned(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// NewOwnedCaller creates a new read-only instance of Owned, bound to a specific deployed contract.
func NewOwnedCaller(address common.Address, caller bind.ContractCaller) (*OwnedCaller, error) {
	contract, err := bindOwned(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedCaller{contract: contract}, nil
}

// NewOwnedTransactor creates a new write-only instance of Owned, bound to a specific deployed contract.
func NewOwnedTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTransactor, error) {
	contract, err := bindOwned(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTransactor{contract: contract}, nil
}

// NewOwnedFilterer creates a new log filterer instance of Owned, bound to a specific deployed contract.
func NewOwnedFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedFilterer, error) {
	contract, err := bindOwned(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedFilterer{contract: contract}, nil
}

// bindOwned binds a generic wrapper to an already deployed contract.
func bindOwned(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Owned.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCallerSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address _newOwner) returns()
func (_Owned *OwnedTransactor) ChangeOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "changeOwner", _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address _newOwner) returns()
func (_Owned *OwnedSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.ChangeOwner(&_Owned.TransactOpts, _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address _newOwner) returns()
func (_Owned *OwnedTransactorSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.ChangeOwner(&_Owned.TransactOpts, _newOwner)
}

// OwnedOwnerUpdateIterator is returned from FilterOwnerUpdate and is used to iterate over the raw logs and unpacked data for OwnerUpdate events raised by the Owned contract.
type OwnedOwnerUpdateIterator struct {
	Event *OwnedOwnerUpdate // Event containing the contract specifics and raw log

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
func (it *OwnedOwnerUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedOwnerUpdate)
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
		it.Event = new(OwnedOwnerUpdate)
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
func (it *OwnedOwnerUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedOwnerUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedOwnerUpdate represents a OwnerUpdate event raised by the Owned contract.
type OwnedOwnerUpdate struct {
	PrevOwner common.Address
	NewOwner  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnerUpdate is a free log retrieval operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: event OwnerUpdate(address _prevOwner, address _newOwner)
func (_Owned *OwnedFilterer) FilterOwnerUpdate(opts *bind.FilterOpts) (*OwnedOwnerUpdateIterator, error) {

	logs, sub, err := _Owned.contract.FilterLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return &OwnedOwnerUpdateIterator{contract: _Owned.contract, event: "OwnerUpdate", logs: logs, sub: sub}, nil
}

// WatchOwnerUpdate is a free log subscription operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: event OwnerUpdate(address _prevOwner, address _newOwner)
func (_Owned *OwnedFilterer) WatchOwnerUpdate(opts *bind.WatchOpts, sink chan<- *OwnedOwnerUpdate) (event.Subscription, error) {

	logs, sub, err := _Owned.contract.WatchLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedOwnerUpdate)
				if err := _Owned.contract.UnpackLog(event, "OwnerUpdate", log); err != nil {
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
