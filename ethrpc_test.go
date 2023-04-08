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
	nodeUrl := "https://mainnet.infura.io/v3/"
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
