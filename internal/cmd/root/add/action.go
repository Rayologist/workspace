package add

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"workspace/internal/config"
	"workspace/internal/git"
)

func add(opt *AddOptions) error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	_, err = c.Workspace(opt.ProjectName)
	if err == nil {
		return fmt.Errorf("workspace '%s' already exists", opt.ProjectName)
	}

	infos := make([]*RepoInfo, len(opt.RepoConfigs))

	for i, repo := range opt.RepoConfigs {
		info, err := parseRepo(repo)
		if err != nil {
			return fmt.Errorf("invalid repo config '%s': %w", repo, err)
		}
		if info.Branch == "" {
			info.Branch = opt.ProjectName
		}
		infos[i] = info
	}

	for _, info := range infos {
		repoConfig := &config.WorkspaceRepoConfig{
			Branch: info.Branch,
		}
		if err := c.AddWorkspaceRepo(opt.ProjectName, info.Alias, repoConfig); err != nil {
			return fmt.Errorf("failed to add repo '%s' to workspace '%s': %w", info.Alias, opt.ProjectName, err)
		}
	}

	w, err := c.Workspace(opt.ProjectName)
	if err != nil {
		return err
	}

	workspaceDir, err := config.WorkspacesDirPath()
	if err != nil {
		return fmt.Errorf("failed to get workspaces dir path: %w", err)
	}

	var rollbacks []func()
	success := false

	defer func() {
		if success {
			return
		}

		fmt.Println("\nErrors occurred during workspace setup.")
		fmt.Printf("Rolling back changes...\n\n")
		// LIFO order for rollbacks
		for i := len(rollbacks) - 1; i >= 0; i-- {
			rollbacks[i]()
		}

		workspacePath := filepath.Join(workspaceDir, opt.ProjectName)

		if _, err := os.Stat(workspacePath); err == nil {
			fmt.Printf("Removing workspace directory '%s'...\n", workspacePath)
			if err := os.Remove(workspacePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Failed to remove workspace directory '%s' during rollback: %v\n", workspacePath, err)
			}
		}

		fmt.Printf("\bRollback completed.\n\n")
	}()

	for alias, workspaceRepoConfig := range w.Repos {
		repo, err := c.RepoByAlias(alias)
		if err != nil {
			return fmt.Errorf("failed to get repo '%s' from config: %w", alias, err)
		}

		worktreePath := filepath.Join(workspaceDir, opt.ProjectName, alias)

		isNewBranch, err := git.AddWorktree(repo.Path, worktreePath, workspaceRepoConfig.Branch)
		if err != nil {
			return fmt.Errorf("failed to add worktree for repo '%s': %w", alias, err)
		}

		if isNewBranch {
			rollbacks = append(rollbacks, func() {
				fmt.Printf("Rolling back new branch '%s' for repo '%s'...\n", workspaceRepoConfig.Branch, alias)
				fmt.Printf("Removing worktree for repo '%s' at path '%s'...\n", alias, worktreePath)
				if err := git.RemoveWorktree(repo.Path, worktreePath, false); err != nil {
					fmt.Printf("Failed to remove worktree for repo '%s' during rollback: %v\n", alias, err)
				}
				fmt.Printf("Worktree for repo '%s' rolled back successfully.\n", alias)

				fmt.Printf("Deleting new branch '%s' for repo '%s'...\n", workspaceRepoConfig.Branch, alias)
				if err := git.DeleteBranch(repo.Path, workspaceRepoConfig.Branch, false); err != nil {
					fmt.Printf("Failed to delete branch '%s' for repo '%s' during rollback: %v\n", workspaceRepoConfig.Branch, alias, err)
				}
				fmt.Printf("Branch '%s' for repo '%s' rolled back successfully.\n", workspaceRepoConfig.Branch, alias)
			})
		} else {
			rollbacks = append(rollbacks, func() {
				fmt.Printf("Rolling back existing branch worktree for repo '%s'...\n", alias)
				if err := git.RemoveWorktree(repo.Path, worktreePath, false); err != nil {
					fmt.Printf("Failed to remove worktree for repo '%s' during rollback: %v\n", alias, err)
				}
				fmt.Printf("Worktree for repo '%s' rolled back successfully.\n", alias)
			})
		}

	}

	for alias := range w.Repos {
		repo, err := c.RepoByAlias(alias)
		if err != nil {
			return fmt.Errorf("failed to get repo '%s' from config: %w", alias, err)
		}

		worktreePath := filepath.Join(workspaceDir, opt.ProjectName, alias)

		for _, hook := range repo.Hooks.Setup {
			cmd := exec.Command("sh", "-c", hook)
			cmd.Env = append(os.Environ(), "WS_TARGET="+worktreePath, "WS_SOURCE="+repo.Path)
			cmd.Dir = worktreePath
			fmt.Printf("Running setup hook for repo '%s': %s\n", alias, hook)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to run setup hook for repo '%s': %w (output: %s)", alias, err, out)
			}
			if len(out) > 0 {
				fmt.Println(string(out))
			}
		}
	}

	if err := c.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	success = true

	fmt.Printf("Workspace '%s' created successfully.\n", opt.ProjectName)

	return nil
}

type RepoInfo struct {
	Alias  string
	Branch string
}

func parseRepo(repo string) (*RepoInfo, error) {
	parts := strings.SplitN(repo, ":", 2)

	if len(parts) == 1 {
		return &RepoInfo{
			Alias:  parts[0],
			Branch: "",
		}, nil
	}

	if len(parts) == 2 {
		return &RepoInfo{
			Alias:  parts[0],
			Branch: parts[1],
		}, nil
	}
	return nil, fmt.Errorf("invalid repo format: %s", repo)
}
