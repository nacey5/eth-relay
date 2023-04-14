package eth_relay

import (
	"math/big"
	"sync"
)

// NonceManager @s map[address]nonce
type NonceManager struct {
	//read write map,must consider of the mul
	lock          sync.Mutex
	nonceMemCache map[string]*big.Int
}

func NewNonceManager() *NonceManager {
	return &NonceManager{
		lock: sync.Mutex{},
	}
}

func (n *NonceManager) SetNonce(address string, nonce *big.Int) {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	n.nonceMemCache[address] = nonce
}

func (n *NonceManager) GetNonce(address string) *big.Int {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.nonceMemCache[address]
}

// PlusNonce @return after pulls nonce
func (n *NonceManager) PlusNonce(address string) *big.Int {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	oldNonce := n.nonceMemCache[address]
	newNonce := oldNonce.Add(oldNonce, big.NewInt(int64(1)))
	n.nonceMemCache[address] = newNonce
	return newNonce
}
