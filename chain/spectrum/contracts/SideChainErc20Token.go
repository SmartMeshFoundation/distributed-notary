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

// AtmosphereTokenABI is the input ABI used to generate the binding from.
const AtmosphereTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"lockout\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lockin_htlc\",\"outputs\":[{\"name\":\"SecretHash\",\"type\":\"bytes32\"},{\"name\":\"Expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"queryLockin\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"cancelLockin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"lockin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"queryLockout\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"secret_hash\",\"type\":\"bytes32\"},{\"name\":\"expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"prepareLockout\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"cancelLockOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"secret_hash\",\"type\":\"bytes32\"},{\"name\":\"expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"prepareLockin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lockout_htlc\",\"outputs\":[{\"name\":\"SecretHash\",\"type\":\"bytes32\"},{\"name\":\"Expiration\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"tokenName\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"PrepareLockin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"LockinSecret\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"PrepareLockout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"Lockout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"CancelLockin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"CancelLockout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// AtmosphereTokenBin is the compiled bytecode used for deploying new contracts.
const AtmosphereTokenBin = `0x60028054600160a060020a03191690556006805460ff1916601217905560c0604052600460808190527f76302e310000000000000000000000000000000000000000000000000000000060a09081526200005d9160079190620000b6565b503480156200006b57600080fd5b50604051620010c6380380620010c683398101604052805160018054600160a060020a03191633178155600055018051620000ae906005906020840190620000b6565b50506200015b565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000f957805160ff191683800117855562000129565b8280016001018555821562000129579182015b82811115620001295782518255916020019190600101906200010c565b50620001379291506200013b565b5090565b6200015891905b8082111562000137576000815560010162000142565b90565b610f5b806200016b6000396000f3006080604052600436106101275763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663043d9180811461013957806306fdde031461015f578063095ea7b3146101e957806318160ddd146102215780631e0ef9a41461024857806323b872dd14610287578063313ce567146102b157806354fd4d50146102dc57806357e1ee59146102f157806370a082311461031257806376188aa51461033357806379ba5097146103545780637fd408d2146103695780638caa80f71461038d5780638da5cb5b146103ae57806392d062cd146103df57806393bd8121146103fd5780639a7165491461041e578063a6f9dae114610448578063a9059cbb14610469578063b85287611461048d578063dd62ed3e146104ae575b34801561013357600080fd5b50600080fd5b34801561014557600080fd5b5061015d600160a060020a03600435166024356104d5565b005b34801561016b57600080fd5b50610174610609565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101ae578181015183820152602001610196565b50505050905090810190601f1680156101db5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101f557600080fd5b5061020d600160a060020a0360043516602435610697565b604080519115158252519081900360200190f35b34801561022d57600080fd5b506102366106fe565b60408051918252519081900360200190f35b34801561025457600080fd5b50610269600160a060020a0360043516610704565b60408051938452602084019290925282820152519081900360600190f35b34801561029357600080fd5b5061020d600160a060020a0360043581169060243516604435610725565b3480156102bd57600080fd5b506102c661082b565b6040805160ff9092168252519081900360200190f35b3480156102e857600080fd5b50610174610834565b3480156102fd57600080fd5b50610269600160a060020a036004351661088f565b34801561031e57600080fd5b50610236600160a060020a03600435166108b7565b34801561033f57600080fd5b5061015d600160a060020a03600435166108d2565b34801561036057600080fd5b5061015d61096a565b34801561037557600080fd5b5061015d600160a060020a0360043516602435610a01565b34801561039957600080fd5b50610269600160a060020a0360043516610b7c565b3480156103ba57600080fd5b506103c3610ba4565b60408051600160a060020a039092168252519081900360200190f35b3480156103eb57600080fd5b5061015d600435602435604435610bb3565b34801561040957600080fd5b5061015d600160a060020a0360043516610c60565b34801561042a57600080fd5b5061015d600160a060020a0360043516602435604435606435610d0e565b34801561045457600080fd5b5061015d600160a060020a0360043516610dd0565b34801561047557600080fd5b5061020d600160a060020a0360043516602435610e31565b34801561049957600080fd5b50610269600160a060020a0360043516610ee3565b3480156104ba57600080fd5b50610236600160a060020a0360043581169060243516610f04565b600160a060020a03821660009081526009602052604081206002810154909180821161050057600080fd5b6001830154431061051057600080fd5b60408051602080820187905282518083038201815291830192839052815191929182918401908083835b602083106105595780518252601f19909201916020918201910161053a565b5181516020939093036101000a60001901801990911692169190911790526040519201829003909120865414925061059391505057600080fd5b508154600060028401819055808455600180850182905581548490039182905511156105be57600080fd5b60408051600160a060020a03871681526020810183905281517f4048ec8cca6761c1cdb0f52fcbe25f7486b68d5db97a6711753891dc53e57b45929181900390910190a15050505050565b6005805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561068f5780601f106106645761010080835404028352916020019161068f565b820191906000526020600020905b81548152906001019060200180831161067257829003601f168201915b505050505081565b336000818152600460209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a35060015b92915050565b60005481565b60086020526000908152604090208054600182015460029092015490919083565b600160a060020a03831660009081526003602052604081205482118015906107705750600160a060020a03841660009081526004602090815260408083203384529091529020548211155b80156107955750600160a060020a038316600090815260036020526040902054828101115b1561082057600160a060020a03808416600081815260036020908152604080832080548801905593881680835284832080548890039055600482528483203384528252918490208054879003905583518681529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a3506001610824565b5060005b9392505050565b60065460ff1681565b6007805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561068f5780601f106106645761010080835404028352916020019161068f565b600160a060020a03166000908152600860205260409020805460018201546002909201549092565b600160a060020a031660009081526003602052604090205490565b600160a060020a0381166000908152600860205260408120600281015490919081106108fd57600080fd5b6001820154431161090d57600080fd5b508054600060028301819055808355600183015560408051600160a060020a03851681526020810183905281517f026b98a8ac743c75f99f54b50949aa5e66574f9b73738858c62935046e4aa6c9929181900390910190a1505050565b600254600160a060020a0316331461098157600080fd5b60015460025460408051600160a060020a03938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a1600280546001805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600160a060020a0382166000908152600860205260408120600281015490919081908110610a2e57600080fd5b60408051602080820187905282518083038201815291830192839052815191929182918401908083835b60208310610a775780518252601f199092019160209182019101610a58565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091208654149250610ab191505057600080fd5b60018301544310610ac157600080fd5b50506002810154600160a060020a038416600090815260036020526040902054808201821115610af057600080fd5b600160a060020a0385166000908152600360205260408120828401905580548301908190558210610b2057600080fd5b600060028401819055808455600184015560408051600160a060020a03871681526020810186905281517ffc3a947b611186b1dbc5b435603c40ce6979cc5821ba68ce3973262fb49eb2e5929181900390910190a15050505050565b600160a060020a03166000908152600960205260409020805460018201546002909201549092565b600154600160a060020a031681565b336000908152600960205260408120908211610bce57600080fd5b600281015415610bdd57600080fd5b33600090815260036020526040902054821115610bf957600080fd5b600281018290558381556001810183905533600081815260036020908152604091829020805486900390558151928352820184905280517fbbae3304c67c8fbb052efa093374fc235534c3d862512a40007e7e35062a04759281900390910190a150505050565b600160a060020a038116600090815260096020526040812060028101549091808211610c8b57600080fd5b60018301544311610c9b57600080fd5b50815460006002840181905580845560018401819055600160a060020a038516808252600360209081526040928390208054860190558251918252810183905281517f625a628f697109c4cbee05890a8e5accf4b75c40503b6d8480b11715d148c2db929181900390910190a150505050565b600154600090600160a060020a03163314610d2857600080fd5b600160a060020a0385161515610d3d57600080fd5b600160a060020a03851660009081526008602052604090206002015415610d6357600080fd5b50600160a060020a03841660008181526008602090815260409182902086815560018101869055600281018590558251938452908301849052815190927f1cc3ff93fb861f5fb2869fc15945f233d14ea7a4afa5721ad3c9804be90f3c6a92908290030190a15050505050565b600154600160a060020a03163314610de757600080fd5b600154600160a060020a0382811691161415610e0257600080fd5b6002805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b336000908152600360205260408120548211801590610e695750600160a060020a038316600090815260036020526040902054828101115b15610edb5733600081815260036020908152604080832080548790039055600160a060020a03871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a35060016106f8565b5060006106f8565b60096020526000908152604090208054600182015460029092015490919083565b600160a060020a039182166000908152600460209081526040808320939094168252919091522054905600a165627a7a7230582065ef12d8f96636fd205da443fe374bb79cba27606ed15c655666954aa04fd3b90029`

