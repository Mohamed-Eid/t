package cmd

import (
	"fmt"

	"t/internal/runner"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     ":stop <task-name-or-pid>",
	Aliases: []string{":kill", ":terminate", ":s"},
	Short:   "Stop a running detached task",
	Long:    "Stop a detached task by task name or process ID (PID).",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		// Load config (we need a runner instance)
		config, err := runner.LoadConfig("tasks.yaml")
		if err != nil {
			// For stopping processes, we don't strictly need a valid config
			config = &runner.Config{} // Empty config
		}

		taskRunner := runner.NewRunner(config)

		// Stop the detached process
		err = taskRunner.StopDetachedProcess(identifier)
		if err != nil {
			fmt.Printf("❌ Error stopping process: %v\n", err)
			fmt.Println("\n💡 Use 't :ps' to see running detached tasks")
			return
		}

		// Success message is printed in StopDetachedProcess
	},
}
