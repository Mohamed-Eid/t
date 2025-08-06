package cmd

import (
	"fmt"
	"time"

	"t/internal/runner"

	"github.com/spf13/cobra"
)

var parallelCmd = &cobra.Command{
	Use:   ":parallel <task-name>",
	Short: "Run task with timing information to show parallel execution",
	Long:  "Execute a task and show timing information to demonstrate parallel execution benefits.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]

		start := time.Now()
		fmt.Printf("⏱️  Starting task '%s' at %s\n", taskName, start.Format("15:04:05.000"))

		// Load config and run task
		config, err := runner.LoadConfig("tasks.yaml")
		if err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			fmt.Println("\n💡 Tip: Run 't :init' to create a tasks.yaml file")
			return
		}

		taskRunner := runner.NewRunner(config)

		if err := taskRunner.RunTask(taskName); err != nil {
			fmt.Printf("❌ Task failed: %v\n", err)
			return
		}

		duration := time.Since(start)
		fmt.Printf("🎉 Task '%s' completed successfully in %v!\n", taskName, duration.Round(time.Millisecond))
	},
}
