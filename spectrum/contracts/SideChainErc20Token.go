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

// EthereumTokenABI is the input ABI used to generate the binding from.
const EthereumTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"secret_hash\",\"type\":\"bytes32\"},{\"name\":\"expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"prePareLockedOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lockin_htlc\",\"outputs\":[{\"name\":\"SecretHash\",\"type\":\"bytes32\"},{\"name\":\"Expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"cancleLockOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"cancelLockin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"lockin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"lockedOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"secret_hash\",\"type\":\"bytes32\"},{\"name\":\"expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"prepareLockin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lockout_htlc\",\"outputs\":[{\"name\":\"SecretHash\",\"type\":\"bytes32\"},{\"name\":\"Expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"Data\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secret_hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"expiration\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"PrepareLockin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"LockinSecret\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"PrePareLockedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// EthereumTokenBin is the compiled bytecode used for deploying new contracts.
const EthereumTokenBin = `0x60028054600160a060020a031916905560c0604052601d60808190527f457468657265756d20546f6b656e20666f722061746d6f73706865726500000060a090815262000050916005919062000115565b506006805460ff1916601217905560408051808201909152600d8082527f457468657265756d546f6b656e000000000000000000000000000000000000006020909201918252620000a49160079162000115565b506040805180820190915260048082527f76302e31000000000000000000000000000000000000000000000000000000006020909201918252620000eb9160089162000115565b50348015620000f957600080fd5b5060018054600160a060020a03191633178155600055620001ba565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200015857805160ff191683800117855562000188565b8280016001018555821562000188579182015b82811115620001885782518255916020019190600101906200016b565b50620001969291506200019a565b5090565b620001b791905b80821115620001965760008155600101620001a1565b90565b610e6680620001ca6000396000f30060806040526004361061011c5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663065edecb811461012e57806306fdde0314610151578063095ea7b3146101db57806318160ddd146102135780631e0ef9a41461023a57806323b872dd14610279578063313ce567146102a357806343ba603e146102ce57806354fd4d50146102e357806370a08231146102f857806376188aa51461031957806379ba50971461033a5780637fd408d21461034f57806384540a28146103735780638da5cb5b1461039757806395d89b41146103c85780639a716549146103dd578063a6f9dae114610407578063a9059cbb14610428578063b85287611461044c578063dd62ed3e14610493575b34801561012857600080fd5b50600080fd5b34801561013a57600080fd5b5061014f6004356024356044356064356104ba565b005b34801561015d57600080fd5b5061016661056a565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101a0578181015183820152602001610188565b50505050905090810190601f1680156101cd5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101e757600080fd5b506101ff600160a060020a03600435166024356105f8565b604080519115158252519081900360200190f35b34801561021f57600080fd5b5061022861065f565b60408051918252519081900360200190f35b34801561024657600080fd5b5061025b600160a060020a0360043516610665565b60408051938452602084019290925282820152519081900360600190f35b34801561028557600080fd5b506101ff600160a060020a0360043581169060243516604435610686565b3480156102af57600080fd5b506102b861078c565b6040805160ff9092168252519081900360200190f35b3480156102da57600080fd5b5061014f610795565b3480156102ef57600080fd5b506101666107f8565b34801561030457600080fd5b50610228600160a060020a0360043516610853565b34801561032557600080fd5b5061014f600160a060020a036004351661086e565b34801561034657600080fd5b5061014f6108bc565b34801561035b57600080fd5b5061014f600160a060020a0360043516602435610953565b34801561037f57600080fd5b5061014f600160a060020a0360043516602435610abd565b3480156103a357600080fd5b506103ac610baf565b60408051600160a060020a039092168252519081900360200190f35b3480156103d457600080fd5b50610166610bbe565b3480156103e957600080fd5b5061014f600160a060020a0360043516602435604435606435610c19565b34801561041357600080fd5b5061014f600160a060020a0360043516610cd5565b34801561043457600080fd5b506101ff600160a060020a0360043516602435610d36565b34801561045857600080fd5b5061046d600160a060020a0360043516610de8565b604080519485526020850193909352838301919091526060830152519081900360800190f35b34801561049f57600080fd5b50610228600160a060020a0360043581169060243516610e0f565b336000908152600b602052604081209083116104d557600080fd5b6002810154156104e457600080fd5b3360009081526003602052604090205483111561050057600080fd5b600281018390558481556001810184905560038082018390553360008181526020928352604090819020805487900390558051868152905191927fecaa134d03b9436d19a5de74adbad850315cb95f6a84b2a343cf97588c5fe19d92918290030190a25050505050565b6005805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156105f05780601f106105c5576101008083540402835291602001916105f0565b820191906000526020600020905b8154815290600101906020018083116105d357829003601f168201915b505050505081565b336000818152600460209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a35060015b92915050565b60005481565b600a6020526000908152604090208054600182015460029092015490919083565b600160a060020a03831660009081526003602052604081205482118015906106d15750600160a060020a03841660009081526004602090815260408083203384529091529020548211155b80156106f65750600160a060020a038316600090815260036020526040902054828101115b1561078157600160a060020a03808416600081815260036020908152604080832080548801905593881680835284832080548890039055600482528483203384528252918490208054879003905583518681529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a3506001610785565b5060005b9392505050565b60065460ff1681565b336000908152600b602052604081206002810154909181116107b657600080fd5b600182015443116107c657600080fd5b600060028301819055808355600183018190556003928301819055338152602092909252604090912080549091019055565b6008805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156105f05780601f106105c5576101008083540402835291602001916105f0565b600160a060020a031660009081526003602052604090205490565b600160a060020a0381166000908152600a60205260408120600281015490911061089757600080fd5b600181015443116108a757600080fd5b60006002820181905580825560019091015550565b600254600160a060020a031633146108d357600080fd5b60015460025460408051600160a060020a03938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a1600280546001805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600160a060020a0382166000908152600a6020526040812060028101549091908190811061098057600080fd5b60408051602080820187905282518083038201815291830192839052815191929182918401908083835b602083106109c95780518252601f1990920191602091820191016109aa565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091208654149250610a0391505057600080fd5b60018301544310610a1357600080fd5b50506002810154600160a060020a038416600090815260036020526040902054808201821115610a4257600080fd5b600160a060020a0385166000908152600360205260408120828401905580548301908190558210610a7257600080fd5b60006002840181905580845560018401556040805185815290517f284a65d036cbbc9dd3c8f871862236ea013cdbf2d086ae0031b86fb27a97c12c9181900360200190a15050505050565b600160a060020a0382166000908152600b60205260408120600281015490918111610ae757600080fd5b60018201544310610af757600080fd5b60408051602080820186905282518083038201815291830192839052815191929182918401908083835b60208310610b405780518252601f199092019160209182019101610b21565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091208554149250610b7a91505057600080fd5b6000600283018190558083556001808401829055600384018290558154839003918290551115610ba957600080fd5b50505050565b600154600160a060020a031681565b6007805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156105f05780601f106105c5576101008083540402835291602001916105f0565b600154600090600160a060020a03163314610c3357600080fd5b600160a060020a0385166000908152600a602052604090206002015415610c5957600080fd5b50600160a060020a0384166000818152600a60209081526040918290208681556001810186905560028101859055825193845290830186905282820185905260608301849052905190917fa9e3e3bf2489963c99696ec61f81185fb4d9d9b09068d0bd189e6f81a2357678919081900360800190a15050505050565b600154600160a060020a03163314610cec57600080fd5b600154600160a060020a0382811691161415610d0757600080fd5b6002805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b336000908152600360205260408120548211801590610d6e5750600160a060020a038316600090815260036020526040902054828101115b15610de05733600081815260036020908152604080832080548790039055600160a060020a03871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a3506001610659565b506000610659565b600b6020526000908152604090208054600182015460028301546003909301549192909184565b600160a060020a039182166000908152600460209081526040808320939094168252919091522054905600a165627a7a723058201b8bd1d0c479449f3461223c3ff313638458c6cc2ca2c6a2d18084441d47cfee0029`

