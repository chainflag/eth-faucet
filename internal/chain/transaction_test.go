package chain

import (
	"context"
	"math/big"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/require"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestTxBuilder(t *testing.T) {
	privateKey, _ := crypto.HexToECDSA("976f9f7772781ff6d1c93941129d417c49a209c674056a3cf5e27e225ee55fa8")
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	simBackend := simulated.NewBackend(types.GenesisAlloc{
		fromAddress: {Balance: big.NewInt(10000000000000000)},
	}, simulated.WithBlockGasLimit(10000000))
	simClient := simBackend.Client()
	defer simBackend.Close()
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
	simBackend.Commit()

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

// ------------------------------------------
// Test 1: Check concurrency with an incrementing on-chain nonce
// ------------------------------------------

type incMockClient struct {
	bind.ContractTransactor
	chainID    *big.Int
	callCount  uint64
	onChainVal uint64
}

func (m *incMockClient) PendingNonceAt(ctx context.Context, _ common.Address) (uint64, error) {
	atomic.AddUint64(&m.callCount, 1)
	// Each call returns 1,2,3,...
	return atomic.AddUint64(&m.onChainVal, 1), nil
}
func (m *incMockClient) ChainID(ctx context.Context) (*big.Int, error) { return m.chainID, nil }
func (m *incMockClient) HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error) {
	return &types.Header{BaseFee: big.NewInt(1)}, nil
}
func (m *incMockClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}
func (m *incMockClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}
func (m *incMockClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}

func TestTxBuilderNonceConcurrency(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	mock := &incMockClient{
		chainID: big.NewInt(1),
	}

	// We'll force a refresh on every transaction, so
	// the code always pulls from the chain (which increments).
	builder := &TxBuild{
		client:            mock,
		privateKey:        privKey,
		signer:            types.NewLondonSigner(mock.chainID),
		fromAddress:       crypto.PubkeyToAddress(privKey.PublicKey),
		supportsEIP1559:   true,
		nonceRefreshEvery: 1,
		lastRefreshTime:   time.Now(),
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
	for n := range results {
		got = append(got, n)
	}

	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })

	// We expect [1..50]
	expected := make([]uint64, total)
	for i := 1; i <= total; i++ {
		expected[i-1] = uint64(i)
	}
	require.Equal(t, expected, got)
}

// ------------------------------------------
// Test 2: Check 'refreshEvery' concurrency
// ------------------------------------------

type refreshMockClient struct {
	bind.ContractTransactor
	pendingNonce uint64
	chainID      *big.Int
	callCount    uint64
}

func (m *refreshMockClient) PendingNonceAt(ctx context.Context, _ common.Address) (uint64, error) {
	atomic.AddUint64(&m.callCount, 1)
	// Always return 100
	return 100, nil
}
func (m *refreshMockClient) ChainID(ctx context.Context) (*big.Int, error) { return m.chainID, nil }
func (m *refreshMockClient) HeaderByNumber(ctx context.Context, _ *big.Int) (*types.Header, error) {
	return &types.Header{BaseFee: big.NewInt(1)}, nil
}
func (m *refreshMockClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}
func (m *refreshMockClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}
func (m *refreshMockClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}

func TestTxBuilderNonceRefreshEvery(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	mock := &refreshMockClient{
		pendingNonce: 100,
		chainID:      big.NewInt(1),
	}
	builder := &TxBuild{
		client:               mock,
		privateKey:           privKey,
		signer:               types.NewLondonSigner(mock.chainID),
		fromAddress:          crypto.PubkeyToAddress(privKey.PublicKey),
		supportsEIP1559:      true,
		nonceRefreshInterval: time.Hour,
		lastRefreshTime:      time.Now(),
		nonceRefreshEvery:    5,
	}

	const total = 20
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

	var got []uint64
	for r := range results {
		got = append(got, r)
	}
	require.Equal(t, total, len(got))

	calls := atomic.LoadUint64(&mock.callCount)
	require.True(t, calls > 1, "Expected multiple calls to PendingNonceAt, got %d", calls)
}
