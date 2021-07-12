package pkg

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
	Transfer(ctx context.Context, to string, value *big.Int) (string, error)
}

type txBuilder struct {
	rpc         *ethclient.Client
	privkey     *ecdsa.PrivateKey
	fromAddress common.Address
}

func NewTxBuilder(provider, privateKeyHex string) ITxBuilder {
	client, err := ethclient.Dial(provider)
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(err)
	}

	publicKeyECDSA, _ := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &txBuilder{
		rpc:         client,
		privkey:     privateKey,
		fromAddress: fromAddress,
	}
}

func (b txBuilder) Transfer(ctx context.Context, to string, value *big.Int) (string, error) {
	nonce, err := b.rpc.PendingNonceAt(ctx, b.fromAddress)
	if err != nil {
		return "", err
	}

	gasLimit := uint64(21000)
	gasPrice, err := b.rpc.SuggestGasPrice(ctx)
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

	chainID, err := b.rpc.ChainID(ctx)
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(chainID), b.privkey)
	if err != nil {
		return unsignedTx.Hash().String(), err
	}

	return signedTx.Hash().String(), b.rpc.SendTransaction(ctx, signedTx)
}
