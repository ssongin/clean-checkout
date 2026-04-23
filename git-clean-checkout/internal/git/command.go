package git

import (
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	Name        string   `json:"name"`
	Args        []string `json:"args"`
	Description string   `json:"description"`
	Destructive bool     `json:"destructive"`
}

type CommandResult struct {
	Command  Command       `json:"command"`
	Duration time.Duration `json:"duration"`
	Success  bool          `json:"success"`
	Error    string        `json:"error,omitempty"`
}

func (c Command) String() string {
	return "git " + strings.Join(append([]string{c.Name}, c.Args...), " ")
}

func (c Command) Run() CommandResult {
	start := time.Now()

	cmd := exec.Command("git", append([]string{c.Name}, c.Args...)...)
	err := cmd.Run()

	duration := time.Since(start)

	if err != nil {
		return CommandResult{
			Command:  c,
			Duration: duration,
			Success:  false,
			Error:    err.Error(),
		}
	}

	return CommandResult{
		Command:  c,
		Duration: duration,
		Success:  true,
	}
}
