package store

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

const (
	qryFindDone = `SELECT id FROM mygrate`
	qryUp       = `INSERT INTO mygrate (id, executed) VALUES (?, ?)`
	qryDown     = `DELETE FROM mygrate WHERE id = ?`
	qryCreate   = `CREATE TABLE IF NOT EXISTS mygrate (
		id VARCHAR(100) NOT NULL,
		executed DATETIME NOT NULL,
		PRIMARY KEY (id)
	)`
)

type SQLStore struct {
	db *sql.DB
	mu sync.Mutex
}

// NewSQLStore returns a neq MySQLStore.
func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

// Init implements mygrate.Store.
func (s *SQLStore) Init() error {
	_, err := s.db.Exec(qryCreate)
	return err
}

// FindDone implements mygrate.Store.
func (s *SQLStore) FindDone() ([]string, error) {
	rows, err := s.db.Query(qryFindDone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var done []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		done = append(done, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return done, nil
}

// Up implements mygrate.Store.
func (s *SQLStore) Up(id string, executed time.Time) error {
	res, err := s.db.Exec(qryUp, id, executed)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// Down implements mygrate.Store.
func (s *SQLStore) Down(id string, executed time.Time) error {
	res, err := s.db.Exec(qryDown, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%s %w", id, ErrIDNotFound)
	}

	return nil
}

// Lock implements mygrate.Locker.
func (s *SQLStore) Lock() error {
	s.mu.Lock()
	return nil
}

// Unlock implements mygrate.Locker.
func (s *SQLStore) Unlock() error {
	s.mu.Unlock()
	return nil
}
