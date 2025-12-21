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

// CharacterNFTCharacterData is an auto generated low-level Go binding around an user-defined struct.
type CharacterNFTCharacterData struct {
	GameCharacterId *big.Int
	CharacterType   string
	Element         string
	Rarity          string
	Level           *big.Int
	MintedAt        *big.Int
}

// CharacterNFTMetaData contains all meta data concerning the CharacterNFT contract.
var CharacterNFTMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"baseURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizeMinter\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedMinters\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"characters\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"gameCharacterId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"characterType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"element\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"rarity\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"level\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"mintedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCharacterData\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCharacterNFT.CharacterData\",\"components\":[{\"name\":\"gameCharacterId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"characterType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"element\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"rarity\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"level\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"mintedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mintCharacter\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gameCharacterId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"characterType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"element\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"rarity\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"level\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeMinter\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"baseURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenByIndex\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenOfOwnerByIndex\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokensOfOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BatchMetadataUpdate\",\"inputs\":[{\"name\":\"_fromTokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_toTokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CharacterMinted\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gameCharacterId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"characterType\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"rarity\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdate\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinterAuthorized\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinterRevoked\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721EnumerableForbiddenBatchMint\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721OutOfBoundsIndex\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// CharacterNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use CharacterNFTMetaData.ABI instead.
var CharacterNFTABI = CharacterNFTMetaData.ABI

// CharacterNFT is an auto generated Go binding around an Ethereum contract.
type CharacterNFT struct {
	CharacterNFTCaller     // Read-only binding to the contract
	CharacterNFTTransactor // Write-only binding to the contract
	CharacterNFTFilterer   // Log filterer for contract events
}

// CharacterNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type CharacterNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CharacterNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CharacterNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CharacterNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CharacterNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CharacterNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CharacterNFTSession struct {
	Contract     *CharacterNFT     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CharacterNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CharacterNFTCallerSession struct {
	Contract *CharacterNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// CharacterNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CharacterNFTTransactorSession struct {
	Contract     *CharacterNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CharacterNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type CharacterNFTRaw struct {
	Contract *CharacterNFT // Generic contract binding to access the raw methods on
}

// CharacterNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CharacterNFTCallerRaw struct {
	Contract *CharacterNFTCaller // Generic read-only contract binding to access the raw methods on
}

// CharacterNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CharacterNFTTransactorRaw struct {
	Contract *CharacterNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCharacterNFT creates a new instance of CharacterNFT, bound to a specific deployed contract.
func NewCharacterNFT(address common.Address, backend bind.ContractBackend) (*CharacterNFT, error) {
	contract, err := bindCharacterNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CharacterNFT{CharacterNFTCaller: CharacterNFTCaller{contract: contract}, CharacterNFTTransactor: CharacterNFTTransactor{contract: contract}, CharacterNFTFilterer: CharacterNFTFilterer{contract: contract}}, nil
}

// NewCharacterNFTCaller creates a new read-only instance of CharacterNFT, bound to a specific deployed contract.
func NewCharacterNFTCaller(address common.Address, caller bind.ContractCaller) (*CharacterNFTCaller, error) {
	contract, err := bindCharacterNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTCaller{contract: contract}, nil
}

// NewCharacterNFTTransactor creates a new write-only instance of CharacterNFT, bound to a specific deployed contract.
func NewCharacterNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*CharacterNFTTransactor, error) {
	contract, err := bindCharacterNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTTransactor{contract: contract}, nil
}

// NewCharacterNFTFilterer creates a new log filterer instance of CharacterNFT, bound to a specific deployed contract.
func NewCharacterNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*CharacterNFTFilterer, error) {
	contract, err := bindCharacterNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTFilterer{contract: contract}, nil
}

