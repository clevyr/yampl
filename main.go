package main

import (
	"github.com/clevyr/go-yampl/cmd"
	"os"
)

//go:generate git config core.hooksPath .githooks

func main() {
	if err := cmd.Command.Execute(); err != nil {
		os.Exit(1)
	}
}
