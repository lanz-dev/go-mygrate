package migrations

import (
	"fmt"
)

func initCreateUp() error {
	fmt.Println("initCreateUp()")
	return nil
}

func initCreateDown() error {
	fmt.Println("initCreateDown()")
	return nil
}
