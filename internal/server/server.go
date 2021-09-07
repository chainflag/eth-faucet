package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/web"
)

const AddressKey string = "address"

type Server struct {
	chain.TxBuilder
	config *Config
	queue  chan string
}

func NewServer(txBuilder chain.TxBuilder, config *Config) *Server {
	return &Server{
		TxBuilder: txBuilder,
		config:    config,
		queue:     make(chan string, config.queueCap),
	}
}

func (s *Server) Run() {
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(web.Dist()))
	router.Handle("/api/claim", negroni.New(NewLimiter(60*time.Second), negroni.Wrap(s.handleClaim())))
	router.Handle("/api/info", s.handleInfo())
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(router)

	go func() {
		for address := range s.queue {
			txHash, err := s.Transfer(context.Background(), address, s.config.payout)
			if err != nil {
				log.WithError(err).Error("Failed to handle transaction in the queue")
			} else {
				log.WithFields(log.Fields{
					"txHash":  txHash,
					"address": address,
				}).Info("Consume from queue successfully")
			}
		}
	}()

	log.Infof("Starting http server %d", s.config.port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.config.port), n))
}

func (s *Server) handleClaim() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}

		address := r.PostFormValue(AddressKey)
		if len(s.queue) != 0 {
			select {
			case s.queue <- address:
				log.WithFields(log.Fields{
					"address": address,
				}).Info("Added to queue successfully")
				fmt.Fprintf(w, "Added %s to the queue", address)
			default:
				log.Warn("Max queue capacity reached")
				http.Error(w, "Faucet queue is too long, please try again later.", http.StatusServiceUnavailable)
			}
			return
		}

		txHash, err := s.Transfer(r.Context(), address, s.config.payout)
		if err != nil {
			log.WithError(err).Error("Could not send transaction")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.WithFields(log.Fields{
			"txHash":  txHash,
			"address": address,
		}).Info("Funded directly successfully")
		fmt.Fprintf(w, txHash.String())
	}
}

func (s *Server) handleInfo() http.HandlerFunc {
	type info struct {
		Account string `json:"account"`
		Payout  string `json:"payout"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info{
			Account: s.Sender().String(),
			Payout:  s.config.payout.String(),
		})
	}
}
