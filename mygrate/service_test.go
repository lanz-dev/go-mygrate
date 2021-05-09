package mygrate_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lanz-dev/go-mygrate/mygrate"
	"github.com/lanz-dev/go-mygrate/store"
)

var (
	errUnitTest = errors.New("unittest")

	errFunc = func() error {
		return errUnitTest
	}
	nilFunc = func() error {
		return nil
	}
)

func buildMock() *store.MockStore {
	mock := &store.MockStore{}
	mock.InitFunc = nilFunc
	mock.LockFunc = nilFunc
	mock.UnlockFunc = nilFunc
	mock.FindDoneFunc = func() ([]string, error) {
		return nil, nil
	}
	return mock
}

func TestService_InitOnlyOnce(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	m := mygrate.New(mygrate.WithStore(mock))

	_, err := m.Migrate(false)
	if !mock.InitCalled {
		t.Fatal(`expected mock.Init() to be called`)
	}
	if err != nil {
		t.Fatalf(`did not expected to receive err '%s'`, err)
	}

	mock.InitFunc = func() error {
		// InitFunc should not be called again, but if so, we will receive an error now
		return errUnitTest
	}
	if _, err := m.Migrate(false); err != nil {
		t.Fatalf(`did not expected to receive err '%s'`, err)
	}
}

func TestService_Migrate_InitErr(t *testing.T) {
	t.Parallel()

	mock := &store.MockStore{}
	mock.InitFunc = func() error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	_, err := m.Migrate(false)

	if !mock.InitCalled {
		t.Fatal(`expected mock.Init() to be called`)
	}

	if !errors.Is(err, mygrate.ErrInitFn) {
		t.Fatalf(`expected err '%s', got '%s'`, mygrate.ErrInitFn, err)
	}
}

func TestService_Migrate_Unlocks(t *testing.T) {
	t.Parallel()

	mock := &store.MockStore{}
	mock.LockFunc = func() error {
		return nil
	}
	mock.UnlockFunc = func() error {
		return nil
	}
	mock.FindDoneFunc = func() ([]string, error) {
		return nil, nil
	}
	m := mygrate.New(mygrate.WithStore(mock))

	_, err := m.Migrate(false)

	if !mock.LockCalled {
		t.Fatal(`expected mock.Lock() to be called`)
	}

	if !mock.UnlockCalled {
		t.Fatal(`expected mock.Unlock() to be called`)
	}

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
}

func TestService_Migrate_LockError(t *testing.T) {
	t.Parallel()

	mock := &store.MockStore{}
	mock.LockFunc = func() error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	_, err := m.Migrate(false)

	if !mock.LockCalled {
		t.Fatal(`expected mock.Lock() to be called`)
	}

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrStore) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrStore)
	}
}

func TestService_Migrate_FindDoneErr(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return nil, errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	_, err := m.Migrate(false)

	if !mock.FindDoneCalled {
		t.Fatal(`expected mock.FindDone() to be called`)
	}

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrStore) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrStore)
	}
}

func TestService_Migrate_MigrationUpErr(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", errFunc, nilFunc)

	_, err := m.Migrate(false)

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrUpFn) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrUpFn)
	}
}

func TestService_Migrate_StoreUpErr(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.UpFunc = func(id string, executed time.Time) error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nilFunc, nilFunc)

	_, err := m.Migrate(false)

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrStore) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrStore)
	}
}

func TestService_Migrate_RedoLast(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1", "2"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	counter := 0
	m.Register(
		"1",
		func() error {
			counter += 100
			return nil
		},
		nilFunc,
	)
	m.Register(
		"2",
		func() error {
			counter++
			return nil
		},
		nilFunc,
	)

	_, err := m.Migrate(true)

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if counter != 1 {
		t.Fatalf(`expected counter to be '%d', got '%d'`, 1, counter)
	}
}

func TestService_Migrate_RedoLastWithErrUp(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1", "2"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nilFunc, nilFunc)
	m.Register("2", errFunc, nilFunc)

	_, err := m.Migrate(true)

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrUpFn) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrUpFn)
	}
}

