package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"kv-ttl/kv"
	"log"
)

// Repository implements the kv.Storage interface and provides storing cache values
// in a Postgres table. Each key-value pair mapped to a row in the table. The key is
// used as a PK, the value is stored as a JSONB type.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// RestoreInto reads rows from the database table and populates the given map.
func (p *Repository) RestoreInto(m *map[string]kv.TtlBox) error {
	rows, err := p.db.Query(`select id, json_value from cache_snapshot`)
	if err != nil {
		return err
	}
	var (
		key   string
		value jsonValue
	)
	for rows.Next() {
		if err = rows.Scan(&key, &value); err != nil {
			continue
		}
		box := kv.TtlBox(value)
		(*m)[key] = box
	}
	return nil
}

// Save clears the table and then inserts new values.
func (p *Repository) Save(m map[string]kv.TtlBox) error {
	_, err := p.db.Exec(`delete from cache_snapshot where true`)
	if err != nil {
		return err
	}
	if len(m) == 0 {
		return nil
	}
	stmt, err := p.db.Prepare(`insert into cache_snapshot values ($1, $2)`)
	if err != nil {
		return err
	}
	for k, v := range m {
		_, err = stmt.Exec(k, jsonValue(v))
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

// For reasons that I don't know the code below doesn't work.
// Origins are taken from https://godoc.org/github.com/lib/pq#hdr-Bulk_imports.
//
// Save clears the table and then inserts new values.
// Insert performed by pq.CopyIn() inside a single transaction for efficient writes.
//func (p *Repository) Save(m map[string]kv.TtlBox) error {
//	_, err := p.db.Exec(`delete from cache_snapshot where true`)
//	if err != nil {
//		return err
//	}
//	if len(m) == 0 {
//		return nil
//	}
//	tx, err := p.db.Begin()
//	if err != nil {
//		return err
//	}
//	txCompleted := false
//	defer func() {
//		if !txCompleted {
//			_ = tx.Rollback()
//		}
//	}()
//	stmt, err := tx.Prepare(pq.CopyIn("cache_snapshot", "id", "json_value"))
//	if err != nil {
//		return err
//	}
//	for k, v := range m {
//		_, err = stmt.Exec(k, jsonValue(v))
//		if err != nil {
//			return err
//		}
//	}
//	if err = tx.Commit(); err != nil {
//		return err
//	}
//	txCompleted = true
//	return nil
//}

// Helper structure used to serialize into and deserialize original values
// from Postgres JSONB type
type jsonValue kv.TtlBox

// Values makes the jsonValue type implement the driver.Valuer interface
func (v jsonValue) Value() (driver.Value, error) {
	return json.Marshal(v)
}

// Scan makes the jsonValue type implement the sql.Scanner interface
func (v *jsonValue) Scan(value interface{}) error {
	bts, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bts, &v)
}
