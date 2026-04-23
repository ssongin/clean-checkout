package refresh

import (
	"fmt"

	"github.com/ssongin/clean-checkout/git-clean-checkout/internal/git"
)

func validateRemote(remote string) error {
	remotes, err := git.ListRemotes()
	if err != nil {
		return fmt.Errorf("failed to list git remotes: %w", err)
	}

	for _, r := range remotes {
		if r == remote {
			return nil
		}
	}

	return fmt.Errorf("invalid remote '%s' (available: %v)", remote, remotes)
}
