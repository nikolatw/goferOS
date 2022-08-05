// Copyright (c) 2017, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
	"time"

	"golang.org/x/term"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

var command = flag.String("c", "", "command to be executed")

func main() {
	flag.Parse()
	err := runAll()
	if e, ok := interp.IsExitStatus(err); ok {
		os.Exit(int(e))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runAll() error {
	r, err := interp.New(interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		return err
	}

	if *command != "" {
		return run(r, strings.NewReader(*command), "")
	}
	if flag.NArg() == 0 {
		if term.IsTerminal(int(os.Stdin.Fd())) {
			return runInteractive(r, os.Stdin, os.Stdout, os.Stderr)
		}
		return run(r, os.Stdin, "")
	}
	for _, path := range flag.Args() {
		if err := runPath(r, path); err != nil {
			return err
		}
	}
	return nil
}

func run(r *interp.Runner, reader io.Reader, name string) error {
	prog, err := syntax.NewParser().Parse(reader, name)
	if err != nil {
		return err
	}
	r.Reset()
	ctx := context.Background()
	return r.Run(ctx, prog)
}

func runPath(r *interp.Runner, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return run(r, f, path)
}

func runInteractive(r *interp.Runner, stdin io.Reader, stdout, stderr io.Writer) error {
	parser := syntax.NewParser()
	prompt(r, stdout, false)
	var runErr error
	fn := func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			prompt(r, stdout, true)
			return true
		}
		ctx := context.Background()
		for _, stmt := range stmts {
			runErr = r.Run(ctx, stmt)
			if r.Exited() {
				return false
			}
		}
		prompt(r, stdout, false)
		return true
	}
	if err := parser.Interactive(stdin, fn); err != nil {
		return err
	}
	return runErr
}

func prompt(r *interp.Runner, stdout io.Writer, incomplete bool) {
	if incomplete {
		fmt.Fprint(stdout, "> ")
	} else {
		username := ""
		user, err := user.Current()
		if err == nil {
			username = user.Username
		}

		hostname, _ := os.Hostname()

		fmt.Fprintf(stdout, "[\033[1;34m%s@%s\033[0m \033[0;36m%s\033[0m] \033[0;30m(%s)\033[0m $ ", username, hostname, r.Dir, time.Now().Format("15:04:05"))
	}
}
