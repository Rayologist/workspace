package config

import "fmt"

type WorkspaceNotFoundError struct {
	Name string
}

func (e *WorkspaceNotFoundError) Error() string {
	return fmt.Sprintf("workspace %q not found", e.Name)
}

type WorkspaceExistsError struct {
	Name string
}

func (e *WorkspaceExistsError) Error() string {
	return fmt.Sprintf("workspace %q already exists", e.Name)
}

type WorkspaceRepoNotFoundError struct {
	SourceAlias string
}

func (e *WorkspaceRepoNotFoundError) Error() string {
	return fmt.Sprintf("repo with source alias %q not found", e.SourceAlias)
}

type WorkspaceRepoExistsError struct {
	SourceAlias string
}

func (e *WorkspaceRepoExistsError) Error() string {
	return fmt.Sprintf("repo with source alias %q already exists", e.SourceAlias)
}

type SourceNotFoundError struct {
	Alias string
}

func (e *SourceNotFoundError) Error() string {
	return fmt.Sprintf("source %q not found", e.Alias)
}

type SourceExistsError struct {
	Alias string
}

func (e *SourceExistsError) Error() string {
	return fmt.Sprintf("source %q already exists", e.Alias)
}
