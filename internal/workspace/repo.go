package workspace

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"workspace/internal/config"
	"workspace/internal/git"
)

type AttachRepoArgs struct {
	WorkspacesDir   string
	TargetWorkspace string
	SourceAlias     string
	SourceBranch    string
}

func AttachRepo(c *config.Config, args AttachRepoArgs) (rollback func(), err error) {
	source, err := c.SourceByAlias(args.SourceAlias)
	if err != nil {
		return nil, err
	}

	repoConfig := &config.WorkspaceRepoConfig{
		Branch: args.SourceBranch,
	}

	if err := c.AddWorkspaceRepo(args.TargetWorkspace, args.SourceAlias, repoConfig); err != nil {
		return nil, fmt.Errorf("failed to add repo '%s' to workspace '%s': %w", args.SourceAlias, args.TargetWorkspace, err)
	}

	worktreePath := filepath.Join(args.WorkspacesDir, args.TargetWorkspace, args.SourceAlias)

	isNewBranch, err := git.AddWorktree(source.Path, worktreePath, args.SourceBranch)
	if err != nil {
		_ = c.RemoveWorkspaceRepo(args.TargetWorkspace, args.SourceAlias)
		return nil, fmt.Errorf("failed to add worktree for '%s': %w", args.SourceAlias, err)
	}

	rollback = func() {
		fmt.Printf("Rolling back worktree for '%s' at '%s'...\n", args.SourceAlias, worktreePath)
		if err := git.RemoveWorktree(source.Path, worktreePath, false); err != nil {
			fmt.Printf("Failed to remove worktree for '%s' during rollback: %v\n", args.SourceAlias, err)
		}
		if isNewBranch {
			fmt.Printf("Deleting new branch '%s' for '%s'...\n", args.SourceBranch, args.SourceAlias)
			if err := git.DeleteBranch(source.Path, args.SourceBranch, false); err != nil {
				fmt.Printf("Failed to delete branch '%s' for '%s' during rollback: %v\n", args.SourceBranch, args.SourceAlias, err)
			}
		}
		if err := c.RemoveWorkspaceRepo(args.TargetWorkspace, args.SourceAlias); err != nil {
			fmt.Printf("Failed to remove repo '%s' from workspace config '%s' during rollback: %v\n", args.SourceAlias, args.TargetWorkspace, err)
		}

		fmt.Printf("Rollback for '%s' completed.\n", args.SourceAlias)
	}

	for _, hook := range source.Hooks.Setup {
		cmd := exec.Command("sh", "-c", hook)
		cmd.Dir = worktreePath
		cmd.Env = append(os.Environ(), "WS_TARGET="+worktreePath, "WS_SOURCE="+source.Path)
		fmt.Printf("Running setup hook for '%s': %s\n", args.SourceAlias, hook)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return rollback, fmt.Errorf("failed to run setup hook for repo '%s': %w (output: %s)", args.SourceAlias, err, out)
		}

		if len(out) > 0 {
			fmt.Printf("Output from setup hook:\n%s\n", string(out))
		}
	}

	return rollback, nil
}
