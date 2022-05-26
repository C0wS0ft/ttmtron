package ttmtron

type Block struct {
	BlockId     string `json:"blockID"`
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
