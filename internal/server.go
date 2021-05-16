package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type request struct {
	Address string `json:"address"`
}

type server struct {
	faucet *faucet
}

func NewServer(faucet *faucet) *server {
	return &server{faucet: faucet}
}

func (s server) Run(port int) {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "$ curl -X POST -d '{\"address\":\"Your ETH address\"}'")
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

	if s.faucet.isEmptyQueue() {
		txHash, err := s.faucet.fundTransfer(req.Address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Fprintf(w, txHash)
		return
	}

	if !s.faucet.tryEnqueue(req.Address) {
		http.Error(w, "Max queue capacity reached", http.StatusServiceUnavailable)
		return
	}

	fmt.Fprintf(w, "Added %s to the queue", req.Address)
}
