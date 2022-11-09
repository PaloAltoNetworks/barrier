package main

import (
	"context"
	"log"

	"github.com/PaloAltoNetworks/barrier"
	_ "github.com/PaloAltoNetworks/barrier/examples/suites/suite1"
	_ "github.com/PaloAltoNetworks/barrier/examples/suites/suite2"
)

func stash(ctx context.Context) interface{} {
	return "stash-string"
}

func main() {
	cmd := barrier.NewCommand("tests", "integration tests", "v0.0.0", stash)
	// Run the command.
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
