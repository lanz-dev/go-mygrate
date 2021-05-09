package store

import (
	"time"
)

// MockStore is a mock implementation of mygrate.Store and mygrate.Locker.
type MockStore struct {
	// InitCalled tracks if the method was called.
	InitCalled bool
	// InitFunc mocks the Init method.
	InitFunc func() error

	// UpCalled tracks if the method was called.
	UpCalled bool
	// UpFunc mocks the Up method.
	UpFunc func(id string, executed time.Time) error

	// DownCalled tracks if the method was called.
	DownCalled bool
	// DownFunc mocks the Down method.
	DownFunc func(id string, executed time.Time) error

	// LockCalled tracks if the method was called.
	LockCalled bool
	// LockFunc mocks the Lock method.
	LockFunc func() error

	// UnlockCalled tracks if the method was called.
	UnlockCalled bool
	// UnlockFunc mocks the Unlock method.
	UnlockFunc func() error

	// FindDoneCalled tracks if the method was called.
	FindDoneCalled bool
	// FindDoneFunc mocks the FindDone method.
	FindDoneFunc func() ([]string, error)
}

// Init calls InitFunc.
func (mock *MockStore) Init() error {
	if mock.InitFunc == nil {
		return nil
	}
	mock.InitCalled = true
	return mock.InitFunc()
}

// Up calls UpFunc.
func (mock *MockStore) Up(id string, executed time.Time) error {
	if mock.UpFunc == nil {
		panic("MockStore.UpFunc: method is nil but Store.Up was just called")
	}
	mock.UpCalled = true
	return mock.UpFunc(id, executed)
}

// Down calls DownFunc.
func (mock *MockStore) Down(id string, executed time.Time) error {
	if mock.DownFunc == nil {
		panic("MockStore.DownFunc: method is nil but Store.Down was just called")
	}
	mock.DownCalled = true
	return mock.DownFunc(id, executed)
}

// Lock calls LockFunc.
func (mock *MockStore) Lock() error {
	if mock.LockFunc == nil {
		panic("MockStore.LockFunc: method is nil but Store.Lock was just called")
	}
	mock.LockCalled = true
	return mock.LockFunc()
}

// Unlock calls UnlockFunc.
func (mock *MockStore) Unlock() error {
	if mock.UnlockFunc == nil {
		panic("MockStore.UnlockFunc: method is nil but Store.Unlock was just called")
	}
	mock.UnlockCalled = true
	return mock.UnlockFunc()
}

// FindDone calls FindDoneFunc.
func (mock *MockStore) FindDone() ([]string, error) {
	if mock.FindDoneFunc == nil {
		panic("MockStore.FindDoneFunc: method is nil but Store.FindDone was just called")
	}
	mock.FindDoneCalled = true
	return mock.FindDoneFunc()
}
