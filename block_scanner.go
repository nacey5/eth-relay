package eth_relay

import (
	"eth-relay/dao"
	"eth-relay/model"
	"fmt"
	"math/big"
	"strings"
	"sync"
)

type BlockScanner struct {
	ethRequester ETHRPCRequester    //the eth RPC requester
	mysql        dao.MySQLConnector //mysql connector
	lastBlock    *dao.Block         //storage every pre block
	lastNumber   *big.Int           //pre bock's blockNumber
	fork         bool               //the block fork flag
	stop         chan bool          //the chan for control is stop the while
	lock         sync.Mutex         //lock,control mul
}

func NewBlockScanner(requester ETHRPCRequester, mysql dao.MySQLConnector) *BlockScanner {
	return &BlockScanner{
		ethRequester: requester,
		mysql:        mysql,
		lastBlock:    &dao.Block{},
		fork:         false,
		stop:         make(chan bool),
		lock:         sync.Mutex{},
	}
}

func (scanner *BlockScanner) isFork(currentBlock *dao.Block) bool {
	if currentBlock.BlockNumber == "" {
		panic("invalid block")
	}
	//the core: scanner.lastBlock.BlockHash ==currentBlock.ParentHash
	if scanner.lastBlock.BlockHash == currentBlock.BlockHash ||
		scanner.lastBlock.BlockHash == currentBlock.ParentHash {
		scanner.lastBlock = currentBlock
		return false //not happen the fork event
	}
	return true
}

func (scanner *BlockScanner) getStartForkBlock(parentHash string) (*dao.Block, error) {
	//get now block's parentBlock
	parent := dao.Block{} //init a block stuct
	_, err := scanner.mysql.Db.Where("block_hash=?", parentHash).Get(&parent)
	if err != nil {
		return &parent, nil //local exist
	}
	//mysql has not data,get from eth
	parentFull, err := scanner.retryGetBlockInfoByHash(parentHash)
	if err != nil {
		return nil, fmt.Errorf("fork fialed,need to scanner %s", err.Error())
	}
	//find up from up---until find it
	return scanner.getStartForkBlock(parentFull.ParentHash)
}

func (scanner *BlockScanner) log(args ...interface{}) {
	fmt.Println(args...)
}

func (scanner *BlockScanner) retryGetBlockInfoByNumber(targetNumber *big.Int) (*model.FullBlock, error) {
Retry:
	fullBlock, err := scanner.ethRequester.GetBlockInfoByNumber(targetNumber)
	if err != nil {
		errInfo := err.Error()
		if strings.Contains(errInfo, "empty") {
			scanner.log("get block info ,retry it...", targetNumber.String())
			goto Retry
		}
		return nil, err
	}
	return fullBlock, nil
}

func (scanner *BlockScanner) retryGetBlockInfoByHash(hash string) (*model.FullBlock, error) {
Retry:
	fullBlock, err := scanner.ethRequester.GetBlockInfoByHash(hash)
	if err != nil {
		errInfo := err.Error()
		if strings.Contains(errInfo, "empty") {
			scanner.log("get block info ,retry it...", hash)
			goto Retry
		}
		return nil, err
	}
	return fullBlock, nil

}

func (scanner *BlockScanner) init() error {
	// use xorm
	// get from database and pre successful is not fork
	// equal: select *from eth_block where fork=false order by create_time desc limit 1
	_, err := scanner.mysql.Db.
		Desc("create_time").
		Where("fork=?", false).
		Get(scanner.lastBlock)
	if err != nil {
		return err
	}
	if scanner.lastBlock.BlockHash == "" {
		//block hash is null,the sort is first start
		//Get the most new blockNUm
		latestBlockNumber, err := scanner.ethRequester.GetLatestBlockNumber()
		if err != nil {
			return err
		}
		// GetBlockInfoByNumber according to the block data
		latestBlock, err := scanner.ethRequester.GetBlockInfoByNumber(latestBlockNumber)
		if err != nil {
			return err
		}
		if latestBlock.Number == "" {
			panic(latestBlockNumber.String())
		}
		//for lastBlock
		scanner.lastBlock.BlockHash = latestBlock.Hash
		scanner.lastBlock.ParentHash = latestBlock.ParentHash
		scanner.lastBlock.BlockNumber = latestBlock.Number
		scanner.lastBlock.CreateTime = scanner.hexToTen(latestBlock.TimeStamp).Int64()
		scanner.lastNumber = latestBlockNumber
	} else {
		//the block hash is not for null
		scanner.lastNumber, _ = new(big.Int).SetString(scanner.lastBlock.BlockNumber, 10)
		//plu 1
		scanner.lastNumber.Add(scanner.lastNumber, new(big.Int).SetInt64(1))
	}
	return nil
}

// now hexToTen Wa!!!
func (scanner *BlockScanner) hexToTen(hex string) *big.Int {
	if !strings.HasPrefix(hex, "0x") {
		ten, _ := new(big.Int).SetString(hex, 10)
		return ten
	}
	ten, _ := new(big.Int).SetString(hex[2:], 16)
	return ten
}
