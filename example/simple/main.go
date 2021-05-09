package main

import (
	"github.com/lanz-dev/go-mygrate/example/simple/migrations"
	"github.com/lanz-dev/go-mygrate/mygrate"
)

func main() {
	myg := mygrate.New( // New will use a Default FileStore with the Path .mygrate
	// custom Store:
	// mygrate.WithStore(store.NewFileStoreWithPath(".migration")),
	)

	migrations.Register(myg)

	const redoLast = true // redoLast migration
	if _, err := myg.Migrate(redoLast); err != nil {
		panic(err)
	}
}
