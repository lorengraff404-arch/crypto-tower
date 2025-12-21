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

// ItemNFTItemData is an auto generated low-level Go binding around an user-defined struct.
type ItemNFTItemData struct {
	ItemType     string
	Name         string
	Rarity       string
	IsConsumable bool
	MaxSupply    *big.Int
	CreatedAt    *big.Int
}

// ItemNFTMetaData contains all meta data concerning the ItemNFT contract.
var ItemNFTMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizeMinter\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedMinters\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOfBatch\",\"inputs\":[{\"name\":\"accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createItem\",\"inputs\":[{\"name\":\"itemType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"rarity\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isConsumable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"exists\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getItemData\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structItemNFT.ItemData\",\"components\":[{\"name\":\"itemType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"rarity\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isConsumable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"items\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"itemType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"rarity\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isConsumable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mintBatch\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeMinter\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeBatchTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uri\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ItemCreated\",\"inputs\":[{\"name\":\"itemId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"itemType\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"rarity\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinterAuthorized\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinterRevoked\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferBatch\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ids\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferSingle\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"URI\",\"inputs\":[{\"name\":\"value\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC1155InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidArrayLength\",\"inputs\":[{\"name\":\"idsLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"valuesLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1155MissingApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// ItemNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use ItemNFTMetaData.ABI instead.
var ItemNFTABI = ItemNFTMetaData.ABI

// ItemNFT is an auto generated Go binding around an Ethereum contract.
type ItemNFT struct {
	ItemNFTCaller     // Read-only binding to the contract
	ItemNFTTransactor // Write-only binding to the contract
	ItemNFTFilterer   // Log filterer for contract events
}

// ItemNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type ItemNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ItemNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ItemNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ItemNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ItemNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ItemNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ItemNFTSession struct {
	Contract     *ItemNFT          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ItemNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ItemNFTCallerSession struct {
	Contract *ItemNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ItemNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ItemNFTTransactorSession struct {
	Contract     *ItemNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ItemNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type ItemNFTRaw struct {
	Contract *ItemNFT // Generic contract binding to access the raw methods on
}

// ItemNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ItemNFTCallerRaw struct {
	Contract *ItemNFTCaller // Generic read-only contract binding to access the raw methods on
}

// ItemNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ItemNFTTransactorRaw struct {
	Contract *ItemNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewItemNFT creates a new instance of ItemNFT, bound to a specific deployed contract.
func NewItemNFT(address common.Address, backend bind.ContractBackend) (*ItemNFT, error) {
	contract, err := bindItemNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ItemNFT{ItemNFTCaller: ItemNFTCaller{contract: contract}, ItemNFTTransactor: ItemNFTTransactor{contract: contract}, ItemNFTFilterer: ItemNFTFilterer{contract: contract}}, nil
}

// NewItemNFTCaller creates a new read-only instance of ItemNFT, bound to a specific deployed contract.
func NewItemNFTCaller(address common.Address, caller bind.ContractCaller) (*ItemNFTCaller, error) {
	contract, err := bindItemNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ItemNFTCaller{contract: contract}, nil
}

// NewItemNFTTransactor creates a new write-only instance of ItemNFT, bound to a specific deployed contract.
func NewItemNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*ItemNFTTransactor, error) {
	contract, err := bindItemNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ItemNFTTransactor{contract: contract}, nil
}

// NewItemNFTFilterer creates a new log filterer instance of ItemNFT, bound to a specific deployed contract.
func NewItemNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*ItemNFTFilterer, error) {
	contract, err := bindItemNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ItemNFTFilterer{contract: contract}, nil
}

