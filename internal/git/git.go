package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func ValidateRepo(path string) error {
	info, err := os.Stat(filepath.Join(path, ".git"))

	if os.IsNotExist(err) {
		return fmt.Errorf("directory not exists: %s", path)
	}

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("not a git repository: %s", path)
	}

	return nil
}

func ValidateBranch(path, branch string) error {
	_, err := git(path, "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)

	return err
}

func DeleteBranch(path, branch string, force bool) error {
	flag := "-d"

	if force {
		flag = "-D"
	}

	_, err := git(path, "--no-optional-locks", "branch", flag, branch)

	return err
}
