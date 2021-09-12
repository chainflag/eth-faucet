package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/chainflag/eth-faucet/internal/chain"
)

type Limiter struct {
	cache      *ttlcache.Cache
	proxyCount int
	ttl        time.Duration
}

func NewLimiter(proxyCount int, ttl time.Duration) *Limiter {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)
	return &Limiter{
		cache:      cache,
		proxyCount: proxyCount,
		ttl:        ttl,
	}
}

func (l *Limiter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	address := r.PostFormValue(AddressKey)
	if !chain.IsValidAddress(address, true) {
		http.Error(w, "invalid address", http.StatusBadRequest)
		return
	}
	clintIP := getClientIPFromRequest(l.proxyCount, r)
	if l.limitByKey(w, address) || l.limitByKey(w, clintIP) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	next.ServeHTTP(w, r.WithContext(ctx))
	if w.(negroni.ResponseWriter).Status() == http.StatusOK {
		l.cache.SetWithTTL(address, true, l.ttl)
		l.cache.SetWithTTL(clintIP, true, l.ttl)
		log.WithFields(log.Fields{
			"address":  address,
			"clientIP": clintIP,
		}).Info("Maximum request limit has been reached")
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

func getClientIPFromRequest(proxyCount int, r *http.Request) string {
	if proxyCount > 0 {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		xRealIP := r.Header.Get("X-Real-Ip")

		if xForwardedFor != "" {
			xForwardedForParts := strings.Split(xForwardedFor, ",")
			// Avoid reading the user's forged request header by configuring the count of reverse proxies
			partIndex := len(xForwardedForParts) - proxyCount
			if partIndex < 0 {
				partIndex = 0
			}
			return strings.TrimSpace(xForwardedForParts[partIndex])
		}

		if xRealIP != "" {
			return strings.TrimSpace(xRealIP)
		}
	}

	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
