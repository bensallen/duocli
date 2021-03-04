package docs

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// Generate writes out man page for the application to the provide path
func Generate(c *cli.Context) error {
	path := c.String("output")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	man, err := c.App.ToMan()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "%s", man)
	return err
}
