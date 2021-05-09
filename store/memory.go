package store

import (
	"fmt"
	"sync"
	"time"
)

// MemoryStore store the migration state in a map.
type MemoryStore struct {
	migrations map[string]time.Time
	mu         sync.Mutex
}

// NewMemoryStore will return a MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		migrations: map[string]time.Time{},
	}
}

// Init implements mygrate.Store.
func (m *MemoryStore) Init() error {
	return nil
}

// FindDone implements mygrate.Store.
func (m *MemoryStore) FindDone() ([]string, error) {
	done := make([]string, 0, len(m.migrations))
	for ID := range m.migrations {
		done = append(done, ID)
	}
	return done, nil
}

// Up implements mygrate.Store.
func (m *MemoryStore) Up(id string, executed time.Time) error {
	m.migrations[id] = executed
	return nil
}

// Down implements mygrate.Store.
func (m *MemoryStore) Down(id string, executed time.Time) error {
	if _, ok := m.migrations[id]; !ok {
		return fmt.Errorf("%s %w", id, ErrIDNotFound)
	}
	delete(m.migrations, id)
	return nil
}

// Lock implements mygrate.Locker.
func (m *MemoryStore) Lock() error {
	m.mu.Lock()
	return nil
}

// Unlock implements mygrate.Locker.
func (m *MemoryStore) Unlock() error {
	m.mu.Unlock()
	return nil
}
