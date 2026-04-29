package add

import (
	"fmt"
	"os"
	"path/filepath"

	"workspace/internal/cli"
	"workspace/internal/config"
	"workspace/internal/workspace"

	"github.com/spf13/cobra"
)

type AddOptions struct {
	Config func() (*config.Config, error)

	ProjectName string
	RepoConfigs []string
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &AddOptions{
		Config: r.Config,
	}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new project to the workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectName = args[0]
			return runAdd(opts)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.RepoConfigs, "repo", "r", []string{}, "Add a repository to the project (branch defaults to 'main', if not specified) - format: <alias>[:<branch>]")
	cmd.MarkFlagRequired("repo")

	return cmd
}

func runAdd(opts *AddOptions) error {
	c, err := opts.Config()
	if err != nil {
		return err
	}

	if _, err = c.Workspace(opts.ProjectName); err == nil {
		return fmt.Errorf("workspace '%s' already exists", opts.ProjectName)
	}

	repoConfigs := make([]cli.RepoConfig, len(opts.RepoConfigs))

	for i, config := range opts.RepoConfigs {
		repoConfig := cli.ParseRepoConfig(config)
		if !repoConfig.HasBranch() {
			repoConfig.TargetBranch = opts.ProjectName
		}
		repoConfigs[i] = repoConfig
	}

	workspaceDir, err := config.WorkspacesDirPath()
	if err != nil {
		return err
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

		workspacePath := filepath.Join(workspaceDir, opts.ProjectName)

		if _, err := os.Stat(workspacePath); err == nil {
			fmt.Printf("Removing workspace directory '%s'...\n", workspacePath)
			if err := os.Remove(workspacePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Failed to remove workspace directory '%s' during rollback: %v\n", workspacePath, err)
			}
		}

		fmt.Printf("\bRollback completed.\n\n")
	}()

	for _, repoConfig := range repoConfigs {
		rollback, err := workspace.AttachRepo(c, workspace.AttachRepoArgs{
			WorkspacesDir: workspaceDir,
			ProjectName:   opts.ProjectName,
			SourceAlias:   repoConfig.SourceAlias,
			SourceBranch:  repoConfig.TargetBranch,
		})
		if rollback != nil {
			rollbacks = append(rollbacks, rollback)
		}
		if err != nil {
			return err
		}
	}

	if err := c.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	success = true

	fmt.Printf("Workspace '%s' created successfully.\n", opts.ProjectName)

	return nil
}
