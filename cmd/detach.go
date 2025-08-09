package cmd

import (
	"fmt"

	"t/internal/runner"

	"github.com/spf13/cobra"
)

var detachCmd = &cobra.Command{
	Use:     ":detach <task-name>",
	Aliases: []string{":d", ":bg", ":background"},
	Short:   "Run task in background (detached mode)",
	Long:    "Start a task in the background and return immediately. Perfect for long-running processes like development servers.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]

		// Load config
		config, err := runner.LoadConfig("tasks.yaml")
		if err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			fmt.Println("\n💡 Tip: Run 't :init' to create a tasks.yaml file")
			return
		}

		taskRunner := runner.NewRunner(config)

		// Run task in detached mode
		detachedProc, err := taskRunner.RunTaskDetached(taskName)
		if err != nil {
			fmt.Printf("❌ Failed to start detached task: %v\n", err)
			return
		}

		// Show success message (already printed in RunTaskDetached)
		_ = detachedProc
	},
}
