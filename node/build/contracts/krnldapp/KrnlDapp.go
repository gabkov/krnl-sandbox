// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package krnldapp

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// KrnldappMetaData contains all meta data concerning the Krnldapp contract.
var KrnldappMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"GetCounter\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"SayHi\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"protectedFunctionality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newAuthority\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unprotectFunctionShouldReturn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_number\",\"type\":\"uint256\"}],\"name\":\"unprotectedFunctionShouldThrow\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// KrnldappABI is the input ABI used to generate the binding from.
// Deprecated: Use KrnldappMetaData.ABI instead.
var KrnldappABI = KrnldappMetaData.ABI

// Krnldapp is an auto generated Go binding around an Ethereum contract.
type Krnldapp struct {
	KrnldappCaller     // Read-only binding to the contract
	KrnldappTransactor // Write-only binding to the contract
	KrnldappFilterer   // Log filterer for contract events
}

// KrnldappCaller is an auto generated read-only Go binding around an Ethereum contract.
type KrnldappCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KrnldappTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KrnldappTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KrnldappFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KrnldappFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KrnldappSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KrnldappSession struct {
	Contract     *Krnldapp         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KrnldappCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KrnldappCallerSession struct {
	Contract *KrnldappCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// KrnldappTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KrnldappTransactorSession struct {
	Contract     *KrnldappTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// KrnldappRaw is an auto generated low-level Go binding around an Ethereum contract.
type KrnldappRaw struct {
	Contract *Krnldapp // Generic contract binding to access the raw methods on
}

// KrnldappCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KrnldappCallerRaw struct {
	Contract *KrnldappCaller // Generic read-only contract binding to access the raw methods on
}

// KrnldappTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KrnldappTransactorRaw struct {
	Contract *KrnldappTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKrnldapp creates a new instance of Krnldapp, bound to a specific deployed contract.
func NewKrnldapp(address common.Address, backend bind.ContractBackend) (*Krnldapp, error) {
	contract, err := bindKrnldapp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Krnldapp{KrnldappCaller: KrnldappCaller{contract: contract}, KrnldappTransactor: KrnldappTransactor{contract: contract}, KrnldappFilterer: KrnldappFilterer{contract: contract}}, nil
}

// NewKrnldappCaller creates a new read-only instance of Krnldapp, bound to a specific deployed contract.
func NewKrnldappCaller(address common.Address, caller bind.ContractCaller) (*KrnldappCaller, error) {
	contract, err := bindKrnldapp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KrnldappCaller{contract: contract}, nil
}

// NewKrnldappTransactor creates a new write-only instance of Krnldapp, bound to a specific deployed contract.
func NewKrnldappTransactor(address common.Address, transactor bind.ContractTransactor) (*KrnldappTransactor, error) {
	contract, err := bindKrnldapp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KrnldappTransactor{contract: contract}, nil
}

// NewKrnldappFilterer creates a new log filterer instance of Krnldapp, bound to a specific deployed contract.
func NewKrnldappFilterer(address common.Address, filterer bind.ContractFilterer) (*KrnldappFilterer, error) {
	contract, err := bindKrnldapp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KrnldappFilterer{contract: contract}, nil
}

// bindKrnldapp binds a generic wrapper to an already deployed contract.
func bindKrnldapp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KrnldappMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Krnldapp *KrnldappRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Krnldapp.Contract.KrnldappCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Krnldapp *KrnldappRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Krnldapp.Contract.KrnldappTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Krnldapp *KrnldappRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Krnldapp.Contract.KrnldappTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Krnldapp *KrnldappCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Krnldapp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Krnldapp *KrnldappTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Krnldapp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Krnldapp *KrnldappTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Krnldapp.Contract.contract.Transact(opts, method, params...)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Krnldapp *KrnldappCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Krnldapp.contract.Call(opts, &out, "authority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Krnldapp *KrnldappSession) Authority() (common.Address, error) {
	return _Krnldapp.Contract.Authority(&_Krnldapp.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Krnldapp *KrnldappCallerSession) Authority() (common.Address, error) {
	return _Krnldapp.Contract.Authority(&_Krnldapp.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Krnldapp *KrnldappCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Krnldapp.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Krnldapp *KrnldappSession) Counter() (*big.Int, error) {
	return _Krnldapp.Contract.Counter(&_Krnldapp.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Krnldapp *KrnldappCallerSession) Counter() (*big.Int, error) {
	return _Krnldapp.Contract.Counter(&_Krnldapp.CallOpts)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x1626ba7e.
//
// Solidity: function isValidSignature(bytes32 _hash, bytes _signature) view returns(bytes4)
func (_Krnldapp *KrnldappCaller) IsValidSignature(opts *bind.CallOpts, _hash [32]byte, _signature []byte) ([4]byte, error) {
	var out []interface{}
	err := _Krnldapp.contract.Call(opts, &out, "isValidSignature", _hash, _signature)

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// IsValidSignature is a free data retrieval call binding the contract method 0x1626ba7e.
//
// Solidity: function isValidSignature(bytes32 _hash, bytes _signature) view returns(bytes4)
func (_Krnldapp *KrnldappSession) IsValidSignature(_hash [32]byte, _signature []byte) ([4]byte, error) {
	return _Krnldapp.Contract.IsValidSignature(&_Krnldapp.CallOpts, _hash, _signature)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x1626ba7e.
//
// Solidity: function isValidSignature(bytes32 _hash, bytes _signature) view returns(bytes4)
func (_Krnldapp *KrnldappCallerSession) IsValidSignature(_hash [32]byte, _signature []byte) ([4]byte, error) {
	return _Krnldapp.Contract.IsValidSignature(&_Krnldapp.CallOpts, _hash, _signature)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Krnldapp *KrnldappCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Krnldapp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Krnldapp *KrnldappSession) Owner() (common.Address, error) {
	return _Krnldapp.Contract.Owner(&_Krnldapp.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Krnldapp *KrnldappCallerSession) Owner() (common.Address, error) {
	return _Krnldapp.Contract.Owner(&_Krnldapp.CallOpts)
}

// UnprotectedFunctionShouldThrow is a free data retrieval call binding the contract method 0x36601440.
//
// Solidity: function unprotectedFunctionShouldThrow(uint256 _number) pure returns()
func (_Krnldapp *KrnldappCaller) UnprotectedFunctionShouldThrow(opts *bind.CallOpts, _number *big.Int) error {
	var out []interface{}
	err := _Krnldapp.contract.Call(opts, &out, "unprotectedFunctionShouldThrow", _number)

	if err != nil {
		return err
	}

	return err

}

// UnprotectedFunctionShouldThrow is a free data retrieval call binding the contract method 0x36601440.
//
// Solidity: function unprotectedFunctionShouldThrow(uint256 _number) pure returns()
func (_Krnldapp *KrnldappSession) UnprotectedFunctionShouldThrow(_number *big.Int) error {
	return _Krnldapp.Contract.UnprotectedFunctionShouldThrow(&_Krnldapp.CallOpts, _number)
}

// UnprotectedFunctionShouldThrow is a free data retrieval call binding the contract method 0x36601440.
//
// Solidity: function unprotectedFunctionShouldThrow(uint256 _number) pure returns()
func (_Krnldapp *KrnldappCallerSession) UnprotectedFunctionShouldThrow(_number *big.Int) error {
	return _Krnldapp.Contract.UnprotectedFunctionShouldThrow(&_Krnldapp.CallOpts, _number)
}

// ProtectedFunctionality is a paid mutator transaction binding the contract method 0xff806f18.
//
// Solidity: function protectedFunctionality(string name, bytes32 _hash, bytes _signature) returns(uint256)
func (_Krnldapp *KrnldappTransactor) ProtectedFunctionality(opts *bind.TransactOpts, name string, _hash [32]byte, _signature []byte) (*types.Transaction, error) {
	return _Krnldapp.contract.Transact(opts, "protectedFunctionality", name, _hash, _signature)
}

// ProtectedFunctionality is a paid mutator transaction binding the contract method 0xff806f18.
//
// Solidity: function protectedFunctionality(string name, bytes32 _hash, bytes _signature) returns(uint256)
func (_Krnldapp *KrnldappSession) ProtectedFunctionality(name string, _hash [32]byte, _signature []byte) (*types.Transaction, error) {
	return _Krnldapp.Contract.ProtectedFunctionality(&_Krnldapp.TransactOpts, name, _hash, _signature)
}

// ProtectedFunctionality is a paid mutator transaction binding the contract method 0xff806f18.
//
// Solidity: function protectedFunctionality(string name, bytes32 _hash, bytes _signature) returns(uint256)
func (_Krnldapp *KrnldappTransactorSession) ProtectedFunctionality(name string, _hash [32]byte, _signature []byte) (*types.Transaction, error) {
	return _Krnldapp.Contract.ProtectedFunctionality(&_Krnldapp.TransactOpts, name, _hash, _signature)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address _newAuthority) returns()
func (_Krnldapp *KrnldappTransactor) SetAuthority(opts *bind.TransactOpts, _newAuthority common.Address) (*types.Transaction, error) {
	return _Krnldapp.contract.Transact(opts, "setAuthority", _newAuthority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address _newAuthority) returns()
func (_Krnldapp *KrnldappSession) SetAuthority(_newAuthority common.Address) (*types.Transaction, error) {
	return _Krnldapp.Contract.SetAuthority(&_Krnldapp.TransactOpts, _newAuthority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address _newAuthority) returns()
func (_Krnldapp *KrnldappTransactorSession) SetAuthority(_newAuthority common.Address) (*types.Transaction, error) {
	return _Krnldapp.Contract.SetAuthority(&_Krnldapp.TransactOpts, _newAuthority)
}

// UnprotectFunctionShouldReturn is a paid mutator transaction binding the contract method 0x39df3fcb.
//
// Solidity: function unprotectFunctionShouldReturn() returns(uint256)
func (_Krnldapp *KrnldappTransactor) UnprotectFunctionShouldReturn(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Krnldapp.contract.Transact(opts, "unprotectFunctionShouldReturn")
}

// UnprotectFunctionShouldReturn is a paid mutator transaction binding the contract method 0x39df3fcb.
//
// Solidity: function unprotectFunctionShouldReturn() returns(uint256)
func (_Krnldapp *KrnldappSession) UnprotectFunctionShouldReturn() (*types.Transaction, error) {
	return _Krnldapp.Contract.UnprotectFunctionShouldReturn(&_Krnldapp.TransactOpts)
}

// UnprotectFunctionShouldReturn is a paid mutator transaction binding the contract method 0x39df3fcb.
//
// Solidity: function unprotectFunctionShouldReturn() returns(uint256)
func (_Krnldapp *KrnldappTransactorSession) UnprotectFunctionShouldReturn() (*types.Transaction, error) {
	return _Krnldapp.Contract.UnprotectFunctionShouldReturn(&_Krnldapp.TransactOpts)
}

// KrnldappGetCounterIterator is returned from FilterGetCounter and is used to iterate over the raw logs and unpacked data for GetCounter events raised by the Krnldapp contract.
type KrnldappGetCounterIterator struct {
	Event *KrnldappGetCounter // Event containing the contract specifics and raw log

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
func (it *KrnldappGetCounterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KrnldappGetCounter)
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
		it.Event = new(KrnldappGetCounter)
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
func (it *KrnldappGetCounterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KrnldappGetCounterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KrnldappGetCounter represents a GetCounter event raised by the Krnldapp contract.
type KrnldappGetCounter struct {
	From    common.Address
	Counter *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGetCounter is a free log retrieval operation binding the contract event 0xc804064b5fcfee5e4b8f29473de7ddc714f71910769442aaad3917db1f13776a.
//
// Solidity: event GetCounter(address from, uint256 counter)
func (_Krnldapp *KrnldappFilterer) FilterGetCounter(opts *bind.FilterOpts) (*KrnldappGetCounterIterator, error) {

	logs, sub, err := _Krnldapp.contract.FilterLogs(opts, "GetCounter")
	if err != nil {
		return nil, err
	}
	return &KrnldappGetCounterIterator{contract: _Krnldapp.contract, event: "GetCounter", logs: logs, sub: sub}, nil
}

// WatchGetCounter is a free log subscription operation binding the contract event 0xc804064b5fcfee5e4b8f29473de7ddc714f71910769442aaad3917db1f13776a.
//
// Solidity: event GetCounter(address from, uint256 counter)
func (_Krnldapp *KrnldappFilterer) WatchGetCounter(opts *bind.WatchOpts, sink chan<- *KrnldappGetCounter) (event.Subscription, error) {

	logs, sub, err := _Krnldapp.contract.WatchLogs(opts, "GetCounter")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KrnldappGetCounter)
				if err := _Krnldapp.contract.UnpackLog(event, "GetCounter", log); err != nil {
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

// ParseGetCounter is a log parse operation binding the contract event 0xc804064b5fcfee5e4b8f29473de7ddc714f71910769442aaad3917db1f13776a.
//
// Solidity: event GetCounter(address from, uint256 counter)
func (_Krnldapp *KrnldappFilterer) ParseGetCounter(log types.Log) (*KrnldappGetCounter, error) {
	event := new(KrnldappGetCounter)
	if err := _Krnldapp.contract.UnpackLog(event, "GetCounter", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KrnldappSayHiIterator is returned from FilterSayHi and is used to iterate over the raw logs and unpacked data for SayHi events raised by the Krnldapp contract.
type KrnldappSayHiIterator struct {
	Event *KrnldappSayHi // Event containing the contract specifics and raw log

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
func (it *KrnldappSayHiIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KrnldappSayHi)
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
		it.Event = new(KrnldappSayHi)
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
func (it *KrnldappSayHiIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KrnldappSayHiIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KrnldappSayHi represents a SayHi event raised by the Krnldapp contract.
type KrnldappSayHi struct {
	Name string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSayHi is a free log retrieval operation binding the contract event 0x9d44799ad3c5bf4bf9587b56402ddbdca8baa14084514bf35a6fcc0daf00ae7b.
//
// Solidity: event SayHi(string name)
func (_Krnldapp *KrnldappFilterer) FilterSayHi(opts *bind.FilterOpts) (*KrnldappSayHiIterator, error) {

	logs, sub, err := _Krnldapp.contract.FilterLogs(opts, "SayHi")
	if err != nil {
		return nil, err
	}
	return &KrnldappSayHiIterator{contract: _Krnldapp.contract, event: "SayHi", logs: logs, sub: sub}, nil
}

// WatchSayHi is a free log subscription operation binding the contract event 0x9d44799ad3c5bf4bf9587b56402ddbdca8baa14084514bf35a6fcc0daf00ae7b.
//
// Solidity: event SayHi(string name)
func (_Krnldapp *KrnldappFilterer) WatchSayHi(opts *bind.WatchOpts, sink chan<- *KrnldappSayHi) (event.Subscription, error) {

	logs, sub, err := _Krnldapp.contract.WatchLogs(opts, "SayHi")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KrnldappSayHi)
				if err := _Krnldapp.contract.UnpackLog(event, "SayHi", log); err != nil {
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

// ParseSayHi is a log parse operation binding the contract event 0x9d44799ad3c5bf4bf9587b56402ddbdca8baa14084514bf35a6fcc0daf00ae7b.
//
// Solidity: event SayHi(string name)
func (_Krnldapp *KrnldappFilterer) ParseSayHi(log types.Log) (*KrnldappSayHi, error) {
	event := new(KrnldappSayHi)
	if err := _Krnldapp.contract.UnpackLog(event, "SayHi", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
