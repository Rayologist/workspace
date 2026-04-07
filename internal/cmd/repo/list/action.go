package list

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"workspace/internal/config"
)

func list() error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, err = fmt.Fprintf(w, "ALIAS\tREPO\tBRANCH\tSETUPS\n")
	if err != nil {
		return err
	}

	for k, v := range c.Repos {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(wd, v.Path)
		if err != nil {
			relPath = v.Path
		}

		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", k, relPath, v.Branch, len(v.Hooks.Setup))
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}
