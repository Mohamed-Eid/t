package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"t/internal/runner"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:     ":logs <task-name-or-pid>",
	Aliases: []string{":log", ":l", ":tail"},
	Short:   "View logs of a detached task",
	Long:    "Display the logs of a running or recently finished detached task.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		// Load config (we need a runner instance)
		config, err := runner.LoadConfig("tasks.yaml")
		if err != nil {
			config = &runner.Config{} // Empty config
		}

		taskRunner := runner.NewRunner(config)

		// Get list of detached processes to find the log file
		processes, err := taskRunner.ListDetachedProcesses()
		if err != nil {
			fmt.Printf("❌ Error listing detached processes: %v\n", err)
			return
		}

		var logFile string
		var taskName string

		// Try to find by PID first
		if pid, err := strconv.Atoi(identifier); err == nil {
			for _, proc := range processes {
				if proc.PID == pid {
					logFile = proc.LogFile
					taskName = proc.TaskName
					break
				}
			}
		} else {
			// Search by task name
			for _, proc := range processes {
				if proc.TaskName == identifier {
					logFile = proc.LogFile
					taskName = proc.TaskName
					break
				}
			}
		}

		if logFile == "" {
			fmt.Printf("❌ No detached task found with identifier: %s\n", identifier)
			fmt.Println("\n💡 Use 't :ps' to see running detached tasks")
			return
		}

		// Check if log file exists
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			fmt.Printf("❌ Log file not found: %s\n", logFile)
			return
		}

		fmt.Printf("📝 Logs for task '%s':\n", taskName)
		fmt.Printf("📄 File: %s\n\n", logFile)

		// Follow flag for tail -f behavior
		follow, _ := cmd.Flags().GetBool("follow")

		// Display logs using appropriate command for the platform
		var tailCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			if follow {
				// PowerShell equivalent of tail -f
				tailCmd = exec.Command("powershell", "-Command",
					fmt.Sprintf("Get-Content '%s' -Wait -Tail 50", logFile))
			} else {
				// Show last 50 lines
				tailCmd = exec.Command("powershell", "-Command",
					fmt.Sprintf("Get-Content '%s' -Tail 50", logFile))
			}
		} else {
			if follow {
				tailCmd = exec.Command("tail", "-f", "-n", "50", logFile)
			} else {
				tailCmd = exec.Command("tail", "-n", "50", logFile)
			}
		}

		tailCmd.Stdout = os.Stdout
		tailCmd.Stderr = os.Stderr

		if follow {
			fmt.Println("📡 Following logs (Press Ctrl+C to exit)...")
			fmt.Println("─────────────────────────────────────────────")
		} else {
			fmt.Println("📋 Last 50 lines:")
			fmt.Println("─────────────────────────────────────────────")
		}

		if err := tailCmd.Run(); err != nil {
			fmt.Printf("❌ Error viewing logs: %v\n", err)
		}
	},
}

func init() {
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output (like tail -f)")
}
