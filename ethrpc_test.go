package eth_relay

import (
	"encoding/json"
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
