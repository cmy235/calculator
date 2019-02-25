package main

import (
	"time"
)

func (ds *cache) new() {
	ds.keyVals = make(map[string]*val)
}

func (ds *cache) get(url string) (bool, val) {
	if out, ok := ds.keyVals[url]; ok {
		return true, *out
	}

	return false, val{}
}

func (ds *cache) set(out *Output, url string, timeIn time.Time) {
	valToStore := val{
		out,
		timeIn,
	}
	ds.keyVals[url] = &valToStore
}

func (ds *cache) setTicker() {
	ticker := time.NewTicker(1 * time.Second)
	datastore.startTicking(*ticker)
}

func (ds *cache) startTicking(ticker time.Ticker) {
	defer ticker.Stop()
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			ds.expireOldKeys()
		}
	}
}

func (ds *cache) expireOldKeys() {
	ds.mu.Lock()
	for url, value := range ds.keyVals {
		if value.expiration.Unix()+60 <= time.Now().Unix() {
			delete(ds.keyVals, url)
		}
	}
	ds.mu.Unlock()
}
