package eth_relay

import (
	"errors"
	"eth-relay/model"
	"eth-relay/tool"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type ETHRPCRequester struct {
	nonceManager *NonceManager
	client       *ETHRPCClient
}

func NewETHRPCRequester(nodeUrl string) *ETHRPCRequester {
	requester := &ETHRPCRequester{}
	requester.nonceManager = NewNonceManager()
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

func (r *ETHRPCRequester) GetLatestBlockNumber() (*big.Int, error) {
	methodName := "eth_blockNumber"
	number := ""
	err := r.client.GetRpc().Call(&number, methodName)
	if err != nil {
		return nil, fmt.Errorf("get the latest BlockNumber failed: %s", err.Error())
	}
	ten, _ := new(big.Int).SetString(number[2:], 16)
	return ten, nil
}

func (r *ETHRPCRequester) GetBlockInfoByNumber(blockNumber *big.Int) (*model.FullBlock, error) {
	number := fmt.Sprintf("%#x", blockNumber)
	methodName := "eth_getBlockByNumber"
	fullBlock := model.FullBlock{}
	err := r.client.GetRpc().Call(&fullBlock, methodName, number, true)
	if err != nil {
		return nil, fmt.Errorf("get block info failed! %s", err.Error())
	}
	if fullBlock.Number == "" {
		return nil, fmt.Errorf("block info is empty %s", blockNumber.String())
	}
	return &fullBlock, nil
}

func (r *ETHRPCRequester) GetBlockInfoByHash(blockHash string) (*model.FullBlock, error) {
	methodName := "eth_getBlockByHash"
	fullBlock := model.FullBlock{}
	err := r.client.GetRpc().Call(&fullBlock, methodName, blockHash, true)
	if err != nil {
		return nil, fmt.Errorf("get block info failed! %s", err.Error())
	}
	if fullBlock.Number == "" {
		return nil, fmt.Errorf("block info is empty %s", blockHash)
	}
	return &fullBlock, nil
}

func (r *ETHRPCRequester) CreateETHWallet(password string) (string, error) {
	if password == "" {
		return "", errors.New("pwd cant empty")
	}
	if len(password) < 6 {
		return "", errors.New("pwd's len must more than 6 words")
	}
	keydir := "./keystore"
	ks := keystore.NewKeyStore(keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	wallet, err := ks.NewAccount(password)
	if err != nil {
		return "0x", err
	}
	return wallet.Address.String(), nil
}

func (r *ETHRPCRequester) SendTransaction(address string, transaction *types.Transaction) (string, error) {
	//signUp the transaction data
	signTx, err := tool.SignETHTransaction(address, transaction)
	if err != nil {
		return "", fmt.Errorf("signUp failed:%s", err.Error())
	}
	//rlp Serialization
	txRlpData, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", fmt.Errorf(" rlp Serialization failed:%s", err.Error())
	}
	// use the eth interface
	txHash := ""
	methodName := "eth_sendRawTransaction"
	err = r.client.GetRpc().Call(&txHash, methodName, common.Bytes2Hex(txRlpData))
	if err != nil {
		return "", fmt.Errorf("send transaction failed:%s", err.Error())
	}
	oldNonce := r.nonceManager.GetNonce(address)
	if oldNonce == nil {
		//get from eth evm,get the most new transaction nonce6 cw....................
		r.nonceManager.SetNonce(address, new(big.Int).SetUint64(transaction.Nonce()))
	}
	return txHash, nil
}

func (r *ETHRPCRequester) GetNonce(address string) (uint64, error) {
	methodName := "eth_getTransactionCount"
	nonce := ""
	//query the most new,according to the etTransactionCount
	err := r.client.GetRpc().Call(&nonce, methodName, address, "pending")
	if err != nil {
		return 0, fmt.Errorf("send transaction failed:%s", err.Error())
	}
	// back to the decimal
	n, _ := new(big.Int).SetString(nonce[2:], 16)
	//the transaction hash
	return n.Uint64(), nil
}

// SendETHTransaction send the eth
func (r *ETHRPCRequester) SendETHTransaction(fromStr, toStr, valueStr string, gasLimit, gasPrice uint64) (string, error) {
	if !common.IsHexAddress(fromStr) || !common.IsHexAddress(toStr) {
		return "", errors.New("invalid address")
	}
	to := common.HexToAddress(toStr)
	gas := new(big.Int).SetUint64(gasPrice)

	realV := tool.GetRealDecimalValue(valueStr, 18)
	if realV == "" {
		return "", errors.New("invalid value")
	}
	amount, _ := new(big.Int).SetString(realV, 10)
	//get the nonce
	nonce := r.nonceManager.GetNonce(fromStr)
	if nonce == nil {
		//nonce is not exist
		n, err := r.GetNonce(fromStr)
		if err != nil {
			return "", fmt.Errorf("get nonce failed %s", err.Error())
		}
		nonce = new(big.Int).SetUint64(n)
		//set the nonce from the fromStr
		r.nonceManager.SetNonce(fromStr, nonce)
	}
	data := []byte("")

	transaction := types.NewTx(&types.LegacyTx{
		Nonce:    nonce.Uint64(),
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gas,
		Data:     data,
	})
	return r.SendTransaction(fromStr, transaction)
}

func (r *ETHRPCRequester) SendERC20Transaction(fromStr, contract, receiver, valueStr string,
	gasLimit, gasPrice uint64,
	decimal int) (string, error) {
	if !common.IsHexAddress(fromStr) ||
		!common.IsHexAddress(contract) ||
		!common.IsHexAddress(receiver) {
		return "", errors.New("invalid address")
	}
	//transfer the contract to the type of the address
	to := common.HexToAddress(contract)
	gasPrice2 := new(big.Int).SetUint64(gasPrice)
	amount := new(big.Int).SetInt64(0)
	nonce := r.nonceManager.GetNonce(fromStr)
	if nonce == nil {
		n, err := r.GetNonce(fromStr)
		if err != nil {
			return "", fmt.Errorf("get nonce failed %s", err.Error())
		}
		nonce = new(big.Int).SetUint64(n)
		r.nonceManager.SetNonce(fromStr, nonce)
	}

	data := tool.BuildERC20TransferData(valueStr, receiver, decimal)
	//transform the data to the byte
	byteData := common.FromHex(data)
	transaction := types.NewTx(&types.LegacyTx{
		Nonce:    nonce.Uint64(),
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice2,
		Data:     byteData,
	})
	return r.SendTransaction(fromStr, transaction)
}
