package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

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
	s.router.HandleFunc("/api/claim", s.handleClaim())
	s.router.HandleFunc("/api/info", s.handleInfo())
}

func (s server) Start(port int) {
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.Use(negroni.NewStatic(web.Dist()))
	n.UseHandler(s.router)

	log.Infof("Starting http server %d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), n))
}

func (s server) handleClaim() http.HandlerFunc {
	type claimReq struct {
		Address string `json:"address"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}

		var req claimReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
		if len(req.Address) == 0 || !re.MatchString(req.Address) {
			http.Error(w, "Invalid address", http.StatusBadRequest)
			return
		}

		if !s.faucet.isEmptyQueue() {
			if s.faucet.tryEnqueue(req.Address) {
				log.WithFields(log.Fields{
					"address": req.Address,
				}).Info("Added to queue successfully")
				fmt.Fprintf(w, "Added %s to the queue", req.Address)
			} else {
				log.Warn("Max queue capacity reached")
				http.Error(w, "Faucet queue is too long, please try again later.", http.StatusServiceUnavailable)
			}
			return
		}

		txHash, err := s.faucet.Transfer(context.Background(), req.Address, s.faucet.GetPayoutWei())
		if err != nil {
			log.WithError(err).Error("Could not send transaction")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.WithFields(log.Fields{
			"txHash":  txHash,
			"address": req.Address,
		}).Info("Funded directly successfully")
		fmt.Fprintf(w, txHash)
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
			Account: s.faucet.FromAddress().String(),
			Payout:  s.faucet.GetPayoutWei().String(),
		})
	}
}
