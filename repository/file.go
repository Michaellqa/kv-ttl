package repository

import (
	"encoding/json"
	"kv-ttl/kv"
	"os"
)

const fileMode = 0666

// Repository implements the kv.Storage interface and provides storing
// cache values as a json file.
type FileRepo struct {
	fileName string
}

func NewFileRepo(fileName string) *FileRepo {
	return &FileRepo{
		fileName: fileName,
	}
}

// RestoreInto reads from the file and tries to parse json content into the map.
func (r *FileRepo) RestoreInto(m *map[string]kv.TtlBox) error {
	f, err := os.OpenFile(r.fileName, os.O_RDONLY, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&m)
}

// Save create a file if needed and dumps json representation of the map into it.
func (r *FileRepo) Save(m map[string]kv.TtlBox) error {
	f, err := os.OpenFile(r.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(m)
}
