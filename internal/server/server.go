package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LK4D4/trylock"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/web"
)

const AddressKey string = "address"

type Server struct {
	chain.TxBuilder
	mutex trylock.Mutex
	cfg   *Config
	queue chan string
}

func NewServer(builder chain.TxBuilder, cfg *Config) *Server {
	return &Server{
		TxBuilder: builder,
		cfg:       cfg,
		queue:     make(chan string, cfg.queueCap),
	}
}

func (s *Server) setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(web.Dist()))
	limiter := NewLimiter(s.cfg.proxyCount, s.cfg.interval*time.Minute)
	router.Handle("/api/claim", negroni.New(limiter, negroni.Wrap(s.handleClaim())))
	router.Handle("/api/info", s.handleInfo())

	return router
}

func (s *Server) Run() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			s.consumeQueue()
		}
	}()

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(s.setupRouter())
	log.Infof("Starting http server %d", s.cfg.httpPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.cfg.httpPort), n))
}

func (s *Server) consumeQueue() {
	if len(s.queue) == 0 {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	for len(s.queue) != 0 {
		address := <-s.queue
		txHash, err := s.Transfer(context.Background(), address, s.cfg.payout)
		if err != nil {
			log.WithError(err).Error("Failed to handle transaction in the queue")
		} else {
			log.WithFields(log.Fields{
				"txHash":  txHash,
				"address": address,
			}).Info("Consume from queue successfully")
		}
	}
}

func (s *Server) handleClaim() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}

		address := r.PostFormValue(AddressKey)
		// Try to lock mutex if the work queue is empty
		if len(s.queue) != 0 || !s.mutex.TryLock() {
			select {
			case s.queue <- address:
				log.WithFields(log.Fields{
					"address": address,
				}).Info("Added to queue successfully")
				fmt.Fprintf(w, "Added %s to the queue", address)
			default:
				log.Warn("Max queue capacity reached")
				errMsg := "Faucet queue is too long, please try again later"
				http.Error(w, errMsg, http.StatusServiceUnavailable)
			}
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		txHash, err := s.Transfer(ctx, address, s.cfg.payout)
		s.mutex.Unlock()
		if err != nil {
			log.WithError(err).Error("Failed to send transaction")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.WithFields(log.Fields{
			"txHash":  txHash,
			"address": address,
		}).Info("Funded directly successfully")
		fmt.Fprintf(w, "Txhash: %s", txHash)
	}
}

func (s *Server) handleInfo() http.HandlerFunc {
	type info struct {
		Account   string `json:"account"`
		ChainName string `json:"network"`
		Payout    string `json:"payout"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info{
			Account:   s.Sender().String(),
			ChainName: s.cfg.chainName,
			Payout:    s.cfg.payout.String(),
		})
	}
}
