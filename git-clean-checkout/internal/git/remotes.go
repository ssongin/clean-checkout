package git

import (
	"os/exec"
	"strings"
)

func ListRemotes() ([]string, error) {
	cmd := exec.Command("git", "remote")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	var remotes []string
	for _, l := range lines {
		if l != "" {
			remotes = append(remotes, strings.TrimSpace(l))
		}
	}

	return remotes, nil
}