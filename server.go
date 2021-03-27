package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/chainflag/eth-faucet/core"
)

var port int

type request struct {
	Address string `json:"address"`
}

type response struct {
	Address string   `json:"address"`
	TxHash  string   `json:"txhash"`
	Amount  *big.Int `json:"amount"`
}

type server struct {
	faucet *core.Faucet
}

func NewServer(faucet *core.Faucet) *server {
	return &server{faucet: faucet}
}

func (s server) Run() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("$ curl -X POST -d '{\"address\":\"Your ETH address\"}'"))
	})
	r.Methods("POST").Path("/").HandlerFunc(s.faucetHandler)

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(r)

	http.ListenAndServe(":"+strconv.Itoa(port), n)
}

func (s server) faucetHandler(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if len(req.Address) == 0 || !re.MatchString(req.Address) {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}

	txHash, err := s.faucet.TransferEther(req.Address)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	resp := &response{
		Address: req.Address,
		TxHash:  txHash,
		Amount:  s.faucet.GetPayoutWei(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
