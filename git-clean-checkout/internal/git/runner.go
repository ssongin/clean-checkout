package git

import (
	"fmt"
	"os/exec"
)

func Run(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %v failed: %w", args, err)
	}
	return nil
}