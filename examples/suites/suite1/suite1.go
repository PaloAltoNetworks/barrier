package suite1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/barrier"
)

func init() {

	s := barrier.RegisterSuite(barrier.Suite{
		Name:        "Basic Suite1",
		Description: "Basic Suite1 description",
		Setup: func(ctx context.Context, si barrier.SuiteInfo) (interface{}, barrier.TearDownFunction, error) {

			suiteURL := "http://localhost:8080"

			return suiteURL, func() { fmt.Println("Suite Done") }, nil
		},
	})

	// Demo: Register a test as a part of a suite.
	s.RegisterTest(barrier.Test{
		Name:        "Basic HTTP test",
		Description: "Basic HTTP test",
		Author:      "Satyam",
		Tags:        []string{"suite=sanity", "feature=basic", "test=get"},
		Setup: func(ctx context.Context, ti barrier.TestInfo) (interface{}, barrier.TearDownFunction, error) {

			// Demo: Setup function creates a URL that could be used in the tests.
			url := "http://localhost:3333"

			// Demo: return a tear down function to be executed at end of test.
			return url, func() { fmt.Println("Test Done") }, nil
		},
		Function: func(ctx context.Context, t barrier.TestInfo) error {

			// Demo: Tests have an ID that can be accessed
			_ = t.TestID()

			// Demo: Use suite setup vars in test
			suiteURL := t.SuiteSetupInfo().(string)

			// Demo: Usage of stash
			fmt.Println("stash is: " + t.Stash().(string))

			barrier.Step(t, "perform a get on setup", func() error {
				_, _ = http.Get(suiteURL) // make http call but ignore results
				return nil
			})

			// Demo: Use setup vars in test
			url := t.SetupInfo().(string)

			barrier.Step(t, "perform a get", func() error {
				_, _ = http.Get(url) // make http call but ignore results
				return nil
			})

			return nil
		},
	})
}
