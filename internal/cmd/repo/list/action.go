package list

import (
	"fmt"
	"os"
	"text/tabwriter"

	"workspace/internal/config"
)

func list() error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, err = fmt.Fprintf(w, "ALIAS\tREPO\tBRANCH\tHOOKS\n")
	if err != nil {
		return err
	}

	for k, v := range c.Repos {
		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", k, v.Path, v.Branch, len(v.Hooks.Setup))
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
