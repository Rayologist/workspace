package cli

import "strings"

type RepoConfig struct {
	SourceAlias  string
	TargetBranch string
}

func ParseRepoConfig(s string) RepoConfig {
	alias, branch, _ := strings.Cut(s, ":")

	return RepoConfig{
		SourceAlias:  alias,
		TargetBranch: branch,
	}
}

func (r RepoConfig) HasBranch() bool {
	return r.TargetBranch != ""
}
