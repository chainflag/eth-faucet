package chain

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type TxBuilder interface {
	Sender() common.Address
	Transfer(ctx context.Context, to string, value *big.Int) (common.Hash, error)
}

type TxBuild struct {
	client      bind.ContractTransactor
	privateKey  *ecdsa.PrivateKey
	signer      types.Signer
	fromAddress common.Address
	nonce       uint64
	nonceMutex  sync.Mutex
}

func NewTxBuilder(provider string, privateKey *ecdsa.PrivateKey, chainID *big.Int) (TxBuilder, error) {
	client, err := ethclient.Dial(provider)
	if err != nil {
		return nil, err
	}

	if chainID == nil {
		chainID, err = client.ChainID(context.Background())
		if err != nil {
			return nil, err
		}
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	return &TxBuild{
		client:      client,
		privateKey:  privateKey,
		signer:      types.NewEIP155Signer(chainID),
		fromAddress: fromAddress,
		nonce:       nonce,
	}, nil
}

func (b *TxBuild) Sender() common.Address {
	return b.fromAddress
}

func (b *TxBuild) Transfer(ctx context.Context, to string, value *big.Int) (common.Hash, error) {
	gasLimit := uint64(21000)
	gasPrice, err := b.client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	toAddress := common.HexToAddress(to)
	unsignedTx := types.NewTx(&types.LegacyTx{
		Nonce:    b.getAndIncrementNonce(),
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})

	signedTx, err := types.SignTx(unsignedTx, b.signer, b.privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	if err = b.client.SendTransaction(ctx, signedTx); err != nil {
		log.Error("failed to send tx", "tx hash", signedTx.Hash().String(), "err", err)
		// check if contain nonce, and reset nonce
		if strings.Contains(err.Error(), "nonce") {
			b.resetNonce(context.Background())
		}
		return common.Hash{}, err
	}

	return signedTx.Hash(), nil
}

func (b *TxBuild) getAndIncrementNonce() uint64 {
	b.nonceMutex.Lock()
	defer b.nonceMutex.Unlock()

	nonce := b.nonce
	b.nonce++
	return nonce
}

func (b *TxBuild) resetNonce(ctx context.Context) {
	nonce, err := b.client.PendingNonceAt(ctx, b.Sender())
	if err != nil {
		log.Warn("failed to reset nonce", "address", b.Sender(), "err", err)
		return
	}

	b.nonce = nonce
}
