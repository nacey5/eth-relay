package eth_relay

import (
	"eth-relay/dao"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestBlockScannerStart(t *testing.T) {
	//init the rpc requester
	//the eth main-net
	mainNet := "https://eth-mainnet.g.alchemy.com/v2/5Qr_VuMZh2dAdvsqacDUIW8ew9LuuLfC"
	requester := NewETHRPCRequester(mainNet)

	//init the database connect
	option := dao.MysqlOptions{
		HostName:           "127.0.0.1",
		Port:               "3306",
		DbName:             "eth_relay",
		User:               "root",
		Password:           "a1160124552",
		TablePrefix:        "eth_",
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
		ConnMaxLifetime:    15,
	}

	// add the data tables
	tables := []interface{}{}
	tables = append(tables, dao.Block{}, dao.Transaction{})
	//according to the config init the mysql connection
	mysql := dao.NewMySQLConnector(&option, tables)

	//init the root scanner
	scanner := NewBlockScanner(*requester, mysql)
	err := scanner.Start()
	if err != nil {
		panic(err)
	}
	//use the select mock the main thread pending
	select {}
}
