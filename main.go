package main

import (
	"errors"
	"fmt"
	"os"
)

var subcommands = map[string]Command{
	"encode": NewEncodeCommand(),
	"decode": NewDecodeCommand(),
}

type Command interface {
	Init([]string) error
	Run() error
	Name() string
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass a sub-command")
	}

	subcommand := os.Args[1]

	if cmd, ok := subcommands[subcommand]; ok {
		err := cmd.Init(os.Args[2:])
		if err != nil {
			return err
		}

		return cmd.Run()
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}
