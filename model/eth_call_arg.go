package model

import "github.com/ethereum/go-ethereum/common"

// CallArg the struct for the eth_call
type CallArg struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      string         `json:"gas"`
	GasPrice string         `json:"gas_price`
	Value    string         `json:"value"`
	Data     string         `json:"data"`
	Nonce    string         `json:"nonce"`
}
