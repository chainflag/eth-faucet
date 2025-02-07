package chain

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"time"
	"sync"
	"sort"
)

func TestTxBuilder(t *testing.T) {
	privateKey, _ := crypto.HexToECDSA("976f9f7772781ff6d1c93941129d417c49a209c674056a3cf5e27e225ee55fa8")
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	simClient := backends.NewSimulatedBackend(
		core.GenesisAlloc{
			fromAddress: {Balance: big.NewInt(10000000000000000)},
		}, 10000000,
	)
	defer simClient.Close()
	var s *backends.SimulatedBackend
	patches := gomonkey.ApplyMethod(reflect.TypeOf(s), "SuggestGasPrice", func(_ *backends.SimulatedBackend, _ context.Context) (*big.Int, error) {
		return big.NewInt(875000000), nil
	})
	defer patches.Reset()

	txBuilder := &TxBuild{
		client:          simClient,
		privateKey:      privateKey,
		signer:          types.NewLondonSigner(big.NewInt(1337)),
		fromAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
		supportsEIP1559: false,
	}
	bgCtx := context.Background()
	toAddress := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	value := big.NewInt(1000)
	txHash, err := txBuilder.Transfer(bgCtx, toAddress.Hex(), value)
	if err != nil {
		t.Errorf("could not add tx to pending block: %v", err)
	}
	simClient.Commit()

	block, err := simClient.BlockByNumber(bgCtx, big.NewInt(1))
	if err != nil {
		t.Errorf("could not get block at height 1: %v", err)
	}
	if txHash != block.Transactions()[0].Hash() {
		t.Errorf("did not commit sent transaction. expected hash %v got hash %v", block.Transactions()[0].Hash(), txHash)
	}

	bal, err := simClient.BalanceAt(bgCtx, toAddress, nil)
	if err != nil {
		t.Error(err)
	}
	if bal.Cmp(value) != 0 {
		t.Errorf("expected balance for to address not received. expected: %v actual: %v", value, bal)
	}
}

type mockEthClient struct {
	bind.ContractTransactor
	pendingNonce uint64
	chainID      *big.Int
}

func (m *mockEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return m.pendingNonce, nil
}
func (m *mockEthClient) ChainID(ctx context.Context) (*big.Int, error) {
	return m.chainID, nil
}
func (m *mockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return &types.Header{BaseFee: big.NewInt(1)}, nil
}
func (m *mockEthClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}
func (m *mockEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}
func (m *mockEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}

func TestTxBuilderNonceConcurrency(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	mock := &mockEthClient{
		pendingNonce: 0,
		chainID:      big.NewInt(1),
	}

	builder := &TxBuild{
		client:          mock,
		privateKey:      privKey,
		signer:          types.NewLondonSigner(mock.chainID),
		fromAddress:     crypto.PubkeyToAddress(privKey.PublicKey),
		supportsEIP1559: true,
		refreshInterval: time.Hour,
		lastRefreshTime: time.Now(),
	}

	const total = 50
	var wg sync.WaitGroup
	wg.Add(total)

	results := make(chan uint64, total)
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			n, err := builder.getNextNonce(context.Background())
			require.NoError(t, err)
			results <- n
		}()
	}

	wg.Wait()
	close(results)

	got := make([]uint64, 0, total)
	for r := range results {
		got = append(got, r)
	}

	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })

	expected := make([]uint64, total)
	for i := 0; i < total; i++ {
		expected[i] = uint64(i + 1)
	}

	require.Equal(t, expected, got)
}
