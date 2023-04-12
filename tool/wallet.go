package tool

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
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