func TestService_Migrate_RedoLastWithErrDown(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1", "2"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nilFunc, nilFunc)
	m.Register("2", nilFunc, errFunc)

	_, err := m.Migrate(true)

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrDownFn) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrDownFn)
	}
}

func TestService_Migrate_DetectOpen(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	counter := 0
	m.Register("1", nil, nil)
	m.Register(
		"2",
		func() error {
			counter++
			return nil
		},
		nilFunc,
	)

	_, err := m.Migrate(true)

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if counter != 1 {
		t.Fatalf(`expected counter to be '%d', got '%d'`, 1, counter)
	}
}

func TestService_Migrate_OrderIsByRegister(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"2"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	counter := 0
	m.Register("2", nilFunc, nilFunc)
	m.Register(
		"1",
		func() error {
			counter++
			return nil
		},
		nilFunc,
	)

	_, err := m.Migrate(true)

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if counter != 1 {
		t.Fatalf(`expected counter to be '%d', got '%d'`, 1, counter)
	}
}

func TestService_Migrate(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nilFunc, nilFunc)
	m.Register("2", nilFunc, nilFunc)
	m.Register("3", nilFunc, nilFunc)

	counter, err := m.Migrate(true)
	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if counter != 3 {
		t.Fatalf(`expected counter to be '%d', got '%d'`, 3, counter)
	}
}

func TestService_Rollback_InitErr(t *testing.T) {
	t.Parallel()

	mock := &store.MockStore{}
	mock.InitFunc = func() error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	err := m.Rollback("1")

	if !mock.InitCalled {
		t.Fatal(`expected mock.Init() to be called`)
	}

	if !errors.Is(err, mygrate.ErrInitFn) {
		t.Fatalf(`expected err '%s', got '%s'`, mygrate.ErrInitFn, err)
	}
}

func TestService_Rollback_Unlocks(t *testing.T) {
	t.Parallel()

	mock := &store.MockStore{}
	mock.LockFunc = func() error {
		return nil
	}
	mock.UnlockFunc = func() error {
		return nil
	}
	mock.FindDoneFunc = func() ([]string, error) {
		return nil, nil
	}
	m := mygrate.New(mygrate.WithStore(mock))

	err := m.Rollback("1")

	if !mock.LockCalled {
		t.Fatal(`expected mock.Lock() to be called`)
	}

	if !mock.UnlockCalled {
		t.Fatal(`expected mock.Unlock() to be called`)
	}

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
}

func TestService_Rollback_LockError(t *testing.T) {
	t.Parallel()

	mock := &store.MockStore{}
	mock.LockFunc = func() error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	err := m.Rollback("1")

	if !mock.LockCalled {
		t.Fatal(`expected mock.Lock() to be called`)
	}

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrStore) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrStore)
	}
}

func TestService_Rollback_FindDoneErr(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return nil, errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	err := m.Rollback("1")

	if !mock.FindDoneCalled {
		t.Fatal(`expected mock.FindDone() to be called`)
	}

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrStore) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrStore)
	}
}

func TestService_Rollback_MigrationDownErr(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1", "2"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}
	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nilFunc, nilFunc)
	m.Register("2", nilFunc, errFunc)

	err := m.Rollback("1")

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrDownFn) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrDownFn)
	}
}

func TestService_Rollback_StoreDownErr(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1", "2"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return errUnitTest
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}
	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nilFunc, nilFunc)
	m.Register("2", nilFunc, nilFunc)

	err := m.Rollback("1")

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrStore) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrStore)
	}
}

func TestService_Rollback_ThreeEntries(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1", "2", "3"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	counter := 3
	m.Register(
		"1",
		nilFunc,
		func() error {
			counter--
			return nil
		},
	)
	m.Register(
		"2",
		nilFunc,
		func() error {
			counter--
			return nil
		},
	)
	m.Register(
		"3",
		nilFunc,
		func() error {
			counter--
			return nil
		},
	)

	err := m.Rollback("1")

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if counter != 0 {
		t.Fatalf(`expected counter to be '%d', got '%d'`, 0, counter)
	}
}

