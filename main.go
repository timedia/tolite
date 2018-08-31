package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"matilda.basement.timedia.co.jp/tolite/tolite"
	"matilda.basement.timedia.co.jp/tolite/version"
)

func main() {
	app := cli.NewApp()

	app.Name = "Tolite"
	app.Usage = "gitolite.conf management tool"
	app.Version = version.String()

	app.Commands = []cli.Command{
		{
			Name:  "convert",
			Usage: "pipe gitolite.yml to gitolite.conf.",
			Action: func(c *cli.Context) error {
				a := c.Args().Get(1)
				v, err := tolite.ParseYaml([]byte(a))
				if err != nil {
					return err
				}
				fmt.Printf(tolite.GenerateConf(v))
				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update gitolite.conf by gitolite.yml.",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "users",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add new user",
					Action: func(c *cli.Context) error {
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove user",
					Action: func(c *cli.Context) error {
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list user",
					Action: func(c *cli.Context) error {
						return nil
					},
				},
				{
					Name:  "update",
					Usage: "update user data",
					Action: func(c *cli.Context) error {
						return nil
					},
				},
			},
		},
		{
			Name: "groups",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "repos",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "admin_repos",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
