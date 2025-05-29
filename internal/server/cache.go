package server

import (
	"math/big"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v2"
)

var (
	mainnetBalanceCache = ttlcache.NewCache()
	cacheMutex          sync.Mutex
)

func init() {
	// Set default TTL to 60 seconds
	mainnetBalanceCache.SetTTL(time.Duration(60) * time.Second)
	// Set maximum number of items to 1000
	mainnetBalanceCache.SetCacheSizeLimit(1000)
}

// GetCachedBalance returns the cached balance for an address if it exists and is not expired
func GetCachedBalance(address string) (*big.Int, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if cached, err := mainnetBalanceCache.Get(address); err == nil {
		if balance, ok := cached.(*big.Int); ok {
			return balance, true
		}
	}
	return nil, false
}

// CacheBalance adds a balance to the cache
func CacheBalance(address string, balance *big.Int) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	mainnetBalanceCache.Set(address, balance)
}
