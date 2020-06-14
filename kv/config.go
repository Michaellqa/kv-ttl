package kv

import "time"

const DefaultBackupInterval = 5 * time.Second

// Configuration defines set of parameters to configure a cache.
type Configuration struct {
	BackupInterval time.Duration
	Storage        Storage
}
