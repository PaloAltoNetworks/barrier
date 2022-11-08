package collection

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/barrier"
)

func init() {
	barrier.RegisterTest(barrier.Test{
		Name:        "Basic HTTP test",
		Description: "Basic HTTP test",
		Author:      "Satyam",
		Tags:        []string{"suite=sanity", "feature=basic", "test=get"},
		Setup: func(ctx context.Context, ti barrier.TestInfo) (interface{}, barrier.TearDownFunction, error) {

			// Demo: Setup function creates a URL that could be used in the tests.
			url := "http://localhost:3333"

			// Demo: return a tear down function to be executed at end of test.
			return url, func() { fmt.Println("Done") }, nil
		},
		Function: func(ctx context.Context, t barrier.TestInfo) error {

			// Demo: Tests have an ID that can be accessed
			_ = t.TestID()

			// Demo: Use setup vars etc in test
			url := t.SetupInfo().(string)

			barrier.Step(t, "perform a get", func() error {
				_, _ = http.Get(url) // make http call but ignore results
				return nil
			})

			return nil
		},
	})
}
