package eth_relay

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
)

type ETHRPCClient struct {
	NodeUrl string      //代表节点的url链接
	client  *rpc.Client //代表RPC客户端句柄
}

func NewETHRPCClient(nodeUrl string) *ETHRPCClient {
	client := &ETHRPCClient{
		NodeUrl: nodeUrl,
	}
	client.initRpc()
	return client
}

func (erc *ETHRPCClient) initRpc() {
	//使用go-ethereum 库中的RPC库来初始化
	rpcClient, err := rpc.DialHTTP(erc.NodeUrl)
	if err != nil {
		//init failed
		errorInfo := fmt.Errorf("init rpcClient failed :%s", err.Error()).Error()
		panic(errorInfo)
	}
	//successful
	erc.client = rpcClient
}

func (erc *ETHRPCClient) GetRpc() *rpc.Client {
	if erc.client == nil {
		erc.initRpc()
	}
	return erc.client
}
