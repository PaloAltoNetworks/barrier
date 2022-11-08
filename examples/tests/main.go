package main

import (
	"log"

	"github.com/PaloAltoNetworks/gotest"
	_ "github.com/PaloAltoNetworks/gotest/examples/tests/collection"
)

func main() {
	cmd := gotest.NewCommand("tests", "integration tests", "v0.0.0", nil)
	// Run the command.
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
