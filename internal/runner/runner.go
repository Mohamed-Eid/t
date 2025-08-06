package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
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
	mutex  sync.RWMutex
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
	return r.runTaskWithSync(taskName)
}

// runTaskWithSync executes a task with proper synchronization
func (r *Runner) runTaskWithSync(taskName string) error {
	// Check if already ran (with read lock)
	r.mutex.RLock()
	if r.Ran[taskName] {
		r.mutex.RUnlock()
		return nil
	}
	r.mutex.RUnlock()

	task, exists := r.Config.Tasks[taskName]
	if !exists {
		return fmt.Errorf("task %s not found", taskName)
	}

	// Run dependencies in parallel if possible
	if len(task.Deps) > 0 {
		if err := r.runDependenciesParallel(task.Deps); err != nil {
			return err
		}
	}

	// Check again if task was run by a dependency (with write lock)
	r.mutex.Lock()
	if r.Ran[taskName] {
		r.mutex.Unlock()
		return nil
	}

	fmt.Printf("🔧 Running task: %s\n", taskName)

	// Mark as running to prevent duplicate execution
	r.Ran[taskName] = true
	r.mutex.Unlock()

	// Run task commands sequentially (commands within a task should be sequential)
	return r.executeCommands(taskName, task.Cmds)
}

// runDependenciesParallel runs dependencies in parallel where possible
func (r *Runner) runDependenciesParallel(deps []string) error {
	if len(deps) == 1 {
		// Single dependency - run directly
		return r.runTaskWithSync(deps[0])
	}

	// Multiple dependencies - run in parallel
	var wg sync.WaitGroup
	errChan := make(chan error, len(deps))

	for _, dep := range deps {
		wg.Add(1)
		go func(depName string) {
			defer wg.Done()
			if err := r.runTaskWithSync(depName); err != nil {
				errChan <- fmt.Errorf("dependency %s failed: %w", depName, err)
			}
		}(dep)
	}

	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// executeCommands runs the commands for a task sequentially
func (r *Runner) executeCommands(taskName string, commands []string) error {
	for _, rawCmd := range commands {
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
