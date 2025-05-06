/*
Copyright © 2025 Dave Savic
*/

package cmd

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/davesavic/rbdd/app"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the cucumber tests",
	Long:  `Run the cucumber tests using the gherkin syntax.`,
	Run: func(cmd *cobra.Command, args []string) {
		directories, err := cmd.Flags().GetStringSlice("directories")
		if err != nil || len(directories) == 0 {
			directories = []string{"features"}
		}

		suite := godog.TestSuite{
			Name:                 "rbdd",
			TestSuiteInitializer: app.InitializeTestSuite,
			Options: &godog.Options{
				Format: "pretty",
				Paths:  directories,
			},
		}

		if suite.Run() != 0 {
			fmt.Println("Test suite failed")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceP("directories", "d", []string{"features"}, "Directories to run the tests in")
}
