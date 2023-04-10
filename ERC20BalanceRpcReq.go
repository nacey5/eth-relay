package eth_relay

type ERC20BalanceRpcReq struct {
	ContractAddress string //the eth contract address
	UserAddress     string
	ContractDecimal int //the aim of decimal
}
