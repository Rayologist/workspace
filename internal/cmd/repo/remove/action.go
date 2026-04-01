package remove

import (
	"fmt"

	"workspace/internal/config"
)

func remove(opts *RemoveOptions) error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	if _, exists := c.Repos[opts.Alias]; !exists {
		return fmt.Errorf("repository '%s' not found", opts.Alias)
	}

	delete(c.Repos, opts.Alias)

	if err := c.Save(); err != nil {
		return err
	}

	fmt.Printf("Repository '%s' removed successfully.\n", opts.Alias)

	return nil
}
