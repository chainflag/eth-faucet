package server

import (
	"errors"
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

var cidrs []*net.IPNet

func init() {
	maxCidrBlocks := []string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}

	cidrs = make([]*net.IPNet, len(maxCidrBlocks))
	for i, maxCidrBlock := range maxCidrBlocks {
		_, cidr, _ := net.ParseCIDR(maxCidrBlock)
		cidrs[i] = cidr
	}
}

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
	clintIP := getClientIPFromRequest(r)
	if l.limitByKey(w, address) || l.limitByKey(w, clintIP) {
		return
	}

	next.ServeHTTP(w, r)
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

// https://en.wikipedia.org/wiki/Private_network
func isPrivateAddress(address string) (bool, error) {
	ipAddress := net.ParseIP(address)
	if ipAddress == nil {
		return false, errors.New("address is not valid")
	}

	for i := range cidrs {
		if cidrs[i].Contains(ipAddress) {
			return true, nil
		}
	}

	return false, nil
}

func getClientIPFromRequest(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	isPrivate, err := isPrivateAddress(remoteIP)
	if !isPrivate && err == nil {
		return remoteIP
	}

	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		maxProxyCount := 1
		xForwardedForList := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(xForwardedForList[len(xForwardedForList)-maxProxyCount])
	} else if xRealIP := r.Header.Get("X-Real-Ip"); xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	return remoteIP
}
