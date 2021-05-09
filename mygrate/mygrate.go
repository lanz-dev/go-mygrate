// Package mygrate provides the programmatically interface for Mμgrate.
package mygrate

import (
	"time"
)

// Registerer provides methods to register a migration.
// Properly you will use this interface (or your own) to pass an instance
// of Mygrate around and just allow to Register new migrations with it.
type Registerer interface {
	// Register a migration.
	Register(string, func() error, func() error)
}

// Locker provides method which lock databases, filesystems, etc.
type Locker interface {
	// Lock will be called before migrating.
	Lock() error
	// Unlock will be called after migrating.
	Unlock() error
}

// Store provides method to save the current state of migrations.
type Store interface {
	// Init should prepare the store for usage.
	Init() error

	// Up will be called after the migrations up func was run.
	Up(id string, executed time.Time) error
	// Down will be called after the migrations down func was run.
	Down(id string, executed time.Time) error

	// FindDone returns the IDs from already ran migrations.
	FindDone() ([]string, error)
}

type mygration struct {
	ID   string
	Up   func() error
	Down func() error
}
