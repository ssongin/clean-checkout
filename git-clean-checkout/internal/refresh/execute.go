package refresh

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ssongin/clean-checkout/git-clean-checkout/internal/git"
)

type ExecutionResult struct {
	Results []git.CommandResult `json:"results"`
}

func Execute(plan *Plan, opts Options) (*ExecutionResult, error) {
	if opts.RequireConfirm && !opts.DryRun {
		fmt.Println("Planned commands:")
		for _, c := range plan.Commands {
			flag := ""
			if c.Destructive {
				flag = " [DESTRUCTIVE]"
			}
			fmt.Println("-", c.String(), flag)
		}

		fmt.Print("\nProceed? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		if input != "y\n" && input != "Y\n" {
			return nil, fmt.Errorf("aborted by user")
		}
	}

	var results []git.CommandResult

	for _, cmd := range plan.Commands {
		if opts.DryRun {
			results = append(results, git.CommandResult{
				Command: cmd,
				Success: true,
			})
			continue
		}

		res := cmd.Run()
		results = append(results, res)

		if !res.Success {
			return &ExecutionResult{Results: results}, fmt.Errorf("command failed")
		}
	}

	return &ExecutionResult{Results: results}, nil
}
