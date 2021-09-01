package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ReneKroon/ttlcache/v2"

	"github.com/chainflag/eth-faucet/internal/chain"
)

type Limiter struct {
	cache *ttlcache.Cache
	ttl   time.Duration
}

func NewLimiter(ttl time.Duration) *Limiter {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)
	return &Limiter{
		cache: cache,
		ttl:   ttl,
	}
}

func (l *Limiter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	address, err := getEthAddress(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ip, err := getIP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ttl, err := l.cache.GetWithTTL(address); err != ttlcache.ErrNotFound {
		http.Error(w, fmt.Sprintf("you have exceeded the rate limit. %v", ttl), http.StatusTooManyRequests)
		return
	}
	if _, ttl, err := l.cache.GetWithTTL(ip); err != ttlcache.ErrNotFound {
		http.Error(w, fmt.Sprintf("you have exceeded the rate limit. %v", ttl), http.StatusTooManyRequests)
		return
	}

	ctx := context.WithValue(r.Context(), AddressKey, address)
	next.ServeHTTP(w, r.WithContext(ctx))
	l.cache.SetWithTTL(address, true, l.ttl)
	l.cache.SetWithTTL(ip, true, l.ttl)
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

func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if netIP := net.ParseIP(ip); netIP != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}
