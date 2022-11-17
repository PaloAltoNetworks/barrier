package barrier

import (
	"fmt"

	"github.com/spf13/viper"
)

func listTests(suites []*suiteInfo) error {

	onlySuites := viper.GetBool("only-suites")
	for _, suite := range suites {
		if onlySuites {
			fmt.Printf("%s\n", suite)
		} else {
			suite.listTests()
		}
	}

	return nil
}
