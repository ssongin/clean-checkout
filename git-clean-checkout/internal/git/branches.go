package git

import (
	"os/exec"
	"strings"
)

func ListBranches() ([]string, error) {
	cmd := exec.Command("git", "branch")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var branches []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		branch := strings.TrimPrefix(line, "* ")
		branches = append(branches, branch)
	}

	return branches, nil
}