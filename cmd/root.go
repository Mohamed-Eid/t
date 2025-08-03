/*
Copyright © 2025 Mohamed Eid <medoeid50@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"t/internal/runner"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "t",
	Short: "A simple and efficient task runner for your projects",
	Long: `t is a lightweight task runner similar to Make but with YAML configuration.
It allows you to define and execute tasks with dependencies, variables, and commands.

Examples:
  t :init         Initialize a new tasks.yaml file
  t :list         List all available tasks
  t :version      Show version information
  t build         Run the build task
  t test          Run the test task
  t <task-name>   Run any task defined in tasks.yaml

Note: Tool commands start with ':' to avoid conflicts with user-defined tasks.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Show help when no task is specified
			cmd.Help()
			return
		}

		taskName := args[0]

		// Load config and run task
		config, err := runner.LoadConfig("tasks.yaml")
		if err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			fmt.Println("\n💡 Tip: Run 't :init' to create a tasks.yaml file")
			os.Exit(1)
		}

		taskRunner := runner.NewRunner(config)

		if err := taskRunner.RunTask(taskName); err != nil {
			fmt.Printf("❌ Task failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🎉 Task '%s' completed successfully!\n", taskName)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}
