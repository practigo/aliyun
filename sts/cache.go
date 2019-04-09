package sts

import (
	"sync"
	"time"
)

// A KeyFunc maps the param to a key string.
type KeyFunc func(*AssumeRoleParam) string

// A cache caches the credentials to descrease requests.
type cache struct {
	g Getter
	k KeyFunc

	// the internal creds map
	mu sync.RWMutex
	v  map[string]Credentials
}

// get from the map
func (c *cache) get(key string) (Credentials, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cred, ok := c.v[key]
	return cred, ok
}

// set to the map
func (c *cache) set(key string, cred Credentials) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v[key] = cred
}

func (c *cache) Get(p *AssumeRoleParam, dur int64) (cred Credentials, err error) {
	key := c.k(p)

	cred, ok := c.get(key)
	if ok {
		// check if it's expired after dur seconds
		if cred.Expiration.After(time.Now().Add(time.Duration(dur) * time.Second)) {
			return
		}
	}

	// not exist or need update
	cred, err = c.g.Get(p, 0)
	if err != nil {
		return
	}

	// update credential
	c.set(key, cred)
	return
}

// Wrap wraps a Getter to enable caching.
// If the provide keyFunc k is nil, the
// DefaultKey is used.
func Wrap(g Getter, k KeyFunc) Getter {
	c := &cache{
		g: g,
		k: k,
		v: make(map[string]Credentials),
	}
	if c.k == nil {
		c.k = DefaultKey
	}
	return c
}

// DefaultKey concatenates the params as a key.
func DefaultKey(p *AssumeRoleParam) string {
	// usually uid is part of roleArn
	return p.RoleArn + p.RoleSessionName + p.Policy
}
