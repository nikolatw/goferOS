package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "git",
		Description: "Pure go git implementation",
		Commands: []*cli.Command{
			{
				Name:  "clone",
				Usage: "Clone a repository into a new directory",
				Action: func(cCtx *cli.Context) error {
					_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
						URL:      cCtx.Args().First(),
						Progress: os.Stdout,
					})

					return err
				},
			},
			{
				Name:  "checkout",
				Usage: "Switch branches or restore working tree files",
				Action: func(cCtx *cli.Context) error {
					ref := plumbing.NewBranchReferenceName(cCtx.Args().First())

					rep, err := git.PlainOpen(".")
					if err != nil {
						return err
					}
					fullRef, err := rep.Reference(ref, true)
					if err != nil {
						return err
					}
					wt, err := rep.Worktree()
					if err != nil {
						return err
					}
					err = wt.Checkout(&git.CheckoutOptions{
						Hash:  fullRef.Hash(),
						Force: true,
						Keep:  false,
					})

					return err
				},
			},
			{
				Name:  "add",
				Usage: "Add file contents to the index",
				Action: func(cCtx *cli.Context) error {
					rep, err := git.PlainOpen(".")
					if err != nil {
						return err
					}

					wt, err := rep.Worktree()
					if err != nil {
						return err
					}

					_, err = wt.Add(cCtx.Args().First())
					if err != nil {
						return err
					}

					return err
				},
			},
			{
				Name:  "commit",
				Usage: "Record changes to the repository",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "message",
						Aliases:  []string{"m"},
						Usage:    "commit message",
						Required: true,
					},
					&cli.BoolFlag{
						Name:     "all",
						Aliases:  []string{"A"},
						Usage:    "adds, modifies, and removes index entries to match the working tree",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					rep, err := git.PlainOpen(".")
					if err != nil {
						return err
					}

					wt, err := rep.Worktree()
					if err != nil {
						return err
					}

					_, err = wt.Commit(cCtx.String("message"), &git.CommitOptions{
						All: cCtx.Bool("all"),
					})
					if err != nil {
						return err
					}

					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
