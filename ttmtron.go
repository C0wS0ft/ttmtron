package ttmtron

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/trustwallet/go-libs/client"
)

// Init TronRequest lib
func Init(baseURL string) *TronRequest {
	return &TronRequest{client.Request{
		BaseURL:          baseURL,
		Headers:          make(map[string]string),
		HttpClient:       DefaultClient,
		HttpErrorHandler: DefaultErrorHandler,
	}}
}

// CurrentBlockNumber return current block number
func (t *TronRequest) CurrentBlockNumber(ctx context.Context) (uint64, error) {
	var block Block
	err := t.Post(&block, "wallet/getnowblock", nil)

	return block.BlockHeader.Data.Number, errors.Wrap(err, "unable to post wallet/getnowblock")
}

// GetBlockByNumber return block by number
func (t *TronRequest) GetBlockByNumber(ctx context.Context, num uint64) (*Block, error) {
	var block Block
	err := t.Post(&block, "wallet/getblockbynum", GetBlockByNumRequest{Num: num})

	return &block, errors.Wrap(err, "unable to get block by number")
}

func (t *TronRequest) GetBlockByLimitNext(ctx context.Context, startnum, endnum uint64) (*Blocks, error) {
	var blocks *Blocks
	err := t.Post(&blocks, "wallet/getblockbylimitnext", GetBlockByLimitNextRequest{StartNum: startnum, EndNum: endnum})

	return blocks, errors.Wrap(err, "unable to get block by limit")
}

// GetAccountBalance due to unknown reason Tron accepts only HEX address, not base58
func (t *TronRequest) GetAccountBalance(ctx context.Context, address string, asset string) (uint64, error) {
	address = Base58ToHex(address)

	var account AccountReply

	err := t.Post(&account, "wallet/getaccount", AccountRequest{Address: address, Visible: false})

	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}

// triggerConstantContract Call smart contract function and return constant_result field
func (t *TronRequest) triggerConstantContract(ctx context.Context, ownerAddress string, smartContractAddress string, function string, params []string) ([]string, error) { //nolint:lll
	var result TriggerConstantContractReply
	var parameters string

	for _, v := range params {
		parameters += v
	}

	err := t.Post(&result, "wallet/triggerconstantcontract",
		TriggerConstantContractRequest{
			OwnerAddress:     ownerAddress,
			ContractAddress:  smartContractAddress,
			FunctionSelector: function,
			Parameter:        parameters,
			Visible:          true, // base58
		},
	)

	if err != nil {
		return nil, err
	}

	return result.ConstantResult, nil
}

func (t *TronRequest) triggerSmartContractTransfer(ctx context.Context, ownerAddress string,
	smartContractAddress string, callValue uint64, parameter string, feeLimit uint64) (*CreateTransactionReply, error) {
	var result TriggerSmartContractTransferReply

	contract := TriggerSmartContractRequest{
		OwnerAddress:     Base58ToHex(ownerAddress),
		ContractAddress:  Base58ToHex(smartContractAddress),
		FunctionSelector: "transfer (address, uint256)",
		CallValue:        callValue,
		Parameter:        parameter,
		FeeLimit:         feeLimit,
		CallTokenValue:   0,
		TokenID:          0,
	}

	err := t.Post(&result, "wallet/triggersmartcontract", contract)

	if err != nil {
		return nil, err
	}

	if !result.Result.Result {
		qq, _ := hex.DecodeString(result.Result.Message)

		return nil, errors.New(string(qq))
	}

	return &result.Transaction, nil
}

// GetTRC20TokenSymbol returns TRC20 smart contract token symbol
func (t *TronRequest) GetTRC20TokenSymbol(ctx context.Context, ownerAddress string, smartContractAddress string) (string, error) { //nolint:lll
	res, err := t.triggerConstantContract(ctx, ownerAddress, smartContractAddress, "symbol()", []string{})

	if err != nil {
		return "", err
	}

	if len(res) < 1 {
		return "", errors.New("unable to get token symbol")
	}

	symbol, err := DecodeConstantToSymbol(res[0])

	if err != nil {
		return "", err
	}

	return symbol, nil
}

