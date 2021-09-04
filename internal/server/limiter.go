package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/urfave/negroni"

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
	address := r.PostFormValue(AddressKey)
	if !chain.IsValidAddress(address, true) {
		http.Error(w, "invalid address", http.StatusBadRequest)
		return
	}
	ip, err := getIP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if l.limitByKey(w, address) || l.limitByKey(w, ip) {
		return
	}

	next.ServeHTTP(w, r)
	if w.(negroni.ResponseWriter).Status() == http.StatusOK {
		l.cache.SetWithTTL(address, true, l.ttl)
		l.cache.SetWithTTL(ip, true, l.ttl)
	}
}

func (l *Limiter) limitByKey(w http.ResponseWriter, key string) bool {
	if _, ttl, err := l.cache.GetWithTTL(key); err == nil {
		errMsg := fmt.Sprintf("You have exceeded the rate limit. Please wait %s before you try again", ttl.Round(time.Second))
		http.Error(w, errMsg, http.StatusTooManyRequests)
		return true
	}
	return false
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