// DeployAtmosphereToken deploys a new Ethereum contract, binding an instance of AtmosphereToken to it.
func DeployAtmosphereToken(auth *bind.TransactOpts, backend bind.ContractBackend, tokenName string) (common.Address, *types.Transaction, *AtmosphereToken, error) {
	parsed, err := abi.JSON(strings.NewReader(AtmosphereTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtmosphereTokenBin), backend, tokenName)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtmosphereToken{AtmosphereTokenCaller: AtmosphereTokenCaller{contract: contract}, AtmosphereTokenTransactor: AtmosphereTokenTransactor{contract: contract}, AtmosphereTokenFilterer: AtmosphereTokenFilterer{contract: contract}}, nil
}

// AtmosphereToken is an auto generated Go binding around an Ethereum contract.
type AtmosphereToken struct {
	AtmosphereTokenCaller     // Read-only binding to the contract
	AtmosphereTokenTransactor // Write-only binding to the contract
	AtmosphereTokenFilterer   // Log filterer for contract events
}

// AtmosphereTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtmosphereTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtmosphereTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtmosphereTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtmosphereTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtmosphereTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtmosphereTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtmosphereTokenSession struct {
	Contract     *AtmosphereToken  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtmosphereTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtmosphereTokenCallerSession struct {
	Contract *AtmosphereTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AtmosphereTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtmosphereTokenTransactorSession struct {
	Contract     *AtmosphereTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AtmosphereTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtmosphereTokenRaw struct {
	Contract *AtmosphereToken // Generic contract binding to access the raw methods on
}

// AtmosphereTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtmosphereTokenCallerRaw struct {
	Contract *AtmosphereTokenCaller // Generic read-only contract binding to access the raw methods on
}

// AtmosphereTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtmosphereTokenTransactorRaw struct {
	Contract *AtmosphereTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtmosphereToken creates a new instance of AtmosphereToken, bound to a specific deployed contract.
func NewAtmosphereToken(address common.Address, backend bind.ContractBackend) (*AtmosphereToken, error) {
	contract, err := bindAtmosphereToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtmosphereToken{AtmosphereTokenCaller: AtmosphereTokenCaller{contract: contract}, AtmosphereTokenTransactor: AtmosphereTokenTransactor{contract: contract}, AtmosphereTokenFilterer: AtmosphereTokenFilterer{contract: contract}}, nil
}

// NewAtmosphereTokenCaller creates a new read-only instance of AtmosphereToken, bound to a specific deployed contract.
func NewAtmosphereTokenCaller(address common.Address, caller bind.ContractCaller) (*AtmosphereTokenCaller, error) {
	contract, err := bindAtmosphereToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenCaller{contract: contract}, nil
}

// NewAtmosphereTokenTransactor creates a new write-only instance of AtmosphereToken, bound to a specific deployed contract.
func NewAtmosphereTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*AtmosphereTokenTransactor, error) {
	contract, err := bindAtmosphereToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenTransactor{contract: contract}, nil
}

