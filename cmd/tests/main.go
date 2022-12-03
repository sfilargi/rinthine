package main

import (
	"fmt"

	"github.com/rinthine/pkg/coreops"
)

func main() {
	a, b, c := coreops.Authorize("stavros_app", "stavros", "secret")

	fmt.Printf("%+v %+v %+v\n", a, b, c)
}
