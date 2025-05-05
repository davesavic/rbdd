/*
Copyright © 2025 Dave Savic
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version details",
	Long:  `Display the version of the application, including the build date and commit hash.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\n", rootCmd.Use, Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("OS/Arch: %s/%s\n", Os, Arch)
		fmt.Printf("Built at: %s\n", Date)
		fmt.Printf("Built by: %s\n", BuiltBy)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
