package collection

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/gotest"
)

func init() {
	gotest.RegisterTest(gotest.Test{
		Name:        "Basic HTTP test",
		Description: "Basic HTTP test",
		Author:      "Satyam",
		Tags:        []string{"suite=sanity", "feature=basic", "test=get"},
		Setup: func(ctx context.Context, ti gotest.TestInfo) (interface{}, gotest.TearDownFunction, error) {

			// Demo: Setup function creates a URL that could be used in the tests.
			url := "http://localhost:3333"

			// Demo: return a tear down function to be executed at end of test.
			return url, func() { fmt.Println("Done") }, nil
		},
		Function: func(ctx context.Context, t gotest.TestInfo) error {

			// Demo: Tests have an ID that can be accessed
			_ = t.TestID()

			// Demo: Use setup vars etc in test
			url := t.SetupInfo().(string)

			gotest.Step(t, "perform a get", func() error {
				_, _ = http.Get(url) // make http call but ignore results
				return nil
			})

			return nil
		},
	})
}
