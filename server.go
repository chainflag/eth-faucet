package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type request struct {
	Address string `json:"address"`
}

type response struct {
	Message string `json:"msg"`
}

type server struct {
	worker *worker
}

func NewServer(worker *worker) *server {
	return &server{worker: worker}
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

	if !s.worker.TryEnqueue(req.Address) {
		http.Error(w, "Max queue capacity reached", 503)
		return
	}

	resp := &response{
		Message: req.Address,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
