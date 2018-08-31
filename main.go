package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli"
	"matilda.basement.timedia.co.jp/tolite/tolitelib"
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
				a := c.Args().Get(0)
				fmt.Println(a)
				v, err := tolitelib.ParseYaml([]byte(a))
				if err != nil {
					return err
				}
				fmt.Printf(tolitelib.GenerateConf(v))
				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update gitolite.conf by gitolite.yml.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input",
					Value: "gitolite.yml",
					Usage: "input yml file",
				},
				cli.StringFlag{
					Name:  "output",
					Value: "gitolite.conf",
					Usage: "output conf file",
				},
			},
			Action: func(c *cli.Context) error {
				buf, err := ioutil.ReadFile(c.String("input"))
				v, err4 := tolitelib.ParseYaml(buf)
				if err4 != nil {
					return err
				}
				outfile, err2 := os.Create(c.String("output"))
				if err2 != nil {
					log.Fatal(err)
				}
				outfile.Write([]byte(tolitelib.GenerateConf(v)))
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
