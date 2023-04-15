package eth_relay

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
)

type ETHRPCClient struct {
	NodeUrl string      //the node url
	client  *rpc.Client //rpc
}

func NewETHRPCClient(nodeUrl string) *ETHRPCClient {
	client := &ETHRPCClient{
		NodeUrl: nodeUrl,
	}
	client.initRpc()
	return client
}

func (erc *ETHRPCClient) initRpc() {
	//use go-ethereum rpc init
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
