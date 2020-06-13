package postgres

import (
	"database/sql"
	"kv-ttl/kv"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (p *Repository) RestoreInto(*map[string]kv.TtlBox) error {
	panic("implement me")
}

func (p *Repository) Save(map[string]kv.TtlBox) error {
	panic("implement me")
}
