package core

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ITxBuilder interface {
	buildTransaction(to string, value *big.Int, data []byte) (*types.Transaction, error)
	signTransaction(tx *types.Transaction) (*types.Transaction, error)
	submitTransaction(tx *types.Transaction) (string, error)
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

func (b txBuilder) buildTransaction(to string, value *big.Int, data []byte) (*types.Transaction, error) {
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

func (b txBuilder) signTransaction(tx *types.Transaction) (*types.Transaction, error) {
	return types.SignTx(tx, types.NewEIP155Signer(b.chainID), b.privkey)
}

func (b txBuilder) submitTransaction(tx *types.Transaction) (string, error) {
	if err := b.rpc.SendTransaction(context.Background(), tx); err != nil {
		return "", err
	}

	receipt, err := bind.WaitMined(context.Background(), b.rpc, tx)
	if err != nil {
		return "", err
	}

	if receipt.Status == 0 {
		return "", fmt.Errorf("tx %x failed", tx.Hash().Hex())
	}

	return tx.Hash().Hex(), nil
}
