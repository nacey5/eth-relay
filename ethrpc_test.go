package eth_relay

import (
	"encoding/json"
	"eth-relay/tool"
	"fmt"
	"testing"
)

func TestNewETHRPCClient(t *testing.T) {
	client2 := NewETHRPCClient("www.baidu.com").GetRpc()
	if client2 == nil {
		fmt.Println("client2 init failed")
	}

	client := NewETHRPCClient("123://456").GetRpc()
	if client == nil {
		fmt.Println("client init failed")
	}
}

func TestGetTransactionByHash(t *testing.T) {
	nodeUrl := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	txHash := "0xd03c50db89188055d05126e6044ae76f2389ca4cbf7dd68309978bcd2846c87f"
	if txHash == "" || len(txHash) != 66 {
		errStr := fmt.Errorf("not egeal").Error()
		fmt.Println(errStr)
		return
	}
	txInfo, err := NewETHRPCRequester(nodeUrl).GetTransactionByHash(txHash)
	if err != nil {
		errStr := fmt.Errorf("not allowed").Error()
		fmt.Println(errStr)
		return
	}
	//successful,json the result
	json, _ := json.Marshal(txInfo)
	fmt.Println(string(json))
}

func Test_GetETHBalance(t *testing.T) {
	nodeUrl := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	address := "0x0D0707963952f2fBA59dD06f2b425ace40b492Fe"
	if address == "" || len(address) != 42 {
		fmt.Println("not egl")
		return
	}
	balance, err := NewETHRPCRequester(nodeUrl).GetETHBalance(address)
	if err != nil {
		//query failed
		fmt.Println("query eth failed,info:", err.Error())
		return
	}
	fmt.Println(balance)
}

func Test_GetETHBalances(t *testing.T) {
	nodeUrl := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	address1 := "0x0D0707963952f2fBA59dD06f2b425ace40b492Fe"
	address2 := "0xf89260db97765A00a343aba8e5682715804769ca"
	addresses := []string{address1, address2}
	balances, err := NewETHRPCRequester(nodeUrl).GetETHBalances(addresses)
	if err != nil {
		//query failed
		fmt.Println("query eth failed,info:", err.Error())
		return
	}
	fmt.Println(balances)
}

func Test_GetERCBalances(t *testing.T) {
	nodeUrl := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	address := "0xe16C1623c1AA7D919cd2241d8b36d9E79C1Be2A2"
	contract1 := "0x78021ABD9b06f0456CB9DB95a846C302c34E8b8D"
	contract2 := "0xB8c77482e45E1E44dE1745F52C74426C631bDD52"
	params := []ERC20BalanceRpcReq{}
	item := ERC20BalanceRpcReq{}
	item.ContractAddress = contract1
	item.UserAddress = address
	item.ContractDecimal = 18
	params = append(params, item)
	item.ContractAddress = contract2
	params = append(params, item)

	balances, err := NewETHRPCRequester(nodeUrl).GetERC20Balances(params)
	if err != nil {
		fmt.Println("query eth failed,info:", err.Error())
		return
	}
	fmt.Println(balances)
}

func TestETHRPCRequester_GetLatestBlockNumber(t *testing.T) {
	nodeUrl := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	number, err := NewETHRPCRequester(nodeUrl).GetLatestBlockNumber()
	if err != nil {
		//query failed
		fmt.Println("get the latest BlockNumber failed,info:", err.Error())
		return
	}
	fmt.Println("decimal:", number.String())
}

func TestGetFullBlockInfo(t *testing.T) {
	nodeUrl := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	requester := NewETHRPCRequester(nodeUrl)
	number, _ := requester.GetLatestBlockNumber()
	fmt.Println("block Number is :", number)
	fullBlock, err := requester.GetBlockInfoByNumber(number)
	if err != nil {
		fmt.Println("get block info failed,info: ", err.Error())
	}
	jsonl, _ := json.Marshal(fullBlock)
	fmt.Println("get the info by blockNumber:\n", string(jsonl))
}

func TestGetFullBlockInfoByHash(t *testing.T) {
	nodeUrl := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	requester := NewETHRPCRequester(nodeUrl)
	blockHash := "0x06286a17cdb1b6a70d79ec6a622a2615708a127ab9ff638c6ab38099bf135acc"
	fmt.Println("blockHash is :", blockHash)
	fullBlock, err := requester.GetBlockInfoByHash(blockHash)
	if err != nil {
		fmt.Println("get block info failed,info: ", err.Error())
	}
	jsonl, _ := json.Marshal(fullBlock)
	fmt.Println("get the info by blockHash:\n", string(jsonl))
}

func TestMakeMethodId(t *testing.T) {
	contractABI := `[ { "constant": true, "inputs": [ { "name": "arg1",
"type": "uint8" }, { "name":
"arg2", "type": "uint8" } ], "name": "add", "outputs": [ {
"name": "", "type":
"uint8" } ], "payable": false, "stateMutability": "pure",
"type": "function" } ]
`
	methodName := "add"
	methodId, err := tool.MakeMethodId(methodName, contractABI)
	if err != nil {
		fmt.Println("create methodId failed", err.Error())
		return
	}
	fmt.Println("create methodId successful", methodId)
}

func TestCreateETHWallet(t *testing.T) {
	nodeUrl := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	address1, err := NewETHRPCRequester(nodeUrl).CreateETHWallet("12345")
	if err != nil {
		fmt.Println("first,failed to create wallet", err.Error())
	} else {
		fmt.Println("first,success,eth address is:", address1)
	}

	address2, err := NewETHRPCRequester(nodeUrl).CreateETHWallet("123456aass")
	if err != nil {
		fmt.Println("second,failed to create wallet", err.Error())
	} else {
		fmt.Println("second,success,eth address is:", address2)
	}
}

func TestUnlockETHWallet(t *testing.T) {
	address := "0x15902acd111a5265e07455fD3B938440A74b465B"
	keyDir := "./keystore"
	err1 := tool.UnlockETHWallet(keyDir, address, "189")
	if err1 != nil {
		fmt.Println("unlock failed", err1.Error())
	} else {
		fmt.Println("unlock successful")

	}
	err2 := tool.UnlockETHWallet(keyDir, address, "123456aass")
	if err2 != nil {
		fmt.Println("unlock failed", err2.Error())
	} else {
		fmt.Println("unlock successful")
	}
}