// GetTRC20TokenDecimals returns TRC20 smart contract token decimals
func (t *TronRequest) GetTRC20TokenDecimals(ctx context.Context, ownerAddress string, smartContractAddress string) (uint64, error) { //nolint:lll
	res, err := t.triggerConstantContract(ctx, ownerAddress, smartContractAddress, "decimals()", []string{})

	if err != nil {
		return 0, err
	}

	if len(res) < 1 {
		return 0, errors.New("unable to get token decimals")
	}

	decimals, err := HexToInt256(res[0])

	if err != nil {
		return 0, err
	}

	return decimals.Uint64(), nil
}

// GetTRC20TokenBalance returns TRC20 smart contract token balance
func (t *TronRequest) GetTRC20TokenBalance(ctx context.Context, ownerAddress string, smartContractAddress string) (uint64, error) { //nolint:lll
	address, err := EncodeAddressToParameter(ownerAddress)

	if err != nil {
		return 0, err
	}

	res, err := t.triggerConstantContract(ctx, ownerAddress, smartContractAddress, "balanceOf(address)", []string{address})

	if err != nil {
		return 0, err
	}

	if len(res) < 1 {
		return 0, errors.New("unable to get balance")
	}

	balance, err := HexToInt256(res[0])

	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}

func (t *TronRequest) GetTRC20SmartContract(ctx context.Context, smartContractAddress string) (*GetContractReply, error) {
	var result GetContractReply

	err := t.Post(&result, "wallet/getcontract",
		GetContractRequest{
			Value:   smartContractAddress,
			Visible: true,
		},
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TronRequest) GetTRC10TokenInfoByID(ctx context.Context, tokenID uint64) (*GetAssetIssueByIDReply, error) {
	var result GetAssetIssueByIDReply

	err := t.Post(&result, "wallet/getassetissuebyid", GetAssetIssueByIDRequest{Value: tokenID})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TronRequest) TransferTRX(ctx context.Context, from, to string, amount uint64) (*CreateTransactionReply, error) {
	var result CreateTransactionReply

	value := TransferValue{
		Amount:       amount,
		OwnerAddress: Base58ToHex(from),
		ToAddress:    Base58ToHex(to),
	}

	err := t.Post(&result, "wallet/createtransaction", value)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TronRequest) TransferTRC10Token(ctx context.Context, from, to string, tokenID uint64, amount uint64) (*CreateTransactionReply, error) { //nolint:lll
	var result CreateTransactionReply

	tokenIDStr := fmt.Sprintf("%d", tokenID)

	value := TransferValue{
		Amount:       amount,
		OwnerAddress: Base58ToHex(from),
		ToAddress:    Base58ToHex(to),
		AssetName:    hex.EncodeToString([]byte(tokenIDStr)),
	}

	err := t.Post(&result, "wallet/transferasset", value)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TronRequest) TransferTRC20Token(ctx context.Context, from, to, smartContractAddress string, amount uint64, feeLimit uint64) (*CreateTransactionReply, error) { //nolint:lll
	parameter := AddParameterAddress(to)
	parameter += AddParameterAmount(amount * 1000000) // check this

	feeLimit *= 1000000

	result, err := t.triggerSmartContractTransfer(ctx, from, smartContractAddress, 0, parameter, feeLimit)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *TronRequest) BroadcastSignedTransaction(rawdata, rawdataHex string) (*BroadcastSignedTransactionReply, error) {
	var result BroadcastSignedTransactionReply

	err := t.Post(&result, "wallet/broadcasttransaction",
		BroadcastSignedTransactionRequest{
			RawData:    rawdata,
			RawDataHex: rawdataHex,
		},
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