// DeployEthereumToken deploys a new Ethereum contract, binding an instance of EthereumToken to it.
func DeployEthereumToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EthereumToken, error) {
	parsed, err := abi.JSON(strings.NewReader(EthereumTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EthereumTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EthereumToken{EthereumTokenCaller: EthereumTokenCaller{contract: contract}, EthereumTokenTransactor: EthereumTokenTransactor{contract: contract}, EthereumTokenFilterer: EthereumTokenFilterer{contract: contract}}, nil
}

// EthereumToken is an auto generated Go binding around an Ethereum contract.
type EthereumToken struct {
	EthereumTokenCaller     // Read-only binding to the contract
	EthereumTokenTransactor // Write-only binding to the contract
	EthereumTokenFilterer   // Log filterer for contract events
}

// EthereumTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthereumTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthereumTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthereumTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthereumTokenSession struct {
	Contract     *EthereumToken    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthereumTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthereumTokenCallerSession struct {
	Contract *EthereumTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// EthereumTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthereumTokenTransactorSession struct {
	Contract     *EthereumTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// EthereumTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthereumTokenRaw struct {
	Contract *EthereumToken // Generic contract binding to access the raw methods on
}

// EthereumTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthereumTokenCallerRaw struct {
	Contract *EthereumTokenCaller // Generic read-only contract binding to access the raw methods on
}

// EthereumTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthereumTokenTransactorRaw struct {
	Contract *EthereumTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthereumToken creates a new instance of EthereumToken, bound to a specific deployed contract.
func NewEthereumToken(address common.Address, backend bind.ContractBackend) (*EthereumToken, error) {
	contract, err := bindEthereumToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthereumToken{EthereumTokenCaller: EthereumTokenCaller{contract: contract}, EthereumTokenTransactor: EthereumTokenTransactor{contract: contract}, EthereumTokenFilterer: EthereumTokenFilterer{contract: contract}}, nil
}

// NewEthereumTokenCaller creates a new read-only instance of EthereumToken, bound to a specific deployed contract.
func NewEthereumTokenCaller(address common.Address, caller bind.ContractCaller) (*EthereumTokenCaller, error) {
	contract, err := bindEthereumToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumTokenCaller{contract: contract}, nil
}

// NewEthereumTokenTransactor creates a new write-only instance of EthereumToken, bound to a specific deployed contract.
func NewEthereumTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*EthereumTokenTransactor, error) {
	contract, err := bindEthereumToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumTokenTransactor{contract: contract}, nil
}

// NewEthereumTokenFilterer creates a new log filterer instance of EthereumToken, bound to a specific deployed contract.
func NewEthereumTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*EthereumTokenFilterer, error) {
	contract, err := bindEthereumToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthereumTokenFilterer{contract: contract}, nil
}

