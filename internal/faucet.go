package internal

import (
	"context"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/chainflag/eth-faucet/internal/chain"
)

type faucet struct {
	chain.ITxBuilder
	payout *big.Int
	queue  chan string
}

func NewFaucet(builder chain.ITxBuilder, queueCap int) *faucet {
	return &faucet{
		ITxBuilder: builder,
		queue:      make(chan string, queueCap),
	}
}

func (f faucet) isEmptyQueue() bool {
	return len(f.queue) == 0
}

func (f *faucet) tryEnqueue(job string) bool {
	select {
	case f.queue <- job:
		return true
	default:
		return false
	}
}

func (f faucet) GetPayoutWei() *big.Int {
	return f.payout
}

func (f *faucet) SetPayoutEther(amount int) {
	ether := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	payoutWei := new(big.Int).Mul(big.NewInt(int64(amount)), ether)
	f.payout = payoutWei
}

func (f *faucet) Run() {
	for address := range f.queue {
		txHash, err := f.Transfer(context.Background(), address, f.GetPayoutWei())
		if err != nil {
			log.WithError(err).Error("Failed to handle transaction in the queue")
		} else {
			log.WithFields(log.Fields{
				"txHash":  txHash,
				"address": address,
			}).Info("Consume from queue successfully")
		}
	}
}

func (f *faucet) Close() {
	close(f.queue)
}
