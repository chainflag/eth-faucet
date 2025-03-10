package chain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
)

var ErrNoPrivateKeys = errors.New("no private keys provided")

// TxBuilderLoadBalancer implements the TxBuilder interface and distributes
// transactions across multiple underlying TxBuilder instances
type TxBuilderLoadBalancer struct {
	builders []TxBuilder
	current  uint32
	mu       sync.RWMutex
}

// NewTxBuilderLoadBalancer creates a new load balancer with multiple TxBuilder instances
func NewTxBuilderLoadBalancer(provider string, privateKeys []*ecdsa.PrivateKey, chainID *big.Int) (TxBuilder, error) {
	if len(privateKeys) == 0 {
		return nil, ErrNoPrivateKeys
	}

	lb := &TxBuilderLoadBalancer{
		builders: make([]TxBuilder, 0, len(privateKeys)),
		current:  0,
	}

	for _, pk := range privateKeys {
		builder, err := NewTxBuilder(provider, pk, chainID)
		if err != nil {
			return nil, err
		}
		lb.builders = append(lb.builders, builder)
	}

	return lb, nil
}

// Sender returns the current address being used for transactions
// Note: This is mainly for informational purposes as we rotate through addresses
func (lb *TxBuilderLoadBalancer) Sender() common.Address {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	idx := atomic.LoadUint32(&lb.current) % uint32(len(lb.builders))
	return lb.builders[idx].Sender()
}

// Transfer sends a transaction using the next available TxBuilder in round-robin fashion
func (lb *TxBuilderLoadBalancer) Transfer(ctx context.Context, to string, value *big.Int) (common.Hash, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// Get the next builder index using atomic operation for thread safety
	idx := atomic.AddUint32(&lb.current, 1) % uint32(len(lb.builders))

	// Use the selected builder to send the transaction
	return lb.builders[idx].Transfer(ctx, to, value)
}

// GetAllSenders returns all addresses managed by this load balancer
func (lb *TxBuilderLoadBalancer) GetAllSenders() []common.Address {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	addresses := make([]common.Address, len(lb.builders))
	for i, builder := range lb.builders {
		addresses[i] = builder.Sender()
	}
	return addresses
}
