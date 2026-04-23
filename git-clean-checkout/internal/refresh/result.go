package refresh

import "github.com/ssongin/clean-checkout/git-clean-checkout/internal/git"

type Result struct {
	Branch          string        `json:"branch"`
	Reset           bool          `json:"reset"`
	DryRun          bool          `json:"dry_run"`
	DeletedBranches []string      `json:"deleted_branches"`
	SkippedBranches []string      `json:"skipped_branches"`
	Commands        []git.Command `json:"commands"`
}
