package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// Task represents a single task configuration
type Task struct {
	Desc        string            `yaml:"desc"`
	Deps        []string          `yaml:"deps"`
	Cmds        []string          `yaml:"cmds"`
	Interactive map[string]Prompt `yaml:"interactive"`
}

// Prompt represents an interactive prompt configuration
type Prompt struct {
	Message  string `yaml:"message"`
	Required bool   `yaml:"required"`
	Default  string `yaml:"default"`
}

// Config represents the entire tasks.yaml configuration
type Config struct {
	Version string            `yaml:"version"`
	Vars    map[string]string `yaml:"vars"`
	Tasks   map[string]Task   `yaml:"tasks"`
}

// DetachedProcess represents a background process
type DetachedProcess struct {
	PID       int       `json:"pid"`
	TaskName  string    `json:"task_name"`
	Command   string    `json:"command"`
	StartedAt time.Time `json:"started_at"`
	LogFile   string    `json:"log_file"`
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

	// Prompt for interactive input if needed
	interactiveInputs, err := r.promptForInput(taskName, task)
	if err != nil {
		r.mutex.Unlock()
		return fmt.Errorf("interactive input failed: %w", err)
	}

	// Mark as running to prevent duplicate execution
	r.Ran[taskName] = true
	r.mutex.Unlock()

	// Run task commands sequentially (commands within a task should be sequential)
	return r.executeCommandsWithInteractive(taskName, task.Cmds, interactiveInputs)
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

// executeCommandsWithInteractive runs the commands for a task sequentially with interactive inputs
func (r *Runner) executeCommandsWithInteractive(taskName string, commands []string, interactiveInputs map[string]string) error {
	for _, rawCmd := range commands {
		// First expand regular variables
		cmdStr, err := r.expandVars(rawCmd)
		if err != nil {
			return err
		}

		// Then expand interactive variables
		cmdStr, err = r.expandVarsWithInteractive(cmdStr, interactiveInputs)
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
} // expandVars replaces variables in commands with their values
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

// promptForInput prompts the user for interactive input
func (r *Runner) promptForInput(taskName string, task Task) (map[string]string, error) {
	inputs := make(map[string]string)

	if len(task.Interactive) == 0 {
		return inputs, nil
	}

	fmt.Printf("🤔 Task '%s' requires interactive input:\n\n", taskName)

	reader := bufio.NewReader(os.Stdin)

	for varName, prompt := range task.Interactive {
		// Show the prompt message
		fmt.Printf("📝 %s", prompt.Message)

		// Show default value if available
		if prompt.Default != "" {
			fmt.Printf(" [%s]", prompt.Default)
		}

		// Show required indicator
		if prompt.Required {
			fmt.Printf(" (required)")
		}

		fmt.Print(": ")

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		// Clean the input
		input = strings.TrimSpace(input)

		// Use default if no input provided
		if input == "" && prompt.Default != "" {
			input = prompt.Default
		}

		// Check if required input is provided
		if prompt.Required && input == "" {
			return nil, fmt.Errorf("required input '%s' not provided", varName)
		}

		inputs[varName] = input
		fmt.Printf("✅ %s: %s\n", varName, input)
	}

	fmt.Println()
	return inputs, nil
}

// expandVarsWithInteractive replaces variables in commands with their values including interactive inputs
func (r *Runner) expandVarsWithInteractive(cmdStr string, interactiveInputs map[string]string) (string, error) {
	result := cmdStr

	// Expand interactive variables using $variable syntax
	for varName, value := range interactiveInputs {
		result = strings.ReplaceAll(result, "$"+varName, value)
	}

	return result, nil
} // RunTaskDetached runs a task in the background and returns immediately
func (r *Runner) RunTaskDetached(taskName string) (*DetachedProcess, error) {
	task, exists := r.Config.Tasks[taskName]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskName)
	}

	// Run dependencies first (synchronously)
	if len(task.Deps) > 0 {
		fmt.Printf("🔧 Running dependencies for detached task: %s\n", taskName)
		if err := r.runDependenciesParallel(task.Deps); err != nil {
			return nil, fmt.Errorf("dependencies failed: %w", err)
		}
	}

	// Create logs directory if it doesn't exist
	logsDir := ".t-logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file for this task
	timestamp := time.Now().Format("20060102-150405")
	logFile := filepath.Join(logsDir, fmt.Sprintf("%s-%s.log", taskName, timestamp))

	// Start the first command in detached mode
	if len(task.Cmds) == 0 {
		return nil, fmt.Errorf("task %s has no commands to run", taskName)
	}

	// For detached mode, we'll run the first command as the main process
	// and any additional commands as setup
	mainCmd := task.Cmds[len(task.Cmds)-1]    // Use last command as main
	setupCmds := task.Cmds[:len(task.Cmds)-1] // Previous commands as setup

	// Run setup commands first (if any)
	if len(setupCmds) > 0 {
		fmt.Printf("🔧 Running setup commands for detached task: %s\n", taskName)
		for _, rawCmd := range setupCmds {
			cmdStr, err := r.expandVars(rawCmd)
			if err != nil {
				return nil, err
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

			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf("setup command failed: %s", cmdStr)
			}
			fmt.Printf("✅ done\n")
		}
	}

	// Expand variables in the main command
	cmdStr, err := r.expandVars(mainCmd)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🚀 Starting detached task: %s\n", taskName)
	fmt.Printf("➡️  %s\n", cmdStr)

	// Create the command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	// Create or open log file
	logFileHandle, err := os.Create(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	// Redirect output to log file
	cmd.Stdout = logFileHandle
	cmd.Stderr = logFileHandle

	// Set up process group for proper cleanup of child processes
	if runtime.GOOS == "windows" {
		// On Windows, create a new process group
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
	} else {
		// On Unix-like systems, we'll handle process groups differently
		// For now, use basic process creation and handle cleanup in stop command
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		logFileHandle.Close()
		return nil, fmt.Errorf("failed to start detached process: %w", err)
	}

	// Create detached process info
	detachedProc := &DetachedProcess{
		PID:       cmd.Process.Pid,
		TaskName:  taskName,
		Command:   cmdStr,
		StartedAt: time.Now(),
		LogFile:   logFile,
	}

	// Save process info to file for later reference
	if err := r.saveDetachedProcess(detachedProc); err != nil {
		fmt.Printf("⚠️  Warning: failed to save process info: %v\n", err)
	}

	fmt.Printf("✅ Task '%s' started in background (PID: %d)\n", taskName, cmd.Process.Pid)
	fmt.Printf("📝 Logs: %s\n", logFile)
	fmt.Printf("🛑 Stop with: t :stop %s (or PID %d)\n", taskName, cmd.Process.Pid)

	// Start a goroutine to wait for the process and clean up
	go func() {
		defer logFileHandle.Close()
		cmd.Wait()
		r.removeDetachedProcess(detachedProc.PID)
	}()

	return detachedProc, nil
}

// saveDetachedProcess saves process info to a file
func (r *Runner) saveDetachedProcess(proc *DetachedProcess) error {
	processesDir := ".t-processes"
	if err := os.MkdirAll(processesDir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(processesDir, fmt.Sprintf("%d.json", proc.PID))
	data, err := json.MarshalIndent(proc, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// removeDetachedProcess removes process info file
func (r *Runner) removeDetachedProcess(pid int) {
	processesDir := ".t-processes"
	filename := filepath.Join(processesDir, fmt.Sprintf("%d.json", pid))
	os.Remove(filename) // Ignore errors
}

// ListDetachedProcesses returns all currently tracked detached processes
func (r *Runner) ListDetachedProcesses() ([]*DetachedProcess, error) {
	processesDir := ".t-processes"

	// Check if directory exists
	if _, err := os.Stat(processesDir); os.IsNotExist(err) {
		return []*DetachedProcess{}, nil
	}

	files, err := filepath.Glob(filepath.Join(processesDir, "*.json"))
	if err != nil {
		return nil, err
	}

	var processes []*DetachedProcess
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip invalid files
		}

		var proc DetachedProcess
		if err := json.Unmarshal(data, &proc); err != nil {
			continue // Skip invalid JSON
		}

		// Check if process is still running
		if r.isProcessRunning(proc.PID) {
			processes = append(processes, &proc)
		} else {
			// Clean up dead process
			os.Remove(file)
		}
	}

	return processes, nil
}

// isProcessRunning checks if a process with the given PID is still running
func (r *Runner) isProcessRunning(pid int) bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
		output, err := cmd.Output()
		if err != nil {
			return false
		}
		return bytes.Contains(output, []byte(strconv.Itoa(pid)))
	} else {
		// Unix-like systems
		cmd := exec.Command("ps", "-p", strconv.Itoa(pid))
		err := cmd.Run()
		return err == nil
	}
}

