package tool

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

var ETHUnlockMap map[string]accounts.Account

var Unlocks *keystore.KeyStore

func MakeMethodId(methodName string, abiStr string) (string, error) {
	abi := &abi.ABI{}
	err := abi.UnmarshalJSON([]byte(abiStr))
	if err != nil {
		return "", err
	}
	//according to methodName get the method instance
	method := abi.Methods[methodName]
	return string(method.ID), nil
}

func UnlockETHWallet(keyDir string, address, password string) error {
	if Unlocks == nil {
		Unlocks = keystore.NewKeyStore(
			keyDir,
			keystore.StandardScryptN,
			keystore.StandardScryptP,
		)
		if Unlocks == nil {
			return errors.New("ks is nil")
		}
	}
	unlock := accounts.Account{Address: common.HexToAddress(address)}
	//ks.Unlock use the keystore.go unlock the func
	if err := Unlocks.Unlock(unlock, password); err != nil {
		return errors.New("unlock err:" + err.Error())
	}
	if ETHUnlockMap == nil {
		ETHUnlockMap = map[string]accounts.Account{}
	}
	ETHUnlockMap[address] = unlock
	return nil
}

type txdata struct {
	AccountNonce uint64          `json:"nonce" gencodec:"required"`    //tran num
	Price        *big.Int        `json:"gasPrice" gencodec:"required"` //gasPrice
	GasLimit     uint64          `json:"gas" gencodec:"required"`      //gasLimit
	Recipient    *common.Address `json:"to" rlp="nil"`                 //to Recipient address
	Amount       *big.Int        `json:"value" gencodec:"required"`    //the value of coin
	Payload      []byte          `json:"input" gencodec:"required"`    //data params

	//signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	Hash *common.Hash `json:"hash" rlp:"-"`
}

func SignETHTransaction(address string, transaction *types.Transaction) (*types.Transaction, error) {
	if Unlocks == nil {
		return nil, errors.New("you need to init keystore first!")
	}
	account := ETHUnlockMap[address]
	if !common.IsHexAddress(account.Address.String()) {
		// jg the address is unlock or not
		return nil, errors.New("account need to unlock first!")
	}
	//Sign
	return Unlocks.SignTx(account, transaction, nil)
}
