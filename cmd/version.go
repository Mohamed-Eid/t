package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version will be set during build
	Version = "dev"
	// BuildDate will be set during build
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   ":version",
	Short: "Print the version information",
	Long:  "Display version and build information for the t task runner.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("t task runner\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Printf("Author: Mohamed Eid\n")
	},
}
