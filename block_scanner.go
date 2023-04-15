package eth_relay

import (
	"encoding/json"
	"errors"
	"eth-relay/dao"
	"eth-relay/model"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
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

// use the getScannerBlockNumber() must accused the init() ---init()--->getScannerBlockNumber()
func (scanner *BlockScanner) getScannerBlockNumber() (*big.Int, error) {
	//use eth request get the most new block number
	newBlockNumber, err := scanner.ethRequester.GetLatestBlockNumber()
	if err != nil {
		return nil, err
	}
	latestNumber := newBlockNumber
	//use the new() init and set the value
	//if not do this according to that,may affect the block num get after
	targetNumber := new(big.Int).Set(scanner.lastNumber)
	//compare the block number
	// -1 if x<y. 0 if x==y.  +1 if x>y
	if latestNumber.Cmp(scanner.lastNumber) < 0 {
		//the most new height must smaller than the set ,then the height of new >= set
	Next:
		for {
			select {
			case <-time.After(time.Duration(4 * time.Second)): //late for 4s retry
				number, err := scanner.ethRequester.GetLatestBlockNumber()
				if err != nil && number.Cmp(scanner.lastNumber) >= 0 {
					break Next //jump for
				}
			}
		}
	}
	return targetNumber, nil // return the height of the block number
}

// scan the block
func (scanner *BlockScanner) scan() error {
	//get the blockNumber you want scan
	targetNumber, err := scanner.getScannerBlockNumber()
	if err != nil {
		return err
	}
	//use the func get the blockNumber that can retry
	fullBlock, err := scanner.retryGetBlockInfoByNumber(targetNumber)
	if err != nil {
		return err
	}
	//the block plu 1
	scanner.lastNumber.Add(scanner.lastNumber, new(big.Int).SetInt64(1))

	//you must accused two tables,must use the transactional for the tables
	tx := scanner.mysql.Db.NewSession() //start affair
	defer tx.Close()

	//prepare the info for the block,must judge the block info exist or not
	block := dao.Block{}
	_, err = tx.Where("block_hash=?", fullBlock.Hash).Get(&block)
	if err == nil && block.Id == 0 {
		//not exists
		block.BlockNumber = scanner.hexToTen(fullBlock.Number).String()
		block.ParentHash = fullBlock.ParentHash
		block.CreateTime = scanner.hexToTen(fullBlock.TimeStamp).Int64()
		block.BlockHash = fullBlock.Hash
		block.Fork = false
		if _, err := tx.Insert(&block); err != nil {
			tx.Rollback()
			return err
		}
	}
	//check the block has fork or not
	if scanner.isFork(&block) {
		data, _ := json.Marshal(fullBlock)
		scanner.log("fork", string(data))
		tx.Commit() //though fork,must save the block transactional commit
		scanner.fork = true
		return errors.New("fork check") //return error ,let your lay occur
	}

	// analytic
	scanner.log(
		"scan block start==>", "number:",
		scanner.hexToTen(fullBlock.Number), "hash:", fullBlock.Hash)
	for index, transaction := range fullBlock.Transactions {
		//the print operation mock the init operation,and to the every tx,we can get the info for that
		scanner.log("tx hash==>", transaction.(map[string]interface{})["hash"])
		//in order to control the print info ,for 5 that is enough
		if index == 5 {
			break
		}
	}
	scanner.log("scan block finish \n==========")
	//save the transaction info
	for _, transaction := range fullBlock.Transactions {
		if _, err := tx.InsertOne(&transaction); err != nil {
			tx.Rollback() //affair rollback
			return err
		}
	}
	return tx.Commit()
}

// Start listen the stop event and start an go routine to scan prevent from pending the main thread
func (scanner *BlockScanner) Start() error {
	scanner.lock.Lock() //lock
	//init the data----mainly init the blockNumber inner
	if err := scanner.init(); err != nil {
		scanner.lock.Unlock()
		return err
	}
	execute := func() {
		//scan func,the func for scan the block
		if err := scanner.scan(); err != nil {
			scanner.log(err.Error())
			return
		}
		time.Sleep(1 * time.Second)
	}
	// start a go routine to for the block
	go func() {
		for {
			select {
			case <-scanner.stop:
				scanner.log("finished block scanner!")
				return
			default:
				if !scanner.fork {
					//inter this branch approve that has not any fork,can execute every circle
					execute()
					continue
				}
				// if fork=true,that has fork branch,re-init
				// re-get from database pre successful and not any fork
				if err := scanner.init(); err != nil {
					scanner.log(err.Error())
					return
				}
				scanner.fork = false
			}

		}
	}()

	return nil
}

func (scanner *BlockScanner) Stop() {
	scanner.lock.Unlock()
	scanner.stop <- true
}
