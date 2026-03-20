package server

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni/v3"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/web"
)

type Server struct {
	txBuilder chain.TxBuilder
	cfg       *Config
	server    *http.Server
	payoutWei *big.Int
}

func NewServer(builder chain.TxBuilder, cfg *Config) *Server {
	return &Server{
		txBuilder: builder,
		cfg:       cfg,
		payoutWei: chain.EtherToWei(cfg.payout),
	}
}

func (s *Server) setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(web.Dist()))
	limiter := NewLimiter(s.cfg.proxyCount, time.Duration(s.cfg.interval)*time.Minute)
	middlewares := []negroni.Handler{limiter}
	if s.cfg.hcaptchaSecret != "" {
		middlewares = append(middlewares, NewCaptcha(s.cfg.hcaptchaSiteKey, s.cfg.hcaptchaSecret))
	}
	middlewares = append(middlewares, negroni.Wrap(s.handleClaim()))
	router.Handle("/api/claim", negroni.New(middlewares...))
	router.Handle("/api/info", s.handleInfo())

	return router
}

func (s *Server) Run() {
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(s.setupRouter())

	s.server = &http.Server{
		Addr:         ":" + strconv.Itoa(s.cfg.httpPort),
		Handler:      n,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Infof("Starting http server on port %d", s.cfg.httpPort)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Server) handleClaim() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		address, ok := r.Context().Value(addressContextKey).(string)
		if !ok || address == "" {
			renderJSON(w, claimResponse{Message: "invalid request"}, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		txHash, err := s.txBuilder.Transfer(ctx, address, new(big.Int).Set(s.payoutWei))
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"address": address,
			}).Error("Failed to send transaction")
			renderJSON(w, claimResponse{Message: fmt.Sprintf("Transaction failed: %v", err)}, http.StatusInternalServerError)
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
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		renderJSON(w, infoResponse{
			Account:         s.txBuilder.Sender().String(),
			Network:         s.cfg.network,
			Symbol:          s.cfg.symbol,
			Payout:          strconv.FormatFloat(s.cfg.payout, 'f', -1, 64),
			HcaptchaSiteKey: s.cfg.hcaptchaSiteKey,
		}, http.StatusOK)
	}
}
