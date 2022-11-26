package main

import (
	"fmt"

	"github.com/rinthine/pkg/coreops"
)

func main() {

	app := coreops.App{
		Name:        "testapp",
		User:        "stavros",
		Description: "",
		HomeUrl:     "",
		RedirectUrl: "",
		Password:    coreops.RandomString(36),
	}

	err := coreops.CreateApp(&app)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Printf("OK")
	}
}
