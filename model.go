package ttmtron

const (
	TRC10 tokenTypes = "TRC10"
	TRC20 tokenTypes = "TRC20"
)

type tokenTypes = string

type Block struct {
	BlockID     string `json:"blockID"`
	Txs         []Tx   `json:"transactions"`
	BlockHeader struct {
		Data BlockData `json:"raw_data"`
	} `json:"block_header"`
}

type BlockData struct {
	Number    int64 `json:"number"`
	Timestamp int64 `json:"timestamp"`
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
	Parameter struct {
		Value TransferValue `json:"value"`
	} `json:"parameter"`
}

type TransferValue struct {
	Amount       int    `json:"amount"`
	OwnerAddress string `json:"owner_address"`
	ToAddress    string `json:"to_address"`
	AssetName    string `json:"asset_name,omitempty"`
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
