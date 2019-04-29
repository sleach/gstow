package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "gstow"
	app.Usage = "gstow is a stow replacement written in Golang"
	app.Author = "Sean Leach <sleach@wiggum.com>"
	app.EnableBashCompletion = true
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Just simulate/print what would happen",
		},
		cli.BoolFlag{
			Name:  "force",
			Usage: "Force overwrite of any existing symlinks",
		},
		cli.StringFlag{
			Name:  "dir, d",
			Usage: "Set gstow dir to DIR (default is current dir)",
		},
		cli.StringFlag{
			Name:  "target, t",
			Usage: "Set target gstow dir to DIR (default is parent of current dir)",
		},
		cli.StringSliceFlag{
			Name:  "ignore",
			Usage: "List of regular expressions of files to ignore",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "stow",
			Aliases: []string{"s"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				f1 := os.Args[1]
				f2 := os.Args[2]
				fs1, _ := os.Stat(f1)
				fs2, _ := os.Stat(f2)
				fmt.Printf("%s == %s => %v\n", f1, f2, os.SameFile(fs1, fs2))
				return nil
			},
		},
		{
			Name:    "unstow",
			Aliases: []string{"u"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				fmt.Println("unstowed: ", c.Args().First())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
