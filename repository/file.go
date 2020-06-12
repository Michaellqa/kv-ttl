package repository

import (
	"encoding/json"
	"kv-ttl/kv"
	"os"
)

const fileMode = 0666

type FileRepo struct {
	fileName string
}

func NewFileRepo(fileName string) *FileRepo {
	return &FileRepo{
		fileName: fileName,
	}
}

func (r *FileRepo) RestoreInto(m *map[string]kv.TtlBox) error {
	f, err := os.OpenFile(r.fileName, os.O_RDONLY, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&m)
}

func (r *FileRepo) Save(m map[string]kv.TtlBox) error {
	f, err := os.OpenFile(r.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(m)
}
