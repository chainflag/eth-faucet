package chain

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ITxBuilder interface {
	Sender() common.Address
	Transfer(ctx context.Context, to string, value *big.Int) (string, error)
}

type TxBuilder struct {
	chainID     *big.Int
	client      *ethclient.Client
	privateKey  *ecdsa.PrivateKey
	fromAddress common.Address
}

func NewTxBuilder(provider string, privateKey *ecdsa.PrivateKey, chainID *big.Int) ITxBuilder {
	client, err := ethclient.Dial(provider)
	if err != nil {
		panic(err)
	}

	if chainID == nil {
		chainID, err = client.ChainID(context.Background())
		if err != nil {
			panic(err)
		}
	}

	return &TxBuilder{
		chainID:     chainID,
		client:      client,
		privateKey:  privateKey,
		fromAddress: crypto.PubkeyToAddress(privateKey.PublicKey),
	}
}

func (b *TxBuilder) Sender() common.Address {
	return b.fromAddress
}

func (b *TxBuilder) Transfer(ctx context.Context, to string, value *big.Int) (string, error) {
	nonce, err := b.client.PendingNonceAt(ctx, b.fromAddress)
	if err != nil {
		return "", err
	}

	gasLimit := uint64(21000)
	gasPrice, err := b.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(to)
	unsignedTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})

	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(b.chainID), b.privateKey)
	if err != nil {
		return unsignedTx.Hash().String(), err
	}

	return signedTx.Hash().String(), b.client.SendTransaction(ctx, signedTx)
}
