package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   ":init",
	Short: "init t file (tasks.yaml)",
	Long:  "Initialize the task file (tasks.yaml) with a default structure.",
	Run: func(cmd *cobra.Command, args []string) {
		initTasksFile()
	},
}

func initTasksFile() {
	if _, err := os.Stat("tasks.yaml"); err == nil {
		fmt.Println("❌ tasks.yaml already exists in current directory")
		fmt.Println("Remove it first or use a different directory")
		return
	}

	// Create the tasks.yaml file with a default structure
	file, err := os.Create("tasks.yaml")
	if err != nil {
		fmt.Printf("❌ Error creating tasks.yaml: %v\n", err)
		return
	}
	defer file.Close()

	// Default tasks.yaml content
	defaultContent := `version: "1"

vars:
  APP_NAME: "myapp"
  BUILD_DIR: "bin"

tasks:
  build:
    desc: "Build the application"
    deps: [clean]
    cmds:
      - "mkdir -p {{.BUILD_DIR}}"
      - "go build -ldflags='-s -w' -o {{.BUILD_DIR}}/{{.APP_NAME}} ."

  test:
    desc: "Run tests"
    cmds:
      - "go test ./..."

  clean:
    desc: "Clean build artifacts"
    cmds:
      - "rm -rf {{.BUILD_DIR}}"
      - "rm -f {{.APP_NAME}} {{.APP_NAME}}.exe"

  dev:
    desc: "Run in development mode"
    cmds:
      - "go run ."

  install:
    desc: "Install dependencies"
    cmds:
      - "go mod download"
      - "go mod tidy"

  lint:
    desc: "Run linter"
    cmds:
      - "go fmt ./..."
      - "go vet ./..."
`

	// Write the content to the file
	_, err = file.WriteString(defaultContent)
	if err != nil {
		fmt.Printf("❌ Error writing to tasks.yaml: %v\n", err)
		return
	}

	fmt.Println("✅ Created tasks.yaml with default tasks:")
	fmt.Println("   • build  - Build the application")
	fmt.Println("   • test   - Run tests")
	fmt.Println("   • clean  - Clean build artifacts")
	fmt.Println("   • dev    - Run in development mode")
	fmt.Println("   • install- Install dependencies")
	fmt.Println("   • lint   - Run linter")
	fmt.Println("")
	fmt.Println("Run 't build' to get started!")
}
