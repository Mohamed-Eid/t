package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

type Task struct {
	Desc string   `yaml:"desc"`
	Deps []string `yaml:"deps"`
	Cmds []string `yaml:"cmds"`
}

type Config struct {
	Version string            `yaml:"version"`
	Vars    map[string]string `yaml:"vars"`
	Tasks   map[string]Task   `yaml:"tasks"`
}

type Runner struct {
	Config *Config
	Ran    map[string]bool
}

func NewRunner(config *Config) *Runner {
	return &Runner{
		Config: config,
		Ran:    make(map[string]bool),
	}
}

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
		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command failed: %s", cmdStr)
		}

		fmt.Printf("✅ Done\n")
	}

	r.Ran[taskName] = true
	return nil
}

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
