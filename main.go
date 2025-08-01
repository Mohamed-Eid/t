package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: t <task>")
		os.Exit(1)
	}

	taskName := os.Args[1]

	config, err := LoadConfig("tasks.yaml")
	if err != nil {
		fmt.Println("❌ Error loading config:", err)
		os.Exit(1)
	}

	runner := NewRunner(config)

	if err := runner.RunTask(taskName); err != nil {
		fmt.Println("❌ Task failed:", err)
		os.Exit(1)
	}
}