// bindCharacterNFT binds a generic wrapper to an already deployed contract.
func bindCharacterNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CharacterNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CharacterNFT *CharacterNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CharacterNFT.Contract.CharacterNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CharacterNFT *CharacterNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CharacterNFT.Contract.CharacterNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CharacterNFT *CharacterNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CharacterNFT.Contract.CharacterNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CharacterNFT *CharacterNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CharacterNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CharacterNFT *CharacterNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CharacterNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CharacterNFT *CharacterNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CharacterNFT.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_CharacterNFT *CharacterNFTCaller) AuthorizedMinters(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "authorizedMinters", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_CharacterNFT *CharacterNFTSession) AuthorizedMinters(arg0 common.Address) (bool, error) {
	return _CharacterNFT.Contract.AuthorizedMinters(&_CharacterNFT.CallOpts, arg0)
}

// AuthorizedMinters is a free data retrieval call binding the contract method 0xaa2fe91b.
//
// Solidity: function authorizedMinters(address ) view returns(bool)
func (_CharacterNFT *CharacterNFTCallerSession) AuthorizedMinters(arg0 common.Address) (bool, error) {
	return _CharacterNFT.Contract.AuthorizedMinters(&_CharacterNFT.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CharacterNFT *CharacterNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CharacterNFT *CharacterNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _CharacterNFT.Contract.BalanceOf(&_CharacterNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CharacterNFT *CharacterNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _CharacterNFT.Contract.BalanceOf(&_CharacterNFT.CallOpts, owner)
}

// Characters is a free data retrieval call binding the contract method 0x4810bc59.
//
// Solidity: function characters(uint256 ) view returns(uint256 gameCharacterId, string characterType, string element, string rarity, uint256 level, uint256 mintedAt)
func (_CharacterNFT *CharacterNFTCaller) Characters(opts *bind.CallOpts, arg0 *big.Int) (struct {
	GameCharacterId *big.Int
	CharacterType   string
	Element         string
	Rarity          string
	Level           *big.Int
	MintedAt        *big.Int
}, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "characters", arg0)

	outstruct := new(struct {
		GameCharacterId *big.Int
		CharacterType   string
		Element         string
		Rarity          string
		Level           *big.Int
		MintedAt        *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameCharacterId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.CharacterType = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Element = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Rarity = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.Level = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.MintedAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Characters is a free data retrieval call binding the contract method 0x4810bc59.
//
// Solidity: function characters(uint256 ) view returns(uint256 gameCharacterId, string characterType, string element, string rarity, uint256 level, uint256 mintedAt)
func (_CharacterNFT *CharacterNFTSession) Characters(arg0 *big.Int) (struct {
	GameCharacterId *big.Int
	CharacterType   string
	Element         string
	Rarity          string
	Level           *big.Int
	MintedAt        *big.Int
}, error) {
	return _CharacterNFT.Contract.Characters(&_CharacterNFT.CallOpts, arg0)
}

// Characters is a free data retrieval call binding the contract method 0x4810bc59.
//
// Solidity: function characters(uint256 ) view returns(uint256 gameCharacterId, string characterType, string element, string rarity, uint256 level, uint256 mintedAt)
func (_CharacterNFT *CharacterNFTCallerSession) Characters(arg0 *big.Int) (struct {
	GameCharacterId *big.Int
	CharacterType   string
	Element         string
	Rarity          string
	Level           *big.Int
	MintedAt        *big.Int
}, error) {
	return _CharacterNFT.Contract.Characters(&_CharacterNFT.CallOpts, arg0)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CharacterNFT *CharacterNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CharacterNFT *CharacterNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _CharacterNFT.Contract.GetApproved(&_CharacterNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CharacterNFT *CharacterNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _CharacterNFT.Contract.GetApproved(&_CharacterNFT.CallOpts, tokenId)
}

// GetCharacterData is a free data retrieval call binding the contract method 0x14dede35.
//
// Solidity: function getCharacterData(uint256 tokenId) view returns((uint256,string,string,string,uint256,uint256))
func (_CharacterNFT *CharacterNFTCaller) GetCharacterData(opts *bind.CallOpts, tokenId *big.Int) (CharacterNFTCharacterData, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "getCharacterData", tokenId)

	if err != nil {
		return *new(CharacterNFTCharacterData), err
	}

	out0 := *abi.ConvertType(out[0], new(CharacterNFTCharacterData)).(*CharacterNFTCharacterData)

	return out0, err

}

// GetCharacterData is a free data retrieval call binding the contract method 0x14dede35.
//
// Solidity: function getCharacterData(uint256 tokenId) view returns((uint256,string,string,string,uint256,uint256))
func (_CharacterNFT *CharacterNFTSession) GetCharacterData(tokenId *big.Int) (CharacterNFTCharacterData, error) {
	return _CharacterNFT.Contract.GetCharacterData(&_CharacterNFT.CallOpts, tokenId)
}

// GetCharacterData is a free data retrieval call binding the contract method 0x14dede35.
//
// Solidity: function getCharacterData(uint256 tokenId) view returns((uint256,string,string,string,uint256,uint256))
func (_CharacterNFT *CharacterNFTCallerSession) GetCharacterData(tokenId *big.Int) (CharacterNFTCharacterData, error) {
	return _CharacterNFT.Contract.GetCharacterData(&_CharacterNFT.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CharacterNFT *CharacterNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CharacterNFT *CharacterNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _CharacterNFT.Contract.IsApprovedForAll(&_CharacterNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CharacterNFT *CharacterNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _CharacterNFT.Contract.IsApprovedForAll(&_CharacterNFT.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CharacterNFT *CharacterNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CharacterNFT *CharacterNFTSession) Name() (string, error) {
	return _CharacterNFT.Contract.Name(&_CharacterNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CharacterNFT *CharacterNFTCallerSession) Name() (string, error) {
	return _CharacterNFT.Contract.Name(&_CharacterNFT.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CharacterNFT *CharacterNFTCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CharacterNFT *CharacterNFTSession) Owner() (common.Address, error) {
	return _CharacterNFT.Contract.Owner(&_CharacterNFT.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CharacterNFT *CharacterNFTCallerSession) Owner() (common.Address, error) {
	return _CharacterNFT.Contract.Owner(&_CharacterNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CharacterNFT *CharacterNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CharacterNFT *CharacterNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _CharacterNFT.Contract.OwnerOf(&_CharacterNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CharacterNFT *CharacterNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _CharacterNFT.Contract.OwnerOf(&_CharacterNFT.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CharacterNFT *CharacterNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CharacterNFT *CharacterNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CharacterNFT.Contract.SupportsInterface(&_CharacterNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CharacterNFT *CharacterNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CharacterNFT.Contract.SupportsInterface(&_CharacterNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CharacterNFT *CharacterNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CharacterNFT *CharacterNFTSession) Symbol() (string, error) {
	return _CharacterNFT.Contract.Symbol(&_CharacterNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CharacterNFT *CharacterNFTCallerSession) Symbol() (string, error) {
	return _CharacterNFT.Contract.Symbol(&_CharacterNFT.CallOpts)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_CharacterNFT *CharacterNFTCaller) TokenByIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "tokenByIndex", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_CharacterNFT *CharacterNFTSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _CharacterNFT.Contract.TokenByIndex(&_CharacterNFT.CallOpts, index)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_CharacterNFT *CharacterNFTCallerSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _CharacterNFT.Contract.TokenByIndex(&_CharacterNFT.CallOpts, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_CharacterNFT *CharacterNFTCaller) TokenOfOwnerByIndex(opts *bind.CallOpts, owner common.Address, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "tokenOfOwnerByIndex", owner, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_CharacterNFT *CharacterNFTSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _CharacterNFT.Contract.TokenOfOwnerByIndex(&_CharacterNFT.CallOpts, owner, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_CharacterNFT *CharacterNFTCallerSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _CharacterNFT.Contract.TokenOfOwnerByIndex(&_CharacterNFT.CallOpts, owner, index)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CharacterNFT *CharacterNFTCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CharacterNFT *CharacterNFTSession) TokenURI(tokenId *big.Int) (string, error) {
	return _CharacterNFT.Contract.TokenURI(&_CharacterNFT.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CharacterNFT *CharacterNFTCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _CharacterNFT.Contract.TokenURI(&_CharacterNFT.CallOpts, tokenId)
}

// TokensOfOwner is a free data retrieval call binding the contract method 0x8462151c.
//
// Solidity: function tokensOfOwner(address owner) view returns(uint256[])
func (_CharacterNFT *CharacterNFTCaller) TokensOfOwner(opts *bind.CallOpts, owner common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "tokensOfOwner", owner)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// TokensOfOwner is a free data retrieval call binding the contract method 0x8462151c.
//
// Solidity: function tokensOfOwner(address owner) view returns(uint256[])
func (_CharacterNFT *CharacterNFTSession) TokensOfOwner(owner common.Address) ([]*big.Int, error) {
	return _CharacterNFT.Contract.TokensOfOwner(&_CharacterNFT.CallOpts, owner)
}

// TokensOfOwner is a free data retrieval call binding the contract method 0x8462151c.
//
// Solidity: function tokensOfOwner(address owner) view returns(uint256[])
func (_CharacterNFT *CharacterNFTCallerSession) TokensOfOwner(owner common.Address) ([]*big.Int, error) {
	return _CharacterNFT.Contract.TokensOfOwner(&_CharacterNFT.CallOpts, owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CharacterNFT *CharacterNFTCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CharacterNFT.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CharacterNFT *CharacterNFTSession) TotalSupply() (*big.Int, error) {
	return _CharacterNFT.Contract.TotalSupply(&_CharacterNFT.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CharacterNFT *CharacterNFTCallerSession) TotalSupply() (*big.Int, error) {
	return _CharacterNFT.Contract.TotalSupply(&_CharacterNFT.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.Contract.Approve(&_CharacterNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.Contract.Approve(&_CharacterNFT.TransactOpts, to, tokenId)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_CharacterNFT *CharacterNFTTransactor) AuthorizeMinter(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "authorizeMinter", minter)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_CharacterNFT *CharacterNFTSession) AuthorizeMinter(minter common.Address) (*types.Transaction, error) {
	return _CharacterNFT.Contract.AuthorizeMinter(&_CharacterNFT.TransactOpts, minter)
}

// AuthorizeMinter is a paid mutator transaction binding the contract method 0x0c984832.
//
// Solidity: function authorizeMinter(address minter) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) AuthorizeMinter(minter common.Address) (*types.Transaction, error) {
	return _CharacterNFT.Contract.AuthorizeMinter(&_CharacterNFT.TransactOpts, minter)
}

// MintCharacter is a paid mutator transaction binding the contract method 0xaf58cf99.
//
// Solidity: function mintCharacter(address to, uint256 gameCharacterId, string characterType, string element, string rarity, uint256 level, string tokenURI) returns(uint256)
func (_CharacterNFT *CharacterNFTTransactor) MintCharacter(opts *bind.TransactOpts, to common.Address, gameCharacterId *big.Int, characterType string, element string, rarity string, level *big.Int, tokenURI string) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "mintCharacter", to, gameCharacterId, characterType, element, rarity, level, tokenURI)
}

// MintCharacter is a paid mutator transaction binding the contract method 0xaf58cf99.
//
// Solidity: function mintCharacter(address to, uint256 gameCharacterId, string characterType, string element, string rarity, uint256 level, string tokenURI) returns(uint256)
func (_CharacterNFT *CharacterNFTSession) MintCharacter(to common.Address, gameCharacterId *big.Int, characterType string, element string, rarity string, level *big.Int, tokenURI string) (*types.Transaction, error) {
	return _CharacterNFT.Contract.MintCharacter(&_CharacterNFT.TransactOpts, to, gameCharacterId, characterType, element, rarity, level, tokenURI)
}

// MintCharacter is a paid mutator transaction binding the contract method 0xaf58cf99.
//
// Solidity: function mintCharacter(address to, uint256 gameCharacterId, string characterType, string element, string rarity, uint256 level, string tokenURI) returns(uint256)
func (_CharacterNFT *CharacterNFTTransactorSession) MintCharacter(to common.Address, gameCharacterId *big.Int, characterType string, element string, rarity string, level *big.Int, tokenURI string) (*types.Transaction, error) {
	return _CharacterNFT.Contract.MintCharacter(&_CharacterNFT.TransactOpts, to, gameCharacterId, characterType, element, rarity, level, tokenURI)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CharacterNFT *CharacterNFTTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CharacterNFT *CharacterNFTSession) RenounceOwnership() (*types.Transaction, error) {
	return _CharacterNFT.Contract.RenounceOwnership(&_CharacterNFT.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CharacterNFT *CharacterNFTTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CharacterNFT.Contract.RenounceOwnership(&_CharacterNFT.TransactOpts)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_CharacterNFT *CharacterNFTTransactor) RevokeMinter(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "revokeMinter", minter)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_CharacterNFT *CharacterNFTSession) RevokeMinter(minter common.Address) (*types.Transaction, error) {
	return _CharacterNFT.Contract.RevokeMinter(&_CharacterNFT.TransactOpts, minter)
}

// RevokeMinter is a paid mutator transaction binding the contract method 0xcfbd4885.
//
// Solidity: function revokeMinter(address minter) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) RevokeMinter(minter common.Address) (*types.Transaction, error) {
	return _CharacterNFT.Contract.RevokeMinter(&_CharacterNFT.TransactOpts, minter)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.Contract.SafeTransferFrom(&_CharacterNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.Contract.SafeTransferFrom(&_CharacterNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_CharacterNFT *CharacterNFTTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_CharacterNFT *CharacterNFTSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _CharacterNFT.Contract.SafeTransferFrom0(&_CharacterNFT.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _CharacterNFT.Contract.SafeTransferFrom0(&_CharacterNFT.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CharacterNFT *CharacterNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CharacterNFT *CharacterNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _CharacterNFT.Contract.SetApprovalForAll(&_CharacterNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _CharacterNFT.Contract.SetApprovalForAll(&_CharacterNFT.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI) returns()
func (_CharacterNFT *CharacterNFTTransactor) SetBaseURI(opts *bind.TransactOpts, baseURI string) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "setBaseURI", baseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI) returns()
func (_CharacterNFT *CharacterNFTSession) SetBaseURI(baseURI string) (*types.Transaction, error) {
	return _CharacterNFT.Contract.SetBaseURI(&_CharacterNFT.TransactOpts, baseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) SetBaseURI(baseURI string) (*types.Transaction, error) {
	return _CharacterNFT.Contract.SetBaseURI(&_CharacterNFT.TransactOpts, baseURI)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.Contract.TransferFrom(&_CharacterNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CharacterNFT.Contract.TransferFrom(&_CharacterNFT.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CharacterNFT *CharacterNFTTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CharacterNFT.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CharacterNFT *CharacterNFTSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CharacterNFT.Contract.TransferOwnership(&_CharacterNFT.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CharacterNFT *CharacterNFTTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CharacterNFT.Contract.TransferOwnership(&_CharacterNFT.TransactOpts, newOwner)
}

// CharacterNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the CharacterNFT contract.
type CharacterNFTApprovalIterator struct {
	Event *CharacterNFTApproval // Event containing the contract specifics and raw log

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
func (it *CharacterNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTApproval)
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
		it.Event = new(CharacterNFTApproval)
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
func (it *CharacterNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTApproval represents a Approval event raised by the CharacterNFT contract.
type CharacterNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_CharacterNFT *CharacterNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*CharacterNFTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTApprovalIterator{contract: _CharacterNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_CharacterNFT *CharacterNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CharacterNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTApproval)
				if err := _CharacterNFT.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_CharacterNFT *CharacterNFTFilterer) ParseApproval(log types.Log) (*CharacterNFTApproval, error) {
	event := new(CharacterNFTApproval)
	if err := _CharacterNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharacterNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the CharacterNFT contract.
type CharacterNFTApprovalForAllIterator struct {
	Event *CharacterNFTApprovalForAll // Event containing the contract specifics and raw log

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
func (it *CharacterNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTApprovalForAll)
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
		it.Event = new(CharacterNFTApprovalForAll)
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
func (it *CharacterNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTApprovalForAll represents a ApprovalForAll event raised by the CharacterNFT contract.
type CharacterNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_CharacterNFT *CharacterNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*CharacterNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTApprovalForAllIterator{contract: _CharacterNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_CharacterNFT *CharacterNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *CharacterNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTApprovalForAll)
				if err := _CharacterNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_CharacterNFT *CharacterNFTFilterer) ParseApprovalForAll(log types.Log) (*CharacterNFTApprovalForAll, error) {
	event := new(CharacterNFTApprovalForAll)
	if err := _CharacterNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharacterNFTBatchMetadataUpdateIterator is returned from FilterBatchMetadataUpdate and is used to iterate over the raw logs and unpacked data for BatchMetadataUpdate events raised by the CharacterNFT contract.
type CharacterNFTBatchMetadataUpdateIterator struct {
	Event *CharacterNFTBatchMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *CharacterNFTBatchMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTBatchMetadataUpdate)
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
		it.Event = new(CharacterNFTBatchMetadataUpdate)
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
func (it *CharacterNFTBatchMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTBatchMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTBatchMetadataUpdate represents a BatchMetadataUpdate event raised by the CharacterNFT contract.
type CharacterNFTBatchMetadataUpdate struct {
	FromTokenId *big.Int
	ToTokenId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchMetadataUpdate is a free log retrieval operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_CharacterNFT *CharacterNFTFilterer) FilterBatchMetadataUpdate(opts *bind.FilterOpts) (*CharacterNFTBatchMetadataUpdateIterator, error) {

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &CharacterNFTBatchMetadataUpdateIterator{contract: _CharacterNFT.contract, event: "BatchMetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchBatchMetadataUpdate is a free log subscription operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_CharacterNFT *CharacterNFTFilterer) WatchBatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *CharacterNFTBatchMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTBatchMetadataUpdate)
				if err := _CharacterNFT.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
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

// ParseBatchMetadataUpdate is a log parse operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_CharacterNFT *CharacterNFTFilterer) ParseBatchMetadataUpdate(log types.Log) (*CharacterNFTBatchMetadataUpdate, error) {
	event := new(CharacterNFTBatchMetadataUpdate)
	if err := _CharacterNFT.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharacterNFTCharacterMintedIterator is returned from FilterCharacterMinted and is used to iterate over the raw logs and unpacked data for CharacterMinted events raised by the CharacterNFT contract.
type CharacterNFTCharacterMintedIterator struct {
	Event *CharacterNFTCharacterMinted // Event containing the contract specifics and raw log

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
func (it *CharacterNFTCharacterMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTCharacterMinted)
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
		it.Event = new(CharacterNFTCharacterMinted)
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
func (it *CharacterNFTCharacterMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTCharacterMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTCharacterMinted represents a CharacterMinted event raised by the CharacterNFT contract.
type CharacterNFTCharacterMinted struct {
	TokenId         *big.Int
	Owner           common.Address
	GameCharacterId *big.Int
	CharacterType   string
	Rarity          string
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCharacterMinted is a free log retrieval operation binding the contract event 0xdf275d816af024e62f941936bea92bfb70589d4df53af1f62cf06f85b47caca9.
//
// Solidity: event CharacterMinted(uint256 indexed tokenId, address indexed owner, uint256 gameCharacterId, string characterType, string rarity)
func (_CharacterNFT *CharacterNFTFilterer) FilterCharacterMinted(opts *bind.FilterOpts, tokenId []*big.Int, owner []common.Address) (*CharacterNFTCharacterMintedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "CharacterMinted", tokenIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTCharacterMintedIterator{contract: _CharacterNFT.contract, event: "CharacterMinted", logs: logs, sub: sub}, nil
}

// WatchCharacterMinted is a free log subscription operation binding the contract event 0xdf275d816af024e62f941936bea92bfb70589d4df53af1f62cf06f85b47caca9.
//
// Solidity: event CharacterMinted(uint256 indexed tokenId, address indexed owner, uint256 gameCharacterId, string characterType, string rarity)
func (_CharacterNFT *CharacterNFTFilterer) WatchCharacterMinted(opts *bind.WatchOpts, sink chan<- *CharacterNFTCharacterMinted, tokenId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "CharacterMinted", tokenIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTCharacterMinted)
				if err := _CharacterNFT.contract.UnpackLog(event, "CharacterMinted", log); err != nil {
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

// ParseCharacterMinted is a log parse operation binding the contract event 0xdf275d816af024e62f941936bea92bfb70589d4df53af1f62cf06f85b47caca9.
//
// Solidity: event CharacterMinted(uint256 indexed tokenId, address indexed owner, uint256 gameCharacterId, string characterType, string rarity)
func (_CharacterNFT *CharacterNFTFilterer) ParseCharacterMinted(log types.Log) (*CharacterNFTCharacterMinted, error) {
	event := new(CharacterNFTCharacterMinted)
	if err := _CharacterNFT.contract.UnpackLog(event, "CharacterMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharacterNFTMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the CharacterNFT contract.
type CharacterNFTMetadataUpdateIterator struct {
	Event *CharacterNFTMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *CharacterNFTMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTMetadataUpdate)
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
		it.Event = new(CharacterNFTMetadataUpdate)
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
func (it *CharacterNFTMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTMetadataUpdate represents a MetadataUpdate event raised by the CharacterNFT contract.
type CharacterNFTMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_CharacterNFT *CharacterNFTFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*CharacterNFTMetadataUpdateIterator, error) {

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &CharacterNFTMetadataUpdateIterator{contract: _CharacterNFT.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_CharacterNFT *CharacterNFTFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *CharacterNFTMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTMetadataUpdate)
				if err := _CharacterNFT.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_CharacterNFT *CharacterNFTFilterer) ParseMetadataUpdate(log types.Log) (*CharacterNFTMetadataUpdate, error) {
	event := new(CharacterNFTMetadataUpdate)
	if err := _CharacterNFT.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharacterNFTMinterAuthorizedIterator is returned from FilterMinterAuthorized and is used to iterate over the raw logs and unpacked data for MinterAuthorized events raised by the CharacterNFT contract.
type CharacterNFTMinterAuthorizedIterator struct {
	Event *CharacterNFTMinterAuthorized // Event containing the contract specifics and raw log

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
func (it *CharacterNFTMinterAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTMinterAuthorized)
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
		it.Event = new(CharacterNFTMinterAuthorized)
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
func (it *CharacterNFTMinterAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTMinterAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTMinterAuthorized represents a MinterAuthorized event raised by the CharacterNFT contract.
type CharacterNFTMinterAuthorized struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterAuthorized is a free log retrieval operation binding the contract event 0x83b05b6735acd4b85e3bded8e72c851d1a87718f81e3c8e6f0c9d9a2baa88e46.
//
// Solidity: event MinterAuthorized(address indexed minter)
func (_CharacterNFT *CharacterNFTFilterer) FilterMinterAuthorized(opts *bind.FilterOpts, minter []common.Address) (*CharacterNFTMinterAuthorizedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "MinterAuthorized", minterRule)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTMinterAuthorizedIterator{contract: _CharacterNFT.contract, event: "MinterAuthorized", logs: logs, sub: sub}, nil
}

// WatchMinterAuthorized is a free log subscription operation binding the contract event 0x83b05b6735acd4b85e3bded8e72c851d1a87718f81e3c8e6f0c9d9a2baa88e46.
//
// Solidity: event MinterAuthorized(address indexed minter)
func (_CharacterNFT *CharacterNFTFilterer) WatchMinterAuthorized(opts *bind.WatchOpts, sink chan<- *CharacterNFTMinterAuthorized, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "MinterAuthorized", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTMinterAuthorized)
				if err := _CharacterNFT.contract.UnpackLog(event, "MinterAuthorized", log); err != nil {
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
func (_CharacterNFT *CharacterNFTFilterer) ParseMinterAuthorized(log types.Log) (*CharacterNFTMinterAuthorized, error) {
	event := new(CharacterNFTMinterAuthorized)
	if err := _CharacterNFT.contract.UnpackLog(event, "MinterAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharacterNFTMinterRevokedIterator is returned from FilterMinterRevoked and is used to iterate over the raw logs and unpacked data for MinterRevoked events raised by the CharacterNFT contract.
type CharacterNFTMinterRevokedIterator struct {
	Event *CharacterNFTMinterRevoked // Event containing the contract specifics and raw log

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
func (it *CharacterNFTMinterRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTMinterRevoked)
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
		it.Event = new(CharacterNFTMinterRevoked)
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
func (it *CharacterNFTMinterRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTMinterRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTMinterRevoked represents a MinterRevoked event raised by the CharacterNFT contract.
type CharacterNFTMinterRevoked struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterRevoked is a free log retrieval operation binding the contract event 0x44f4322f8daa225d5f4877ad0f7d3dfba248a774396f3ca99405ed40a044fe81.
//
// Solidity: event MinterRevoked(address indexed minter)
func (_CharacterNFT *CharacterNFTFilterer) FilterMinterRevoked(opts *bind.FilterOpts, minter []common.Address) (*CharacterNFTMinterRevokedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "MinterRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTMinterRevokedIterator{contract: _CharacterNFT.contract, event: "MinterRevoked", logs: logs, sub: sub}, nil
}

// WatchMinterRevoked is a free log subscription operation binding the contract event 0x44f4322f8daa225d5f4877ad0f7d3dfba248a774396f3ca99405ed40a044fe81.
//
// Solidity: event MinterRevoked(address indexed minter)
func (_CharacterNFT *CharacterNFTFilterer) WatchMinterRevoked(opts *bind.WatchOpts, sink chan<- *CharacterNFTMinterRevoked, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "MinterRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTMinterRevoked)
				if err := _CharacterNFT.contract.UnpackLog(event, "MinterRevoked", log); err != nil {
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
func (_CharacterNFT *CharacterNFTFilterer) ParseMinterRevoked(log types.Log) (*CharacterNFTMinterRevoked, error) {
	event := new(CharacterNFTMinterRevoked)
	if err := _CharacterNFT.contract.UnpackLog(event, "MinterRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharacterNFTOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CharacterNFT contract.
type CharacterNFTOwnershipTransferredIterator struct {
	Event *CharacterNFTOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CharacterNFTOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTOwnershipTransferred)
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
		it.Event = new(CharacterNFTOwnershipTransferred)
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
func (it *CharacterNFTOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTOwnershipTransferred represents a OwnershipTransferred event raised by the CharacterNFT contract.
type CharacterNFTOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CharacterNFT *CharacterNFTFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CharacterNFTOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTOwnershipTransferredIterator{contract: _CharacterNFT.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CharacterNFT *CharacterNFTFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CharacterNFTOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTOwnershipTransferred)
				if err := _CharacterNFT.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_CharacterNFT *CharacterNFTFilterer) ParseOwnershipTransferred(log types.Log) (*CharacterNFTOwnershipTransferred, error) {
	event := new(CharacterNFTOwnershipTransferred)
	if err := _CharacterNFT.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharacterNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the CharacterNFT contract.
type CharacterNFTTransferIterator struct {
	Event *CharacterNFTTransfer // Event containing the contract specifics and raw log

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
func (it *CharacterNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharacterNFTTransfer)
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
		it.Event = new(CharacterNFTTransfer)
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
func (it *CharacterNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharacterNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharacterNFTTransfer represents a Transfer event raised by the CharacterNFT contract.
type CharacterNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_CharacterNFT *CharacterNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*CharacterNFTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _CharacterNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &CharacterNFTTransferIterator{contract: _CharacterNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_CharacterNFT *CharacterNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CharacterNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _CharacterNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharacterNFTTransfer)
				if err := _CharacterNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_CharacterNFT *CharacterNFTFilterer) ParseTransfer(log types.Log) (*CharacterNFTTransfer, error) {
	event := new(CharacterNFTTransfer)
	if err := _CharacterNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
