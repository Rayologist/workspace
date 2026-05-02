package git

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func IsRepo(path string) (bool, error) {
	info, err := os.Stat(filepath.Join(path, ".git"))

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func BranchExists(path, branch string) (bool, error) {
	_, err := git(path, "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)
	if err != nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}

	return false, err
}

func DeleteBranch(path, branch string, force bool) error {
	flag := "-d"

	if force {
		flag = "-D"
	}

	_, err := git(path, "--no-optional-locks", "branch", flag, branch)

	return err
}