// bindItemNFT binds a generic wrapper to an already deployed contract.
func bindItemNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ItemNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ItemNFT *ItemNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ItemNFT.Contract.ItemNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ItemNFT *ItemNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ItemNFT.Contract.ItemNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ItemNFT *ItemNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ItemNFT.Contract.ItemNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ItemNFT *ItemNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ItemNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ItemNFT *ItemNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ItemNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ItemNFT *ItemNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ItemNFT.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_ItemNFT *ItemNFTCaller) AuthorizedMinters(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "authorizedMinters", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_ItemNFT *ItemNFTSession) AuthorizedMinters(arg0 common.Address) (bool, error) {
	return _ItemNFT.Contract.AuthorizedMinters(&_ItemNFT.CallOpts, arg0)
}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_ItemNFT *ItemNFTCallerSession) AuthorizedMinters(arg0 common.Address) (bool, error) {
	return _ItemNFT.Contract.AuthorizedMinters(&_ItemNFT.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_ItemNFT *ItemNFTCaller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_ItemNFT *ItemNFTSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _ItemNFT.Contract.BalanceOf(&_ItemNFT.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_ItemNFT *ItemNFTCallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _ItemNFT.Contract.BalanceOf(&_ItemNFT.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_ItemNFT *ItemNFTCaller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_ItemNFT *ItemNFTSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _ItemNFT.Contract.BalanceOfBatch(&_ItemNFT.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_ItemNFT *ItemNFTCallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _ItemNFT.Contract.BalanceOfBatch(&_ItemNFT.CallOpts, accounts, ids)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_ItemNFT *ItemNFTCaller) Exists(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "exists", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_ItemNFT *ItemNFTSession) Exists(id *big.Int) (bool, error) {
	return _ItemNFT.Contract.Exists(&_ItemNFT.CallOpts, id)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_ItemNFT *ItemNFTCallerSession) Exists(id *big.Int) (bool, error) {
	return _ItemNFT.Contract.Exists(&_ItemNFT.CallOpts, id)
}

// GetItemData is a free data retrieval call binding the contract method 0x8bc6976e.
//
// Solidity: function getItemData(uint256 id) view returns((string,string,string,bool,uint256,uint256))
func (_ItemNFT *ItemNFTCaller) GetItemData(opts *bind.CallOpts, id *big.Int) (ItemNFTItemData, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "getItemData", id)

	if err != nil {
		return *new(ItemNFTItemData), err
	}

	out0 := *abi.ConvertType(out[0], new(ItemNFTItemData)).(*ItemNFTItemData)

	return out0, err

}

// GetItemData is a free data retrieval call binding the contract method 0x8bc6976e.
//
// Solidity: function getItemData(uint256 id) view returns((string,string,string,bool,uint256,uint256))
func (_ItemNFT *ItemNFTSession) GetItemData(id *big.Int) (ItemNFTItemData, error) {
	return _ItemNFT.Contract.GetItemData(&_ItemNFT.CallOpts, id)
}

// GetItemData is a free data retrieval call binding the contract method 0x8bc6976e.
//
// Solidity: function getItemData(uint256 id) view returns((string,string,string,bool,uint256,uint256))
func (_ItemNFT *ItemNFTCallerSession) GetItemData(id *big.Int) (ItemNFTItemData, error) {
	return _ItemNFT.Contract.GetItemData(&_ItemNFT.CallOpts, id)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_ItemNFT *ItemNFTCaller) IsApprovedForAll(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "isApprovedForAll", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_ItemNFT *ItemNFTSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _ItemNFT.Contract.IsApprovedForAll(&_ItemNFT.CallOpts, account, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_ItemNFT *ItemNFTCallerSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _ItemNFT.Contract.IsApprovedForAll(&_ItemNFT.CallOpts, account, operator)
}

// Items is a free data retrieval call binding the contract method 0xbfb231d2.
//
// Solidity: function items(uint256 ) view returns(string itemType, string name, string rarity, bool isConsumable, uint256 maxSupply, uint256 createdAt)
func (_ItemNFT *ItemNFTCaller) Items(opts *bind.CallOpts, arg0 *big.Int) (struct {
	ItemType     string
	Name         string
	Rarity       string
	IsConsumable bool
	MaxSupply    *big.Int
	CreatedAt    *big.Int
}, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "items", arg0)

	outstruct := new(struct {
		ItemType     string
		Name         string
		Rarity       string
		IsConsumable bool
		MaxSupply    *big.Int
		CreatedAt    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ItemType = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Rarity = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.IsConsumable = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.MaxSupply = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.CreatedAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Items is a free data retrieval call binding the contract method 0xbfb231d2.
//
// Solidity: function items(uint256 ) view returns(string itemType, string name, string rarity, bool isConsumable, uint256 maxSupply, uint256 createdAt)
func (_ItemNFT *ItemNFTSession) Items(arg0 *big.Int) (struct {
	ItemType     string
	Name         string
	Rarity       string
	IsConsumable bool
	MaxSupply    *big.Int
	CreatedAt    *big.Int
}, error) {
	return _ItemNFT.Contract.Items(&_ItemNFT.CallOpts, arg0)
}

// Items is a free data retrieval call binding the contract method 0xbfb231d2.
//
// Solidity: function items(uint256 ) view returns(string itemType, string name, string rarity, bool isConsumable, uint256 maxSupply, uint256 createdAt)
func (_ItemNFT *ItemNFTCallerSession) Items(arg0 *big.Int) (struct {
	ItemType     string
	Name         string
	Rarity       string
	IsConsumable bool
	MaxSupply    *big.Int
	CreatedAt    *big.Int
}, error) {
	return _ItemNFT.Contract.Items(&_ItemNFT.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ItemNFT *ItemNFTCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ItemNFT *ItemNFTSession) Owner() (common.Address, error) {
	return _ItemNFT.Contract.Owner(&_ItemNFT.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ItemNFT *ItemNFTCallerSession) Owner() (common.Address, error) {
	return _ItemNFT.Contract.Owner(&_ItemNFT.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ItemNFT *ItemNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ItemNFT *ItemNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ItemNFT.Contract.SupportsInterface(&_ItemNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ItemNFT *ItemNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ItemNFT.Contract.SupportsInterface(&_ItemNFT.CallOpts, interfaceId)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ItemNFT *ItemNFTCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ItemNFT *ItemNFTSession) TotalSupply() (*big.Int, error) {
	return _ItemNFT.Contract.TotalSupply(&_ItemNFT.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ItemNFT *ItemNFTCallerSession) TotalSupply() (*big.Int, error) {
	return _ItemNFT.Contract.TotalSupply(&_ItemNFT.CallOpts)
}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_ItemNFT *ItemNFTCaller) TotalSupply0(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "totalSupply0", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_ItemNFT *ItemNFTSession) TotalSupply0(id *big.Int) (*big.Int, error) {
	return _ItemNFT.Contract.TotalSupply0(&_ItemNFT.CallOpts, id)
}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_ItemNFT *ItemNFTCallerSession) TotalSupply0(id *big.Int) (*big.Int, error) {
	return _ItemNFT.Contract.TotalSupply0(&_ItemNFT.CallOpts, id)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_ItemNFT *ItemNFTCaller) Uri(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _ItemNFT.contract.Call(opts, &out, "uri", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_ItemNFT *ItemNFTSession) Uri(tokenId *big.Int) (string, error) {
	return _ItemNFT.Contract.Uri(&_ItemNFT.CallOpts, tokenId)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_ItemNFT *ItemNFTCallerSession) Uri(tokenId *big.Int) (string, error) {
	return _ItemNFT.Contract.Uri(&_ItemNFT.CallOpts, tokenId)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_ItemNFT *ItemNFTTransactor) AuthorizeMinter(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "authorizeMinter", minter)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_ItemNFT *ItemNFTSession) AuthorizeMinter(minter common.Address) (*types.Transaction, error) {
	return _ItemNFT.Contract.AuthorizeMinter(&_ItemNFT.TransactOpts, minter)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_ItemNFT *ItemNFTTransactorSession) AuthorizeMinter(minter common.Address) (*types.Transaction, error) {
	return _ItemNFT.Contract.AuthorizeMinter(&_ItemNFT.TransactOpts, minter)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 id, uint256 amount) returns()
func (_ItemNFT *ItemNFTTransactor) Burn(opts *bind.TransactOpts, from common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "burn", from, id, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 id, uint256 amount) returns()
func (_ItemNFT *ItemNFTSession) Burn(from common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ItemNFT.Contract.Burn(&_ItemNFT.TransactOpts, from, id, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 id, uint256 amount) returns()
func (_ItemNFT *ItemNFTTransactorSession) Burn(from common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ItemNFT.Contract.Burn(&_ItemNFT.TransactOpts, from, id, amount)
}

// CreateItem is a paid mutator transaction binding the contract method 0x62f41b3c.
//
// Solidity: function createItem(string itemType, string name, string rarity, bool isConsumable, uint256 maxSupply, string tokenURI) returns(uint256)
func (_ItemNFT *ItemNFTTransactor) CreateItem(opts *bind.TransactOpts, itemType string, name string, rarity string, isConsumable bool, maxSupply *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "createItem", itemType, name, rarity, isConsumable, maxSupply, tokenURI)
}

// CreateItem is a paid mutator transaction binding the contract method 0x62f41b3c.
//
// Solidity: function createItem(string itemType, string name, string rarity, bool isConsumable, uint256 maxSupply, string tokenURI) returns(uint256)
func (_ItemNFT *ItemNFTSession) CreateItem(itemType string, name string, rarity string, isConsumable bool, maxSupply *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ItemNFT.Contract.CreateItem(&_ItemNFT.TransactOpts, itemType, name, rarity, isConsumable, maxSupply, tokenURI)
}

// CreateItem is a paid mutator transaction binding the contract method 0x62f41b3c.
//
// Solidity: function createItem(string itemType, string name, string rarity, bool isConsumable, uint256 maxSupply, string tokenURI) returns(uint256)
func (_ItemNFT *ItemNFTTransactorSession) CreateItem(itemType string, name string, rarity string, isConsumable bool, maxSupply *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ItemNFT.Contract.CreateItem(&_ItemNFT.TransactOpts, itemType, name, rarity, isConsumable, maxSupply, tokenURI)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 id, uint256 amount) returns()
func (_ItemNFT *ItemNFTTransactor) Mint(opts *bind.TransactOpts, to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "mint", to, id, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 id, uint256 amount) returns()
func (_ItemNFT *ItemNFTSession) Mint(to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ItemNFT.Contract.Mint(&_ItemNFT.TransactOpts, to, id, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 id, uint256 amount) returns()
func (_ItemNFT *ItemNFTTransactorSession) Mint(to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ItemNFT.Contract.Mint(&_ItemNFT.TransactOpts, to, id, amount)
}

// MintBatch is a paid mutator transaction binding the contract method 0xd81d0a15.
//
// Solidity: function mintBatch(address to, uint256[] ids, uint256[] amounts) returns()
func (_ItemNFT *ItemNFTTransactor) MintBatch(opts *bind.TransactOpts, to common.Address, ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "mintBatch", to, ids, amounts)
}

// MintBatch is a paid mutator transaction binding the contract method 0xd81d0a15.
//
// Solidity: function mintBatch(address to, uint256[] ids, uint256[] amounts) returns()
func (_ItemNFT *ItemNFTSession) MintBatch(to common.Address, ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _ItemNFT.Contract.MintBatch(&_ItemNFT.TransactOpts, to, ids, amounts)
}

// MintBatch is a paid mutator transaction binding the contract method 0xd81d0a15.
//
// Solidity: function mintBatch(address to, uint256[] ids, uint256[] amounts) returns()
func (_ItemNFT *ItemNFTTransactorSession) MintBatch(to common.Address, ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _ItemNFT.Contract.MintBatch(&_ItemNFT.TransactOpts, to, ids, amounts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ItemNFT *ItemNFTTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ItemNFT *ItemNFTSession) RenounceOwnership() (*types.Transaction, error) {
	return _ItemNFT.Contract.RenounceOwnership(&_ItemNFT.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ItemNFT *ItemNFTTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ItemNFT.Contract.RenounceOwnership(&_ItemNFT.TransactOpts)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_ItemNFT *ItemNFTTransactor) RevokeMinter(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "revokeMinter", minter)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_ItemNFT *ItemNFTSession) RevokeMinter(minter common.Address) (*types.Transaction, error) {
	return _ItemNFT.Contract.RevokeMinter(&_ItemNFT.TransactOpts, minter)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_ItemNFT *ItemNFTTransactorSession) RevokeMinter(minter common.Address) (*types.Transaction, error) {
	return _ItemNFT.Contract.RevokeMinter(&_ItemNFT.TransactOpts, minter)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_ItemNFT *ItemNFTTransactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_ItemNFT *ItemNFTSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _ItemNFT.Contract.SafeBatchTransferFrom(&_ItemNFT.TransactOpts, from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_ItemNFT *ItemNFTTransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _ItemNFT.Contract.SafeBatchTransferFrom(&_ItemNFT.TransactOpts, from, to, ids, values, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_ItemNFT *ItemNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "safeTransferFrom", from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_ItemNFT *ItemNFTSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ItemNFT.Contract.SafeTransferFrom(&_ItemNFT.TransactOpts, from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_ItemNFT *ItemNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ItemNFT.Contract.SafeTransferFrom(&_ItemNFT.TransactOpts, from, to, id, value, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ItemNFT *ItemNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ItemNFT *ItemNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ItemNFT.Contract.SetApprovalForAll(&_ItemNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ItemNFT *ItemNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ItemNFT.Contract.SetApprovalForAll(&_ItemNFT.TransactOpts, operator, approved)
}

// SetURI is a paid mutator transaction binding the contract method 0x862440e2.
//
// Solidity: function setURI(uint256 tokenId, string tokenURI) returns()
func (_ItemNFT *ItemNFTTransactor) SetURI(opts *bind.TransactOpts, tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "setURI", tokenId, tokenURI)
}

// SetURI is a paid mutator transaction binding the contract method 0x862440e2.
//
// Solidity: function setURI(uint256 tokenId, string tokenURI) returns()
func (_ItemNFT *ItemNFTSession) SetURI(tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ItemNFT.Contract.SetURI(&_ItemNFT.TransactOpts, tokenId, tokenURI)
}

// SetURI is a paid mutator transaction binding the contract method 0x862440e2.
//
// Solidity: function setURI(uint256 tokenId, string tokenURI) returns()
func (_ItemNFT *ItemNFTTransactorSession) SetURI(tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ItemNFT.Contract.SetURI(&_ItemNFT.TransactOpts, tokenId, tokenURI)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ItemNFT *ItemNFTTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ItemNFT.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ItemNFT *ItemNFTSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ItemNFT.Contract.TransferOwnership(&_ItemNFT.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ItemNFT *ItemNFTTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ItemNFT.Contract.TransferOwnership(&_ItemNFT.TransactOpts, newOwner)
}

// ItemNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ItemNFT contract.
type ItemNFTApprovalForAllIterator struct {
	Event *ItemNFTApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ItemNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ItemNFTApprovalForAll)
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
		it.Event = new(ItemNFTApprovalForAll)
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
func (it *ItemNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ItemNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ItemNFTApprovalForAll represents a ApprovalForAll event raised by the ItemNFT contract.
type ItemNFTApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_ItemNFT *ItemNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*ItemNFTApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ItemNFT.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ItemNFTApprovalForAllIterator{contract: _ItemNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_ItemNFT *ItemNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ItemNFTApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ItemNFT.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ItemNFTApprovalForAll)
				if err := _ItemNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_ItemNFT *ItemNFTFilterer) ParseApprovalForAll(log types.Log) (*ItemNFTApprovalForAll, error) {
	event := new(ItemNFTApprovalForAll)
	if err := _ItemNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ItemNFTItemCreatedIterator is returned from FilterItemCreated and is used to iterate over the raw logs and unpacked data for ItemCreated events raised by the ItemNFT contract.
type ItemNFTItemCreatedIterator struct {
	Event *ItemNFTItemCreated // Event containing the contract specifics and raw log

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
func (it *ItemNFTItemCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ItemNFTItemCreated)
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
		it.Event = new(ItemNFTItemCreated)
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
func (it *ItemNFTItemCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ItemNFTItemCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ItemNFTItemCreated represents a ItemCreated event raised by the ItemNFT contract.
type ItemNFTItemCreated struct {
	ItemId   *big.Int
	ItemType string
	Name     string
	Rarity   string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterItemCreated is a free log retrieval operation binding the contract event 0x95e04d1a7ad6b08119551a486a6e61d5be7969c038b84999326cffe707c717cd.
//
// Solidity: event ItemCreated(uint256 indexed itemId, string itemType, string name, string rarity)
func (_ItemNFT *ItemNFTFilterer) FilterItemCreated(opts *bind.FilterOpts, itemId []*big.Int) (*ItemNFTItemCreatedIterator, error) {

	var itemIdRule []interface{}
	for _, itemIdItem := range itemId {
		itemIdRule = append(itemIdRule, itemIdItem)
	}

	logs, sub, err := _ItemNFT.contract.FilterLogs(opts, "ItemCreated", itemIdRule)
	if err != nil {
		return nil, err
	}
	return &ItemNFTItemCreatedIterator{contract: _ItemNFT.contract, event: "ItemCreated", logs: logs, sub: sub}, nil
}

// WatchItemCreated is a free log subscription operation binding the contract event 0x95e04d1a7ad6b08119551a486a6e61d5be7969c038b84999326cffe707c717cd.
//
// Solidity: event ItemCreated(uint256 indexed itemId, string itemType, string name, string rarity)
func (_ItemNFT *ItemNFTFilterer) WatchItemCreated(opts *bind.WatchOpts, sink chan<- *ItemNFTItemCreated, itemId []*big.Int) (event.Subscription, error) {

	var itemIdRule []interface{}
	for _, itemIdItem := range itemId {
		itemIdRule = append(itemIdRule, itemIdItem)
	}

	logs, sub, err := _ItemNFT.contract.WatchLogs(opts, "ItemCreated", itemIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ItemNFTItemCreated)
				if err := _ItemNFT.contract.UnpackLog(event, "ItemCreated", log); err != nil {
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

// ParseItemCreated is a log parse operation binding the contract event 0x95e04d1a7ad6b08119551a486a6e61d5be7969c038b84999326cffe707c717cd.
//
// Solidity: event ItemCreated(uint256 indexed itemId, string itemType, string name, string rarity)
func (_ItemNFT *ItemNFTFilterer) ParseItemCreated(log types.Log) (*ItemNFTItemCreated, error) {
	event := new(ItemNFTItemCreated)
	if err := _ItemNFT.contract.UnpackLog(event, "ItemCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ItemNFTMinterAuthorizedIterator is returned from FilterMinterAuthorized and is used to iterate over the raw logs and unpacked data for MinterAuthorized events raised by the ItemNFT contract.
type ItemNFTMinterAuthorizedIterator struct {
	Event *ItemNFTMinterAuthorized // Event containing the contract specifics and raw log

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
func (it *ItemNFTMinterAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ItemNFTMinterAuthorized)
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
		it.Event = new(ItemNFTMinterAuthorized)
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
func (it *ItemNFTMinterAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ItemNFTMinterAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ItemNFTMinterAuthorized represents a MinterAuthorized event raised by the ItemNFT contract.
type ItemNFTMinterAuthorized struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterAuthorized is a free log retrieval operation binding the contract event 0x83b05b6735acd4b85e3bded8e72c851d1a87718f81e3c8e6f0c9d9a2baa88e46.
//
// Solidity: event MinterAuthorized(address indexed minter)
func (_ItemNFT *ItemNFTFilterer) FilterMinterAuthorized(opts *bind.FilterOpts, minter []common.Address) (*ItemNFTMinterAuthorizedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _ItemNFT.contract.FilterLogs(opts, "MinterAuthorized", minterRule)
	if err != nil {
		return nil, err
	}
	return &ItemNFTMinterAuthorizedIterator{contract: _ItemNFT.contract, event: "MinterAuthorized", logs: logs, sub: sub}, nil
}

// WatchMinterAuthorized is a free log subscription operation binding the contract event 0x83b05b6735acd4b85e3bded8e72c851d1a87718f81e3c8e6f0c9d9a2baa88e46.
//
// Solidity: event MinterAuthorized(address indexed minter)
func (_ItemNFT *ItemNFTFilterer) WatchMinterAuthorized(opts *bind.WatchOpts, sink chan<- *ItemNFTMinterAuthorized, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _ItemNFT.contract.WatchLogs(opts, "MinterAuthorized", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ItemNFTMinterAuthorized)
				if err := _ItemNFT.contract.UnpackLog(event, "MinterAuthorized", log); err != nil {
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
func (_ItemNFT *ItemNFTFilterer) ParseMinterAuthorized(log types.Log) (*ItemNFTMinterAuthorized, error) {
	event := new(ItemNFTMinterAuthorized)
	if err := _ItemNFT.contract.UnpackLog(event, "MinterAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ItemNFTMinterRevokedIterator is returned from FilterMinterRevoked and is used to iterate over the raw logs and unpacked data for MinterRevoked events raised by the ItemNFT contract.
type ItemNFTMinterRevokedIterator struct {
	Event *ItemNFTMinterRevoked // Event containing the contract specifics and raw log

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
func (it *ItemNFTMinterRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ItemNFTMinterRevoked)
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
		it.Event = new(ItemNFTMinterRevoked)
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
func (it *ItemNFTMinterRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ItemNFTMinterRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ItemNFTMinterRevoked represents a MinterRevoked event raised by the ItemNFT contract.
type ItemNFTMinterRevoked struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterRevoked is a free log retrieval operation binding the contract event 0x44f4322f8daa225d5f4877ad0f7d3dfba248a774396f3ca99405ed40a044fe81.
//
// Solidity: event MinterRevoked(address indexed minter)
func (_ItemNFT *ItemNFTFilterer) FilterMinterRevoked(opts *bind.FilterOpts, minter []common.Address) (*ItemNFTMinterRevokedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _ItemNFT.contract.FilterLogs(opts, "MinterRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return &ItemNFTMinterRevokedIterator{contract: _ItemNFT.contract, event: "MinterRevoked", logs: logs, sub: sub}, nil
}

// WatchMinterRevoked is a free log subscription operation binding the contract event 0x44f4322f8daa225d5f4877ad0f7d3dfba248a774396f3ca99405ed40a044fe81.
//
// Solidity: event MinterRevoked(address indexed minter)
func (_ItemNFT *ItemNFTFilterer) WatchMinterRevoked(opts *bind.WatchOpts, sink chan<- *ItemNFTMinterRevoked, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _ItemNFT.contract.WatchLogs(opts, "MinterRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ItemNFTMinterRevoked)
				if err := _ItemNFT.contract.UnpackLog(event, "MinterRevoked", log); err != nil {
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
func (_ItemNFT *ItemNFTFilterer) ParseMinterRevoked(log types.Log) (*ItemNFTMinterRevoked, error) {
	event := new(ItemNFTMinterRevoked)
	if err := _ItemNFT.contract.UnpackLog(event, "MinterRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ItemNFTOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ItemNFT contract.
type ItemNFTOwnershipTransferredIterator struct {
	Event *ItemNFTOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ItemNFTOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ItemNFTOwnershipTransferred)
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
		it.Event = new(ItemNFTOwnershipTransferred)
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
func (it *ItemNFTOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ItemNFTOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ItemNFTOwnershipTransferred represents a OwnershipTransferred event raised by the ItemNFT contract.
type ItemNFTOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ItemNFT *ItemNFTFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ItemNFTOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ItemNFT.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ItemNFTOwnershipTransferredIterator{contract: _ItemNFT.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ItemNFT *ItemNFTFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ItemNFTOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ItemNFT.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ItemNFTOwnershipTransferred)
				if err := _ItemNFT.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ItemNFT *ItemNFTFilterer) ParseOwnershipTransferred(log types.Log) (*ItemNFTOwnershipTransferred, error) {
	event := new(ItemNFTOwnershipTransferred)
	if err := _ItemNFT.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ItemNFTTransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the ItemNFT contract.
type ItemNFTTransferBatchIterator struct {
	Event *ItemNFTTransferBatch // Event containing the contract specifics and raw log

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
func (it *ItemNFTTransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ItemNFTTransferBatch)
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
		it.Event = new(ItemNFTTransferBatch)
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
func (it *ItemNFTTransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ItemNFTTransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ItemNFTTransferBatch represents a TransferBatch event raised by the ItemNFT contract.
type ItemNFTTransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ItemNFT *ItemNFTFilterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*ItemNFTTransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ItemNFT.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ItemNFTTransferBatchIterator{contract: _ItemNFT.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ItemNFT *ItemNFTFilterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *ItemNFTTransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ItemNFT.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ItemNFTTransferBatch)
				if err := _ItemNFT.contract.UnpackLog(event, "TransferBatch", log); err != nil {
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

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ItemNFT *ItemNFTFilterer) ParseTransferBatch(log types.Log) (*ItemNFTTransferBatch, error) {
	event := new(ItemNFTTransferBatch)
	if err := _ItemNFT.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ItemNFTTransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the ItemNFT contract.
type ItemNFTTransferSingleIterator struct {
	Event *ItemNFTTransferSingle // Event containing the contract specifics and raw log

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
func (it *ItemNFTTransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ItemNFTTransferSingle)
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
		it.Event = new(ItemNFTTransferSingle)
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
func (it *ItemNFTTransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ItemNFTTransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ItemNFTTransferSingle represents a TransferSingle event raised by the ItemNFT contract.
type ItemNFTTransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ItemNFT *ItemNFTFilterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*ItemNFTTransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ItemNFT.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ItemNFTTransferSingleIterator{contract: _ItemNFT.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ItemNFT *ItemNFTFilterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *ItemNFTTransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ItemNFT.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ItemNFTTransferSingle)
				if err := _ItemNFT.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ItemNFT *ItemNFTFilterer) ParseTransferSingle(log types.Log) (*ItemNFTTransferSingle, error) {
	event := new(ItemNFTTransferSingle)
	if err := _ItemNFT.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ItemNFTURIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the ItemNFT contract.
type ItemNFTURIIterator struct {
	Event *ItemNFTURI // Event containing the contract specifics and raw log

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
func (it *ItemNFTURIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ItemNFTURI)
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
		it.Event = new(ItemNFTURI)
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
func (it *ItemNFTURIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ItemNFTURIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ItemNFTURI represents a URI event raised by the ItemNFT contract.
type ItemNFTURI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ItemNFT *ItemNFTFilterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*ItemNFTURIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ItemNFT.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &ItemNFTURIIterator{contract: _ItemNFT.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ItemNFT *ItemNFTFilterer) WatchURI(opts *bind.WatchOpts, sink chan<- *ItemNFTURI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ItemNFT.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ItemNFTURI)
				if err := _ItemNFT.contract.UnpackLog(event, "URI", log); err != nil {
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

// ParseURI is a log parse operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ItemNFT *ItemNFTFilterer) ParseURI(log types.Log) (*ItemNFTURI, error) {
	event := new(ItemNFTURI)
	if err := _ItemNFT.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
