package git

func AddWorktree(repoPath, worktreePath, branchName string) (bool, error) {
	err := ValidateBranch(repoPath, branchName)

	isNewBranch := err != nil

	args := []string{
		"--no-optional-locks",
		"worktree",
		"add",
		"--quiet",
	}

	if isNewBranch {
		args = append(args, "-b", branchName, worktreePath)
	} else {
		args = append(args, worktreePath, branchName)
	}

	if _, err = git(repoPath, args...); err != nil {
		return false, err
	}

	return isNewBranch, nil
}

func RemoveWorktree(repoPath, worktreePath string, force bool) error {
	args := []string{
		"worktree",
		"remove",
	}

	if force {
		args = append(args, "--force")
	}

	args = append(args, worktreePath)

	_, err := git(repoPath, args...)

	return err
}
