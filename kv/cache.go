package kv

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const defaultCleanInterval = time.Second

// Auxiliary struct to take care of TTL.
type TtlBox struct {
	CreatedAt time.Time
	Expired   *time.Time
	Content   T
}

// T holds user's values.
type T struct {
	V string
}

//
type Cache struct {
	config Configuration
	mu     sync.RWMutex
	values map[string]TtlBox
}

func NewCache(config Configuration) *Cache {
	c := &Cache{
		mu:     sync.RWMutex{},
		values: make(map[string]TtlBox),
	}
	c.configure(config)
	c.startCleaner(defaultCleanInterval)
	return c
}

// configure applies configuration, starts background processes if needed.
func (c *Cache) configure(config Configuration) {
	c.config = config
	if c.config.Storage == nil {
		c.config.Storage = &NotImplementedStorage{}
	}
	err := c.config.Storage.RestoreInto(&c.values)
	if err != nil {
		log.Println(err)
	}
	if c.config.BackupInterval != 0 {
		c.startAutoBackup()
	}
}

// startCleaner initiates background process that deletes expired pairs from cache.
func (c *Cache) startCleaner(delta time.Duration) {
	tick := time.Tick(delta)
	go func() {
		for range tick {
			c.mu.Lock()
			now := time.Now()
			for k, v := range c.values {
				if v.Expired != nil && now.After(*v.Expired) {
					fmt.Printf("deleted by cleaner: %s %v\n", k, v)
					delete(c.values, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}

// startAutoBackup initiates background process that makes snapshots of cache data.
func (c *Cache) startAutoBackup() {
	tick := time.Tick(c.config.BackupInterval)
	go func() {
		for range tick {
			c.makeSnapshot()
		}
	}()
}

func (c *Cache) makeSnapshot() {
	c.mu.RLock()
	mapCopy := make(map[string]TtlBox, len(c.values))
	for k, v := range c.values {
		mapCopy[k] = v
	}
	c.mu.RUnlock()
	if err := c.config.Storage.Save(mapCopy); err != nil {
		log.Println(err)
	}
}

// Add sets value for a key without TTL. If the key existed in the cache
// the new value overwrites the old one.
func (c *Cache) Add(key string, value T) bool {
	return c.add(key, value, nil)
}

// AddWithTtl sets value for a key and stores the expiration date for it.
// If the key existed in the cache the new value overwrites the old one.
func (c *Cache) AddWithTtl(key string, value T, ttl time.Duration) bool {
	expired := time.Now().Add(ttl)
	return c.add(key, value, &expired)
}

// Get returns the value for a given key.
// The boolean value indicates the existence of the key in the cache.
func (c *Cache) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.values[key]
	return value.Content, ok
}

func (c *Cache) GetAll() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	results := make([]T, 0, len(c.values))
	for _, b := range c.values {
		results = append(results, b.Content)
	}
	return results
}

// Remove removes value for a given key.
func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.values, key)
}

// todo: rename
// GetTtl returns the duration of how long ago the value was added to the cache.
// The boolean value indicates the existence of the key in the cache.
func (c *Cache) GetTtl(key string) (time.Duration, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.values[key]
	if !ok {
		return 0, false
	}
	return time.Now().Sub(value.CreatedAt), true
}

// SetTtl changes previous expiration time for the key if it is in the cache.
// Otherwise false is returned
func (c *Cache) SetTtl(key string, ttl *time.Time) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.values[key]
	if !ok {
		return false
	}
	value.Expired = ttl
	c.values[key] = value
	return true
}

func (c *Cache) add(key string, value T, ttl *time.Time) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.values[key]; ok {
		return false
	}
	c.values[key] = TtlBox{
		CreatedAt: time.Now(),
		Expired:   ttl,
		Content:   value,
	}
	return true
}
