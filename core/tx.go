package core

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ITxBuilder interface {
	BuildUnsignedTx(to string, value *big.Int, data []byte) (*types.Transaction, error)
	SubmitSignedTx(tx *types.Transaction) (*types.Transaction, error)
}

type txBuilder struct {
	rpc         *ethclient.Client
	chainID     *big.Int
	privkey     *ecdsa.PrivateKey
	fromAddress common.Address
}

func NewTxBuilder(provider, hexkey string) *txBuilder {
	client, err := ethclient.Dial(provider)
	if err != nil {
		panic(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		panic(err)
	}

	publicKeyECDSA, _ := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &txBuilder{
		rpc:         client,
		chainID:     chainID,
		privkey:     privateKey,
		fromAddress: fromAddress,
	}
}

func (b txBuilder) BuildUnsignedTx(to string, value *big.Int, data []byte) (*types.Transaction, error) {
	nonce, err := b.rpc.PendingNonceAt(context.Background(), b.fromAddress)
	if err != nil {
		return nil, err
	}

	toAddress := common.HexToAddress(to)
	gasLimit, err := b.rpc.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})

	gasPrice, err := b.rpc.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	return tx, nil
}

func (b txBuilder) SubmitSignedTx(tx *types.Transaction) (*types.Transaction, error) {
	tx, err := types.SignTx(tx, types.NewEIP155Signer(b.chainID), b.privkey)
	if err != nil {
		return nil, err
	}

	if err := b.rpc.SendTransaction(context.Background(), tx); err != nil {
		return nil, err
	}

	return tx, nil
}
