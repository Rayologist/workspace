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

	err = c.RemoveRepo(opts.Alias)
	if err != nil {
		return err
	}

	if err := c.Save(); err != nil {
		return err
	}

	fmt.Printf("Repository '%s' removed successfully.\n", opts.Alias)

	return nil
}
