package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"

	"github.com/chainflag/eth-faucet/internal/chain"
)

type MockTxBuilder struct {
	mock.Mock
}

func (m *MockTxBuilder) Sender() common.Address {
	args := m.Called()
	return args.Get(0).(common.Address)
}

func (m *MockTxBuilder) Transfer(ctx context.Context, to string, value *big.Int) (common.Hash, error) {
	args := m.Called(ctx, to, value)
	return args.Get(0).(common.Hash), args.Error(1)
}

func setupTestServer(mockBuilder chain.TxBuilder) *Server {
	cfg := &Config{
		httpPort:   8080,
		proxyCount: 0,
		interval:   0,
		network:    "testnet",
		symbol:     "ETH",
		payout:     1.0,
	}
	return NewServer(mockBuilder, cfg)
}

func TestHandleClaim(t *testing.T) {
	mockBuilder := new(MockTxBuilder)
	expectedAddress := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	expectedAmount := chain.EtherToWei(1.0)
	mockBuilder.On("Transfer", mock.Anything, expectedAddress, expectedAmount).Return(common.Hash{1}, nil)

	server := setupTestServer(mockBuilder)
	reqBody := strings.NewReader(fmt.Sprintf(`{"address": "%s"}`, expectedAddress))
	req, err := http.NewRequest("POST", "/api/claim", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := server.handleClaim()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, rr.Code)
	}

	var resp claimResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	mockBuilder.AssertExpectations(t)

}

func TestHandleInfo(t *testing.T) {
	mockBuilder := new(MockTxBuilder)
	mockBuilder.On("Sender").Return(common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"))

	server := setupTestServer(mockBuilder)
	req, err := http.NewRequest("GET", "/api/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := server.handleInfo()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, rr.Code)
	}

	var resp infoResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	mockBuilder.AssertExpectations(t)
}
