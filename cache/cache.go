package cache

import (
	"sync"
	"time"
)

const (
	// DefaultTTL             = 5 * time.Minute
	DefaultCleanupInterval = 30 * time.Second
	DefaultMaxSize         = 1000
)

type item struct {
	value      any
	expiration int64
}

type Cache struct {
	items           map[string]item
	rw              sync.RWMutex
	ttl             time.Duration
	cleanupInterval time.Duration
	maxSize         int
	stopChan        chan struct{}
}

func NewCache(defaultTTL time.Duration, cleanupInterval time.Duration, maxSize int) *Cache {
	// может быть необходим ttl = 0
	/*if defaultTTL <= 0 { 
		defaultTTL = DefaultTTL

	}*/

	if cleanupInterval <= 0 {
		cleanupInterval = DefaultCleanupInterval
	}

	if maxSize <= 0 {
		maxSize = DefaultMaxSize
	}

	c := &Cache{
		items:           make(map[string]item),
		ttl:             defaultTTL,
		cleanupInterval: cleanupInterval,
		stopChan:        make(chan struct{}),
		maxSize:         maxSize,
	}
	go c.backgroundCleanup()
	return c
}

func (c *Cache) Stop() {
	close(c.stopChan)
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if c.maxSize > 0 && len(c.items) >= c.maxSize {
		// log.Println("кэш переполнен")
		c.deleteOldest()
	}

	if ttl == 0 {
		ttl = c.ttl
	}

	var exp int64
	if ttl > 0 {
		exp = time.Now().UnixNano() + int64(ttl)
	} else {
		exp = 0
	}

	c.items[key] = item{
		value:      value,
		expiration: exp,
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	item, isExist := c.items[key]

	if !isExist {
		return nil, false
	}

	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		delete(c.items, key)
		return nil, false
	}
	return item.value, isExist
}

func (c *Cache) Delete(key string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	delete(c.items, key)
}

func (c *Cache) Cleanup() {
	c.rw.Lock()
	defer c.rw.Unlock()

	now := time.Now().UnixNano()

	for key, item := range c.items {
		if item.expiration > 0 && now > item.expiration {
			delete(c.items, key)
		}
	}
}

func (c *Cache) Exists(key string) bool {
	c.rw.Lock()
	defer c.rw.Unlock()

	item, exists := c.items[key]

	if !exists {
		return false
	}

	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		delete(c.items, key)
		return false
	}
	
	return true
}

func (c *Cache) Keys() []string {
	c.rw.Lock()
	defer c.rw.Unlock()

	keys := make([]string, 0, len(c.items))
	
	for key, item := range c.items {
		if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		delete(c.items, key)
		continue
		}
		keys = append(keys, key)
	}
	return keys
}

// метод для тестов
func (c *Cache) Size() int {
    	c.rw.RLock()
   	defer c.rw.RUnlock()
   	
	return len(c.items)
}

func (c *Cache) backgroundCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Cleanup()
		case <-c.stopChan:
			return
		}
	}
}

// для действий в случае переполнения кэша
func (c *Cache) deleteOldest() {
	if len(c.items) == 0 { 
		return
	}

	for key, item := range c.items {
		if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
			delete(c.items, key)
			return
		}
	}

	var targetKey string
	var earliestExp int64 = 0
	found := false

	for key, item := range c.items {
        	if item.expiration > 0 {
            
           		if !found || item.expiration < earliestExp {
                	earliestExp = item.expiration
                	targetKey = key
                	found = true
            		}
        	}
    	}

	if found && targetKey != "" {
        	delete(c.items, targetKey)
        	return
    	}

	for key := range c.items {
       		delete(c.items, key)
        	return
    	}	
}
