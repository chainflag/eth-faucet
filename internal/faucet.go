package internal

import (
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/chainflag/eth-faucet/internal/pkg"
)

type faucet struct {
	payout    *big.Int
	queue     chan string
	txBuilder pkg.ITxBuilder
}

func NewFaucet(capacity int, builder pkg.ITxBuilder) *faucet {
	return &faucet{
		queue:     make(chan string, capacity),
		txBuilder: builder,
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

func (f *faucet) SetPayoutEther(amount int64) {
	payoutWei := new(big.Int).Mul(big.NewInt(amount), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	f.payout = payoutWei
}

func (f faucet) fundTransfer(to string) (string, error) {
	unsignedTx, err := f.txBuilder.BuildUnsignedTx(to, f.payout, nil)
	if err != nil {
		return "", err
	}

	return f.txBuilder.SignAndSubmitTx(unsignedTx)
}

func (f *faucet) Run() {
	for address := range f.queue {
		txHash, err := f.fundTransfer(address)
		if err != nil {
			log.WithError(err).Error("Failed to handle transaction in the queue")
		}
		log.WithFields(log.Fields{
			"txHash":  txHash,
			"address": address,
		}).Info("Consume from queue successfully")
	}
}

func (f *faucet) Close() {
	close(f.queue)
}
