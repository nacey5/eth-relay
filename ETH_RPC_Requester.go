package eth_relay

import "eth-relay/model"

type ETHRPCRequester struct {
	client *ETHRPCClient
}

func NewETHRPCRequester(nodeUrl string) *ETHRPCRequester {
	requester := &ETHRPCRequester{}
	requester.client = NewETHRPCClient(nodeUrl)
	return requester
}

func (r *ETHRPCRequester) GetTransactionByHash(txHash string) (model.Transaction, error) {
	methodName := "eth_getTransactionByHash"
	result := model.Transaction{}
	// 下面call函数的result参数传入的是model.Transaction结构体的引用
	// 这样内部所设置的值在函数执行完之后才能有效果
	err := r.client.GetRpc().Call(&result, methodName, txHash)
	return result, err
}
