package refresh

import "github.com/ssongin/clean-checkout/git-clean-checkout/internal/git"

type Plan struct {
	Branch   string        `json:"branch"`
	Commands []git.Command `json:"commands"`
}

type Options struct {
	Reset           bool
	DryRun          bool
	OnlyDestructive bool
	RequireConfirm  bool
	Remote          string
}

var protectedBranches = map[string]bool{
	"main":    true,
	"develop": true,
}

func PlanRefresh(branch string, opts Options) (*Plan, error) {
	var cmds []git.Command

	if err := validateRemote(opts.Remote); err != nil {
		return nil, err
	}

	cmds = append(cmds, git.Checkout(branch))
	cmds = append(cmds, git.Fetch(opts.Remote))

	if opts.Reset {
		cmds = append(cmds, git.ResetHard(opts.Remote+"/"+branch))
	} else {
		cmds = append(cmds, git.PullFFOnly())
	}

	branches, err := git.ListBranches()
	if err != nil {
		return nil, err
	}

	for _, b := range branches {
		if b == branch || protectedBranches[b] {
			continue
		}
		cmds = append(cmds, git.DeleteBranch(b))
	}

	if opts.OnlyDestructive {
		var filtered []git.Command
		for _, c := range cmds {
			if c.Destructive {
				filtered = append(filtered, c)
			}
		}
		cmds = filtered
	}

	return &Plan{
		Branch:   branch,
		Commands: cmds,
	}, nil
}
