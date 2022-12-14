package main

import (
	"log"

	"github.com/PaloAltoNetworks/barrier"
	_ "github.com/PaloAltoNetworks/barrier/examples/tests/collection"
)

func main() {
	cmd := barrier.NewCommand("tests", "integration tests", "v0.0.0", nil)
	// Run the command.
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