// NewAtmosphereTokenFilterer creates a new log filterer instance of AtmosphereToken, bound to a specific deployed contract.
func NewAtmosphereTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*AtmosphereTokenFilterer, error) {
	contract, err := bindAtmosphereToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenFilterer{contract: contract}, nil
}

// bindAtmosphereToken binds a generic wrapper to an already deployed contract.
func bindAtmosphereToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtmosphereTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtmosphereToken *AtmosphereTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtmosphereToken.Contract.AtmosphereTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtmosphereToken *AtmosphereTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.AtmosphereTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtmosphereToken *AtmosphereTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.AtmosphereTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtmosphereToken *AtmosphereTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtmosphereToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtmosphereToken *AtmosphereTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtmosphereToken *AtmosphereTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_AtmosphereToken *AtmosphereTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AtmosphereToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_AtmosphereToken *AtmosphereTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _AtmosphereToken.Contract.Allowance(&_AtmosphereToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_AtmosphereToken *AtmosphereTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _AtmosphereToken.Contract.Allowance(&_AtmosphereToken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_AtmosphereToken *AtmosphereTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AtmosphereToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_AtmosphereToken *AtmosphereTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _AtmosphereToken.Contract.BalanceOf(&_AtmosphereToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_AtmosphereToken *AtmosphereTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _AtmosphereToken.Contract.BalanceOf(&_AtmosphereToken.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_AtmosphereToken *AtmosphereTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _AtmosphereToken.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_AtmosphereToken *AtmosphereTokenSession) Decimals() (uint8, error) {
	return _AtmosphereToken.Contract.Decimals(&_AtmosphereToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_AtmosphereToken *AtmosphereTokenCallerSession) Decimals() (uint8, error) {
	return _AtmosphereToken.Contract.Decimals(&_AtmosphereToken.CallOpts)
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_AtmosphereToken *AtmosphereTokenCaller) LockinHtlc(opts *bind.CallOpts, arg0 common.Address) (struct {
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
	err := _AtmosphereToken.contract.Call(opts, out, "lockin_htlc", arg0)
	return *ret, err
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_AtmosphereToken *AtmosphereTokenSession) LockinHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _AtmosphereToken.Contract.LockinHtlc(&_AtmosphereToken.CallOpts, arg0)
}

// LockinHtlc is a free data retrieval call binding the contract method 0x1e0ef9a4.
//
// Solidity: function lockin_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_AtmosphereToken *AtmosphereTokenCallerSession) LockinHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _AtmosphereToken.Contract.LockinHtlc(&_AtmosphereToken.CallOpts, arg0)
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_AtmosphereToken *AtmosphereTokenCaller) LockoutHtlc(opts *bind.CallOpts, arg0 common.Address) (struct {
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
	err := _AtmosphereToken.contract.Call(opts, out, "lockout_htlc", arg0)
	return *ret, err
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_AtmosphereToken *AtmosphereTokenSession) LockoutHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _AtmosphereToken.Contract.LockoutHtlc(&_AtmosphereToken.CallOpts, arg0)
}

// LockoutHtlc is a free data retrieval call binding the contract method 0xb8528761.
//
// Solidity: function lockout_htlc( address) constant returns(SecretHash bytes32, Expiration uint256, value uint256)
func (_AtmosphereToken *AtmosphereTokenCallerSession) LockoutHtlc(arg0 common.Address) (struct {
	SecretHash [32]byte
	Expiration *big.Int
	Value      *big.Int
}, error) {
	return _AtmosphereToken.Contract.LockoutHtlc(&_AtmosphereToken.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_AtmosphereToken *AtmosphereTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AtmosphereToken.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_AtmosphereToken *AtmosphereTokenSession) Name() (string, error) {
	return _AtmosphereToken.Contract.Name(&_AtmosphereToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_AtmosphereToken *AtmosphereTokenCallerSession) Name() (string, error) {
	return _AtmosphereToken.Contract.Name(&_AtmosphereToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AtmosphereToken *AtmosphereTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AtmosphereToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AtmosphereToken *AtmosphereTokenSession) Owner() (common.Address, error) {
	return _AtmosphereToken.Contract.Owner(&_AtmosphereToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AtmosphereToken *AtmosphereTokenCallerSession) Owner() (common.Address, error) {
	return _AtmosphereToken.Contract.Owner(&_AtmosphereToken.CallOpts)
}

// QueryLockin is a free data retrieval call binding the contract method 0x57e1ee59.
//
// Solidity: function queryLockin(account address) constant returns(bytes32, uint256, uint256)
func (_AtmosphereToken *AtmosphereTokenCaller) QueryLockin(opts *bind.CallOpts, account common.Address) ([32]byte, *big.Int, *big.Int, error) {
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
	err := _AtmosphereToken.contract.Call(opts, out, "queryLockin", account)
	return *ret0, *ret1, *ret2, err
}

// QueryLockin is a free data retrieval call binding the contract method 0x57e1ee59.
//
// Solidity: function queryLockin(account address) constant returns(bytes32, uint256, uint256)
func (_AtmosphereToken *AtmosphereTokenSession) QueryLockin(account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	return _AtmosphereToken.Contract.QueryLockin(&_AtmosphereToken.CallOpts, account)
}

// QueryLockin is a free data retrieval call binding the contract method 0x57e1ee59.
//
// Solidity: function queryLockin(account address) constant returns(bytes32, uint256, uint256)
func (_AtmosphereToken *AtmosphereTokenCallerSession) QueryLockin(account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	return _AtmosphereToken.Contract.QueryLockin(&_AtmosphereToken.CallOpts, account)
}

// QueryLockout is a free data retrieval call binding the contract method 0x8caa80f7.
//
// Solidity: function queryLockout(account address) constant returns(bytes32, uint256, uint256)
func (_AtmosphereToken *AtmosphereTokenCaller) QueryLockout(opts *bind.CallOpts, account common.Address) ([32]byte, *big.Int, *big.Int, error) {
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
	err := _AtmosphereToken.contract.Call(opts, out, "queryLockout", account)
	return *ret0, *ret1, *ret2, err
}

// QueryLockout is a free data retrieval call binding the contract method 0x8caa80f7.
//
// Solidity: function queryLockout(account address) constant returns(bytes32, uint256, uint256)
func (_AtmosphereToken *AtmosphereTokenSession) QueryLockout(account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	return _AtmosphereToken.Contract.QueryLockout(&_AtmosphereToken.CallOpts, account)
}

// QueryLockout is a free data retrieval call binding the contract method 0x8caa80f7.
//
// Solidity: function queryLockout(account address) constant returns(bytes32, uint256, uint256)
func (_AtmosphereToken *AtmosphereTokenCallerSession) QueryLockout(account common.Address) ([32]byte, *big.Int, *big.Int, error) {
	return _AtmosphereToken.Contract.QueryLockout(&_AtmosphereToken.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_AtmosphereToken *AtmosphereTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AtmosphereToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_AtmosphereToken *AtmosphereTokenSession) TotalSupply() (*big.Int, error) {
	return _AtmosphereToken.Contract.TotalSupply(&_AtmosphereToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_AtmosphereToken *AtmosphereTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _AtmosphereToken.Contract.TotalSupply(&_AtmosphereToken.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_AtmosphereToken *AtmosphereTokenCaller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AtmosphereToken.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_AtmosphereToken *AtmosphereTokenSession) Version() (string, error) {
	return _AtmosphereToken.Contract.Version(&_AtmosphereToken.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_AtmosphereToken *AtmosphereTokenCallerSession) Version() (string, error) {
	return _AtmosphereToken.Contract.Version(&_AtmosphereToken.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AtmosphereToken *AtmosphereTokenTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AtmosphereToken *AtmosphereTokenSession) AcceptOwnership() (*types.Transaction, error) {
	return _AtmosphereToken.Contract.AcceptOwnership(&_AtmosphereToken.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AtmosphereToken *AtmosphereTokenTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AtmosphereToken.Contract.AcceptOwnership(&_AtmosphereToken.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.Approve(&_AtmosphereToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.Approve(&_AtmosphereToken.TransactOpts, _spender, _value)
}

// CancelLockOut is a paid mutator transaction binding the contract method 0x93bd8121.
//
// Solidity: function cancelLockOut(account address) returns()
func (_AtmosphereToken *AtmosphereTokenTransactor) CancelLockOut(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "cancelLockOut", account)
}

// CancelLockOut is a paid mutator transaction binding the contract method 0x93bd8121.
//
// Solidity: function cancelLockOut(account address) returns()
func (_AtmosphereToken *AtmosphereTokenSession) CancelLockOut(account common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.CancelLockOut(&_AtmosphereToken.TransactOpts, account)
}

// CancelLockOut is a paid mutator transaction binding the contract method 0x93bd8121.
//
// Solidity: function cancelLockOut(account address) returns()
func (_AtmosphereToken *AtmosphereTokenTransactorSession) CancelLockOut(account common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.CancelLockOut(&_AtmosphereToken.TransactOpts, account)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(account address) returns()
func (_AtmosphereToken *AtmosphereTokenTransactor) CancelLockin(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "cancelLockin", account)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(account address) returns()
func (_AtmosphereToken *AtmosphereTokenSession) CancelLockin(account common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.CancelLockin(&_AtmosphereToken.TransactOpts, account)
}

// CancelLockin is a paid mutator transaction binding the contract method 0x76188aa5.
//
// Solidity: function cancelLockin(account address) returns()
func (_AtmosphereToken *AtmosphereTokenTransactorSession) CancelLockin(account common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.CancelLockin(&_AtmosphereToken.TransactOpts, account)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_AtmosphereToken *AtmosphereTokenTransactor) ChangeOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "changeOwner", _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_AtmosphereToken *AtmosphereTokenSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.ChangeOwner(&_AtmosphereToken.TransactOpts, _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_AtmosphereToken *AtmosphereTokenTransactorSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.ChangeOwner(&_AtmosphereToken.TransactOpts, _newOwner)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(account address, secret bytes32) returns()
func (_AtmosphereToken *AtmosphereTokenTransactor) Lockin(opts *bind.TransactOpts, account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "lockin", account, secret)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(account address, secret bytes32) returns()
func (_AtmosphereToken *AtmosphereTokenSession) Lockin(account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.Lockin(&_AtmosphereToken.TransactOpts, account, secret)
}

// Lockin is a paid mutator transaction binding the contract method 0x7fd408d2.
//
// Solidity: function lockin(account address, secret bytes32) returns()
func (_AtmosphereToken *AtmosphereTokenTransactorSession) Lockin(account common.Address, secret [32]byte) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.Lockin(&_AtmosphereToken.TransactOpts, account, secret)
}

// Lockout is a paid mutator transaction binding the contract method 0x043d9180.
//
// Solidity: function lockout(from address, secret bytes32) returns()
func (_AtmosphereToken *AtmosphereTokenTransactor) Lockout(opts *bind.TransactOpts, from common.Address, secret [32]byte) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "lockout", from, secret)
}

// Lockout is a paid mutator transaction binding the contract method 0x043d9180.
//
// Solidity: function lockout(from address, secret bytes32) returns()
func (_AtmosphereToken *AtmosphereTokenSession) Lockout(from common.Address, secret [32]byte) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.Lockout(&_AtmosphereToken.TransactOpts, from, secret)
}

// Lockout is a paid mutator transaction binding the contract method 0x043d9180.
//
// Solidity: function lockout(from address, secret bytes32) returns()
func (_AtmosphereToken *AtmosphereTokenTransactorSession) Lockout(from common.Address, secret [32]byte) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.Lockout(&_AtmosphereToken.TransactOpts, from, secret)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0x9a716549.
//
// Solidity: function prepareLockin(account address, secret_hash bytes32, expiration uint256, value uint256) returns()
func (_AtmosphereToken *AtmosphereTokenTransactor) PrepareLockin(opts *bind.TransactOpts, account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "prepareLockin", account, secret_hash, expiration, value)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0x9a716549.
//
// Solidity: function prepareLockin(account address, secret_hash bytes32, expiration uint256, value uint256) returns()
func (_AtmosphereToken *AtmosphereTokenSession) PrepareLockin(account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.PrepareLockin(&_AtmosphereToken.TransactOpts, account, secret_hash, expiration, value)
}

// PrepareLockin is a paid mutator transaction binding the contract method 0x9a716549.
//
// Solidity: function prepareLockin(account address, secret_hash bytes32, expiration uint256, value uint256) returns()
func (_AtmosphereToken *AtmosphereTokenTransactorSession) PrepareLockin(account common.Address, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.PrepareLockin(&_AtmosphereToken.TransactOpts, account, secret_hash, expiration, value)
}

// PrepareLockout is a paid mutator transaction binding the contract method 0x92d062cd.
//
// Solidity: function prepareLockout(secret_hash bytes32, expiration uint256, value uint256) returns()
func (_AtmosphereToken *AtmosphereTokenTransactor) PrepareLockout(opts *bind.TransactOpts, secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "prepareLockout", secret_hash, expiration, value)
}

// PrepareLockout is a paid mutator transaction binding the contract method 0x92d062cd.
//
// Solidity: function prepareLockout(secret_hash bytes32, expiration uint256, value uint256) returns()
func (_AtmosphereToken *AtmosphereTokenSession) PrepareLockout(secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.PrepareLockout(&_AtmosphereToken.TransactOpts, secret_hash, expiration, value)
}

// PrepareLockout is a paid mutator transaction binding the contract method 0x92d062cd.
//
// Solidity: function prepareLockout(secret_hash bytes32, expiration uint256, value uint256) returns()
func (_AtmosphereToken *AtmosphereTokenTransactorSession) PrepareLockout(secret_hash [32]byte, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.PrepareLockout(&_AtmosphereToken.TransactOpts, secret_hash, expiration, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.Transfer(&_AtmosphereToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.Transfer(&_AtmosphereToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.TransferFrom(&_AtmosphereToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_AtmosphereToken *AtmosphereTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _AtmosphereToken.Contract.TransferFrom(&_AtmosphereToken.TransactOpts, _from, _to, _value)
}

// AtmosphereTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AtmosphereToken contract.
type AtmosphereTokenApprovalIterator struct {
	Event *AtmosphereTokenApproval // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenApproval)
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
		it.Event = new(AtmosphereTokenApproval)
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
func (it *AtmosphereTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenApproval represents a Approval event raised by the AtmosphereToken contract.
type AtmosphereTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*AtmosphereTokenApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenApprovalIterator{contract: _AtmosphereToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenApproval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenApproval)
				if err := _AtmosphereToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// AtmosphereTokenCancelLockinIterator is returned from FilterCancelLockin and is used to iterate over the raw logs and unpacked data for CancelLockin events raised by the AtmosphereToken contract.
type AtmosphereTokenCancelLockinIterator struct {
	Event *AtmosphereTokenCancelLockin // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenCancelLockinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenCancelLockin)
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
		it.Event = new(AtmosphereTokenCancelLockin)
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
func (it *AtmosphereTokenCancelLockinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenCancelLockinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenCancelLockin represents a CancelLockin event raised by the AtmosphereToken contract.
type AtmosphereTokenCancelLockin struct {
	Account    common.Address
	SecretHash [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancelLockin is a free log retrieval operation binding the contract event 0x026b98a8ac743c75f99f54b50949aa5e66574f9b73738858c62935046e4aa6c9.
//
// Solidity: e CancelLockin(account address, secretHash bytes32)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterCancelLockin(opts *bind.FilterOpts) (*AtmosphereTokenCancelLockinIterator, error) {

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "CancelLockin")
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenCancelLockinIterator{contract: _AtmosphereToken.contract, event: "CancelLockin", logs: logs, sub: sub}, nil
}

// WatchCancelLockin is a free log subscription operation binding the contract event 0x026b98a8ac743c75f99f54b50949aa5e66574f9b73738858c62935046e4aa6c9.
//
// Solidity: e CancelLockin(account address, secretHash bytes32)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchCancelLockin(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenCancelLockin) (event.Subscription, error) {

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "CancelLockin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenCancelLockin)
				if err := _AtmosphereToken.contract.UnpackLog(event, "CancelLockin", log); err != nil {
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

// AtmosphereTokenCancelLockoutIterator is returned from FilterCancelLockout and is used to iterate over the raw logs and unpacked data for CancelLockout events raised by the AtmosphereToken contract.
type AtmosphereTokenCancelLockoutIterator struct {
	Event *AtmosphereTokenCancelLockout // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenCancelLockoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenCancelLockout)
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
		it.Event = new(AtmosphereTokenCancelLockout)
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
func (it *AtmosphereTokenCancelLockoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenCancelLockoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenCancelLockout represents a CancelLockout event raised by the AtmosphereToken contract.
type AtmosphereTokenCancelLockout struct {
	Account    common.Address
	SecretHash [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancelLockout is a free log retrieval operation binding the contract event 0x625a628f697109c4cbee05890a8e5accf4b75c40503b6d8480b11715d148c2db.
//
// Solidity: e CancelLockout(account address, secretHash bytes32)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterCancelLockout(opts *bind.FilterOpts) (*AtmosphereTokenCancelLockoutIterator, error) {

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "CancelLockout")
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenCancelLockoutIterator{contract: _AtmosphereToken.contract, event: "CancelLockout", logs: logs, sub: sub}, nil
}

// WatchCancelLockout is a free log subscription operation binding the contract event 0x625a628f697109c4cbee05890a8e5accf4b75c40503b6d8480b11715d148c2db.
//
// Solidity: e CancelLockout(account address, secretHash bytes32)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchCancelLockout(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenCancelLockout) (event.Subscription, error) {

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "CancelLockout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenCancelLockout)
				if err := _AtmosphereToken.contract.UnpackLog(event, "CancelLockout", log); err != nil {
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

// AtmosphereTokenLockinSecretIterator is returned from FilterLockinSecret and is used to iterate over the raw logs and unpacked data for LockinSecret events raised by the AtmosphereToken contract.
type AtmosphereTokenLockinSecretIterator struct {
	Event *AtmosphereTokenLockinSecret // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenLockinSecretIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenLockinSecret)
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
		it.Event = new(AtmosphereTokenLockinSecret)
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
func (it *AtmosphereTokenLockinSecretIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenLockinSecretIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenLockinSecret represents a LockinSecret event raised by the AtmosphereToken contract.
type AtmosphereTokenLockinSecret struct {
	Account common.Address
	Secret  [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLockinSecret is a free log retrieval operation binding the contract event 0xfc3a947b611186b1dbc5b435603c40ce6979cc5821ba68ce3973262fb49eb2e5.
//
// Solidity: e LockinSecret(account address, secret bytes32)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterLockinSecret(opts *bind.FilterOpts) (*AtmosphereTokenLockinSecretIterator, error) {

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "LockinSecret")
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenLockinSecretIterator{contract: _AtmosphereToken.contract, event: "LockinSecret", logs: logs, sub: sub}, nil
}

// WatchLockinSecret is a free log subscription operation binding the contract event 0xfc3a947b611186b1dbc5b435603c40ce6979cc5821ba68ce3973262fb49eb2e5.
//
// Solidity: e LockinSecret(account address, secret bytes32)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchLockinSecret(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenLockinSecret) (event.Subscription, error) {

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "LockinSecret")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenLockinSecret)
				if err := _AtmosphereToken.contract.UnpackLog(event, "LockinSecret", log); err != nil {
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

// AtmosphereTokenLockoutIterator is returned from FilterLockout and is used to iterate over the raw logs and unpacked data for Lockout events raised by the AtmosphereToken contract.
type AtmosphereTokenLockoutIterator struct {
	Event *AtmosphereTokenLockout // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenLockoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenLockout)
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
		it.Event = new(AtmosphereTokenLockout)
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
func (it *AtmosphereTokenLockoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenLockoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenLockout represents a Lockout event raised by the AtmosphereToken contract.
type AtmosphereTokenLockout struct {
	Account    common.Address
	SecretHash [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLockout is a free log retrieval operation binding the contract event 0x4048ec8cca6761c1cdb0f52fcbe25f7486b68d5db97a6711753891dc53e57b45.
//
// Solidity: e Lockout(account address, secretHash bytes32)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterLockout(opts *bind.FilterOpts) (*AtmosphereTokenLockoutIterator, error) {

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "Lockout")
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenLockoutIterator{contract: _AtmosphereToken.contract, event: "Lockout", logs: logs, sub: sub}, nil
}

// WatchLockout is a free log subscription operation binding the contract event 0x4048ec8cca6761c1cdb0f52fcbe25f7486b68d5db97a6711753891dc53e57b45.
//
// Solidity: e Lockout(account address, secretHash bytes32)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchLockout(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenLockout) (event.Subscription, error) {

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "Lockout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenLockout)
				if err := _AtmosphereToken.contract.UnpackLog(event, "Lockout", log); err != nil {
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

// AtmosphereTokenOwnerUpdateIterator is returned from FilterOwnerUpdate and is used to iterate over the raw logs and unpacked data for OwnerUpdate events raised by the AtmosphereToken contract.
type AtmosphereTokenOwnerUpdateIterator struct {
	Event *AtmosphereTokenOwnerUpdate // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenOwnerUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenOwnerUpdate)
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
		it.Event = new(AtmosphereTokenOwnerUpdate)
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
func (it *AtmosphereTokenOwnerUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenOwnerUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenOwnerUpdate represents a OwnerUpdate event raised by the AtmosphereToken contract.
type AtmosphereTokenOwnerUpdate struct {
	PrevOwner common.Address
	NewOwner  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnerUpdate is a free log retrieval operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterOwnerUpdate(opts *bind.FilterOpts) (*AtmosphereTokenOwnerUpdateIterator, error) {

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenOwnerUpdateIterator{contract: _AtmosphereToken.contract, event: "OwnerUpdate", logs: logs, sub: sub}, nil
}

// WatchOwnerUpdate is a free log subscription operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchOwnerUpdate(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenOwnerUpdate) (event.Subscription, error) {

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenOwnerUpdate)
				if err := _AtmosphereToken.contract.UnpackLog(event, "OwnerUpdate", log); err != nil {
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

// AtmosphereTokenPrepareLockinIterator is returned from FilterPrepareLockin and is used to iterate over the raw logs and unpacked data for PrepareLockin events raised by the AtmosphereToken contract.
type AtmosphereTokenPrepareLockinIterator struct {
	Event *AtmosphereTokenPrepareLockin // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenPrepareLockinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenPrepareLockin)
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
		it.Event = new(AtmosphereTokenPrepareLockin)
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
func (it *AtmosphereTokenPrepareLockinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenPrepareLockinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenPrepareLockin represents a PrepareLockin event raised by the AtmosphereToken contract.
type AtmosphereTokenPrepareLockin struct {
	Account common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPrepareLockin is a free log retrieval operation binding the contract event 0x1cc3ff93fb861f5fb2869fc15945f233d14ea7a4afa5721ad3c9804be90f3c6a.
//
// Solidity: e PrepareLockin(account address, value uint256)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterPrepareLockin(opts *bind.FilterOpts) (*AtmosphereTokenPrepareLockinIterator, error) {

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "PrepareLockin")
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenPrepareLockinIterator{contract: _AtmosphereToken.contract, event: "PrepareLockin", logs: logs, sub: sub}, nil
}

// WatchPrepareLockin is a free log subscription operation binding the contract event 0x1cc3ff93fb861f5fb2869fc15945f233d14ea7a4afa5721ad3c9804be90f3c6a.
//
// Solidity: e PrepareLockin(account address, value uint256)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchPrepareLockin(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenPrepareLockin) (event.Subscription, error) {

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "PrepareLockin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenPrepareLockin)
				if err := _AtmosphereToken.contract.UnpackLog(event, "PrepareLockin", log); err != nil {
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

// AtmosphereTokenPrepareLockoutIterator is returned from FilterPrepareLockout and is used to iterate over the raw logs and unpacked data for PrepareLockout events raised by the AtmosphereToken contract.
type AtmosphereTokenPrepareLockoutIterator struct {
	Event *AtmosphereTokenPrepareLockout // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenPrepareLockoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenPrepareLockout)
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
		it.Event = new(AtmosphereTokenPrepareLockout)
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
func (it *AtmosphereTokenPrepareLockoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenPrepareLockoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenPrepareLockout represents a PrepareLockout event raised by the AtmosphereToken contract.
type AtmosphereTokenPrepareLockout struct {
	Account common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPrepareLockout is a free log retrieval operation binding the contract event 0xbbae3304c67c8fbb052efa093374fc235534c3d862512a40007e7e35062a0475.
//
// Solidity: e PrepareLockout(account address, _value uint256)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterPrepareLockout(opts *bind.FilterOpts) (*AtmosphereTokenPrepareLockoutIterator, error) {

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "PrepareLockout")
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenPrepareLockoutIterator{contract: _AtmosphereToken.contract, event: "PrepareLockout", logs: logs, sub: sub}, nil
}

// WatchPrepareLockout is a free log subscription operation binding the contract event 0xbbae3304c67c8fbb052efa093374fc235534c3d862512a40007e7e35062a0475.
//
// Solidity: e PrepareLockout(account address, _value uint256)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchPrepareLockout(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenPrepareLockout) (event.Subscription, error) {

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "PrepareLockout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenPrepareLockout)
				if err := _AtmosphereToken.contract.UnpackLog(event, "PrepareLockout", log); err != nil {
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

// AtmosphereTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AtmosphereToken contract.
type AtmosphereTokenTransferIterator struct {
	Event *AtmosphereTokenTransfer // Event containing the contract specifics and raw log

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
func (it *AtmosphereTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtmosphereTokenTransfer)
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
		it.Event = new(AtmosphereTokenTransfer)
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
func (it *AtmosphereTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtmosphereTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtmosphereTokenTransfer represents a Transfer event raised by the AtmosphereToken contract.
type AtmosphereTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_AtmosphereToken *AtmosphereTokenFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*AtmosphereTokenTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _AtmosphereToken.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &AtmosphereTokenTransferIterator{contract: _AtmosphereToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_AtmosphereToken *AtmosphereTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AtmosphereTokenTransfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _AtmosphereToken.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtmosphereTokenTransfer)
				if err := _AtmosphereToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
const OwnedBin = `0x608060405260018054600160a060020a031916905534801561002057600080fd5b5060008054600160a060020a031916331790556101f7806100426000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166379ba5097811461005b5780638da5cb5b14610072578063a6f9dae1146100a3575b600080fd5b34801561006757600080fd5b506100706100c4565b005b34801561007e57600080fd5b5061008761015b565b60408051600160a060020a039092168252519081900360200190f35b3480156100af57600080fd5b50610070600160a060020a036004351661016a565b600154600160a060020a031633146100db57600080fd5b60005460015460408051600160a060020a03938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a1600180546000805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600054600160a060020a031681565b600054600160a060020a0316331461018157600080fd5b600054600160a060020a038281169116141561019c57600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a7230582040c844b07381ce9608df79423225ec48740b9b58ff9e8bc65a95931d7fdf0c7d0029`

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
const StandardTokenBin = `0x608060405260028054600160a060020a03199081169091556001805490911633179055610599806100316000396000f3006080604052600436106100985763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b3811461009d57806318160ddd146100d557806323b872dd146100fc57806370a082311461012657806379ba5097146101475780638da5cb5b1461015e578063a6f9dae11461018f578063a9059cbb146101b0578063dd62ed3e146101d4575b600080fd5b3480156100a957600080fd5b506100c1600160a060020a03600435166024356101fb565b604080519115158252519081900360200190f35b3480156100e157600080fd5b506100ea610262565b60408051918252519081900360200190f35b34801561010857600080fd5b506100c1600160a060020a0360043581169060243516604435610268565b34801561013257600080fd5b506100ea600160a060020a036004351661036e565b34801561015357600080fd5b5061015c610389565b005b34801561016a57600080fd5b50610173610420565b60408051600160a060020a039092168252519081900360200190f35b34801561019b57600080fd5b5061015c600160a060020a036004351661042f565b3480156101bc57600080fd5b506100c1600160a060020a0360043516602435610490565b3480156101e057600080fd5b506100ea600160a060020a0360043581169060243516610542565b336000818152600460209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a35060015b92915050565b60005481565b600160a060020a03831660009081526003602052604081205482118015906102b35750600160a060020a03841660009081526004602090815260408083203384529091529020548211155b80156102d85750600160a060020a038316600090815260036020526040902054828101115b1561036357600160a060020a03808416600081815260036020908152604080832080548801905593881680835284832080548890039055600482528483203384528252918490208054879003905583518681529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a3506001610367565b5060005b9392505050565b600160a060020a031660009081526003602052604090205490565b600254600160a060020a031633146103a057600080fd5b60015460025460408051600160a060020a03938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a1600280546001805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600154600160a060020a031681565b600154600160a060020a0316331461044657600080fd5b600154600160a060020a038281169116141561046157600080fd5b6002805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b3360009081526003602052604081205482118015906104c85750600160a060020a038316600090815260036020526040902054828101115b1561053a5733600081815260036020908152604080832080548790039055600160a060020a03871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a350600161025c565b50600061025c565b600160a060020a039182166000908152600460209081526040808320939094168252919091522054905600a165627a7a72305820e6399a6641fba180bcdeed17909bdcc1884b7bf5569b218ac22cf4f986c950b40029`

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
