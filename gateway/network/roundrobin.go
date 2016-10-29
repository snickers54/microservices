package network

import (
	"crypto/sha1"
	"encoding/hex"
	"sync"
	"time"
)

type cachedEndpointRR struct {
	LastUsed     time.Time
	Size         uint
	CurrentIndex uint
}

var cache map[string]cachedEndpointRR
var cacheMutex = &sync.Mutex{}

func invalidateRRStates() {
	cacheMutex.Lock()
	for key, value := range cache {
		if time.Now().Sub(value.LastUsed) > 15*time.Minute {
			delete(cache, key)
		}
	}
	cacheMutex.Unlock()
}

func init() {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				invalidateRRStates()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func createEntry(endpoints Endpoints) cachedEndpointRR {
	csRR := cachedEndpointRR{
		LastUsed:     time.Now(),
		CurrentIndex: 0,
		Size:         uint(len(endpoints)),
	}
	return csRR
}

func createHash(endpoints Endpoints) string {
	aggregate := ""
	for _, endpoint := range endpoints {
		aggregate += endpoint.Service.String()
	}
	h := sha1.New()
	h.Write([]byte(aggregate))
	return hex.EncodeToString(h.Sum(nil))
}

func RoundRobin(endpoints Endpoints) *Endpoint {
	cacheMutex.Lock()
	hash := createHash(endpoints)
	var index uint
	csRR, ok := cache[hash]
	if ok {
		csRR.LastUsed = time.Now()
		csRR.CurrentIndex = (csRR.CurrentIndex + 1) % csRR.Size
	} else {
		csRR = createEntry(endpoints)
	}
	index = csRR.CurrentIndex
	cache[hash] = csRR
	cacheMutex.Unlock()
	return endpoints[index]
}
