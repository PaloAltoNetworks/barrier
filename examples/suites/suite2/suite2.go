package suite2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/gotest"
)

func init() {

	s := gotest.RegisterSuite(gotest.Suite{
		Name:        "Basic Suite2",
		Description: "Basic Suite2 description",
		Setup: func(ctx context.Context, si gotest.SuiteInfo) (interface{}, gotest.TearDownFunction, error) {

			suiteURL := "http://localhost:8080"

			return suiteURL, func() { fmt.Println("Suite Done") }, nil
		},
	})

	// Demo: Register a test as a part of a suite.
	s.RegisterTest(gotest.Test{
		Name:        "Basic HTTP test",
		Description: "Basic HTTP test",
		Author:      "Satyam",
		Tags:        []string{"suite=sanity", "feature=basic", "test=get"},
		Setup: func(ctx context.Context, ti gotest.TestInfo) (interface{}, gotest.TearDownFunction, error) {

			// Demo: Setup function creates a URL that could be used in the tests.
			url := "http://localhost:3333"

			// Demo: return a tear down function to be executed at end of test.
			return url, func() { fmt.Println("Test Done") }, nil
		},
		Function: func(ctx context.Context, t gotest.TestInfo) error {

			// Demo: Tests have an ID that can be accessed
			_ = t.TestID()

			// Demo: Use suite setup vars in test
			suiteURL := t.SuiteSetupInfo().(string)

			gotest.Step(t, "perform a get on setup", func() error {
				_, _ = http.Get(suiteURL) // make http call but ignore results
				return nil
			})

			// Demo: Use setup vars in test
			url := t.SetupInfo().(string)

			gotest.Step(t, "perform a get", func() error {
				_, _ = http.Get(url) // make http call but ignore results
				return nil
			})

			return nil
		},
	})
}
