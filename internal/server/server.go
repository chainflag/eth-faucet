package server

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni/v3"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/web"
)

type Server struct {
	chain.TxBuilder
	cfg *Config
}

func NewServer(builder chain.TxBuilder, cfg *Config) *Server {
	return &Server{
		TxBuilder: builder,
		cfg:       cfg,
	}
}

func (s *Server) setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(web.Dist()))
	limiter := NewLimiter(s.cfg.proxyCount, time.Duration(s.cfg.interval)*time.Minute)
	hcaptcha := NewCaptcha(s.cfg.hcaptchaSiteKey, s.cfg.hcaptchaSecret)
	router.Handle("/api/claim", negroni.New(limiter, hcaptcha, negroni.Wrap(s.handleClaim())))
	router.Handle("/api/info", s.handleInfo())

	return router
}

func (s *Server) Run() {
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(s.setupRouter())
	log.Infof("Starting http server %d", s.cfg.httpPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.cfg.httpPort), n))
}

func (s *Server) handleClaim() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}

		// The error always be nil since it has already been handled in limiter
		address, _ := readAddress(r)

		// Check mainnet balance if configured
		if s.cfg.minMainnetBalance != nil && s.cfg.minMainnetBalance.Cmp(big.NewInt(0)) > 0 {
			var balance *big.Int

			// Try to get balance from cache first
			if cachedBalance, exists := GetCachedBalance(address); exists {
				balance = cachedBalance
			} else {
				// If not in cache, query the mainnet node
				mainnetClient, err := ethclient.Dial(s.cfg.mainnetProvider)
				if err != nil {
					log.WithError(err).Error("Failed to connect to mainnet provider")
					renderJSON(w, claimResponse{Message: "Failed to check mainnet balance"}, http.StatusInternalServerError)
					return
				}
				defer mainnetClient.Close()

				balance, err = mainnetClient.BalanceAt(r.Context(), common.HexToAddress(address), nil)
				if err != nil {
					log.WithError(err).Error("Failed to get mainnet balance")
					renderJSON(w, claimResponse{Message: "Failed to check mainnet balance"}, http.StatusInternalServerError)
					return
				}

				// Cache the balance
				CacheBalance(address, balance)
			}

			if balance.Cmp(s.cfg.minMainnetBalance) < 0 {
				renderJSON(w, claimResponse{Message: fmt.Sprintf("Destination address must have at least %s ETH on mainnet", chain.WeiToEther(s.cfg.minMainnetBalance))}, http.StatusBadRequest)
				return
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		txHash, err := s.Transfer(ctx, address, chain.GweiToWei(s.cfg.payout))
		if err != nil {
			log.WithError(err).Error("Failed to send transaction")
			renderJSON(w, claimResponse{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		log.WithFields(log.Fields{
			"txHash":  txHash,
			"address": address,
		}).Info("Transaction sent successfully")
		resp := claimResponse{Message: fmt.Sprintf("Txhash: %s", txHash)}
		renderJSON(w, resp, http.StatusOK)
	}
}

func (s *Server) handleInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}
		renderJSON(w, infoResponse{
			Account:         s.Sender().String(),
			Network:         s.cfg.network,
			Symbol:          s.cfg.symbol,
			Payout:          strconv.FormatInt(s.cfg.payout, 10),
			HcaptchaSiteKey: s.cfg.hcaptchaSiteKey,
			LogoURL:         s.cfg.logoURL,
			BackgroundURL:   s.cfg.backgroundURL,
			FrontendType:    s.cfg.frontendType,
			PaidCustomer:    s.cfg.paidCustomer,
		}, http.StatusOK)
	}
}
