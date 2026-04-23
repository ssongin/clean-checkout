package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ssongin/clean-checkout/git-clean-checkout/internal/refresh"
)

func main() {
	reset := flag.Bool("reset", false, "Hard reset to origin/<branch>")
	dryRun := flag.Bool("dry-run", false, "Show what would be executed without running")
	jsonOut := flag.Bool("json", false, "Output structured JSON (for IDE integration)")
	onlyDestructive := flag.Bool("only-destructive", false, "Run only destructive commands")
	confirm := flag.Bool("confirm", false, "Ask for confirmation before execution")
	remote := flag.String("remote", "origin", "Git remote to use (origin, upstream, etc.)")

	flag.Usage = func() {
		fmt.Println("git-cc - safe git branch refresh utility")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  git-cc [options] <branch>")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	branch := flag.Arg(0)

	opts := refresh.Options{
		Reset:           *reset,
		DryRun:          *dryRun,
		OnlyDestructive: *onlyDestructive,
		RequireConfirm:  *confirm,
		Remote:          *remote,
	}

	plan, err := refresh.PlanRefresh(branch, opts)
	if err != nil {
		exit(err, *jsonOut)
	}

	result, err := refresh.Execute(plan, opts)
	if err != nil {
		exit(err, *jsonOut, result)
	}

	if *jsonOut {
		out := map[string]any{
			"plan":    plan,
			"result":  result,
			"success": err == nil,
		}

		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return
	}

	fmt.Println("✔ git-cc completed successfully")
}

func exit(err error, jsonOut bool, extra ...any) {
	if jsonOut {
		out := map[string]any{
			"error": err.Error(),
		}
		if len(extra) > 0 {
			out["partial"] = extra[0]
		}

		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		os.Exit(1)
	}

	fmt.Println("Error:", err)
	os.Exit(1)
}
