package kv

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
