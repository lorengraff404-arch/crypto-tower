// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// TowerTokenMetaData contains all meta data concerning the TowerToken contract.
var TowerTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_treasury\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_buybackWallet\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_SUPPLY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizeMinter\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedMinters\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"buybackWallet\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"platformTreasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverTokens\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeMinter\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateBuybackWallet\",\"inputs\":[{\"name\":\"newWallet\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTreasury\",\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BuybackWalletUpdated\",\"inputs\":[{\"name\":\"newWallet\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinterAuthorized\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinterRevoked\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TreasuryUpdated\",\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// TowerTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use TowerTokenMetaData.ABI instead.
var TowerTokenABI = TowerTokenMetaData.ABI

// TowerToken is an auto generated Go binding around an Ethereum contract.
type TowerToken struct {
	TowerTokenCaller     // Read-only binding to the contract
	TowerTokenTransactor // Write-only binding to the contract
	TowerTokenFilterer   // Log filterer for contract events
}

// TowerTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TowerTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TowerTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TowerTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TowerTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TowerTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TowerTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TowerTokenSession struct {
	Contract     *TowerToken       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TowerTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TowerTokenCallerSession struct {
	Contract *TowerTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// TowerTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TowerTokenTransactorSession struct {
	Contract     *TowerTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TowerTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TowerTokenRaw struct {
	Contract *TowerToken // Generic contract binding to access the raw methods on
}

// TowerTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TowerTokenCallerRaw struct {
	Contract *TowerTokenCaller // Generic read-only contract binding to access the raw methods on
}

// TowerTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TowerTokenTransactorRaw struct {
	Contract *TowerTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTowerToken creates a new instance of TowerToken, bound to a specific deployed contract.
func NewTowerToken(address common.Address, backend bind.ContractBackend) (*TowerToken, error) {
	contract, err := bindTowerToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TowerToken{TowerTokenCaller: TowerTokenCaller{contract: contract}, TowerTokenTransactor: TowerTokenTransactor{contract: contract}, TowerTokenFilterer: TowerTokenFilterer{contract: contract}}, nil
}

// NewTowerTokenCaller creates a new read-only instance of TowerToken, bound to a specific deployed contract.
func NewTowerTokenCaller(address common.Address, caller bind.ContractCaller) (*TowerTokenCaller, error) {
	contract, err := bindTowerToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TowerTokenCaller{contract: contract}, nil
}

// NewTowerTokenTransactor creates a new write-only instance of TowerToken, bound to a specific deployed contract.
func NewTowerTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TowerTokenTransactor, error) {
	contract, err := bindTowerToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TowerTokenTransactor{contract: contract}, nil
}

// NewTowerTokenFilterer creates a new log filterer instance of TowerToken, bound to a specific deployed contract.
func NewTowerTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TowerTokenFilterer, error) {
	contract, err := bindTowerToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TowerTokenFilterer{contract: contract}, nil
}

