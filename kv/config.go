package kv

import "time"

type Configuration struct {
	BackupInterval time.Duration
	Storage        Storage
}

var DefaultConfig = Configuration{
	BackupInterval: time.Duration(time.Second),
	Storage:        &NotImplementedStorage{},
}

type configuration struct {
	backupInterval time.Duration
	storage        Storage
}
