package migrations

import (
	"fmt"
)

func doSomethingUp() error {
	fmt.Println("doSomethingUp()")
	return nil
}

func doSomethingDown() error {
	fmt.Println("doSomethingDown()")
	return nil
}
