package core

import "math/big"

type Faucet struct {
	payout    *big.Int
	txBuilder ITxBuilder
}

func NewFaucet(builder ITxBuilder) *Faucet {
	return &Faucet{txBuilder: builder}
}

func (f Faucet) TransferEther(to string) (string, error) {
	tx, err := f.txBuilder.buildTransaction(to, f.payout, nil)
	if err != nil {
		return "", err
	}

	if tx, err = f.txBuilder.signTransaction(tx); err != nil {
		return "", err
	}

	return f.txBuilder.submitTransaction(tx)
}

func (f Faucet) GetPayoutWei() *big.Int {
	return f.payout
}

func (f *Faucet) SetPayoutEther(amount int64) {
	payoutWei := new(big.Int).Mul(big.NewInt(amount), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	f.payout = payoutWei
}
