package kv

import "time"

// Configuration defines set of parameters to configure a cache.
type Configuration struct {
	BackupInterval time.Duration
	Storage        Storage
}
