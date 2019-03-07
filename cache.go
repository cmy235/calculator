package main

import (
	"sync"
	"time"
)

// cache is the in memory key/value store
type cache struct {
	mu sync.RWMutex
	// keyVals map[string]*val
	keyVals map[string]val
}

func (c *cache) new() {
	c.keyVals = make(map[string]val)
}

// TODO read lock here?
func (c *cache) get(url string) (bool, val) {
	c.mu.RLock()
	out, ok := c.keyVals[url]
	c.mu.RUnlock()
	return ok, out
}

func (c *cache) set(out Output, url string) {
	valToStore := val{
		out,
		time.Now().Add(time.Second * 60),
	}
	c.mu.Lock()
	c.keyVals[url] = valToStore
	c.mu.Unlock()
}

func (c *cache) setTicker() {
	ticker := time.NewTicker(1 * time.Second)
	datastore.startTicking(*ticker)
}

func (c *cache) startTicking(ticker time.Ticker) {
	defer ticker.Stop()
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			c.expireOldKeys()
		}
	}
}

// TODO
// fix with mutex

// func (c *cache) expireOldKeys() {
// 	c.mu.Lock()
// 	for url, value := range c.keyVals {
// 		if value.expiration.Unix()+60 <= time.Now().Unix() {
// 			delete(c.keyVals, url)
// 		}
// 	}
// 	c.mu.Unlock()
// }

func (c *cache) expireOldKeys() {
	for url, value := range c.keyVals {
		c.mu.RLock()
		expired := value.expiration.Before(time.Now())
		c.mu.RUnlock()

		if expired {
			c.mu.Lock()
			delete(c.keyVals, url)
			c.mu.Unlock()
		}
	}
}
