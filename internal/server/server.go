package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/chainflag/eth-faucet/web"
)

const AddressKey string = "address"

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
	s.router.Handle("/api/claim", negroni.New(NewLimiter(60*time.Second), negroni.Wrap(s.handleClaim())))
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

		address := r.PostFormValue(AddressKey)
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

		txHash, err := s.faucet.Transfer(r.Context(), address, s.faucet.GetPayoutWei())
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
			Account: s.faucet.Sender().String(),
			Payout:  s.faucet.GetPayoutWei().String(),
		})
	}
}