func TestService_Rollback_TwoEntries(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1", "2", "3"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	counter := 3
	m.Register(
		"1",
		nilFunc,
		func() error {
			counter--
			return nil
		},
	)
	m.Register(
		"2",
		nilFunc,
		func() error {
			counter--
			return nil
		},
	)
	m.Register(
		"3",
		nilFunc,
		func() error {
			counter--
			return nil
		},
	)

	err := m.Rollback("2")

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if counter != 1 {
		t.Fatalf(`expected counter to be '%d', got '%d'`, 1, counter)
	}
}

func TestService_Rollback_OneEntry(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	counter := 3
	m.Register(
		"1",
		nilFunc,
		func() error {
			counter--
			return nil
		},
	)

	err := m.Rollback("1")

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if counter != 2 {
		t.Fatalf(`expected counter to be '%d', got '%d'`, 2, counter)
	}
}

func TestService_Rollback_NotAllRegisteredAreDone(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.FindDoneFunc = func() ([]string, error) {
		return []string{"1", "2"}, nil
	}
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}

	m := mygrate.New(mygrate.WithStore(mock))

	counter := 3
	m.Register(
		"1",
		nil,
		func() error {
			counter--
			return nil
		},
	)
	m.Register(
		"2",
		nil,
		func() error {
			counter--
			return nil
		},
	)
	m.Register(
		"3",
		nil,
		func() error {
			counter--
			return nil
		},
	)

	err := m.Rollback("1")

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if counter != 1 {
		t.Fatalf(`expected counter to be '%d', got '%d'`, 1, counter)
	}
}

func TestService_Reset_InitErr(t *testing.T) {
	t.Parallel()

	mock := &store.MockStore{}
	mock.InitFunc = func() error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	err := m.Reset()

	if !mock.InitCalled {
		t.Fatal(`expected mock.Init() to be called`)
	}

	if !errors.Is(err, mygrate.ErrInitFn) {
		t.Fatalf(`expected err '%s', got '%s'`, mygrate.ErrInitFn, err)
	}
}

func TestService_Reset_WithNoMigrationsRegistered(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	m := mygrate.New(mygrate.WithStore(mock))

	err := m.Reset()

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
	if mock.LockCalled {
		t.Fatal(`did not expected mock.Lock() to be called`)
	}
}

func TestService_Reset_Error(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.LockFunc = func() error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nil, nil)

	err := m.Reset()

	if !mock.LockCalled {
		t.Fatal(`expected mock.Lock() to be called`)
	}

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrStore) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrStore)
	}
}

func TestService_Refresh_InitErr(t *testing.T) {
	t.Parallel()

	mock := &store.MockStore{}
	mock.InitFunc = func() error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	err := m.Refresh()

	if !mock.InitCalled {
		t.Fatal(`expected mock.Init() to be called`)
	}

	if !errors.Is(err, mygrate.ErrInitFn) {
		t.Fatalf(`expected err '%s', got '%s'`, mygrate.ErrInitFn, err)
	}
}

func TestService_Refresh_ResetError(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.LockFunc = func() error {
		return errUnitTest
	}
	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nilFunc, nilFunc)

	err := m.Refresh()

	if !mock.LockCalled {
		t.Fatal(`expected mock.Lock() to be called`)
	}

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrStore) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrStore)
	}
}

func TestService_Refresh_MigrateError(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", errFunc, nilFunc)

	err := m.Refresh()

	if !errors.Is(err, errUnitTest) {
		t.Fatalf(`expected err to be '%s'`, errUnitTest)
	}

	if !errors.Is(err, mygrate.ErrUpFn) {
		t.Fatalf(`expected err to be '%s'`, mygrate.ErrUpFn)
	}
}

func TestService_Refresh(t *testing.T) {
	t.Parallel()

	mock := buildMock()
	mock.DownFunc = func(id string, executed time.Time) error {
		return nil
	}
	mock.UpFunc = func(id string, executed time.Time) error {
		return nil
	}
	m := mygrate.New(mygrate.WithStore(mock))

	m.Register("1", nilFunc, nilFunc)

	err := m.Refresh()

	if err != nil {
		t.Fatalf(`did not expected err '%s'`, err)
	}
}
