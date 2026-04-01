package init

import (
	"fmt"
	"os"

	"workspace/internal/config"
)

func initialize() error {
	path, err := config.ConfigPath()
	if err != nil {
		return err
	}

	f, err := os.Stat(path)

	if f != nil {
		fmt.Printf("Workspace already initialized.\n")
		return nil
	}

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	c := config.New()
	err = c.Save()
	if err != nil {
		return err
	}

	fmt.Printf("Workspace initialized successfully.\n")

	return nil
}
