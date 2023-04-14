package dao

import (
	"fmt"
	"testing"
)

func TestNewMySQLConnector(t *testing.T) {
	options := &MysqlOptions{
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
	tables := []interface{}{}
	tables = append(tables, Block{}, Transaction{}) //add the table struct
	NewMySQLConnector(options, tables)
	fmt.Println("create tables successful")
}
