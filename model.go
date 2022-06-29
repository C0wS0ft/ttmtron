package ttmtron

const (
	TRC10       TokenTypes = "TRC10"
	TRC20       TokenTypes = "TRC20"
	SigTransfer            = "a9059cbb"
	SigSwap                = "a8f85ca6"
)

type TokenTypes = string

type Block struct {
	BlockID     string `json:"blockID"`
	Txs         []Tx   `json:"transactions"`
	BlockHeader struct {
		Data BlockData `json:"raw_data"`
	} `json:"block_header"`
}

type Blocks struct {
	Blocks []Block `json:"block"`
}

type BlockData struct {
	Number    uint64 `json:"number"`
	Timestamp int64  `json:"timestamp"`
}

type Tx struct {
	ID        string `json:"txID"`
	BlockTime int64  `json:"block_timestamp"`
	Data      TxData `json:"raw_data"`
}

type TxData struct {
	Timestamp int64      `json:"timestamp"`
	Contracts []Contract `json:"contract"`
}

type ContractType = string

type Contract struct {
	Type      ContractType `json:"type"`
	Parameter Parameter    `json:"parameter"`
}

type Parameter struct {
	Value   TransferValue `json:"value"`
	TypeURL string        `json:"type_url"`
}

type TransferValue struct {
	Data            string `json:"data,omitempty"`             // TRC20
	ContractAddress string `json:"contract_address,omitempty"` // TRC20
	CallValue       uint64 `json:"call_value,omitempty"`       // TRC20
	Amount          uint64 `json:"amount,omitempty"`           // TRC10 & TRX
	OwnerAddress    string `json:"owner_address"`              // TRC10 & TRC20 & TRX
	ToAddress       string `json:"to_address,omitempty"`       // TRC10 & TRX
	AssetName       string `json:"asset_name,omitempty"`       // TRC10 & TRX
}

type AccountRequest struct {
	Address string `json:"address"`
	Visible bool   `json:"visible"`
}

type AccountReply struct {
	AccountName string `json:"account_name"`
	Address     string `json:"address"`
	Balance     uint64 `json:"balance"`
}

type TriggerConstantContractRequest struct {
	OwnerAddress     string `json:"owner_address"`
	ContractAddress  string `json:"contract_address"`
	FunctionSelector string `json:"function_selector"`
	Parameter        string `json:"parameter"`
	Visible          bool   `json:"visible"`
}

type TriggerSmartContractRequest struct {
	OwnerAddress     string `json:"owner_address"`     // Address that triggers the contract, converted to a hex string.
	ContractAddress  string `json:"contract_address"`  // Contract address, converted to a hex string
	FunctionSelector string `json:"function_selector"` // Function call, must not be left blank
	CallValue        uint64 `json:"call_value"`        // Amount of TRX transferred with this transaction, measured in SUN (1 TRX = 1,000,000 SUN)
	Parameter        string `json:"parameter"`         // Parameter encoding needs to be in accordance with the ABI rule
	FeeLimit         uint64 `json:"fee_limit"`         // Maximum TRX consumption, measured in SUN (1 TRX = 1,000,000 SUN)
	//	Visible          bool   `json:"visible"`           // Optional. Whehter the address is in base58check format.
	CallTokenValue uint64 `json:"call_token_value"`
	TokenID        uint64 `json:"token_id"`
}

type TriggerConstantContractReply struct {
	EnergyUsed     uint64   `json:"energy_used"`
	ConstantResult []string `json:"constant_result"`
}

type GetContractRequest struct {
	Value   string `json:"value"`
	Visible bool   `json:"visible"`
}

type GetContractReply struct {
	Bytecode                    string `json:"bytecode"`
	ConsumeUsersResourcePercent uint64 `json:"consume_users_resource_percent"`
	Name                        string `json:"name"`
	OriginAddress               string `json:"origin_address"`
	OriginEnergyLimit           uint64 `json:"origin_energy_limit"`
	ContractAddress             string `json:"contract_address"`
	CodeHash                    string `json:"code_hash"`
}

type GetBlockByNumRequest struct {
	Num uint64 `json:"num"`
}

type GetBlockByLimitNextRequest struct {
	StartNum uint64 `json:"startNum"`
	EndNum   uint64 `json:"endNum"`
}

// TRC10

type GetAssetIssueByIDRequest struct {
	Value uint64 `json:"value"`
}

type GetAssetIssueByIDReply struct {
	OwnerAddress string `json:"owner_address"`
	Name         string `json:"name"`
	Abbr         string `json:"abbr"`
	TotalSupply  int64  `json:"total_supply"`
	TrxNum       int    `json:"trx_num"`
	Num          int    `json:"num"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	Description  string `json:"description"`
	URL          string `json:"url"`
	ID           string `json:"id"`
}

type TriggerSmartContractTransferReply struct {
	Result struct {
		Result  bool   `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Transaction CreateTransactionReply `json:"transaction"`
}

type CreateTransactionReply struct {
	Visible bool   `json:"visible"`
	TxID    string `json:"txID"`
	RawData struct {
		Contract      []Contract `json:"contract"`
		RefBlockBytes string     `json:"ref_block_bytes"`
		RefBlockHash  string     `json:"ref_block_hash"`
		Expiration    uint64     `json:"expiration"`
		Timestamp     uint64     `json:"timestamp"`
		FeeLimit      uint64     `json:"fee_limit"`
	} `json:"raw_data"`
	RawDataHex string `json:"raw_data_hex"`
}

type BroadcastSignedTransactionRequest struct {
	RawData    string `json:"raw_data"`
	RawDataHex string `json:"raw_data_hex"`
}

type BroadcastSignedTransactionReply struct{}
