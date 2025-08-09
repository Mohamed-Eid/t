package cmd

import (
	"fmt"
	"time"

	"t/internal/runner"

	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:     ":ps",
	Aliases: []string{":p", ":processes", ":status"},
	Short:   "List running detached tasks",
	Long:    "Show all currently running detached tasks with their PIDs, start times, and log files.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config (we need a runner instance to access the methods)
		config, err := runner.LoadConfig("tasks.yaml")
		if err != nil {
			// For listing processes, we don't strictly need a valid config
			// But we need a runner instance
			fmt.Printf("⚠️  Warning: Could not load tasks.yaml, showing tracked processes only\n")
			config = &runner.Config{} // Empty config
		}

		taskRunner := runner.NewRunner(config)

		// Get list of detached processes
		processes, err := taskRunner.ListDetachedProcesses()
		if err != nil {
			fmt.Printf("❌ Error listing detached processes: %v\n", err)
			return
		}

		if len(processes) == 0 {
			fmt.Println("📭 No detached tasks are currently running")
			fmt.Println("\n💡 Start a detached task with: t :detach <task-name>")
			return
		}

		fmt.Printf("🔧 Running detached tasks (%d):\n\n", len(processes))

		for _, proc := range processes {
			duration := time.Since(proc.StartedAt).Round(time.Second)
			fmt.Printf("  📋 Task: %s\n", proc.TaskName)
			fmt.Printf("     🆔 PID: %d\n", proc.PID)
			fmt.Printf("     ⏰ Running for: %v\n", duration)
			fmt.Printf("     📝 Log file: %s\n", proc.LogFile)
			fmt.Printf("     🛑 Stop with: t :stop %s\n\n", proc.TaskName)
		}

		fmt.Printf("💡 Use 't :stop <task-name>' or 't :stop <pid>' to stop a task\n")
		fmt.Printf("💡 Use 't :logs <task-name>' to view logs\n")
	},
}
