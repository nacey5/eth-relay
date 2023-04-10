package eth_relay

import (
	"errors"
	"eth-relay/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

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

func (r *ETHRPCRequester) GetETHBalance(address string) (string, error) {
	name := "eth_getBalance"
	result := ""
	err := r.client.GetRpc().Call(&result, name, address, "latest")
	if err != nil {
		return "", nil
	}
	if result == "" {
		return "", errors.New("eth balance is null")
	}

	//the result format is 0x16
	//transfer to the ten
	//prevent form the bit overflow
	ten, _ := new(big.Int).SetString(result[2:], 16)
	return ten.String(), nil
}

// GetETHBalances pathQuery
func (r *ETHRPCRequester) GetETHBalances(addresses []string) ([]string, error) {
	name := "eth_getBalance"
	rets := []*string{}
	size := len(addresses)
	reqs := []rpc.BatchElem{}
	for i := 0; i < size; i++ {
		ret := ""
		//instantiate every elem
		req := rpc.BatchElem{
			Method: name,
			Args:   []interface{}{addresses[i], "latest"},
			Result: &ret,
		}
		reqs = append(reqs, req)
		rets = append(rets, &ret)
	}
	err := r.client.GetRpc().BatchCall(reqs)
	if err != nil {
		return nil, err
	}
	for _, req := range reqs {
		if req.Error != nil {
			return nil, req.Error
		}
	}
	finalRet := []string{}
	for _, item := range rets {
		ten, _ := new(big.Int).SetString((*item)[2:], 16)
		finalRet = append(finalRet, ten.String())
	}
	return finalRet, err
}

func (r *ETHRPCRequester) GetERC20Balances(paramArr []ERC20BalanceRpcReq) ([]string, error) {
	name := "eth_call"
	methodId := "0x70a08231" //the balanceOf methodId
	rets := []*string{}
	size := len(paramArr)
	reqs := []rpc.BatchElem{}
	for i := 0; i < size; i++ {
		ret := ""
		arg := &model.CallArg{}
		userAddress := paramArr[i].UserAddress
		//query args,the query not need the gas,dont set the gas fee
		arg.To = common.HexToAddress(paramArr[i].ContractAddress)
		arg.Data = methodId + "000000000000000000000000" + userAddress[2:]
		//instance every ele
		req := rpc.BatchElem{
			Method: name,
			Args:   []interface{}{arg, "latest"},
			Result: &ret,
		}
		reqs = append(reqs, req)
		rets = append(rets, &ret)
	}
	err := r.client.GetRpc().BatchCall(reqs)
	if err != nil {
		return nil, err
	}
	for _, req := range reqs {
		if req.Error != nil {
			return nil, req.Error
		}
	}
	finalRet := []string{}
	for _, item := range rets {
		if *item == "" {
			continue
		}
		ten, _ := new(big.Int).SetString((*item)[2:], 16)
		finalRet = append(finalRet, ten.String())
	}
	return finalRet, err
}
