package mygrate

import (
	"time"

	"github.com/lanz-dev/go-mygrate/store"
)

// Service provides methods for Mμgrate.
type Service struct {
	initDone   bool
	migrations []mygration
	store      Store
}

// New will create a new Service instance with a default FileStore.
func New(opts ...Option) *Service {
	s := &Service{
		store: store.NewFileStore(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) up(myg mygration) error {
	if err := myg.Up(); err != nil {
		return errUp(myg.ID, err)
	}

	if err := s.store.Up(myg.ID, time.Now().UTC()); err != nil {
		return errStore(myg.ID, err)
	}

	return nil
}

func (s *Service) down(myg mygration) error {
	if err := myg.Down(); err != nil {
		return errDown(myg.ID, err)
	}

	if err := s.store.Down(myg.ID, time.Now().UTC()); err != nil {
		return errStore(myg.ID, err)
	}

	return nil
}

func (s *Service) redo(key int) error {
	if err := s.down(s.migrations[key]); err != nil {
		return err
	}

	if err := s.up(s.migrations[key]); err != nil {
		return err
	}

	return nil
}

func (s *Service) findOpen() ([]mygration, error) {
	doneIDs, err := s.store.FindDone()
	if err != nil {
		return nil, errStore("", err)
	}

	var todo []mygration
	for _, s1 := range s.migrations {
		found := false
		for _, ID := range doneIDs {
			if s1.ID == ID {
				found = true
				break
			}
		}
		if !found {
			todo = append(todo, s1)
		}
	}

	return todo, nil
}

func (s *Service) findRevert(targetID string) ([]mygration, error) {
	doneIDs, err := s.store.FindDone()
	if err != nil {
		return nil, errStore("", err)
	}

	reversed := make([]mygration, len(s.migrations))
	copy(reversed, s.migrations)
	for i, j := 0, len(s.migrations)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	var revert []mygration
	for _, s1 := range reversed {
		for _, ID := range doneIDs {
			if s1.ID == ID {
				revert = append(revert, s1)
				break
			}
		}
		if s1.ID == targetID {
			break
		}
	}

	return revert, nil
}

func (s *Service) init() error {
	if s.initDone {
		return nil
	}

	if err := s.store.Init(); err != nil {
		return errInit(err)
	}

	s.initDone = true

	return nil
}

// Migrate will execute all outstanding migrations.
func (s *Service) Migrate(redoLast bool) (int, error) {
	if err := s.init(); err != nil {
		return 0, err
	}

	locker, ok := s.store.(Locker)
	if ok {
		if err := locker.Lock(); err != nil {
			return 0, errStore("", err)
		}
		defer locker.Unlock()
	}

	todo, err := s.findOpen()
	if err != nil {
		return 0, err
	}

	for _, myg := range todo {
		if err := s.up(myg); err != nil {
			return 0, err
		}
	}

	changes := len(todo)
	if changes == 0 && redoLast && len(s.migrations) >= 1 {
		lastKey := len(s.migrations) - 1
		if err := s.redo(lastKey); err != nil {
			return 0, err
		}
	}

	return changes, nil
}

// Rollback will rollback migrations to (including) the given id.
func (s *Service) Rollback(id string) error {
	if err := s.init(); err != nil {
		return err
	}

	locker, ok := s.store.(Locker)
	if ok {
		if err := locker.Lock(); err != nil {
			return errStore("", err)
		}
		defer locker.Unlock()
	}

	todo, err := s.findRevert(id)
	if err != nil {
		return err
	}

	for _, myg := range todo {
		if err := s.down(myg); err != nil {
			return err
		}
	}

	return nil
}

// Reset will rollback all migrations.
func (s *Service) Reset() error {
	if err := s.init(); err != nil {
		return err
	}

	if len(s.migrations) == 0 {
		return nil
	}

	if err := s.Rollback(s.migrations[0].ID); err != nil {
		return err
	}

	return nil
}

// Refresh will rollback all migrations and execute them again.
func (s *Service) Refresh() error {
	if err := s.init(); err != nil {
		return err
	}

	if err := s.Reset(); err != nil {
		return err
	}

	if _, err := s.Migrate(false); err != nil {
		return err
	}

	return nil
}

// Register will register a migration.
func (s *Service) Register(id string, up func() error, down func() error) {
	s.migrations = append(s.migrations, mygration{
		ID:   id,
		Up:   up,
		Down: down,
	})
}