// StopDetachedProcess stops a detached process by PID or task name
func (r *Runner) StopDetachedProcess(identifier string) error {
	processes, err := r.ListDetachedProcesses()
	if err != nil {
		return err
	}

	var targetPID int
	var targetProc *DetachedProcess

	// Try to parse as PID first
	if pid, err := strconv.Atoi(identifier); err == nil {
		targetPID = pid
		// Find the process info
		for _, proc := range processes {
			if proc.PID == pid {
				targetProc = proc
				break
			}
		}
	} else {
		// Search by task name
		for _, proc := range processes {
			if proc.TaskName == identifier {
				targetPID = proc.PID
				targetProc = proc
				break
			}
		}
	}

	if targetPID == 0 {
		return fmt.Errorf("no detached process found with identifier: %s", identifier)
	}

	// Kill the process and its children
	if runtime.GOOS == "windows" {
		// On Windows, use taskkill with /T flag to kill the process tree
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(targetPID))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to kill process tree %d: %w", targetPID, err)
		}
	} else {
		// On Unix-like systems, try to kill the process group first, then the process
		// First try to kill the process group (negative PID)
		killGroupCmd := exec.Command("kill", fmt.Sprintf("-%d", targetPID))
		killGroupErr := killGroupCmd.Run()

		// Also kill the main process directly
		killCmd := exec.Command("kill", strconv.Itoa(targetPID))
		killErr := killCmd.Run()

		// If both fail, try a more aggressive approach
		if killGroupErr != nil && killErr != nil {
			// Try SIGKILL
			killForceCmd := exec.Command("kill", "-9", strconv.Itoa(targetPID))
			if err := killForceCmd.Run(); err != nil {
				return fmt.Errorf("failed to kill process %d: %w", targetPID, err)
			}
		}
	}

	// Clean up process info
	r.removeDetachedProcess(targetPID)

	if targetProc != nil {
		fmt.Printf("🛑 Stopped detached task '%s' (PID: %d)\n", targetProc.TaskName, targetPID)
	} else {
		fmt.Printf("🛑 Stopped process (PID: %d)\n", targetPID)
	}

	return nil
}
