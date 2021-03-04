package group

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Add(c *cli.Context) error {

	groupNames := c.StringSlice("name")
	fmt.Printf("%#v\n", groupNames)

	return nil
}

func Remove(c *cli.Context) error {

	groupNames := c.StringSlice("name")
	fmt.Printf("%#v\n", groupNames)

	return nil
}
