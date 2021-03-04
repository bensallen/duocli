package main

import (
	"log"
	"os"

	"github.com/bensallen/duocli/pkg/cli/docs"
	"github.com/bensallen/duocli/pkg/cli/user"
	"github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {

	app := &cli.App{

		Name:                 "duocli",
		Usage:                "CLI Interface to Duo Admin API",
		Version:              version,
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "load configuration from `FILE`",
				DefaultText: ".duocli.json",
			},
		},
		Commands: []*cli.Command{
			{
				Name:            "user",
				Usage:           "manage users",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
					{
						Name:   "add",
						Usage:  "add a user",
						Action: user.Add,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "username", Aliases: []string{"u"}, Required: true, Usage: "username"},
							&cli.StringSliceFlag{Name: "group", Aliases: []string{"g"}, Usage: "add user to group, can be specified multiple times to add user to multiple groups"},
							&cli.StringFlag{Name: "email", Aliases: []string{"e"}, Usage: "email address of user"},
							&cli.StringFlag{Name: "firstName", Aliases: []string{"f"}, Usage: "first name of user"},
							&cli.StringFlag{Name: "lastName", Aliases: []string{"l"}, Usage: "last name of user"},
							&cli.StringFlag{Name: "status", Aliases: []string{"s"}, Usage: "status of user: active, disabled, or bypass", Value: "active"},
						},
					},
					{
						Name:   "get",
						Usage:  "get one or more users and display as JSON",
						Action: user.Get,
						Flags: []cli.Flag{
							&cli.StringSliceFlag{Name: "username", Aliases: []string{"u"}, Required: true, Usage: "username, can be specified multiple times"},
						},
					},
					{
						Name:   "modify",
						Usage:  "modify a user's attributes, add or remove group membership",
						Action: user.Modify,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "username", Aliases: []string{"u"}, Required: true, Usage: "username"},
							&cli.StringSliceFlag{Name: "addgroup", Aliases: []string{"g"}, Usage: "add user to groups, adds to existing memberships, and can be specified multiple times to add user to multiple groups"},
							&cli.StringSliceFlag{Name: "delgroup", Aliases: []string{"G"}, Usage: "remove user from groups, removes from existing memberships, and can be specified multiple times to remove user from multiple groups"},
							&cli.StringFlag{Name: "email", Aliases: []string{"e"}, Usage: "email address of user"},
							&cli.StringFlag{Name: "firstName", Aliases: []string{"f"}, Usage: "first name of user"},
							&cli.StringFlag{Name: "lastName", Aliases: []string{"l"}, Usage: "last name of user"},
							&cli.StringFlag{Name: "status", Aliases: []string{"s"}, Usage: "status of user: active, disabled, or bypass"},
						},
					},
					{
						Name:   "remove",
						Usage:  "remove user and any attached phones or tokens",
						Action: user.Remove,
						Flags: []cli.Flag{
							&cli.StringSliceFlag{Name: "username", Aliases: []string{"u"}, Required: true, Usage: "username, can be specified multiple times"},
							&cli.BoolFlag{Name: "phone", Aliases: []string{"P"}, Usage: "remove any phones found attached to the user before removing the user", Value: true},
						},
					},
				},
			},
			{
				Name:   "docs",
				Hidden: true,
				Action: docs.Generate,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Required: true, Usage: "path to write man page"},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Error, %v", err)
	}
}
