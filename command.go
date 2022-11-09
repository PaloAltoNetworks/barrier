package barrier

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand generates a new CLI
func NewCommand(
	name string,
	description string,
	version string,
	stash GenerateStash,
) *cobra.Command {

	cobra.OnInitialize(func() {
		viper.SetEnvPrefix(name)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	})

	var rootCmd = &cobra.Command{
		Use:   name,
		Short: description,
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version and exit.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	var cmdListTests = &cobra.Command{
		Use:           "list",
		Aliases:       []string{"ls"},
		Short:         "List registered tests.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			suites := filterSuites()
			return listTests(suites)
		},
	}

	cmdListTests.Flags().StringSliceP("id", "i", nil, "Only run tests with the given identifier")
	cmdListTests.Flags().StringSliceP("tag", "t", nil, "Only run tests with the given tags")

	var cmdRunTests = &cobra.Command{
		Use:           "test",
		Aliases:       []string{"run"},
		Short:         "Run the registered tests",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			// TODO: add argument check.
			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("limit"))
			defer cancel()

			s := stash(ctx)

			suites := filterSuites()

			for _, suite := range suites {

				if stash != nil {
					// Store caller stash per suite
					suite.stash = s
				}

				err := newTestRunner(
					viper.GetString("build-id"),
					suite,
					viper.GetDuration("limit"),
					viper.GetInt("concurrent"),
					viper.GetInt("stress"),
					viper.GetBool("verbose"),
					viper.GetBool("skip-teardown"),
					viper.GetBool("stop-on-failure"),
				).run(ctx, suite)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	// Parameters to configure suite behaviors
	cmdRunTests.Flags().StringSliceP("suite", "Z", nil, "Only run suites specified")

	// Parameters to configure test behaviors
	cmdRunTests.Flags().BoolP("verbose", "V", false, "Show logs even on success")
	cmdRunTests.Flags().DurationP("limit", "l", 20*time.Minute, "Execution time limit")
	cmdRunTests.Flags().IntP("concurrent", "c", 20, "Max number of concurrent tests")
	cmdRunTests.Flags().IntP("stress", "s", 1, "Number of times to run each test in parallel")
	cmdRunTests.Flags().StringSliceP("id", "i", nil, "Only run tests with the given identifier")
	cmdRunTests.Flags().StringSliceP("tag", "t", nil, "Only run tests with the given tags")
	cmdRunTests.Flags().BoolP("match-all", "M", false, "Match all tags specified")
	cmdRunTests.Flags().BoolP("skip-teardown", "S", false, "Skip teardown step")
	cmdRunTests.Flags().BoolP("stop-on-failure", "X", false, "Stop on the first failed test")

	rootCmd.AddCommand(
		versionCmd,
		cmdListTests,
		cmdRunTests,
	)

	return rootCmd
}

// runSuite returns true if we should consider the suite for running
func runSuite(s *suiteInfo, names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, name := range names {
		if name == s.Name {
			return true
		}
	}
	return false
}

// filterSuites filters the suite based on ids and/or tags
func filterSuites() []*suiteInfo {
	s := []*suiteInfo{}

	names := viper.GetStringSlice("suite")
	for _, suite := range mainSuites.sorted() {

		// Filter Suites
		if !runSuite(suite, names) {
			continue
		}

		// Filter Tests in a suite
		ids := viper.GetStringSlice("id")
		if len(ids) > 0 {
			suite = suite.testsWithIDs(viper.GetBool("verbose"), ids)
		} else {
			tags := viper.GetStringSlice("tag")
			if len(tags) > 0 {
				suite = suite.testsWithArgs(viper.GetBool("verbose"), viper.GetBool("match-all"), tags)
			}
		}
		if len(suite.tests) > 0 {
			s = append(s, suite)
		}
	}
	return s
}
