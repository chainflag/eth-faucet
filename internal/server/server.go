package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/web"
)

type server struct {
	faucet *faucet
	router *http.ServeMux
}

func NewServer(faucet *faucet) *server {
	server := &server{
		faucet: faucet,
		router: http.NewServeMux(),
	}
	server.routes()
	return server
}

func (s *server) routes() {
	s.router.Handle("/", http.FileServer(web.Dist()))
	s.router.Handle("/api/claim", s.handleClaim())
	s.router.Handle("/api/info", s.handleInfo())
}

func (s server) Start(port int) {
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(s.router)

	log.Infof("Starting http server %d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), n))
}

func (s server) handleClaim() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}

		address, err := getEthAddress(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !s.faucet.isEmptyQueue() {
			if s.faucet.tryEnqueue(address) {
				log.WithFields(log.Fields{
					"address": address,
				}).Info("Added to queue successfully")
				fmt.Fprintf(w, "Added %s to the queue", address)
			} else {
				log.Warn("Max queue capacity reached")
				http.Error(w, "Faucet queue is too long, please try again later.", http.StatusServiceUnavailable)
			}
			return
		}

		txHash, err := s.faucet.Transfer(context.Background(), address, s.faucet.GetPayoutWei())
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

func (s server) handleInfo() http.HandlerFunc {
	type infoResp struct {
		Account string `json:"account"`
		Payout  string `json:"payout"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(infoResp{
			Account: s.faucet.Sender().String(),
			Payout:  s.faucet.GetPayoutWei().String(),
		})
	}
}

func getEthAddress(r *http.Request) (string, error) {
	type claimReq struct {
		Address string `json:"address"`
	}
	var req claimReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", err
	}
	if !chain.IsValidAddress(req.Address) {
		return "", errors.New("invalid address")
	}

	return chain.ToCheckSumAddress(req.Address), nil
}
