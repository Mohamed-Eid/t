package cmd

import (
	"fmt"

	"t/internal/runner"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     ":list",
	Short:   "List all available tasks",
	Long:    "Display all tasks defined in the tasks.yaml file with their descriptions.",
	Aliases: []string{":ls", ":tasks"},
	Run: func(cmd *cobra.Command, args []string) {
		listTasks()
	},
}

func listTasks() {
	// Load config
	config, err := runner.LoadConfig("tasks.yaml")
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		fmt.Println("\n💡 Tip: Run 't :init' to create a tasks.yaml file")
		return
	}

	if len(config.Tasks) == 0 {
		fmt.Println("No tasks found in tasks.yaml")
		return
	}

	fmt.Println("📋 Available tasks:")
	fmt.Println()

	for taskName, task := range config.Tasks {
		fmt.Printf("  🔧 %s", taskName)

		if task.Desc != "" {
			fmt.Printf(" - %s", task.Desc)
		}

		if len(task.Deps) > 0 {
			fmt.Printf(" (depends on: %v)", task.Deps)
		}

		fmt.Println()
	}

	fmt.Println()
	fmt.Println("💡 Run 't <task-name>' to execute a task")
}