// bindEthereumToken binds a generic wrapper to an already deployed contract.
func bindEthereumToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthereumTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthereumToken *EthereumTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EthereumToken.Contract.EthereumTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthereumToken *EthereumTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumToken.Contract.EthereumTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthereumToken *EthereumTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthereumToken.Contract.EthereumTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthereumToken *EthereumTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EthereumToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthereumToken *EthereumTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthereumToken *EthereumTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthereumToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EthereumToken *EthereumTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EthereumToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EthereumToken *EthereumTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EthereumToken.Contract.Allowance(&_EthereumToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EthereumToken *EthereumTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EthereumToken.Contract.Allowance(&_EthereumToken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EthereumToken *EthereumTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EthereumToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EthereumToken *EthereumTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EthereumToken.Contract.BalanceOf(&_EthereumToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EthereumToken *EthereumTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EthereumToken.Contract.BalanceOf(&_EthereumToken.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_EthereumToken *EthereumTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _EthereumToken.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_EthereumToken *EthereumTokenSession) Decimals() (uint8, error) {
	return _EthereumToken.Contract.Decimals(&_EthereumToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_EthereumToken *EthereumTokenCallerSession) Decimals() (uint8, error) {
	return _EthereumToken.Contract.Decimals(&_EthereumToken.CallOpts)
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_EthereumToken *EthereumTokenCaller) LockinHtlc(opts *bind.CallOpts, arg0 common.Address) (struct {
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
	err := _EthereumToken.contract.Call(opts, out, "lockin_htlc", arg0)
	return *ret, err
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_EthereumToken *EthereumTokenSession) LockinHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _EthereumToken.Contract.LockinHtlc(&_EthereumToken.CallOpts, arg0)
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_EthereumToken *EthereumTokenCallerSession) LockinHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _EthereumToken.Contract.LockinHtlc(&_EthereumToken.CallOpts, arg0)
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256, Data bytes32)
func (_EthereumToken *EthereumTokenCaller) LockoutHtlc(opts *bind.CallOpts, arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
	Data       [32]byte
}, error) {
	ret := new(struct {
		SecretHash [32]byte
		Expiration *big.Int
		Value      *big.Int
		Data       [32]byte
	})
	out := ret
	err := _EthereumToken.contract.Call(opts, out, "lockout_htlc", arg0)
	return *ret, err
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256, Data bytes32)
func (_EthereumToken *EthereumTokenSession) LockoutHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
	Data       [32]byte
}, error) {
	return _EthereumToken.Contract.LockoutHtlc(&_EthereumToken.CallOpts, arg0)
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256, Data bytes32)
func (_EthereumToken *EthereumTokenCallerSession) LockoutHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
	Data       [32]byte
}, error) {
	return _EthereumToken.Contract.LockoutHtlc(&_EthereumToken.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_EthereumToken *EthereumTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _EthereumToken.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_EthereumToken *EthereumTokenSession) Name() (string, error) {
	return _EthereumToken.Contract.Name(&_EthereumToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_EthereumToken *EthereumTokenCallerSession) Name() (string, error) {
	return _EthereumToken.Contract.Name(&_EthereumToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_EthereumToken *EthereumTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _EthereumToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_EthereumToken *EthereumTokenSession) Owner() (common.Address, error) {
	return _EthereumToken.Contract.Owner(&_EthereumToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_EthereumToken *EthereumTokenCallerSession) Owner() (common.Address, error) {
	return _EthereumToken.Contract.Owner(&_EthereumToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_EthereumToken *EthereumTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _EthereumToken.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_EthereumToken *EthereumTokenSession) Symbol() (string, error) {
	return _EthereumToken.Contract.Symbol(&_EthereumToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_EthereumToken *EthereumTokenCallerSession) Symbol() (string, error) {
	return _EthereumToken.Contract.Symbol(&_EthereumToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EthereumToken *EthereumTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EthereumToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EthereumToken *EthereumTokenSession) TotalSupply() (*big.Int, error) {
	return _EthereumToken.Contract.TotalSupply(&_EthereumToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EthereumToken *EthereumTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _EthereumToken.Contract.TotalSupply(&_EthereumToken.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_EthereumToken *EthereumTokenCaller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _EthereumToken.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_EthereumToken *EthereumTokenSession) Version() (string, error) {
	return _EthereumToken.Contract.Version(&_EthereumToken.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_EthereumToken *EthereumTokenCallerSession) Version() (string, error) {
	return _EthereumToken.Contract.Version(&_EthereumToken.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_EthereumToken *EthereumTokenTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_EthereumToken *EthereumTokenSession) AcceptOwnership() (*types.Transaction, error) {
	return _EthereumToken.Contract.AcceptOwnership(&_EthereumToken.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_EthereumToken *EthereumTokenTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _EthereumToken.Contract.AcceptOwnership(&_EthereumToken.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.Contract.Approve(&_EthereumToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.Contract.Approve(&_EthereumToken.TransactOpts, _spender, _value)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(account address) returns()
func (_EthereumToken *EthereumTokenTransactor) CancelLockin(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "cancelLockin", account)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(account address) returns()
func (_EthereumToken *EthereumTokenSession) CancelLockin(account common.Address) (*types.Transaction, error) {
	return _EthereumToken.Contract.CancelLockin(&_EthereumToken.TransactOpts, account)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(account address) returns()
func (_EthereumToken *EthereumTokenTransactorSession) CancelLockin(account common.Address) (*types.Transaction, error) {
	return _EthereumToken.Contract.CancelLockin(&_EthereumToken.TransactOpts, account)
}

// CancleLockOut is a paid mutator transaction binding the contract method 0x43ba603e.
//
// Solidity: function cancleLockOut() returns()
func (_EthereumToken *EthereumTokenTransactor) CancleLockOut(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "cancleLockOut")
}

// CancleLockOut is a paid mutator transaction binding the contract method 0x43ba603e.
//
// Solidity: function cancleLockOut() returns()
func (_EthereumToken *EthereumTokenSession) CancleLockOut() (*types.Transaction, error) {
	return _EthereumToken.Contract.CancleLockOut(&_EthereumToken.TransactOpts)
}

// CancleLockOut is a paid mutator transaction binding the contract method 0x43ba603e.
//
// Solidity: function cancleLockOut() returns()
func (_EthereumToken *EthereumTokenTransactorSession) CancleLockOut() (*types.Transaction, error) {
	return _EthereumToken.Contract.CancleLockOut(&_EthereumToken.TransactOpts)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_EthereumToken *EthereumTokenTransactor) ChangeOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "changeOwner", _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_EthereumToken *EthereumTokenSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _EthereumToken.Contract.ChangeOwner(&_EthereumToken.TransactOpts, _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_EthereumToken *EthereumTokenTransactorSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _EthereumToken.Contract.ChangeOwner(&_EthereumToken.TransactOpts, _newOwner)
}

// LockedOut is a paid mutator transaction binding the contract method 0x84540a28.
//
// Solidity: function lockedOut(from address, secret bytes32) returns()
func (_EthereumToken *EthereumTokenTransactor) LockedOut(opts *bind.TransactOpts, from common.Address, secret [32]byte) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "lockedOut", from, secret)
}

// LockedOut is a paid mutator transaction binding the contract method 0x84540a28.
//
// Solidity: function lockedOut(from address, secret bytes32) returns()
func (_EthereumToken *EthereumTokenSession) LockedOut(from common.Address, secret [32]byte) (*types.Transaction, error) {
	return _EthereumToken.Contract.LockedOut(&_EthereumToken.TransactOpts, from, secret)
}

// LockedOut is a paid mutator transaction binding the contract method 0x84540a28.
//
// Solidity: function lockedOut(from address, secret bytes32) returns()
func (_EthereumToken *EthereumTokenTransactorSession) LockedOut(from common.Address, secret [32]byte) (*types.Transaction, error) {
	return _EthereumToken.Contract.LockedOut(&_EthereumToken.TransactOpts, from, secret)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(account address, secret bytes32) returns()
func (_EthereumToken *EthereumTokenTransactor) Lockin(opts *bind.TransactOpts, account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "lockin", account, secret)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(account address, secret bytes32) returns()
func (_EthereumToken *EthereumTokenSession) Lockin(account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _EthereumToken.Contract.Lockin(&_EthereumToken.TransactOpts, account, secret)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(account address, secret bytes32) returns()
func (_EthereumToken *EthereumTokenTransactorSession) Lockin(account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _EthereumToken.Contract.Lockin(&_EthereumToken.TransactOpts, account, secret)
}

// PrePareLockedOut is a paid mutator transaction binding the contract method 0x065edecb.
//
// Solidity: function prePareLockedOut(secret_hash bytes32, expiration uint256, value uint256, data bytes32) returns()
func (_EthereumToken *EthereumTokenTransactor) PrePareLockedOut(opts *bind.TransactOpts, secret_hash [32]byte, expiration *big.Int, value *big.Int, data [32]byte) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "prePareLockedOut", secret_hash, expiration, value, data)
}

// PrePareLockedOut is a paid mutator transaction binding the contract method 0x065edecb.
//
// Solidity: function prePareLockedOut(secret_hash bytes32, expiration uint256, value uint256, data bytes32) returns()
func (_EthereumToken *EthereumTokenSession) PrePareLockedOut(secret_hash [32]byte, expiration *big.Int, value *big.Int, data [32]byte) (*types.Transaction, error) {
	return _EthereumToken.Contract.PrePareLockedOut(&_EthereumToken.TransactOpts, secret_hash, expiration, value, data)
}

// PrePareLockedOut is a paid mutator transaction binding the contract method 0x065edecb.
//
// Solidity: function prePareLockedOut(secret_hash bytes32, expiration uint256, value uint256, data bytes32) returns()
func (_EthereumToken *EthereumTokenTransactorSession) PrePareLockedOut(secret_hash [32]byte, expiration *big.Int, value *big.Int, data [32]byte) (*types.Transaction, error) {
	return _EthereumToken.Contract.PrePareLockedOut(&_EthereumToken.TransactOpts, secret_hash, expiration, value, data)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0x9a716549.
//
// Solidity: function prepareLockin(account address, secret_hash bytes32, expiration uint256, value uint256) returns()
func (_EthereumToken *EthereumTokenTransactor) PrepareLockin(opts *bind.TransactOpts, account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "prepareLockin", account, secret_hash, expiration, value)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0x9a716549.
//
// Solidity: function prepareLockin(account address, secret_hash bytes32, expiration uint256, value uint256) returns()
func (_EthereumToken *EthereumTokenSession) PrepareLockin(account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.Contract.PrepareLockin(&_EthereumToken.TransactOpts, account, secret_hash, expiration, value)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0x9a716549.
//
// Solidity: function prepareLockin(account address, secret_hash bytes32, expiration uint256, value uint256) returns()
func (_EthereumToken *EthereumTokenTransactorSession) PrepareLockin(account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.Contract.PrepareLockin(&_EthereumToken.TransactOpts, account, secret_hash, expiration, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.Contract.Transfer(&_EthereumToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.Contract.Transfer(&_EthereumToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.Contract.TransferFrom(&_EthereumToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EthereumToken *EthereumTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthereumToken.Contract.TransferFrom(&_EthereumToken.TransactOpts, _from, _to, _value)
}

// EthereumTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the EthereumToken contract.
type EthereumTokenApprovalIterator struct {
	Event *EthereumTokenApproval // Event containing the contract specifics and raw log

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
func (it *EthereumTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumTokenApproval)
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
		it.Event = new(EthereumTokenApproval)
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
func (it *EthereumTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumTokenApproval represents a Approval event raised by the EthereumToken contract.
type EthereumTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_EthereumToken *EthereumTokenFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*EthereumTokenApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _EthereumToken.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &EthereumTokenApprovalIterator{contract: _EthereumToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_EthereumToken *EthereumTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *EthereumTokenApproval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _EthereumToken.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumTokenApproval)
				if err := _EthereumToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// EthereumTokenLockinSecretIterator is returned from FilterLockinSecret and is used to iterate over the raw logs and unpacked data for LockinSecret events raised by the EthereumToken contract.
type EthereumTokenLockinSecretIterator struct {
	Event *EthereumTokenLockinSecret // Event containing the contract specifics and raw log

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
func (it *EthereumTokenLockinSecretIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumTokenLockinSecret)
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
		it.Event = new(EthereumTokenLockinSecret)
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
func (it *EthereumTokenLockinSecretIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumTokenLockinSecretIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumTokenLockinSecret represents a LockinSecret event raised by the EthereumToken contract.
type EthereumTokenLockinSecret struct {
	Secret [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLockinSecret is a free log retrieval operation binding the contract event 0x284a65d036cbbc9dd3c8f871862236ea013cdbf2d086ae0031b86fb27a97c12c.
//
// Solidity: e LockinSecret(secret bytes32)
func (_EthereumToken *EthereumTokenFilterer) FilterLockinSecret(opts *bind.FilterOpts) (*EthereumTokenLockinSecretIterator, error) {

	logs, sub, err := _EthereumToken.contract.FilterLogs(opts, "LockinSecret")
	if err != nil {
		return nil, err
	}
	return &EthereumTokenLockinSecretIterator{contract: _EthereumToken.contract, event: "LockinSecret", logs: logs, sub: sub}, nil
}

// WatchLockinSecret is a free log subscription operation binding the contract event 0x284a65d036cbbc9dd3c8f871862236ea013cdbf2d086ae0031b86fb27a97c12c.
//
// Solidity: e LockinSecret(secret bytes32)
func (_EthereumToken *EthereumTokenFilterer) WatchLockinSecret(opts *bind.WatchOpts, sink chan<- *EthereumTokenLockinSecret) (event.Subscription, error) {

	logs, sub, err := _EthereumToken.contract.WatchLogs(opts, "LockinSecret")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumTokenLockinSecret)
				if err := _EthereumToken.contract.UnpackLog(event, "LockinSecret", log); err != nil {
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

// EthereumTokenOwnerUpdateIterator is returned from FilterOwnerUpdate and is used to iterate over the raw logs and unpacked data for OwnerUpdate events raised by the EthereumToken contract.
type EthereumTokenOwnerUpdateIterator struct {
	Event *EthereumTokenOwnerUpdate // Event containing the contract specifics and raw log

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
func (it *EthereumTokenOwnerUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumTokenOwnerUpdate)
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
		it.Event = new(EthereumTokenOwnerUpdate)
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
func (it *EthereumTokenOwnerUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumTokenOwnerUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumTokenOwnerUpdate represents a OwnerUpdate event raised by the EthereumToken contract.
type EthereumTokenOwnerUpdate struct {
	PrevOwner common.Address
	NewOwner  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnerUpdate is a free log retrieval operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_EthereumToken *EthereumTokenFilterer) FilterOwnerUpdate(opts *bind.FilterOpts) (*EthereumTokenOwnerUpdateIterator, error) {

	logs, sub, err := _EthereumToken.contract.FilterLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return &EthereumTokenOwnerUpdateIterator{contract: _EthereumToken.contract, event: "OwnerUpdate", logs: logs, sub: sub}, nil
}

// WatchOwnerUpdate is a free log subscription operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_EthereumToken *EthereumTokenFilterer) WatchOwnerUpdate(opts *bind.WatchOpts, sink chan<- *EthereumTokenOwnerUpdate) (event.Subscription, error) {

	logs, sub, err := _EthereumToken.contract.WatchLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumTokenOwnerUpdate)
				if err := _EthereumToken.contract.UnpackLog(event, "OwnerUpdate", log); err != nil {
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

// EthereumTokenPrePareLockedOutIterator is returned from FilterPrePareLockedOut and is used to iterate over the raw logs and unpacked data for PrePareLockedOut events raised by the EthereumToken contract.
type EthereumTokenPrePareLockedOutIterator struct {
	Event *EthereumTokenPrePareLockedOut // Event containing the contract specifics and raw log

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
func (it *EthereumTokenPrePareLockedOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumTokenPrePareLockedOut)
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
		it.Event = new(EthereumTokenPrePareLockedOut)
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
func (it *EthereumTokenPrePareLockedOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumTokenPrePareLockedOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumTokenPrePareLockedOut represents a PrePareLockedOut event raised by the EthereumToken contract.
type EthereumTokenPrePareLockedOut struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterPrePareLockedOut is a free log retrieval operation binding the contract event 0xecaa134d03b9436d19a5de74adbad850315cb95f6a84b2a343cf97588c5fe19d.
//
// Solidity: e PrePareLockedOut(_from indexed address, _value uint256)
func (_EthereumToken *EthereumTokenFilterer) FilterPrePareLockedOut(opts *bind.FilterOpts, _from []common.Address) (*EthereumTokenPrePareLockedOutIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _EthereumToken.contract.FilterLogs(opts, "PrePareLockedOut", _fromRule)
	if err != nil {
		return nil, err
	}
	return &EthereumTokenPrePareLockedOutIterator{contract: _EthereumToken.contract, event: "PrePareLockedOut", logs: logs, sub: sub}, nil
}

// WatchPrePareLockedOut is a free log subscription operation binding the contract event 0xecaa134d03b9436d19a5de74adbad850315cb95f6a84b2a343cf97588c5fe19d.
//
// Solidity: e PrePareLockedOut(_from indexed address, _value uint256)
func (_EthereumToken *EthereumTokenFilterer) WatchPrePareLockedOut(opts *bind.WatchOpts, sink chan<- *EthereumTokenPrePareLockedOut, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _EthereumToken.contract.WatchLogs(opts, "PrePareLockedOut", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumTokenPrePareLockedOut)
				if err := _EthereumToken.contract.UnpackLog(event, "PrePareLockedOut", log); err != nil {
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

// EthereumTokenPrepareLockinIterator is returned from FilterPrepareLockin and is used to iterate over the raw logs and unpacked data for PrepareLockin events raised by the EthereumToken contract.
type EthereumTokenPrepareLockinIterator struct {
	Event *EthereumTokenPrepareLockin // Event containing the contract specifics and raw log

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
func (it *EthereumTokenPrepareLockinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumTokenPrepareLockin)
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
		it.Event = new(EthereumTokenPrepareLockin)
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
func (it *EthereumTokenPrepareLockinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumTokenPrepareLockinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumTokenPrepareLockin represents a PrepareLockin event raised by the EthereumToken contract.
type EthereumTokenPrepareLockin struct {
	Account    common.Address
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPrepareLockin is a free log retrieval operation binding the contract event 0xa9e3e3bf2489963c99696ec61f81185fb4d9d9b09068d0bd189e6f81a2357678.
//
// Solidity: e PrepareLockin(account address, secret_hash bytes32, expiration uint256, value uint256)
func (_EthereumToken *EthereumTokenFilterer) FilterPrepareLockin(opts *bind.FilterOpts) (*EthereumTokenPrepareLockinIterator, error) {

	logs, sub, err := _EthereumToken.contract.FilterLogs(opts, "PrepareLockin")
	if err != nil {
		return nil, err
	}
	return &EthereumTokenPrepareLockinIterator{contract: _EthereumToken.contract, event: "PrepareLockin", logs: logs, sub: sub}, nil
}

// WatchPrepareLockin is a free log subscription operation binding the contract event 0xa9e3e3bf2489963c99696ec61f81185fb4d9d9b09068d0bd189e6f81a2357678.
//
// Solidity: e PrepareLockin(account address, secret_hash bytes32, expiration uint256, value uint256)
func (_EthereumToken *EthereumTokenFilterer) WatchPrepareLockin(opts *bind.WatchOpts, sink chan<- *EthereumTokenPrepareLockin) (event.Subscription, error) {

	logs, sub, err := _EthereumToken.contract.WatchLogs(opts, "PrepareLockin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumTokenPrepareLockin)
				if err := _EthereumToken.contract.UnpackLog(event, "PrepareLockin", log); err != nil {
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

// EthereumTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the EthereumToken contract.
type EthereumTokenTransferIterator struct {
	Event *EthereumTokenTransfer // Event containing the contract specifics and raw log

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
func (it *EthereumTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumTokenTransfer)
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
		it.Event = new(EthereumTokenTransfer)
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
func (it *EthereumTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumTokenTransfer represents a Transfer event raised by the EthereumToken contract.
type EthereumTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_EthereumToken *EthereumTokenFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*EthereumTokenTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _EthereumToken.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &EthereumTokenTransferIterator{contract: _EthereumToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_EthereumToken *EthereumTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *EthereumTokenTransfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _EthereumToken.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumTokenTransfer)
				if err := _EthereumToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
const OwnedBin = `0x608060405260018054600160a060020a031916905534801561002057600080fd5b5060008054600160a060020a031916331790556101f7806100426000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166379ba5097811461005b5780638da5cb5b14610072578063a6f9dae1146100a3575b600080fd5b34801561006757600080fd5b506100706100c4565b005b34801561007e57600080fd5b5061008761015b565b60408051600160a060020a039092168252519081900360200190f35b3480156100af57600080fd5b50610070600160a060020a036004351661016a565b600154600160a060020a031633146100db57600080fd5b60005460015460408051600160a060020a03938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a1600180546000805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600054600160a060020a031681565b600054600160a060020a0316331461018157600080fd5b600054600160a060020a038281169116141561019c57600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a7230582078b9a86f958becf29dde5e2d342215fb11fccbe0c1f46d1c664315b2f1bc4aef0029`

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
// Solidity: function changeOwner(_newOwner address) returns()
func (_Owned *OwnedTransactor) ChangeOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "changeOwner", _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_Owned *OwnedSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.ChangeOwner(&_Owned.TransactOpts, _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
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
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_Owned *OwnedFilterer) FilterOwnerUpdate(opts *bind.FilterOpts) (*OwnedOwnerUpdateIterator, error) {

	logs, sub, err := _Owned.contract.FilterLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return &OwnedOwnerUpdateIterator{contract: _Owned.contract, event: "OwnerUpdate", logs: logs, sub: sub}, nil
}

// WatchOwnerUpdate is a free log subscription operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
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

// StandardTokenABI is the input ABI used to generate the binding from.
const StandardTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// StandardTokenBin is the compiled bytecode used for deploying new contracts.
const StandardTokenBin = `0x608060405260028054600160a060020a03199081169091556001805490911633179055610599806100316000396000f3006080604052600436106100985763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b3811461009d57806318160ddd146100d557806323b872dd146100fc57806370a082311461012657806379ba5097146101475780638da5cb5b1461015e578063a6f9dae11461018f578063a9059cbb146101b0578063dd62ed3e146101d4575b600080fd5b3480156100a957600080fd5b506100c1600160a060020a03600435166024356101fb565b604080519115158252519081900360200190f35b3480156100e157600080fd5b506100ea610262565b60408051918252519081900360200190f35b34801561010857600080fd5b506100c1600160a060020a0360043581169060243516604435610268565b34801561013257600080fd5b506100ea600160a060020a036004351661036e565b34801561015357600080fd5b5061015c610389565b005b34801561016a57600080fd5b50610173610420565b60408051600160a060020a039092168252519081900360200190f35b34801561019b57600080fd5b5061015c600160a060020a036004351661042f565b3480156101bc57600080fd5b506100c1600160a060020a0360043516602435610490565b3480156101e057600080fd5b506100ea600160a060020a0360043581169060243516610542565b336000818152600460209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a35060015b92915050565b60005481565b600160a060020a03831660009081526003602052604081205482118015906102b35750600160a060020a03841660009081526004602090815260408083203384529091529020548211155b80156102d85750600160a060020a038316600090815260036020526040902054828101115b1561036357600160a060020a03808416600081815260036020908152604080832080548801905593881680835284832080548890039055600482528483203384528252918490208054879003905583518681529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a3506001610367565b5060005b9392505050565b600160a060020a031660009081526003602052604090205490565b600254600160a060020a031633146103a057600080fd5b60015460025460408051600160a060020a03938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a1600280546001805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600154600160a060020a031681565b600154600160a060020a0316331461044657600080fd5b600154600160a060020a038281169116141561046157600080fd5b6002805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b3360009081526003602052604081205482118015906104c85750600160a060020a038316600090815260036020526040902054828101115b1561053a5733600081815260036020908152604080832080548790039055600160a060020a03871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a350600161025c565b50600061025c565b600160a060020a039182166000908152600460209081526040808320939094168252919091522054905600a165627a7a72305820b6b5191ccc1749e37ec0919c5358258f6297b7e61695e4620e88e3961db880650029`

// DeployStandardToken deploys a new Ethereum contract, binding an instance of StandardToken to it.
func DeployStandardToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StandardToken, error) {
	parsed, err := abi.JSON(strings.NewReader(StandardTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StandardTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StandardToken{StandardTokenCaller: StandardTokenCaller{contract: contract}, StandardTokenTransactor: StandardTokenTransactor{contract: contract}, StandardTokenFilterer: StandardTokenFilterer{contract: contract}}, nil
}

// StandardToken is an auto generated Go binding around an Ethereum contract.
type StandardToken struct {
	StandardTokenCaller     // Read-only binding to the contract
	StandardTokenTransactor // Write-only binding to the contract
	StandardTokenFilterer   // Log filterer for contract events
}

// StandardTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type StandardTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StandardTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StandardTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StandardTokenSession struct {
	Contract     *StandardToken    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StandardTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StandardTokenCallerSession struct {
	Contract *StandardTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StandardTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StandardTokenTransactorSession struct {
	Contract     *StandardTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StandardTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type StandardTokenRaw struct {
	Contract *StandardToken // Generic contract binding to access the raw methods on
}

// StandardTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StandardTokenCallerRaw struct {
	Contract *StandardTokenCaller // Generic read-only contract binding to access the raw methods on
}

// StandardTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StandardTokenTransactorRaw struct {
	Contract *StandardTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStandardToken creates a new instance of StandardToken, bound to a specific deployed contract.
func NewStandardToken(address common.Address, backend bind.ContractBackend) (*StandardToken, error) {
	contract, err := bindStandardToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StandardToken{StandardTokenCaller: StandardTokenCaller{contract: contract}, StandardTokenTransactor: StandardTokenTransactor{contract: contract}, StandardTokenFilterer: StandardTokenFilterer{contract: contract}}, nil
}

// NewStandardTokenCaller creates a new read-only instance of StandardToken, bound to a specific deployed contract.
func NewStandardTokenCaller(address common.Address, caller bind.ContractCaller) (*StandardTokenCaller, error) {
	contract, err := bindStandardToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StandardTokenCaller{contract: contract}, nil
}

// NewStandardTokenTransactor creates a new write-only instance of StandardToken, bound to a specific deployed contract.
func NewStandardTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*StandardTokenTransactor, error) {
	contract, err := bindStandardToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StandardTokenTransactor{contract: contract}, nil
}

// NewStandardTokenFilterer creates a new log filterer instance of StandardToken, bound to a specific deployed contract.
func NewStandardTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*StandardTokenFilterer, error) {
	contract, err := bindStandardToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StandardTokenFilterer{contract: contract}, nil
}

// bindStandardToken binds a generic wrapper to an already deployed contract.
func bindStandardToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StandardTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StandardToken *StandardTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StandardToken.Contract.StandardTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StandardToken *StandardTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StandardToken.Contract.StandardTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StandardToken *StandardTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StandardToken.Contract.StandardTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StandardToken *StandardTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StandardToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StandardToken *StandardTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StandardToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StandardToken *StandardTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StandardToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_StandardToken *StandardTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StandardToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_StandardToken *StandardTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StandardToken.Contract.Allowance(&_StandardToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_StandardToken *StandardTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StandardToken.Contract.Allowance(&_StandardToken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_StandardToken *StandardTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StandardToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_StandardToken *StandardTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _StandardToken.Contract.BalanceOf(&_StandardToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_StandardToken *StandardTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _StandardToken.Contract.BalanceOf(&_StandardToken.CallOpts, _owner)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_StandardToken *StandardTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StandardToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_StandardToken *StandardTokenSession) Owner() (common.Address, error) {
	return _StandardToken.Contract.Owner(&_StandardToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_StandardToken *StandardTokenCallerSession) Owner() (common.Address, error) {
	return _StandardToken.Contract.Owner(&_StandardToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_StandardToken *StandardTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StandardToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_StandardToken *StandardTokenSession) TotalSupply() (*big.Int, error) {
	return _StandardToken.Contract.TotalSupply(&_StandardToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_StandardToken *StandardTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _StandardToken.Contract.TotalSupply(&_StandardToken.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_StandardToken *StandardTokenTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_StandardToken *StandardTokenSession) AcceptOwnership() (*types.Transaction, error) {
	return _StandardToken.Contract.AcceptOwnership(&_StandardToken.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_StandardToken *StandardTokenTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _StandardToken.Contract.AcceptOwnership(&_StandardToken.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.Approve(&_StandardToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.Approve(&_StandardToken.TransactOpts, _spender, _value)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_StandardToken *StandardTokenTransactor) ChangeOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "changeOwner", _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_StandardToken *StandardTokenSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _StandardToken.Contract.ChangeOwner(&_StandardToken.TransactOpts, _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_StandardToken *StandardTokenTransactorSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _StandardToken.Contract.ChangeOwner(&_StandardToken.TransactOpts, _newOwner)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.Transfer(&_StandardToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.Transfer(&_StandardToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.TransferFrom(&_StandardToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_StandardToken *StandardTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.TransferFrom(&_StandardToken.TransactOpts, _from, _to, _value)
}

// StandardTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the StandardToken contract.
type StandardTokenApprovalIterator struct {
	Event *StandardTokenApproval // Event containing the contract specifics and raw log

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
func (it *StandardTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StandardTokenApproval)
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
		it.Event = new(StandardTokenApproval)
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
func (it *StandardTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StandardTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StandardTokenApproval represents a Approval event raised by the StandardToken contract.
type StandardTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_StandardToken *StandardTokenFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*StandardTokenApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _StandardToken.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &StandardTokenApprovalIterator{contract: _StandardToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_StandardToken *StandardTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StandardTokenApproval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _StandardToken.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StandardTokenApproval)
				if err := _StandardToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// StandardTokenOwnerUpdateIterator is returned from FilterOwnerUpdate and is used to iterate over the raw logs and unpacked data for OwnerUpdate events raised by the StandardToken contract.
type StandardTokenOwnerUpdateIterator struct {
	Event *StandardTokenOwnerUpdate // Event containing the contract specifics and raw log

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
func (it *StandardTokenOwnerUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StandardTokenOwnerUpdate)
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
		it.Event = new(StandardTokenOwnerUpdate)
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
func (it *StandardTokenOwnerUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StandardTokenOwnerUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StandardTokenOwnerUpdate represents a OwnerUpdate event raised by the StandardToken contract.
type StandardTokenOwnerUpdate struct {
	PrevOwner common.Address
	NewOwner  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnerUpdate is a free log retrieval operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_StandardToken *StandardTokenFilterer) FilterOwnerUpdate(opts *bind.FilterOpts) (*StandardTokenOwnerUpdateIterator, error) {

	logs, sub, err := _StandardToken.contract.FilterLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return &StandardTokenOwnerUpdateIterator{contract: _StandardToken.contract, event: "OwnerUpdate", logs: logs, sub: sub}, nil
}

// WatchOwnerUpdate is a free log subscription operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_StandardToken *StandardTokenFilterer) WatchOwnerUpdate(opts *bind.WatchOpts, sink chan<- *StandardTokenOwnerUpdate) (event.Subscription, error) {

	logs, sub, err := _StandardToken.contract.WatchLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StandardTokenOwnerUpdate)
				if err := _StandardToken.contract.UnpackLog(event, "OwnerUpdate", log); err != nil {
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

// StandardTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the StandardToken contract.
type StandardTokenTransferIterator struct {
	Event *StandardTokenTransfer // Event containing the contract specifics and raw log

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
func (it *StandardTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StandardTokenTransfer)
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
		it.Event = new(StandardTokenTransfer)
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
func (it *StandardTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StandardTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StandardTokenTransfer represents a Transfer event raised by the StandardToken contract.
type StandardTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_StandardToken *StandardTokenFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*StandardTokenTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _StandardToken.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &StandardTokenTransferIterator{contract: _StandardToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_StandardToken *StandardTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StandardTokenTransfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _StandardToken.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StandardTokenTransfer)
				if err := _StandardToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// TokenABI is the input ABI used to generate the binding from.
const TokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// TokenBin is the compiled bytecode used for deploying new contracts.
const TokenBin = `0x`

// DeployToken deploys a new Ethereum contract, binding an instance of Token to it.
func DeployToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Token, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// Token is an auto generated Go binding around an Ethereum contract.
type Token struct {
	TokenCaller     // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
	TokenFilterer   // Log filterer for contract events
}

// TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenSession struct {
	Contract     *Token            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCallerSession struct {
	Contract *TokenCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenTransactorSession struct {
	Contract     *TokenTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRaw struct {
	Contract *Token // Generic contract binding to access the raw methods on
}

// TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCallerRaw struct {
	Contract *TokenCaller // Generic read-only contract binding to access the raw methods on
}

// TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenTransactorRaw struct {
	Contract *TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewToken creates a new instance of Token, bound to a specific deployed contract.
func NewToken(address common.Address, backend bind.ContractBackend) (*Token, error) {
	contract, err := bindToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// NewTokenCaller creates a new read-only instance of Token, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contract: contract}, nil
}

// NewTokenTransactor creates a new write-only instance of Token, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}

// NewTokenFilterer creates a new log filterer instance of Token, bound to a specific deployed contract.
func NewTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenFilterer, error) {
	contract, err := bindToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenFilterer{contract: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token.Contract.TokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Token *TokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Token *TokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Token.Contract.Allowance(&_Token.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Token *TokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Token.Contract.Allowance(&_Token.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Token *TokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Token *TokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Token.Contract.BalanceOf(&_Token.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Token *TokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Token.Contract.BalanceOf(&_Token.CallOpts, _owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token *TokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token *TokenSession) TotalSupply() (*big.Int, error) {
	return _Token.Contract.TotalSupply(&_Token.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token *TokenCallerSession) TotalSupply() (*big.Int, error) {
	return _Token.Contract.TotalSupply(&_Token.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_Token *TokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_Token *TokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Approve(&_Token.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_Token *TokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Approve(&_Token.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_Token *TokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_Token *TokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Transfer(&_Token.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_Token *TokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Transfer(&_Token.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_Token *TokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_Token *TokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferFrom(&_Token.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_Token *TokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferFrom(&_Token.TransactOpts, _from, _to, _value)
}

// TokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Token contract.
type TokenApprovalIterator struct {
	Event *TokenApproval // Event containing the contract specifics and raw log

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
func (it *TokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenApproval)
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
		it.Event = new(TokenApproval)
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
func (it *TokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenApproval represents a Approval event raised by the Token contract.
type TokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_Token *TokenFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*TokenApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &TokenApprovalIterator{contract: _Token.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_Token *TokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TokenApproval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenApproval)
				if err := _Token.contract.UnpackLog(event, "Approval", log); err != nil {
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

// TokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Token contract.
type TokenTransferIterator struct {
	Event *TokenTransfer // Event containing the contract specifics and raw log

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
func (it *TokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenTransfer)
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
		it.Event = new(TokenTransfer)
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
func (it *TokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenTransfer represents a Transfer event raised by the Token contract.
type TokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_Token *TokenFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*TokenTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &TokenTransferIterator{contract: _Token.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_Token *TokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TokenTransfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenTransfer)
				if err := _Token.contract.UnpackLog(event, "Transfer", log); err != nil {
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
