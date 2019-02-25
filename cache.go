package main

import (
	"time"
)

func (c *cache) new() {
	c.keyVals = make(map[string]*val)
}

func (c *cache) get(url string) (bool, *val) {
	if out, ok := c.keyVals[url]; ok {
		return true, out
	}
	return false, &val{}
}

func (c *cache) set(out *Output, url string, timeIn time.Time) {
	valToStore := val{
		out,
		timeIn,
	}
	c.keyVals[url] = &valToStore
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

func (c *cache) expireOldKeys() {
	c.mu.Lock()
	for url, value := range c.keyVals {
		if value.expiration.Unix()+60 <= time.Now().Unix() {
			delete(c.keyVals, url)
		}
	}
	c.mu.Unlock()
}
