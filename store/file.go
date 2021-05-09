package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type entry struct {
	ID       string    `json:"id"`
	Executed time.Time `json:"executed"`
}

// FileStore store the migration state in a json based file.
type FileStore struct {
	path       string
	Migrations []entry `json:"migrations"`
	mu         sync.Mutex
}

// NewFileStoreWithPath will return a FileStore with a custom path.
func NewFileStoreWithPath(path string) *FileStore {
	return &FileStore{path: path}
}

// NewFileStore will return a FileStore with the default path ".mygrate".
func NewFileStore() *FileStore {
	return NewFileStoreWithPath(".mygrate")
}

func (f *FileStore) save() error {
	buf, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(f.path, buf, 0600); err != nil {
		return err
	}

	return nil
}

// Init implements mygrate.Store.
func (f *FileStore) Init() error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0600); err != nil {
		return err
	}

	buf, err := os.ReadFile(f.path)
	if err != nil {
		return nil
	}

	if err := json.Unmarshal(buf, f); err != nil {
		return err
	}

	return nil
}

// FindDone implements mygrate.Store.
func (f *FileStore) FindDone() ([]string, error) {
	done := make([]string, 0, len(f.Migrations))
	for _, v := range f.Migrations {
		done = append(done, v.ID)
	}
	return done, nil
}

// Up implements mygrate.Store.
func (f *FileStore) Up(id string, executed time.Time) error {
	f.Migrations = append(f.Migrations, entry{
		ID:       id,
		Executed: executed,
	})
	return f.save()
}

// Down implements mygrate.Store.
func (f *FileStore) Down(id string, executed time.Time) error {
	index := -1
	for i, e := range f.Migrations {
		if e.ID == id {
			index = i
			break
		}
	}

	if index < 0 {
		return fmt.Errorf("%s %w", id, ErrIDNotFound)
	}

	f.Migrations = append(f.Migrations[:index], f.Migrations[index+1:]...)

	return f.save()
}

// Lock implements mygrate.Locker.
func (f *FileStore) Lock() error {
	f.mu.Lock()
	return nil
}

// Unlock implements mygrate.Locker.
func (f *FileStore) Unlock() error {
	f.mu.Unlock()
	return nil
}
