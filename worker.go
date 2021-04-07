package main

import (
	"log"

	"github.com/chainflag/eth-faucet/core"
)

type worker struct {
	faucet *core.Faucet
	queue  chan string
}

func NewWorker(capacity int, faucet *core.Faucet) *worker {
	return &worker{
		faucet: faucet,
		queue:  make(chan string, capacity),
	}
}

func (w *worker) TryEnqueue(job string) bool {
	select {
	case w.queue <- job:
		return true
	default:
		return false
	}
}

func (w *worker) Run() {
	for address := range w.queue {
		txHash, err := w.faucet.TransferEther(address)
		if err != nil {
			log.Println(err)
		}
		log.Println(txHash)
	}
}

func (w *worker) Close() {
	close(w.queue)
}
