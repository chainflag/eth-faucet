package chain

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxBuilder interface {
	Sender() common.Address
	Transfer(ctx context.Context, to string, value *big.Int) (common.Hash, error)
}

type TxBuild struct {
	client               bind.ContractTransactor
	privateKey           *ecdsa.PrivateKey
	signer               types.Signer
	fromAddress          common.Address
	nonce                uint64
	supportsEIP1559      bool
	nonceRefreshEvery    uint64
	nonceRefreshInterval time.Duration
	lastRefreshTime      time.Time
	nonceMu              sync.Mutex
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

	supportsEIP1559, err := checkEIP1559Support(client)
	if err != nil {
		return nil, err
	}

	txBuilder := &TxBuild{
		client:               client,
		privateKey:           privateKey,
		signer:               types.NewLondonSigner(chainID),
		fromAddress:          crypto.PubkeyToAddress(privateKey.PublicKey),
		supportsEIP1559:      supportsEIP1559,
		lastRefreshTime:      time.Time{},
		nonceMu:              sync.Mutex{},
		nonceRefreshInterval: time.Minute * 1,
		nonceRefreshEvery:    100,
	}

	return txBuilder, nil
}

func (b *TxBuild) Sender() common.Address {
	return b.fromAddress
}

func (b *TxBuild) Transfer(ctx context.Context, to string, value *big.Int) (common.Hash, error) {
	gasLimit := uint64(21000)
	toAddress := common.HexToAddress(to)

	var err error
	var unsignedTx *types.Transaction

	nonce, err := b.getNextNonce(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	if b.supportsEIP1559 {
		unsignedTx, err = b.buildEIP1559Tx(ctx, &toAddress, value, gasLimit, nonce)
	} else {
		unsignedTx, err = b.buildLegacyTx(ctx, &toAddress, value, gasLimit, nonce)
	}

	if err != nil {
		return common.Hash{}, err
	}

	signedTx, err := types.SignTx(unsignedTx, b.signer, b.privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	if err = b.client.SendTransaction(ctx, signedTx); err != nil {
		return common.Hash{}, err
	}

	return signedTx.Hash(), nil
}

func (b *TxBuild) buildEIP1559Tx(ctx context.Context, to *common.Address, value *big.Int, gasLimit uint64, nonce uint64) (*types.Transaction, error) {
	header, err := b.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	gasTipCap, err := b.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}

	// gasFeeCap = baseFee * 2 + gasTipCap
	gasFeeCap := new(big.Int).Mul(header.BaseFee, big.NewInt(2))
	gasFeeCap = new(big.Int).Add(gasFeeCap, gasTipCap)

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   b.signer.ChainID(),
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        to,
		Value:     value,
	}), nil
}

func (b *TxBuild) buildLegacyTx(ctx context.Context, to *common.Address, value *big.Int, gasLimit uint64, nonce uint64) (*types.Transaction, error) {
	gasPrice, err := b.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       to,
		Value:    value,
	}), nil
}

func (b *TxBuild) getNextNonce(ctx context.Context) (uint64, error) {
	b.nonceMu.Lock()
	defer b.nonceMu.Unlock()
	b.nonce++
	// fetch from RPC every n txs, or after refresh interval - whichever is hit first
	if time.Since(b.lastRefreshTime) > b.nonceRefreshInterval || b.nonce%b.nonceRefreshEvery == 0 {
		n, err := b.client.PendingNonceAt(ctx, b.fromAddress)
		if err != nil {
			return 0, err
		}
		b.nonce = n
		b.lastRefreshTime = time.Now()
	}
	nonce := b.nonce
	return nonce, nil
}

func checkEIP1559Support(client *ethclient.Client) (bool, error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return false, err
	}

	return header.BaseFee != nil && header.BaseFee.Cmp(big.NewInt(0)) > 0, nil
}
