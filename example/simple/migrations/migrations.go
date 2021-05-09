package migrations

import (
	"github.com/lanz-dev/go-mygrate/mygrate"
)

func Register(r mygrate.Registerer) {
	r.Register("init_create", initCreateUp, initCreateDown)
	r.Register("do_something", doSomethingUp, doSomethingDown)
}