// bindTowerToken binds a generic wrapper to an already deployed contract.
func bindTowerToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TowerTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TowerToken *TowerTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TowerToken.Contract.TowerTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TowerToken *TowerTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TowerToken.Contract.TowerTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TowerToken *TowerTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TowerToken.Contract.TowerTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TowerToken *TowerTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TowerToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TowerToken *TowerTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TowerToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TowerToken *TowerTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TowerToken.Contract.contract.Transact(opts, method, params...)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_TowerToken *TowerTokenCaller) MAXSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "MAX_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_TowerToken *TowerTokenSession) MAXSUPPLY() (*big.Int, error) {
	return _TowerToken.Contract.MAXSUPPLY(&_TowerToken.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_TowerToken *TowerTokenCallerSession) MAXSUPPLY() (*big.Int, error) {
	return _TowerToken.Contract.MAXSUPPLY(&_TowerToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TowerToken *TowerTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TowerToken *TowerTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TowerToken.Contract.Allowance(&_TowerToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TowerToken *TowerTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TowerToken.Contract.Allowance(&_TowerToken.CallOpts, owner, spender)
}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_TowerToken *TowerTokenCaller) AuthorizedMinters(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "authorizedMinters", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_TowerToken *TowerTokenSession) AuthorizedMinters(arg0 common.Address) (bool, error) {
	return _TowerToken.Contract.AuthorizedMinters(&_TowerToken.CallOpts, arg0)
}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_TowerToken *TowerTokenCallerSession) AuthorizedMinters(arg0 common.Address) (bool, error) {
	return _TowerToken.Contract.AuthorizedMinters(&_TowerToken.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TowerToken *TowerTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TowerToken *TowerTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TowerToken.Contract.BalanceOf(&_TowerToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TowerToken *TowerTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TowerToken.Contract.BalanceOf(&_TowerToken.CallOpts, account)
}

// BuybackWallet is a free data retrieval call binding the contract method 0xdeab8aea.
//
// Solidity: function buybackWallet() view returns(address)
func (_TowerToken *TowerTokenCaller) BuybackWallet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "buybackWallet")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BuybackWallet is a free data retrieval call binding the contract method 0xdeab8aea.
//
// Solidity: function buybackWallet() view returns(address)
func (_TowerToken *TowerTokenSession) BuybackWallet() (common.Address, error) {
	return _TowerToken.Contract.BuybackWallet(&_TowerToken.CallOpts)
}

// BuybackWallet is a free data retrieval call binding the contract method 0xdeab8aea.
//
// Solidity: function buybackWallet() view returns(address)
func (_TowerToken *TowerTokenCallerSession) BuybackWallet() (common.Address, error) {
	return _TowerToken.Contract.BuybackWallet(&_TowerToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TowerToken *TowerTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TowerToken *TowerTokenSession) Decimals() (uint8, error) {
	return _TowerToken.Contract.Decimals(&_TowerToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TowerToken *TowerTokenCallerSession) Decimals() (uint8, error) {
	return _TowerToken.Contract.Decimals(&_TowerToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TowerToken *TowerTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TowerToken *TowerTokenSession) Name() (string, error) {
	return _TowerToken.Contract.Name(&_TowerToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TowerToken *TowerTokenCallerSession) Name() (string, error) {
	return _TowerToken.Contract.Name(&_TowerToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TowerToken *TowerTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TowerToken *TowerTokenSession) Owner() (common.Address, error) {
	return _TowerToken.Contract.Owner(&_TowerToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TowerToken *TowerTokenCallerSession) Owner() (common.Address, error) {
	return _TowerToken.Contract.Owner(&_TowerToken.CallOpts)
}

// PlatformTreasury is a free data retrieval call binding the contract method 0xe138818c.
//
// Solidity: function platformTreasury() view returns(address)
func (_TowerToken *TowerTokenCaller) PlatformTreasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "platformTreasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PlatformTreasury is a free data retrieval call binding the contract method 0xe138818c.
//
// Solidity: function platformTreasury() view returns(address)
func (_TowerToken *TowerTokenSession) PlatformTreasury() (common.Address, error) {
	return _TowerToken.Contract.PlatformTreasury(&_TowerToken.CallOpts)
}

// PlatformTreasury is a free data retrieval call binding the contract method 0xe138818c.
//
// Solidity: function platformTreasury() view returns(address)
func (_TowerToken *TowerTokenCallerSession) PlatformTreasury() (common.Address, error) {
	return _TowerToken.Contract.PlatformTreasury(&_TowerToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TowerToken *TowerTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TowerToken *TowerTokenSession) Symbol() (string, error) {
	return _TowerToken.Contract.Symbol(&_TowerToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TowerToken *TowerTokenCallerSession) Symbol() (string, error) {
	return _TowerToken.Contract.Symbol(&_TowerToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TowerToken *TowerTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TowerToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TowerToken *TowerTokenSession) TotalSupply() (*big.Int, error) {
	return _TowerToken.Contract.TotalSupply(&_TowerToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TowerToken *TowerTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _TowerToken.Contract.TotalSupply(&_TowerToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TowerToken *TowerTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TowerToken *TowerTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.Approve(&_TowerToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TowerToken *TowerTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.Approve(&_TowerToken.TransactOpts, spender, value)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_TowerToken *TowerTokenTransactor) AuthorizeMinter(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "authorizeMinter", minter)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_TowerToken *TowerTokenSession) AuthorizeMinter(minter common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.AuthorizeMinter(&_TowerToken.TransactOpts, minter)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_TowerToken *TowerTokenTransactorSession) AuthorizeMinter(minter common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.AuthorizeMinter(&_TowerToken.TransactOpts, minter)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_TowerToken *TowerTokenTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_TowerToken *TowerTokenSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.Burn(&_TowerToken.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_TowerToken *TowerTokenTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.Burn(&_TowerToken.TransactOpts, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_TowerToken *TowerTokenTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "burnFrom", account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_TowerToken *TowerTokenSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.BurnFrom(&_TowerToken.TransactOpts, account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_TowerToken *TowerTokenTransactorSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.BurnFrom(&_TowerToken.TransactOpts, account, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TowerToken *TowerTokenTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TowerToken *TowerTokenSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.Mint(&_TowerToken.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TowerToken *TowerTokenTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.Mint(&_TowerToken.TransactOpts, to, amount)
}

// RecoverTokens is a paid mutator transaction binding the contract method 0x069c9fae.
//
// Solidity: function recoverTokens(address tokenAddress, uint256 amount) returns()
func (_TowerToken *TowerTokenTransactor) RecoverTokens(opts *bind.TransactOpts, tokenAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "recoverTokens", tokenAddress, amount)
}

// RecoverTokens is a paid mutator transaction binding the contract method 0x069c9fae.
//
// Solidity: function recoverTokens(address tokenAddress, uint256 amount) returns()
func (_TowerToken *TowerTokenSession) RecoverTokens(tokenAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.RecoverTokens(&_TowerToken.TransactOpts, tokenAddress, amount)
}

// RecoverTokens is a paid mutator transaction binding the contract method 0x069c9fae.
//
// Solidity: function recoverTokens(address tokenAddress, uint256 amount) returns()
func (_TowerToken *TowerTokenTransactorSession) RecoverTokens(tokenAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.RecoverTokens(&_TowerToken.TransactOpts, tokenAddress, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TowerToken *TowerTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TowerToken *TowerTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _TowerToken.Contract.RenounceOwnership(&_TowerToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TowerToken *TowerTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TowerToken.Contract.RenounceOwnership(&_TowerToken.TransactOpts)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_TowerToken *TowerTokenTransactor) RevokeMinter(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "revokeMinter", minter)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_TowerToken *TowerTokenSession) RevokeMinter(minter common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.RevokeMinter(&_TowerToken.TransactOpts, minter)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_TowerToken *TowerTokenTransactorSession) RevokeMinter(minter common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.RevokeMinter(&_TowerToken.TransactOpts, minter)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TowerToken *TowerTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TowerToken *TowerTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.Transfer(&_TowerToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TowerToken *TowerTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.Transfer(&_TowerToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TowerToken *TowerTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TowerToken *TowerTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.TransferFrom(&_TowerToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TowerToken *TowerTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TowerToken.Contract.TransferFrom(&_TowerToken.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TowerToken *TowerTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TowerToken *TowerTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.TransferOwnership(&_TowerToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TowerToken *TowerTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.TransferOwnership(&_TowerToken.TransactOpts, newOwner)
}

// UpdateBuybackWallet is a paid mutator transaction binding the contract method 0x04dacd50.
//
// Solidity: function updateBuybackWallet(address newWallet) returns()
func (_TowerToken *TowerTokenTransactor) UpdateBuybackWallet(opts *bind.TransactOpts, newWallet common.Address) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "updateBuybackWallet", newWallet)
}

// UpdateBuybackWallet is a paid mutator transaction binding the contract method 0x04dacd50.
//
// Solidity: function updateBuybackWallet(address newWallet) returns()
func (_TowerToken *TowerTokenSession) UpdateBuybackWallet(newWallet common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.UpdateBuybackWallet(&_TowerToken.TransactOpts, newWallet)
}

// UpdateBuybackWallet is a paid mutator transaction binding the contract method 0x04dacd50.
//
// Solidity: function updateBuybackWallet(address newWallet) returns()
func (_TowerToken *TowerTokenTransactorSession) UpdateBuybackWallet(newWallet common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.UpdateBuybackWallet(&_TowerToken.TransactOpts, newWallet)
}

// UpdateTreasury is a paid mutator transaction binding the contract method 0x7f51bb1f.
//
// Solidity: function updateTreasury(address newTreasury) returns()
func (_TowerToken *TowerTokenTransactor) UpdateTreasury(opts *bind.TransactOpts, newTreasury common.Address) (*types.Transaction, error) {
	return _TowerToken.contract.Transact(opts, "updateTreasury", newTreasury)
}

// UpdateTreasury is a paid mutator transaction binding the contract method 0x7f51bb1f.
//
// Solidity: function updateTreasury(address newTreasury) returns()
func (_TowerToken *TowerTokenSession) UpdateTreasury(newTreasury common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.UpdateTreasury(&_TowerToken.TransactOpts, newTreasury)
}

// UpdateTreasury is a paid mutator transaction binding the contract method 0x7f51bb1f.
//
// Solidity: function updateTreasury(address newTreasury) returns()
func (_TowerToken *TowerTokenTransactorSession) UpdateTreasury(newTreasury common.Address) (*types.Transaction, error) {
	return _TowerToken.Contract.UpdateTreasury(&_TowerToken.TransactOpts, newTreasury)
}

// TowerTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TowerToken contract.
type TowerTokenApprovalIterator struct {
	Event *TowerTokenApproval // Event containing the contract specifics and raw log

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
func (it *TowerTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TowerTokenApproval)
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
		it.Event = new(TowerTokenApproval)
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
func (it *TowerTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TowerTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TowerTokenApproval represents a Approval event raised by the TowerToken contract.
type TowerTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TowerToken *TowerTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TowerTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TowerToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TowerTokenApprovalIterator{contract: _TowerToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TowerToken *TowerTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TowerTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TowerToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TowerTokenApproval)
				if err := _TowerToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TowerToken *TowerTokenFilterer) ParseApproval(log types.Log) (*TowerTokenApproval, error) {
	event := new(TowerTokenApproval)
	if err := _TowerToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TowerTokenBuybackWalletUpdatedIterator is returned from FilterBuybackWalletUpdated and is used to iterate over the raw logs and unpacked data for BuybackWalletUpdated events raised by the TowerToken contract.
type TowerTokenBuybackWalletUpdatedIterator struct {
	Event *TowerTokenBuybackWalletUpdated // Event containing the contract specifics and raw log

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
func (it *TowerTokenBuybackWalletUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TowerTokenBuybackWalletUpdated)
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
		it.Event = new(TowerTokenBuybackWalletUpdated)
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
func (it *TowerTokenBuybackWalletUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TowerTokenBuybackWalletUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TowerTokenBuybackWalletUpdated represents a BuybackWalletUpdated event raised by the TowerToken contract.
type TowerTokenBuybackWalletUpdated struct {
	NewWallet common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBuybackWalletUpdated is a free log retrieval operation binding the contract event 0x2e5665e81af8928ff3f79e82650b68c9459cde711420effbff19b3ccee749f81.
//
// Solidity: event BuybackWalletUpdated(address indexed newWallet)
func (_TowerToken *TowerTokenFilterer) FilterBuybackWalletUpdated(opts *bind.FilterOpts, newWallet []common.Address) (*TowerTokenBuybackWalletUpdatedIterator, error) {

	var newWalletRule []interface{}
	for _, newWalletItem := range newWallet {
		newWalletRule = append(newWalletRule, newWalletItem)
	}

	logs, sub, err := _TowerToken.contract.FilterLogs(opts, "BuybackWalletUpdated", newWalletRule)
	if err != nil {
		return nil, err
	}
	return &TowerTokenBuybackWalletUpdatedIterator{contract: _TowerToken.contract, event: "BuybackWalletUpdated", logs: logs, sub: sub}, nil
}

// WatchBuybackWalletUpdated is a free log subscription operation binding the contract event 0x2e5665e81af8928ff3f79e82650b68c9459cde711420effbff19b3ccee749f81.
//
// Solidity: event BuybackWalletUpdated(address indexed newWallet)
func (_TowerToken *TowerTokenFilterer) WatchBuybackWalletUpdated(opts *bind.WatchOpts, sink chan<- *TowerTokenBuybackWalletUpdated, newWallet []common.Address) (event.Subscription, error) {

	var newWalletRule []interface{}
	for _, newWalletItem := range newWallet {
		newWalletRule = append(newWalletRule, newWalletItem)
	}

	logs, sub, err := _TowerToken.contract.WatchLogs(opts, "BuybackWalletUpdated", newWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TowerTokenBuybackWalletUpdated)
				if err := _TowerToken.contract.UnpackLog(event, "BuybackWalletUpdated", log); err != nil {
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

// ParseBuybackWalletUpdated is a log parse operation binding the contract event 0x2e5665e81af8928ff3f79e82650b68c9459cde711420effbff19b3ccee749f81.
//
// Solidity: event BuybackWalletUpdated(address indexed newWallet)
func (_TowerToken *TowerTokenFilterer) ParseBuybackWalletUpdated(log types.Log) (*TowerTokenBuybackWalletUpdated, error) {
	event := new(TowerTokenBuybackWalletUpdated)
	if err := _TowerToken.contract.UnpackLog(event, "BuybackWalletUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TowerTokenMinterAuthorizedIterator is returned from FilterMinterAuthorized and is used to iterate over the raw logs and unpacked data for MinterAuthorized events raised by the TowerToken contract.
type TowerTokenMinterAuthorizedIterator struct {
	Event *TowerTokenMinterAuthorized // Event containing the contract specifics and raw log

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
func (it *TowerTokenMinterAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TowerTokenMinterAuthorized)
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
		it.Event = new(TowerTokenMinterAuthorized)
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
func (it *TowerTokenMinterAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TowerTokenMinterAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TowerTokenMinterAuthorized represents a MinterAuthorized event raised by the TowerToken contract.
type TowerTokenMinterAuthorized struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterAuthorized is a free log retrieval operation binding the contract event 0x83b05b6735acd4b85e3bded8e72c851d1a87718f81e3c8e6f0c9d9a2baa88e46.
//
// Solidity: event MinterAuthorized(address indexed minter)
func (_TowerToken *TowerTokenFilterer) FilterMinterAuthorized(opts *bind.FilterOpts, minter []common.Address) (*TowerTokenMinterAuthorizedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _TowerToken.contract.FilterLogs(opts, "MinterAuthorized", minterRule)
	if err != nil {
		return nil, err
	}
	return &TowerTokenMinterAuthorizedIterator{contract: _TowerToken.contract, event: "MinterAuthorized", logs: logs, sub: sub}, nil
}

// WatchMinterAuthorized is a free log subscription operation binding the contract event 0x83b05b6735acd4b85e3bded8e72c851d1a87718f81e3c8e6f0c9d9a2baa88e46.
//
// Solidity: event MinterAuthorized(address indexed minter)
func (_TowerToken *TowerTokenFilterer) WatchMinterAuthorized(opts *bind.WatchOpts, sink chan<- *TowerTokenMinterAuthorized, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _TowerToken.contract.WatchLogs(opts, "MinterAuthorized", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TowerTokenMinterAuthorized)
				if err := _TowerToken.contract.UnpackLog(event, "MinterAuthorized", log); err != nil {
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

// ParseMinterAuthorized is a log parse operation binding the contract event 0x83b05b6735acd4b85e3bded8e72c851d1a87718f81e3c8e6f0c9d9a2baa88e46.
//
// Solidity: event MinterAuthorized(address indexed minter)
func (_TowerToken *TowerTokenFilterer) ParseMinterAuthorized(log types.Log) (*TowerTokenMinterAuthorized, error) {
	event := new(TowerTokenMinterAuthorized)
	if err := _TowerToken.contract.UnpackLog(event, "MinterAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TowerTokenMinterRevokedIterator is returned from FilterMinterRevoked and is used to iterate over the raw logs and unpacked data for MinterRevoked events raised by the TowerToken contract.
type TowerTokenMinterRevokedIterator struct {
	Event *TowerTokenMinterRevoked // Event containing the contract specifics and raw log

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
func (it *TowerTokenMinterRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TowerTokenMinterRevoked)
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
		it.Event = new(TowerTokenMinterRevoked)
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
func (it *TowerTokenMinterRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TowerTokenMinterRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TowerTokenMinterRevoked represents a MinterRevoked event raised by the TowerToken contract.
type TowerTokenMinterRevoked struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterRevoked is a free log retrieval operation binding the contract event 0x44f4322f8daa225d5f4877ad0f7d3dfba248a774396f3ca99405ed40a044fe81.
//
// Solidity: event MinterRevoked(address indexed minter)
func (_TowerToken *TowerTokenFilterer) FilterMinterRevoked(opts *bind.FilterOpts, minter []common.Address) (*TowerTokenMinterRevokedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _TowerToken.contract.FilterLogs(opts, "MinterRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return &TowerTokenMinterRevokedIterator{contract: _TowerToken.contract, event: "MinterRevoked", logs: logs, sub: sub}, nil
}

// WatchMinterRevoked is a free log subscription operation binding the contract event 0x44f4322f8daa225d5f4877ad0f7d3dfba248a774396f3ca99405ed40a044fe81.
//
// Solidity: event MinterRevoked(address indexed minter)
func (_TowerToken *TowerTokenFilterer) WatchMinterRevoked(opts *bind.WatchOpts, sink chan<- *TowerTokenMinterRevoked, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _TowerToken.contract.WatchLogs(opts, "MinterRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TowerTokenMinterRevoked)
				if err := _TowerToken.contract.UnpackLog(event, "MinterRevoked", log); err != nil {
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

// ParseMinterRevoked is a log parse operation binding the contract event 0x44f4322f8daa225d5f4877ad0f7d3dfba248a774396f3ca99405ed40a044fe81.
//
// Solidity: event MinterRevoked(address indexed minter)
func (_TowerToken *TowerTokenFilterer) ParseMinterRevoked(log types.Log) (*TowerTokenMinterRevoked, error) {
	event := new(TowerTokenMinterRevoked)
	if err := _TowerToken.contract.UnpackLog(event, "MinterRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TowerTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TowerToken contract.
type TowerTokenOwnershipTransferredIterator struct {
	Event *TowerTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TowerTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TowerTokenOwnershipTransferred)
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
		it.Event = new(TowerTokenOwnershipTransferred)
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
func (it *TowerTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TowerTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TowerTokenOwnershipTransferred represents a OwnershipTransferred event raised by the TowerToken contract.
type TowerTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TowerToken *TowerTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TowerTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TowerToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TowerTokenOwnershipTransferredIterator{contract: _TowerToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TowerToken *TowerTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TowerTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TowerToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TowerTokenOwnershipTransferred)
				if err := _TowerToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TowerToken *TowerTokenFilterer) ParseOwnershipTransferred(log types.Log) (*TowerTokenOwnershipTransferred, error) {
	event := new(TowerTokenOwnershipTransferred)
	if err := _TowerToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TowerTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TowerToken contract.
type TowerTokenTransferIterator struct {
	Event *TowerTokenTransfer // Event containing the contract specifics and raw log

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
func (it *TowerTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TowerTokenTransfer)
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
		it.Event = new(TowerTokenTransfer)
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
func (it *TowerTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TowerTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TowerTokenTransfer represents a Transfer event raised by the TowerToken contract.
type TowerTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TowerToken *TowerTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TowerTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TowerToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TowerTokenTransferIterator{contract: _TowerToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TowerToken *TowerTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TowerTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TowerToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TowerTokenTransfer)
				if err := _TowerToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TowerToken *TowerTokenFilterer) ParseTransfer(log types.Log) (*TowerTokenTransfer, error) {
	event := new(TowerTokenTransfer)
	if err := _TowerToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TowerTokenTreasuryUpdatedIterator is returned from FilterTreasuryUpdated and is used to iterate over the raw logs and unpacked data for TreasuryUpdated events raised by the TowerToken contract.
type TowerTokenTreasuryUpdatedIterator struct {
	Event *TowerTokenTreasuryUpdated // Event containing the contract specifics and raw log

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
func (it *TowerTokenTreasuryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TowerTokenTreasuryUpdated)
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
		it.Event = new(TowerTokenTreasuryUpdated)
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
func (it *TowerTokenTreasuryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TowerTokenTreasuryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TowerTokenTreasuryUpdated represents a TreasuryUpdated event raised by the TowerToken contract.
type TowerTokenTreasuryUpdated struct {
	NewTreasury common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTreasuryUpdated is a free log retrieval operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address indexed newTreasury)
func (_TowerToken *TowerTokenFilterer) FilterTreasuryUpdated(opts *bind.FilterOpts, newTreasury []common.Address) (*TowerTokenTreasuryUpdatedIterator, error) {

	var newTreasuryRule []interface{}
	for _, newTreasuryItem := range newTreasury {
		newTreasuryRule = append(newTreasuryRule, newTreasuryItem)
	}

	logs, sub, err := _TowerToken.contract.FilterLogs(opts, "TreasuryUpdated", newTreasuryRule)
	if err != nil {
		return nil, err
	}
	return &TowerTokenTreasuryUpdatedIterator{contract: _TowerToken.contract, event: "TreasuryUpdated", logs: logs, sub: sub}, nil
}

// WatchTreasuryUpdated is a free log subscription operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address indexed newTreasury)
func (_TowerToken *TowerTokenFilterer) WatchTreasuryUpdated(opts *bind.WatchOpts, sink chan<- *TowerTokenTreasuryUpdated, newTreasury []common.Address) (event.Subscription, error) {

	var newTreasuryRule []interface{}
	for _, newTreasuryItem := range newTreasury {
		newTreasuryRule = append(newTreasuryRule, newTreasuryItem)
	}

	logs, sub, err := _TowerToken.contract.WatchLogs(opts, "TreasuryUpdated", newTreasuryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TowerTokenTreasuryUpdated)
				if err := _TowerToken.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
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

// ParseTreasuryUpdated is a log parse operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address indexed newTreasury)
func (_TowerToken *TowerTokenFilterer) ParseTreasuryUpdated(log types.Log) (*TowerTokenTreasuryUpdated, error) {
	event := new(TowerTokenTreasuryUpdated)
	if err := _TowerToken.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
