package git

func Checkout(branch string) Command {
	return Command{
		Name:        "checkout",
		Args:        []string{branch},
		Description: "Switch to branch",
	}
}

func Fetch(remote string) Command {
	return Command{
		Name:        "fetch",
		Args:        []string{remote},
		Description: "Fetch latest from remote",
	}
}

func ResetHard(ref string) Command {
	return Command{
		Name:        "reset",
		Args:        []string{"--hard", ref},
		Description: "Hard reset to remote state",
		Destructive: true,
	}
}

func PullFFOnly() Command {
	return Command{
		Name:        "pull",
		Args:        []string{"--ff-only"},
		Description: "Fast-forward pull",
	}
}

func DeleteBranch(branch string) Command {
	return Command{
		Name:        "branch",
		Args:        []string{"-D", branch},
		Description: "Delete local branch",
		Destructive: true,
	}
}