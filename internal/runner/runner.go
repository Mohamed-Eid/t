package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Task represents a single task configuration
type Task struct {
	Desc string   `yaml:"desc"`
	Deps []string `yaml:"deps"`
	Cmds []string `yaml:"cmds"`
}

// Config represents the entire tasks.yaml configuration
type Config struct {
	Version string            `yaml:"version"`
	Vars    map[string]string `yaml:"vars"`
	Tasks   map[string]Task   `yaml:"tasks"`
}

// Runner handles task execution
type Runner struct {
	Config *Config
	Ran    map[string]bool
}

// LoadConfig loads the tasks.yaml configuration from the specified filename
func LoadConfig(filename string) (*Config, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	// Construct full path to the config file in current directory
	configPath := filepath.Join(cwd, filename)

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("tasks.yaml not found in current directory: %s", cwd)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", configPath, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in %s: %w", filename, err)
	}

	return &config, nil
}

// NewRunner creates a new task runner instance
func NewRunner(config *Config) *Runner {
	return &Runner{
		Config: config,
		Ran:    make(map[string]bool),
	}
}

// RunTask executes a task and its dependencies
func (r *Runner) RunTask(taskName string) error {
	if r.Ran[taskName] {
		return nil
	}

	task, exists := r.Config.Tasks[taskName]
	if !exists {
		return fmt.Errorf("task %s not found", taskName)
	}

	// Run dependencies first
	for _, dep := range task.Deps {
		if err := r.RunTask(dep); err != nil {
			return err
		}
	}

	fmt.Printf("🔧 Running task: %s\n", taskName)

	// Run task commands
	for _, rawCmd := range task.Cmds {
		cmdStr, err := r.expandVars(rawCmd)
		if err != nil {
			return err
		}

		fmt.Printf("➡️  %s\n", cmdStr)

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", "-Command", cmdStr)
		} else {
			cmd = exec.Command("sh", "-c", cmdStr)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command failed: %s", cmdStr)
		}

		fmt.Printf("✅ done\n")
	}

	r.Ran[taskName] = true
	return nil
}

// expandVars replaces variables in commands with their values
func (r *Runner) expandVars(command string) (string, error) {
	tmpl, err := template.New("cmd").Parse(command)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r.Config.Vars); err != nil {
		return "", err
	}

	return buf.String(), nil
}
