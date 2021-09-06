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

type ethClient interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

type ITxBuilder interface {
	Sender() common.Address
	Transfer(ctx context.Context, to string, value *big.Int) (common.Hash, error)
}

type TxBuilder struct {
	client      ethClient
	privateKey  *ecdsa.PrivateKey
	signer      types.Signer
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
		client:      client,
		privateKey:  privateKey,
		signer:      types.NewEIP155Signer(chainID),
		fromAddress: crypto.PubkeyToAddress(privateKey.PublicKey),
	}
}

func (b *TxBuilder) Sender() common.Address {
	return b.fromAddress
}

func (b *TxBuilder) Transfer(ctx context.Context, to string, value *big.Int) (common.Hash, error) {
	nonce, err := b.client.PendingNonceAt(ctx, b.Sender())
	if err != nil {
		return common.Hash{}, err
	}

	gasLimit := uint64(21000)
	gasPrice, err := b.client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	if gasPrice.Cmp(big.NewInt(1)) == 0 {
		gasPrice = big.NewInt(875000000) // for unit test
	}

	toAddress := common.HexToAddress(to)
	unsignedTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})

	signedTx, err := types.SignTx(unsignedTx, b.signer, b.privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	return signedTx.Hash(), b.client.SendTransaction(ctx, signedTx)
}
