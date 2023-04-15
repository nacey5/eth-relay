package dao

type Transaction struct {
	Id               int64  `json:"id"`
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"block_hash"`
	BlockNumber      string `json:"block_number"`
	TransactionIndex string `json:"transaction_index"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasPrice         string `json:"gas_price"`
	Gas              string `json:"gas"`
	Input            string `xorm:"text" json:"input"` //data
}
