package kv

// Storage defines the interface that should be implemented to load initial values
// into cache and periodically update snapshots of the cache data.
type Storage interface {
	RestoreInto(*map[string]TtlBox) error
	Save(map[string]TtlBox) error
}

type NotImplementedStorage struct{}

func (s *NotImplementedStorage) RestoreInto(*map[string]TtlBox) error {
	return nil
}

func (s *NotImplementedStorage) Save(map[string]TtlBox) error {
	return nil
}
