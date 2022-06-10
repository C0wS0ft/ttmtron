package ttmtron

import (
	"context"

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
	// log := logger.FromContext(ctx).WithField("m", "Client::CurrentBlockNumber")
	// log.Debugf("Client::CurrentBlockNumber:: ")

	var block Block
	err := t.Post(&block, "wallet/getnowblock", nil)

	return block.BlockHeader.Data.Number, errors.Wrap(err, "unable to post wallet/getnowblock")
}

// GetBlockByNumber return block by number
func (t *TronRequest) GetBlockByNumber(ctx context.Context, num uint64) (*Block, error) {
	// log := logger.FromContext(ctx).WithField("m", "GetBlockByNumber")
	// log.Debugf("GetBlockByNumber:: num: %v", num)

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
	// log := logger.FromContext(ctx).WithField("m", "Client::GetAccountBalance")
	// log.Debugf("Client::GetAccountBalance:: ")

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
	// log := logger.FromContext(ctx).WithField("m", "Client::triggerConstantContract")
	// log.Debugf("Client::triggerConstantContract:: ")

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

// GetTRC20TokenSymbol returns TRC20 smart contract token symbol
func (t *TronRequest) GetTRC20TokenSymbol(ctx context.Context, ownerAddress string, smartContractAddress string) (string, error) { //nolint:lll
	// log := logger.FromContext(ctx).WithField("m", "Client::GetTRC20TokenSymbol")
	// log.Debugf("Client::GetTRC20TokenSymbol:: ")

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
	// log := logger.FromContext(ctx).WithField("m", "Client::GetTRC20TokenDecimals")
	// log.Debugf("Client::GetTRC20TokenDecimals:: ")

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
	// log := logger.FromContext(ctx).WithField("m", "Client::GetTRC20TokenBalance")
	// log.Debugf("Client::GetTRC20TokenBalance:: ")

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
	// log := logger.FromContext(ctx).WithField("m", "Client::GetTRC20SmartContract")
	// log.Debugf("Client::GetTRC20SmartContract:: ")

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
