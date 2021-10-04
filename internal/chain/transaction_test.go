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
		client:      simClient,
		privateKey:  privateKey,
		signer:      types.NewEIP155Signer(big.NewInt(1337)),
		fromAddress: crypto.PubkeyToAddress(privateKey.PublicKey),
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
