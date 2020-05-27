package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

var (
	shortened  = false
	version    = "dev"
	commit     = "none"
	date       = "unknown"
	output     = "json"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
			fmt.Print(resp)
		},
	}
)

func init() {
	versionCmd.Flags().BoolVar(&shortened, "short", false, "Print just the version number.")
	versionCmd.Flags().StringVar(&output, "output", "json", "Output format. One of 'yaml' or 'json'.")
	rootCmd.AddCommand(versionCmd)
}
